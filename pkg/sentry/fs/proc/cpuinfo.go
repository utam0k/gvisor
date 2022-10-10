// Copyright 2018 The gVisor Authors.
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

package proc

import (
	"bytes"

	"github.com/utam0k/gvisor/pkg/context"
	"github.com/utam0k/gvisor/pkg/sentry/fs"
	"github.com/utam0k/gvisor/pkg/sentry/kernel"
)

// LINT.IfChange

func newCPUInfo(ctx context.Context, msrc *fs.MountSource) *fs.Inode {
	k := kernel.KernelFromContext(ctx)
	features := k.FeatureSet()
	var buf bytes.Buffer
	for i, max := uint(0), k.ApplicationCores(); i < max; i++ {
		features.WriteCPUInfoTo(i, &buf)
	}
	return newStaticProcInode(ctx, msrc, buf.Bytes())
}

// LINT.ThenChange(../../fsimpl/proc/tasks.go)
