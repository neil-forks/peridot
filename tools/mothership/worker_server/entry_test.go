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
	"github.com/stretchr/testify/require"
	mothership_db "go.resf.org/peridot/tools/mothership/db"
	mothershippb "go.resf.org/peridot/tools/mothership/pb"
	"testing"
	"time"
)

func TestWorker_CreateEntry(t *testing.T) {
	require.Nil(t, q[mothership_db.Entry]().Delete())
	defer func() {
		require.Nil(t, q[mothership_db.Entry]().Delete())
	}()

	args := &mothershippb.ProcessRPMArgs{
		Request: &mothershippb.ProcessRPMRequest{
			RpmUri:     "memory://efi-rpm-macros-3-3.el8.src.rpm",
			OsRelease:  "Rocky Linux release 8.8 (Green Obsidian)",
			Checksum:   "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
			Repository: "BaseOS",
		},
		InternalRequest: &mothershippb.ProcessRPMInternalRequest{
			WorkerId: "test-worker",
		},
	}
	entry, err := testW.CreateEntry(args)
	require.Nil(t, err)
	require.NotNil(t, entry)
	require.Equal(t, "Rocky Linux release 8.8 (Green Obsidian)", entry.OsRelease)
	require.Equal(t, "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28", entry.Sha256Sum)
	c, err := q[mothership_db.Entry]().F("name", entry.Name).Count()
	require.Nil(t, err)
	require.Equal(t, c, 1)
}

func TestWorker_SetEntryIDFromRPM(t *testing.T) {
	require.Nil(t, q[mothership_db.Entry]().Delete())
	defer func() {
		require.Nil(t, q[mothership_db.Entry]().Delete())
	}()

	args := &mothershippb.ProcessRPMArgs{
		Request: &mothershippb.ProcessRPMRequest{
			RpmUri:     "memory://efi-rpm-macros-3-3.el8.src.rpm",
			OsRelease:  "Rocky Linux release 8.8 (Green Obsidian)",
			Checksum:   "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
			Repository: "BaseOS",
		},
		InternalRequest: &mothershippb.ProcessRPMInternalRequest{
			WorkerId: "test-worker",
		},
	}
	entry, err := testW.CreateEntry(args)
	require.Nil(t, err)
	require.NotNil(t, entry)

	entry, err = testW.SetEntryIDFromRPM(entry.Name, "memory://efi-rpm-macros-3-3.el8.src.rpm", entry.Sha256Sum)
	require.Nil(t, err)
	require.NotNil(t, entry)
	require.Equal(t, "efi-rpm-macros-3-3.el8.src", entry.EntryId)
}

func TestWorker_SetEntryIDFromRPM_FailedToDownload(t *testing.T) {
	require.Nil(t, q[mothership_db.Entry]().Delete())
	defer func() {
		require.Nil(t, q[mothership_db.Entry]().Delete())
	}()

	args := &mothershippb.ProcessRPMArgs{
		Request: &mothershippb.ProcessRPMRequest{
			RpmUri:     "memory://not-found.rpm",
			OsRelease:  "Rocky Linux release 8.8 (Green Obsidian)",
			Checksum:   "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
			Repository: "BaseOS",
		},
		InternalRequest: &mothershippb.ProcessRPMInternalRequest{
			WorkerId: "test-worker",
		},
	}
	entry, err := testW.CreateEntry(args)
	require.Nil(t, err)
	require.NotNil(t, entry)

	entry, err = testW.SetEntryIDFromRPM(entry.Name, "memory://not-found.rpm", entry.Sha256Sum)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "failed to download resource")
}

