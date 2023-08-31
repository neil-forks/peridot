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

package mothership_rpc

import (
	"context"
	"database/sql"
	base "go.resf.org/peridot/base/go"
	mothership_db "go.resf.org/peridot/tools/mothership/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (s *Server) WorkerPing(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	worker, err := s.getWorkerIdentity(ctx)
	if err != nil {
		return nil, err
	}

	worker.LastCheckinTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if err := base.Q[mothership_db.Worker](s.db).U(worker); err != nil {
		return nil, status.Error(codes.Internal, "failed to update worker")
	}

	return &emptypb.Empty{}, nil
}