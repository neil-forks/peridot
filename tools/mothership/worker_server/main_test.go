// Copyright 2023 Peridot Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mothership_worker_server

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	base "go.resf.org/peridot/base/go"
	storage_memory "go.resf.org/peridot/base/go/storage/memory"
	mothership_db "go.resf.org/peridot/tools/mothership/db"
	mothership_migrations "go.resf.org/peridot/tools/mothership/migrations"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/testsuite"
	"golang.org/x/crypto/openpgp"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	testW        *Worker
	testWRolling *Worker
	//go:embed testdata/RPM-GPG-KEY-Rocky-8
	rocky8GpgKey []byte
	inmf         *inMemoryForge
	tempDirForge string
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

type noopLogger struct{}

func (n *noopLogger) Debug(string, ...any) {}
func (n *noopLogger) Info(string, ...any)  {}
func (n *noopLogger) Warn(string, ...any)  {}
func (n *noopLogger) Error(string, ...any) {}
func (n *noopLogger) With(...any) log.Logger {
	return n
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func TestUnitTestSuite(t *testing.T) {
	ts := new(UnitTestSuite)
	ts.SetLogger(&noopLogger{})
	suite.Run(t, ts)
}

func TestMain(m *testing.M) {
	// Create temporary file
	dir, err := os.MkdirTemp("", "test-db-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	scripts, err := base.EmbedFSToOSFS(dir, mothership_migrations.UpSQLs, ".")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(scripts...),
		postgres.WithDatabase("mshiptest"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.
				ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		panic(err)
	}
	defer pgContainer.Terminate(ctx)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	db, err := base.NewDB(connStr)
	if err != nil {
		panic(err)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	lookasideFS := osfs.New("/")
	inMemStorage := storage_memory.New(lookasideFS, filepath.Join(cwd, "testdata"))

	var gpgKeys openpgp.EntityList
	keyRing, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(rocky8GpgKey))
	if err != nil {
		panic(err)
	}

	gpgKeys = append(gpgKeys, keyRing...)

	tempDirForge, err = os.MkdirTemp("", "test-forge-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDirForge)

	inmf = &inMemoryForge{
		remoteBaseURL: "https://testforge.resf.org",
		localTempDir:  tempDirForge,
		repos:         map[string]bool{},
	}
	testW = New(db, inMemStorage, gpgKeys, inmf, false)
	testWRolling = New(db, inMemStorage, gpgKeys, inmf, true)

	if err := q[mothership_db.Worker]().Create(&mothership_db.Worker{
		Name:      base.NameGen("workers"),
		WorkerID:  "test-worker",
		ApiSecret: "test-secret",
	}); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func q[T any]() base.Pika[T] {
	return base.Q[T](testW.db)
}