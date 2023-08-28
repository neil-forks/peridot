package github_forge

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	transport_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"go.resf.org/peridot/tools/mothership/worker_server/forge"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type Forge struct {
	organization         string
	appId                string
	appPrivateKey        *rsa.PrivateKey
	shouldMakeRepoPublic bool
}

type installationToken struct {
	Token   string
	AppSlug string
}

func New(organization string, appId string, appPrivateKey []byte, shouldMakeRepoPublic bool) (*Forge, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(appPrivateKey)
	if err != nil {
		return nil, err
	}

	return &Forge{
		organization:         organization,
		appId:                appId,
		appPrivateKey:        privateKey,
		shouldMakeRepoPublic: shouldMakeRepoPublic,
	}, nil
}

func (f *Forge) getInstallationToken(jwt string) (*installationToken, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	endpoint := "https://api.github.com/" + filepath.Join("orgs", f.organization, "installation")
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)

	// We can just marshal into a map as we only need the first "id" field (int value)
	respBody := map[string]any{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, err
	}

	installationIdInt := respBody["id"].(float64)
	installationId := strconv.FormatFloat(installationIdInt, 'f', 0, 64)
	appSlug := respBody["app_slug"].(string)

	// Get the installation token
	endpoint = "https://api.github.com/" + filepath.Join("app/installations", installationId, "access_tokens")
	req, err = http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)

	respBody = map[string]any{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	// We need the token string
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, err
	}

	return &installationToken{
		Token:   respBody["token"].(string),
		AppSlug: appSlug,
	}, nil
}

func (f *Forge) GetMeID(installation string, appSlug string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	endpoint := "https://api.github.com/" + filepath.Join("users", fmt.Sprintf("%s[bot]", appSlug))
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+installation)

	respBody := map[string]any{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// We need the token string
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", err
	}

	idInt := respBody["id"].(float64)
	id := strconv.FormatFloat(idInt, 'f', 0, 64)

	return id, nil
}

func (f *Forge) GetAuthenticator() (*forge.Authenticator, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		"iss": f.appId,
		"alg": "RS256",
	})

	tokenString, err := token.SignedString(f.appPrivateKey)
	if err != nil {
		return nil, err
	}

	installation, err := f.getInstallationToken(tokenString)
	if err != nil {
		return nil, err
	}
	// It expires in an hour, but we'll just use 45 minutes
	expires := time.Now().Add(time.Minute * 45)

	meID, err := f.GetMeID(installation.Token, installation.AppSlug)
	if err != nil {
		return nil, err
	}

	transporter := &transport_http.BasicAuth{
		Username: f.appId,
		Password: installation.Token,
	}

	return &forge.Authenticator{
		AuthMethod:  transporter,
		AuthorName:  fmt.Sprintf("%s[bot]", installation.AppSlug),
		AuthorEmail: fmt.Sprintf("%s+%s[bot]@users.noreply.github.com", meID, installation.AppSlug),
		Expires:     expires,
	}, nil
}

func (f *Forge) GetRemote() string {
	return fmt.Sprintf("https://github.com/%s", f.organization)
}

func (f *Forge) EnsureRepositoryExists(auth *forge.Authenticator, repo string) error {
	// Cast AuthMethod to BasicAuth
	basicAuth := auth.AuthMethod.(*transport_http.BasicAuth)
	token := basicAuth.Password

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// First let's check if the repo exists
	endpoint := "https://api.github.com/" + filepath.Join("repos", f.organization, repo)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		// Repo exists, we're done
		return nil
	}

	mapBody := map[string]any{
		"name":         repo,
		"private":      !f.shouldMakeRepoPublic,
		"has_issues":   false,
		"has_projects": false,
		"has_wiki":     false,
	}
	body, err := json.Marshal(mapBody)
	if err != nil {
		return err
	}

	endpoint = "https://api.github.com/" + filepath.Join("orgs", f.organization, "repos")
	req, err = http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("got status code %d", resp.StatusCode)
	}

	return nil
}