func TestWorker_SetEntryState(t *testing.T) {
	require.Nil(t, q[mothership_db.Entry]().Delete())
	defer func() {
		require.Nil(t, q[mothership_db.Entry]().Delete())
	}()

	args := &mothershippb.ProcessRPMArgs{
		Request: &mothershippb.ProcessRPMRequest{
			RpmUri:     "memory://efi-rpm-macros-3-3.el8.src.rpm",
			OsRelease:  "Rocky Linux release 8.8 (Green Obsidian)",
			Checksum:   "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
			Repository: "BaseOS",
		},
		InternalRequest: &mothershippb.ProcessRPMInternalRequest{
			WorkerId: "test-worker",
		},
	}
	entry, err := testW.CreateEntry(args)
	require.Nil(t, err)
	require.NotNil(t, entry)

	importRpmRes := &mothershippb.ImportRPMResponse{
		CommitHash:   "123",
		CommitUri:    "https://testforge.resf.org/peridot/efi-rpm-macros/commit/123",
		CommitBranch: "el-8.8",
		CommitTag:    "imports/el-8.8/efi-rpm-macros-3-3.el8",
		Nevra:        "efi-rpm-macros-0:3-3.el8.aarch64",
		Pkg:          "efi-rpm-macros",
	}
	entry, err = testW.SetEntryState(entry.Name, mothershippb.Entry_ARCHIVED, importRpmRes)
	require.Nil(t, err)
	require.NotNil(t, entry)
	require.Equal(t, mothershippb.Entry_ARCHIVED, entry.State)
	require.Equal(t, "123", entry.CommitHash)
	require.Equal(t, "https://testforge.resf.org/peridot/efi-rpm-macros/commit/123", entry.CommitUri)
	require.Equal(t, "el-8.8", entry.CommitBranch)
	require.Equal(t, "imports/el-8.8/efi-rpm-macros-3-3.el8", entry.CommitTag)
	require.Equal(t, "efi-rpm-macros", entry.Pkg)
}

func TestWorker_SetEntryState_NoRes(t *testing.T) {
	require.Nil(t, q[mothership_db.Entry]().Delete())
	defer func() {
		require.Nil(t, q[mothership_db.Entry]().Delete())
	}()

	args := &mothershippb.ProcessRPMArgs{
		Request: &mothershippb.ProcessRPMRequest{
			RpmUri:     "memory://efi-rpm-macros-3-3.el8.src.rpm",
			OsRelease:  "Rocky Linux release 8.8 (Green Obsidian)",
			Checksum:   "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
			Repository: "BaseOS",
		},
		InternalRequest: &mothershippb.ProcessRPMInternalRequest{
			WorkerId: "test-worker",
		},
	}
	entry, err := testW.CreateEntry(args)
	require.Nil(t, err)
	require.NotNil(t, entry)

	entry, err = testW.SetEntryState(entry.Name, mothershippb.Entry_ON_HOLD, nil)
	require.Nil(t, err)
	require.NotNil(t, entry)
	require.Equal(t, mothershippb.Entry_ON_HOLD, entry.State)
	require.Equal(t, "", entry.CommitHash)
	require.Equal(t, "", entry.CommitUri)
	require.Equal(t, "", entry.CommitBranch)
	require.Equal(t, "", entry.CommitTag)
	require.Equal(t, "", entry.Pkg)
}

func TestWorker_SetEntryState_NoEntry(t *testing.T) {
	require.Nil(t, q[mothership_db.Entry]().Delete())
	defer func() {
		require.Nil(t, q[mothership_db.Entry]().Delete())
	}()

	entry, err := testW.SetEntryState("entries/123", mothershippb.Entry_ON_HOLD, nil)
	require.Nil(t, entry)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "entry does not exist")
}

func TestWorker_SetWorkerLastCheckinTime(t *testing.T) {
	require.Nil(t, testW.SetWorkerLastCheckinTime("test-worker"))
	// Verify that the worker last checkin time is at most 15 seconds ago.
	w, err := q[mothership_db.Worker]().F("worker_id", "test-worker").GetOrNil()
	require.Nil(t, err)
	require.NotNil(t, w)
	require.True(t, w.LastCheckinTime.Valid)
	require.WithinDuration(t, w.LastCheckinTime.Time, time.Now(), 15*time.Second)
}

func TestWorker_SetWorkerLastCheckinTime_NotFound(t *testing.T) {
	err := testW.SetWorkerLastCheckinTime("not-found")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "worker does not exist")
}