// Copyright 2019 The gVisor Authors.
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

package testutil

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/utam0k/gvisor/pkg/abi/linux"
	"github.com/utam0k/gvisor/pkg/context"
	"github.com/utam0k/gvisor/pkg/cpuid"
	"github.com/utam0k/gvisor/pkg/fspath"
	"github.com/utam0k/gvisor/pkg/memutil"
	"github.com/utam0k/gvisor/pkg/sentry/fsbridge"
	"github.com/utam0k/gvisor/pkg/sentry/fsimpl/tmpfs"
	"github.com/utam0k/gvisor/pkg/sentry/kernel"
	"github.com/utam0k/gvisor/pkg/sentry/kernel/auth"
	"github.com/utam0k/gvisor/pkg/sentry/kernel/sched"
	"github.com/utam0k/gvisor/pkg/sentry/limits"
	"github.com/utam0k/gvisor/pkg/sentry/loader"
	"github.com/utam0k/gvisor/pkg/sentry/mm"
	"github.com/utam0k/gvisor/pkg/sentry/pgalloc"
	"github.com/utam0k/gvisor/pkg/sentry/platform"
	"github.com/utam0k/gvisor/pkg/sentry/time"
	"github.com/utam0k/gvisor/pkg/sentry/vfs"

	// Platforms are plugable.
	_ "github.com/utam0k/gvisor/pkg/sentry/platform/kvm"
	_ "github.com/utam0k/gvisor/pkg/sentry/platform/ptrace"
)

var (
	platformFlag           = flag.String("platform", "ptrace", "specify which platform to use")
	platformDevicePathFlag = flag.String("platform_device_path", "", "path to a platform-specific device file (e.g. /dev/kvm for KVM platform). If unset, will use a sane platform-specific default.")
)

// Boot initializes a new bare bones kernel for test.
func Boot() (*kernel.Kernel, error) {
	platformCtr, err := platform.Lookup(*platformFlag)
	if err != nil {
		return nil, fmt.Errorf("platform not found: %v", err)
	}
	deviceFile, err := platformCtr.OpenDevice(*platformDevicePathFlag)
	if err != nil {
		return nil, fmt.Errorf("creating platform: %v", err)
	}
	plat, err := platformCtr.New(deviceFile)
	if err != nil {
		return nil, fmt.Errorf("creating platform: %v", err)
	}

	kernel.VFS2Enabled = true
	k := &kernel.Kernel{
		Platform: plat,
	}

	mf, err := createMemoryFile()
	if err != nil {
		return nil, err
	}
	k.SetMemoryFile(mf)

	// Pass k as the platform since it is savable, unlike the actual platform.
	vdso, err := loader.PrepareVDSO(k)
	if err != nil {
		return nil, fmt.Errorf("creating vdso: %v", err)
	}

	// Create timekeeper.
	tk := kernel.NewTimekeeper(k, vdso.ParamPage.FileRange())
	tk.SetClocks(time.NewCalibratedClocks())

	creds := auth.NewRootCredentials(auth.NewRootUserNamespace())

	// Initiate the Kernel object, which is required by the Context passed
	// to createVFS in order to mount (among other things) procfs.
	if err = k.Init(kernel.InitKernelArgs{
		ApplicationCores:            uint(runtime.GOMAXPROCS(-1)),
		FeatureSet:                  cpuid.HostFeatureSet(),
		Timekeeper:                  tk,
		RootUserNamespace:           creds.UserNamespace,
		Vdso:                        vdso,
		RootUTSNamespace:            kernel.NewUTSNamespace("hostname", "domain", creds.UserNamespace),
		RootIPCNamespace:            kernel.NewIPCNamespace(creds.UserNamespace),
		RootAbstractSocketNamespace: kernel.NewAbstractSocketNamespace(),
		PIDNamespace:                kernel.NewRootPIDNamespace(creds.UserNamespace),
	}); err != nil {
		return nil, fmt.Errorf("initializing kernel: %v", err)
	}

	k.VFS().MustRegisterFilesystemType(tmpfs.Name, &tmpfs.FilesystemType{}, &vfs.RegisterFilesystemTypeOptions{
		AllowUserMount: true,
		AllowUserList:  true,
	})

	ls, err := limits.NewLinuxLimitSet()
	if err != nil {
		return nil, err
	}
	tg := k.NewThreadGroup(nil, k.RootPIDNamespace(), kernel.NewSignalHandlers(), linux.SIGCHLD, ls)
	k.TestOnlySetGlobalInit(tg)

	return k, nil
}

// CreateTask creates a new bare bones task for tests.
func CreateTask(ctx context.Context, name string, tc *kernel.ThreadGroup, mntns *vfs.MountNamespace, root, cwd vfs.VirtualDentry) (*kernel.Task, error) {
	k := kernel.KernelFromContext(ctx)
	if k == nil {
		return nil, fmt.Errorf("cannot find kernel from context")
	}

	exe, err := newFakeExecutable(ctx, k.VFS(), auth.CredentialsFromContext(ctx), root)
	if err != nil {
		return nil, err
	}
	m := mm.NewMemoryManager(k, k, k.SleepForAddressSpaceActivation)
	m.SetExecutable(ctx, fsbridge.NewVFSFile(exe))

	creds := auth.CredentialsFromContext(ctx)
	config := &kernel.TaskConfig{
		Kernel:                  k,
		ThreadGroup:             tc,
		TaskImage:               &kernel.TaskImage{Name: name, MemoryManager: m},
		Credentials:             auth.CredentialsFromContext(ctx),
		NetworkNamespace:        k.RootNetworkNamespace(),
		AllowedCPUMask:          sched.NewFullCPUSet(k.ApplicationCores()),
		UTSNamespace:            kernel.UTSNamespaceFromContext(ctx),
		IPCNamespace:            kernel.IPCNamespaceFromContext(ctx),
		AbstractSocketNamespace: kernel.NewAbstractSocketNamespace(),
		MountNamespaceVFS2:      mntns,
		FSContext:               kernel.NewFSContextVFS2(root, cwd, 0022),
		FDTable:                 k.NewFDTable(),
		UserCounters:            k.GetUserCounters(creds.RealKUID),
	}
	config.NetworkNamespace.IncRef()
	t, err := k.TaskSet().NewTask(ctx, config)
	if err != nil {
		config.ThreadGroup.Release(ctx)
		return nil, err
	}
	return t, nil
}

func newFakeExecutable(ctx context.Context, vfsObj *vfs.VirtualFilesystem, creds *auth.Credentials, root vfs.VirtualDentry) (*vfs.FileDescription, error) {
	const name = "executable"
	pop := &vfs.PathOperation{
		Root:  root,
		Start: root,
		Path:  fspath.Parse(name),
	}
	opts := &vfs.OpenOptions{
		Flags: linux.O_RDONLY | linux.O_CREAT,
		Mode:  0777,
	}
	return vfsObj.OpenAt(ctx, creds, pop, opts)
}

func createMemoryFile() (*pgalloc.MemoryFile, error) {
	const memfileName = "test-memory"
	memfd, err := memutil.CreateMemFD(memfileName, 0)
	if err != nil {
		return nil, fmt.Errorf("error creating memfd: %v", err)
	}
	memfile := os.NewFile(uintptr(memfd), memfileName)
	mf, err := pgalloc.NewMemoryFile(memfile, pgalloc.MemoryFileOpts{})
	if err != nil {
		_ = memfile.Close()
		return nil, fmt.Errorf("error creating pgalloc.MemoryFile: %v", err)
	}
	return mf, nil
}
