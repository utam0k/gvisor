// Copyright 2020 The gVisor Authors.
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

// Package vfs2 provides syscall implementations that use VFS2.
package vfs2

import (
	"github.com/utam0k/gvisor/pkg/sentry/syscalls"
	"github.com/utam0k/gvisor/pkg/sentry/syscalls/linux"
)

// Override syscall table to add syscalls implementations from this package.
func Override() {
	// Override AMD64.
	s := linux.AMD64
	s.Table[0] = syscalls.SupportedPoint("read", Read, linux.PointRead)
	s.Table[1] = syscalls.Supported("write", Write)
	s.Table[2] = syscalls.SupportedPoint("open", Open, linux.PointOpen)
	s.Table[3] = syscalls.SupportedPoint("close", Close, linux.PointClose)
	s.Table[4] = syscalls.Supported("stat", Stat)
	s.Table[5] = syscalls.Supported("fstat", Fstat)
	s.Table[6] = syscalls.Supported("lstat", Lstat)
	s.Table[7] = syscalls.Supported("poll", Poll)
	s.Table[8] = syscalls.Supported("lseek", Lseek)
	s.Table[9] = syscalls.Supported("mmap", Mmap)
	s.Table[16] = syscalls.Supported("ioctl", Ioctl)
	s.Table[17] = syscalls.Supported("pread64", Pread64)
	s.Table[18] = syscalls.Supported("pwrite64", Pwrite64)
	s.Table[19] = syscalls.Supported("readv", Readv)
	s.Table[20] = syscalls.Supported("writev", Writev)
	s.Table[21] = syscalls.Supported("access", Access)
	s.Table[22] = syscalls.SupportedPoint("pipe", Pipe, linux.PointPipe)
	s.Table[23] = syscalls.Supported("select", Select)
	s.Table[32] = syscalls.SupportedPoint("dup", Dup, linux.PointDup)
	s.Table[33] = syscalls.SupportedPoint("dup2", Dup2, linux.PointDup2)
	s.Table[40] = syscalls.Supported("sendfile", Sendfile)
	s.Table[41] = syscalls.SupportedPoint("socket", Socket, linux.PointSocket)
	s.Table[42] = syscalls.SupportedPoint("connect", Connect, linux.PointConnect)
	s.Table[43] = syscalls.SupportedPoint("accept", Accept, linux.PointAccept)
	s.Table[44] = syscalls.Supported("sendto", SendTo)
	s.Table[45] = syscalls.Supported("recvfrom", RecvFrom)
	s.Table[46] = syscalls.Supported("sendmsg", SendMsg)
	s.Table[47] = syscalls.Supported("recvmsg", RecvMsg)
	s.Table[48] = syscalls.Supported("shutdown", Shutdown)
	s.Table[49] = syscalls.SupportedPoint("bind", Bind, linux.PointBind)
	s.Table[50] = syscalls.Supported("listen", Listen)
	s.Table[51] = syscalls.Supported("getsockname", GetSockName)
	s.Table[52] = syscalls.Supported("getpeername", GetPeerName)
	s.Table[53] = syscalls.SupportedPoint("socketpair", SocketPair, linux.PointSocketpair)
	s.Table[54] = syscalls.Supported("setsockopt", SetSockOpt)
	s.Table[55] = syscalls.Supported("getsockopt", GetSockOpt)
	s.Table[59] = syscalls.SupportedPoint("execve", Execve, linux.PointExecve)
	s.Table[72] = syscalls.SupportedPoint("fcntl", Fcntl, linux.PointFcntl)
	s.Table[73] = syscalls.Supported("flock", Flock)
	s.Table[74] = syscalls.Supported("fsync", Fsync)
	s.Table[75] = syscalls.Supported("fdatasync", Fdatasync)
	s.Table[76] = syscalls.Supported("truncate", Truncate)
	s.Table[77] = syscalls.Supported("ftruncate", Ftruncate)
	s.Table[78] = syscalls.Supported("getdents", Getdents)
	s.Table[79] = syscalls.Supported("getcwd", Getcwd)
	s.Table[80] = syscalls.SupportedPoint("chdir", Chdir, linux.PointChdir)
	s.Table[81] = syscalls.SupportedPoint("fchdir", Fchdir, linux.PointFchdir)
	s.Table[82] = syscalls.Supported("rename", Rename)
	s.Table[83] = syscalls.Supported("mkdir", Mkdir)
	s.Table[84] = syscalls.Supported("rmdir", Rmdir)
	s.Table[85] = syscalls.SupportedPoint("creat", Creat, linux.PointCreat)
	s.Table[86] = syscalls.Supported("link", Link)
	s.Table[87] = syscalls.Supported("unlink", Unlink)
	s.Table[88] = syscalls.Supported("symlink", Symlink)
	s.Table[89] = syscalls.Supported("readlink", Readlink)
	s.Table[90] = syscalls.Supported("chmod", Chmod)
	s.Table[91] = syscalls.Supported("fchmod", Fchmod)
	s.Table[92] = syscalls.Supported("chown", Chown)
	s.Table[93] = syscalls.Supported("fchown", Fchown)
	s.Table[94] = syscalls.Supported("lchown", Lchown)
	s.Table[132] = syscalls.Supported("utime", Utime)
	s.Table[133] = syscalls.Supported("mknod", Mknod)
	s.Table[137] = syscalls.Supported("statfs", Statfs)
	s.Table[138] = syscalls.Supported("fstatfs", Fstatfs)
	s.Table[155] = syscalls.Supported("pivot_root", PivotRoot)
	s.Table[161] = syscalls.SupportedPoint("chroot", Chroot, linux.PointChroot)
	s.Table[162] = syscalls.Supported("sync", Sync)
	s.Table[165] = syscalls.Supported("mount", Mount)
	s.Table[166] = syscalls.Supported("umount2", Umount2)
	s.Table[187] = syscalls.Supported("readahead", Readahead)
	s.Table[188] = syscalls.Supported("setxattr", SetXattr)
	s.Table[189] = syscalls.Supported("lsetxattr", Lsetxattr)
	s.Table[190] = syscalls.Supported("fsetxattr", Fsetxattr)
	s.Table[191] = syscalls.Supported("getxattr", GetXattr)
	s.Table[192] = syscalls.Supported("lgetxattr", Lgetxattr)
	s.Table[193] = syscalls.Supported("fgetxattr", Fgetxattr)
	s.Table[194] = syscalls.Supported("listxattr", ListXattr)
	s.Table[195] = syscalls.Supported("llistxattr", Llistxattr)
	s.Table[196] = syscalls.Supported("flistxattr", Flistxattr)
	s.Table[197] = syscalls.Supported("removexattr", RemoveXattr)
	s.Table[198] = syscalls.Supported("lremovexattr", Lremovexattr)
	s.Table[199] = syscalls.Supported("fremovexattr", Fremovexattr)
	s.Table[209] = syscalls.PartiallySupported("io_submit", IoSubmit, "Generally supported with exceptions. User ring optimizations are not implemented.", []string{"gvisor.dev/issue/204"})
	s.Table[213] = syscalls.Supported("epoll_create", EpollCreate)
	s.Table[217] = syscalls.Supported("getdents64", Getdents64)
	s.Table[221] = syscalls.PartiallySupported("fadvise64", Fadvise64, "The syscall is 'supported', but ignores all provided advice.", nil)
	s.Table[232] = syscalls.Supported("epoll_wait", EpollWait)
	s.Table[233] = syscalls.Supported("epoll_ctl", EpollCtl)
	s.Table[235] = syscalls.Supported("utimes", Utimes)
	s.Table[240] = syscalls.Supported("mq_open", MqOpen)
	s.Table[241] = syscalls.Supported("mq_unlink", MqUnlink)
	s.Table[253] = syscalls.PartiallySupportedPoint("inotify_init", InotifyInit, linux.PointInotifyInit, "inotify events are only available inside the sandbox.", nil)
	s.Table[254] = syscalls.PartiallySupportedPoint("inotify_add_watch", InotifyAddWatch, linux.PointInotifyAddWatch, "inotify events are only available inside the sandbox.", nil)
	s.Table[255] = syscalls.PartiallySupportedPoint("inotify_rm_watch", InotifyRmWatch, linux.PointInotifyRmWatch, "inotify events are only available inside the sandbox.", nil)
	s.Table[257] = syscalls.SupportedPoint("openat", Openat, linux.PointOpenat)
	s.Table[258] = syscalls.Supported("mkdirat", Mkdirat)
	s.Table[259] = syscalls.Supported("mknodat", Mknodat)
	s.Table[260] = syscalls.Supported("fchownat", Fchownat)
	s.Table[261] = syscalls.Supported("futimesat", Futimesat)
	s.Table[262] = syscalls.Supported("newfstatat", Newfstatat)
	s.Table[263] = syscalls.Supported("unlinkat", Unlinkat)
	s.Table[264] = syscalls.Supported("renameat", Renameat)
	s.Table[265] = syscalls.Supported("linkat", Linkat)
	s.Table[266] = syscalls.Supported("symlinkat", Symlinkat)
	s.Table[267] = syscalls.Supported("readlinkat", Readlinkat)
	s.Table[268] = syscalls.Supported("fchmodat", Fchmodat)
	s.Table[269] = syscalls.Supported("faccessat", Faccessat)
	s.Table[270] = syscalls.Supported("pselect", Pselect)
	s.Table[271] = syscalls.Supported("ppoll", Ppoll)
	s.Table[275] = syscalls.Supported("splice", Splice)
	s.Table[276] = syscalls.Supported("tee", Tee)
	s.Table[277] = syscalls.Supported("sync_file_range", SyncFileRange)
	s.Table[280] = syscalls.Supported("utimensat", Utimensat)
	s.Table[281] = syscalls.Supported("epoll_pwait", EpollPwait)
	s.Table[282] = syscalls.SupportedPoint("signalfd", Signalfd, linux.PointSignalfd)
	s.Table[283] = syscalls.SupportedPoint("timerfd_create", TimerfdCreate, linux.PointTimerfdCreate)
	s.Table[284] = syscalls.SupportedPoint("eventfd", Eventfd, linux.PointEventfd)
	s.Table[285] = syscalls.PartiallySupported("fallocate", Fallocate, "Not all options are supported.", nil)
	s.Table[286] = syscalls.SupportedPoint("timerfd_settime", TimerfdSettime, linux.PointTimerfdSettime)
	s.Table[287] = syscalls.SupportedPoint("timerfd_gettime", TimerfdGettime, linux.PointTimerfdGettime)
	s.Table[288] = syscalls.SupportedPoint("accept4", Accept4, linux.PointAccept4)
	s.Table[289] = syscalls.SupportedPoint("signalfd4", Signalfd4, linux.PointSignalfd4)
	s.Table[290] = syscalls.SupportedPoint("eventfd2", Eventfd2, linux.PointEventfd2)
	s.Table[291] = syscalls.Supported("epoll_create1", EpollCreate1)
	s.Table[292] = syscalls.SupportedPoint("dup3", Dup3, linux.PointDup3)
	s.Table[293] = syscalls.SupportedPoint("pipe2", Pipe2, linux.PointPipe2)
	s.Table[294] = syscalls.PartiallySupportedPoint("inotify_init1", InotifyInit1, linux.PointInotifyInit1, "inotify events are only available inside the sandbox.", nil)
	s.Table[295] = syscalls.Supported("preadv", Preadv)
	s.Table[296] = syscalls.Supported("pwritev", Pwritev)
	s.Table[299] = syscalls.Supported("recvmmsg", RecvMMsg)
	s.Table[306] = syscalls.Supported("syncfs", Syncfs)
	s.Table[307] = syscalls.Supported("sendmmsg", SendMMsg)
	// FIXME(zkoopmans): Re-enable calls for process_vm_(read/write)v.
	s.Table[316] = syscalls.Supported("renameat2", Renameat2)
	s.Table[319] = syscalls.Supported("memfd_create", MemfdCreate)
	s.Table[322] = syscalls.SupportedPoint("execveat", Execveat, linux.PointExecveat)
	s.Table[327] = syscalls.Supported("preadv2", Preadv2)
	s.Table[328] = syscalls.Supported("pwritev2", Pwritev2)
	s.Table[332] = syscalls.Supported("statx", Statx)
	s.Table[425] = syscalls.PartiallySupported("io_uring_setup", IOUringSetup, "Not all flags and functionality supported.", nil)
	s.Table[436] = syscalls.Supported("close_range", CloseRange)
	s.Table[439] = syscalls.Supported("faccessat2", Faccessat2)
	s.Table[441] = syscalls.Supported("epoll_pwait2", EpollPwait2)
	s.Init()

	// Override ARM64.
	s = linux.ARM64
	s.Table[2] = syscalls.PartiallySupported("io_submit", IoSubmit, "Generally supported with exceptions. User ring optimizations are not implemented.", []string{"gvisor.dev/issue/204"})
	s.Table[5] = syscalls.Supported("setxattr", SetXattr)
	s.Table[6] = syscalls.Supported("lsetxattr", Lsetxattr)
	s.Table[7] = syscalls.Supported("fsetxattr", Fsetxattr)
	s.Table[8] = syscalls.Supported("getxattr", GetXattr)
	s.Table[9] = syscalls.Supported("lgetxattr", Lgetxattr)
	s.Table[10] = syscalls.Supported("fgetxattr", Fgetxattr)
	s.Table[11] = syscalls.Supported("listxattr", ListXattr)
	s.Table[12] = syscalls.Supported("llistxattr", Llistxattr)
	s.Table[13] = syscalls.Supported("flistxattr", Flistxattr)
	s.Table[14] = syscalls.Supported("removexattr", RemoveXattr)
	s.Table[15] = syscalls.Supported("lremovexattr", Lremovexattr)
	s.Table[16] = syscalls.Supported("fremovexattr", Fremovexattr)
	s.Table[17] = syscalls.Supported("getcwd", Getcwd)
	s.Table[19] = syscalls.SupportedPoint("eventfd2", Eventfd2, linux.PointEventfd2)
	s.Table[20] = syscalls.Supported("epoll_create1", EpollCreate1)
	s.Table[21] = syscalls.Supported("epoll_ctl", EpollCtl)
	s.Table[22] = syscalls.Supported("epoll_pwait", EpollPwait)
	s.Table[23] = syscalls.SupportedPoint("dup", Dup, linux.PointDup)
	s.Table[24] = syscalls.SupportedPoint("dup3", Dup3, linux.PointDup3)
	s.Table[25] = syscalls.SupportedPoint("fcntl", Fcntl, linux.PointFcntl)
	s.Table[26] = syscalls.PartiallySupportedPoint("inotify_init1", InotifyInit1, linux.PointInotifyInit1, "inotify events are only available inside the sandbox.", nil)
	s.Table[27] = syscalls.PartiallySupportedPoint("inotify_add_watch", InotifyAddWatch, linux.PointInotifyAddWatch, "inotify events are only available inside the sandbox.", nil)
	s.Table[28] = syscalls.PartiallySupportedPoint("inotify_rm_watch", InotifyRmWatch, linux.PointInotifyRmWatch, "inotify events are only available inside the sandbox.", nil)
	s.Table[29] = syscalls.Supported("ioctl", Ioctl)
	s.Table[32] = syscalls.Supported("flock", Flock)
	s.Table[33] = syscalls.Supported("mknodat", Mknodat)
	s.Table[34] = syscalls.Supported("mkdirat", Mkdirat)
	s.Table[35] = syscalls.Supported("unlinkat", Unlinkat)
	s.Table[36] = syscalls.Supported("symlinkat", Symlinkat)
	s.Table[37] = syscalls.Supported("linkat", Linkat)
	s.Table[38] = syscalls.Supported("renameat", Renameat)
	s.Table[39] = syscalls.Supported("umount2", Umount2)
	s.Table[40] = syscalls.Supported("mount", Mount)
	s.Table[41] = syscalls.Supported("pivot_root", PivotRoot)
	s.Table[43] = syscalls.Supported("statfs", Statfs)
	s.Table[44] = syscalls.Supported("fstatfs", Fstatfs)
	s.Table[45] = syscalls.Supported("truncate", Truncate)
	s.Table[46] = syscalls.Supported("ftruncate", Ftruncate)
	s.Table[47] = syscalls.PartiallySupported("fallocate", Fallocate, "Not all options are supported.", nil)
	s.Table[48] = syscalls.Supported("faccessat", Faccessat)
	s.Table[49] = syscalls.SupportedPoint("chdir", Chdir, linux.PointChdir)
	s.Table[50] = syscalls.SupportedPoint("fchdir", Fchdir, linux.PointFchdir)
	s.Table[51] = syscalls.SupportedPoint("chroot", Chroot, linux.PointChroot)
	s.Table[52] = syscalls.Supported("fchmod", Fchmod)
	s.Table[53] = syscalls.Supported("fchmodat", Fchmodat)
	s.Table[54] = syscalls.Supported("fchownat", Fchownat)
	s.Table[55] = syscalls.Supported("fchown", Fchown)
	s.Table[56] = syscalls.SupportedPoint("openat", Openat, linux.PointOpenat)
	s.Table[57] = syscalls.SupportedPoint("close", Close, linux.PointClose)
	s.Table[59] = syscalls.SupportedPoint("pipe2", Pipe2, linux.PointPipe2)
	s.Table[61] = syscalls.Supported("getdents64", Getdents64)
	s.Table[62] = syscalls.Supported("lseek", Lseek)
	s.Table[63] = syscalls.SupportedPoint("read", Read, linux.PointRead)
	s.Table[64] = syscalls.Supported("write", Write)
	s.Table[65] = syscalls.Supported("readv", Readv)
	s.Table[66] = syscalls.Supported("writev", Writev)
	s.Table[67] = syscalls.Supported("pread64", Pread64)
	s.Table[68] = syscalls.Supported("pwrite64", Pwrite64)
	s.Table[69] = syscalls.Supported("preadv", Preadv)
	s.Table[70] = syscalls.Supported("pwritev", Pwritev)
	s.Table[71] = syscalls.Supported("sendfile", Sendfile)
	s.Table[72] = syscalls.Supported("pselect", Pselect)
	s.Table[73] = syscalls.Supported("ppoll", Ppoll)
	s.Table[74] = syscalls.SupportedPoint("signalfd4", Signalfd4, linux.PointSignalfd4)
	s.Table[76] = syscalls.Supported("splice", Splice)
	s.Table[77] = syscalls.Supported("tee", Tee)
	s.Table[78] = syscalls.Supported("readlinkat", Readlinkat)
	s.Table[79] = syscalls.Supported("newfstatat", Newfstatat)
	s.Table[80] = syscalls.Supported("fstat", Fstat)
	s.Table[81] = syscalls.Supported("sync", Sync)
	s.Table[82] = syscalls.Supported("fsync", Fsync)
	s.Table[83] = syscalls.Supported("fdatasync", Fdatasync)
	s.Table[84] = syscalls.Supported("sync_file_range", SyncFileRange)
	s.Table[85] = syscalls.SupportedPoint("timerfd_create", TimerfdCreate, linux.PointTimerfdCreate)
	s.Table[86] = syscalls.SupportedPoint("timerfd_settime", TimerfdSettime, linux.PointTimerfdSettime)
	s.Table[87] = syscalls.SupportedPoint("timerfd_gettime", TimerfdGettime, linux.PointTimerfdGettime)
	s.Table[88] = syscalls.Supported("utimensat", Utimensat)
	s.Table[180] = syscalls.Supported("mq_open", MqOpen)
	s.Table[181] = syscalls.Supported("mq_unlink", MqUnlink)
	s.Table[198] = syscalls.SupportedPoint("socket", Socket, linux.PointSocket)
	s.Table[199] = syscalls.SupportedPoint("socketpair", SocketPair, linux.PointSocketpair)
	s.Table[200] = syscalls.SupportedPoint("bind", Bind, linux.PointBind)
	s.Table[201] = syscalls.Supported("listen", Listen)
	s.Table[202] = syscalls.SupportedPoint("accept", Accept, linux.PointAccept)
	s.Table[203] = syscalls.SupportedPoint("connect", Connect, linux.PointConnect)
	s.Table[204] = syscalls.Supported("getsockname", GetSockName)
	s.Table[205] = syscalls.Supported("getpeername", GetPeerName)
	s.Table[206] = syscalls.Supported("sendto", SendTo)
	s.Table[207] = syscalls.Supported("recvfrom", RecvFrom)
	s.Table[208] = syscalls.Supported("setsockopt", SetSockOpt)
	s.Table[209] = syscalls.Supported("getsockopt", GetSockOpt)
	s.Table[210] = syscalls.Supported("shutdown", Shutdown)
	s.Table[211] = syscalls.Supported("sendmsg", SendMsg)
	s.Table[212] = syscalls.Supported("recvmsg", RecvMsg)
	s.Table[213] = syscalls.Supported("readahead", Readahead)
	s.Table[221] = syscalls.SupportedPoint("execve", Execve, linux.PointExecve)
	s.Table[222] = syscalls.Supported("mmap", Mmap)
	s.Table[223] = syscalls.PartiallySupported("fadvise64", Fadvise64, "Not all options are supported.", nil)
	s.Table[242] = syscalls.SupportedPoint("accept4", Accept4, linux.PointAccept4)
	s.Table[243] = syscalls.Supported("recvmmsg", RecvMMsg)
	s.Table[267] = syscalls.Supported("syncfs", Syncfs)
	s.Table[269] = syscalls.Supported("sendmmsg", SendMMsg)
	s.Table[276] = syscalls.Supported("renameat2", Renameat2)
	s.Table[279] = syscalls.Supported("memfd_create", MemfdCreate)
	s.Table[281] = syscalls.SupportedPoint("execveat", Execveat, linux.PointExecveat)
	s.Table[286] = syscalls.Supported("preadv2", Preadv2)
	s.Table[287] = syscalls.Supported("pwritev2", Pwritev2)
	s.Table[291] = syscalls.Supported("statx", Statx)
	s.Table[425] = syscalls.PartiallySupported("io_uring_setup", IOUringSetup, "Not all flags and functionality supported.", nil)
	s.Table[436] = syscalls.Supported("close_range", CloseRange)
	s.Table[439] = syscalls.Supported("faccessat2", Faccessat2)
	s.Table[441] = syscalls.Supported("epoll_pwait2", EpollPwait2)
	s.Init()
}
