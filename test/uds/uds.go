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

// Package uds contains helpers for testing external UDS functionality.
package uds

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"
	"github.com/utam0k/gvisor/pkg/log"
	"github.com/utam0k/gvisor/pkg/unet"
)

func doEcho(s *unet.Socket) error {
	buf := make([]byte, 512)
	n, err := s.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read: %d, %w", n, err)
	}

	n, err = s.Write(buf[:n])
	if err != nil {
		return fmt.Errorf("failed to write: %d, %w", n, err)
	}
	return nil
}

// createEchoSocket creates a socket that echoes back anything received.
//
// Only works for stream, seqpacket sockets.
func createEchoSocket(path string, protocol int) (cleanup func(), err error) {
	fd, err := unix.Socket(unix.AF_UNIX, protocol, 0)
	if err != nil {
		return nil, fmt.Errorf("error creating echo(%d) socket: %v", protocol, err)
	}

	if err := unix.Bind(fd, &unix.SockaddrUnix{Name: path}); err != nil {
		return nil, fmt.Errorf("error binding echo(%d) socket: %v", protocol, err)
	}

	if err := unix.Listen(fd, 0); err != nil {
		return nil, fmt.Errorf("error listening echo(%d) socket: %v", protocol, err)
	}

	server, err := unet.NewServerSocket(fd)
	if err != nil {
		return nil, fmt.Errorf("error creating echo(%d) unet socket: %v", protocol, err)
	}

	acceptAndEchoOne := func() error {
		s, err := server.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept: %v", err)
		}
		defer s.Close()

		for {
			if err := doEcho(s); err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
		}
	}

	go func() {
		for {
			if err := acceptAndEchoOne(); err != nil {
				log.Warningf("Failed to handle echo(%d) socket: %v", protocol, err)
				return
			}
		}
	}()

	cleanup = func() {
		if err := server.Close(); err != nil {
			log.Warningf("Failed to close echo(%d) socket: %v", protocol, err)
		}
	}

	return cleanup, nil
}

// connectAndBecomeEcho connects to the given socket and turns into an echo server.
func connectAndBecomeEcho(path string, protocol int) (cleanup func(), err error) {
	usePacket := protocol == unix.SOCK_SEQPACKET
	go func() {
		for {
			sock, err := unet.Connect(path, usePacket)
			log.Infof("Connecting to UDS at %q, got %v", path, err)
			if err != nil {
				// Wait and try again.
				time.Sleep(500 * time.Millisecond)
				continue
			}
			defer sock.Close()
			for {
				log.Infof("Connected to UDS at %q, running echo server", path)
				if err := doEcho(sock); err != nil {
					return
				}
			}
		}
	}()

	return func() {}, nil
}

// createNonListeningSocket creates a socket that is bound but not listening.
//
// Only relevant for stream, seqpacket sockets.
func createNonListeningSocket(path string, protocol int) (cleanup func(), err error) {
	fd, err := unix.Socket(unix.AF_UNIX, protocol, 0)
	if err != nil {
		return nil, fmt.Errorf("error creating nonlistening(%d) socket: %v", protocol, err)
	}

	if err := unix.Bind(fd, &unix.SockaddrUnix{Name: path}); err != nil {
		return nil, fmt.Errorf("error binding nonlistening(%d) socket: %v", protocol, err)
	}

	cleanup = func() {
		if err := unix.Close(fd); err != nil {
			log.Warningf("Failed to close nonlistening(%d) socket: %v", protocol, err)
		}
	}

	return cleanup, nil
}

// createNullSocket creates a socket that reads anything received.
//
// Only works for dgram sockets.
func createNullSocket(path string, protocol int) (cleanup func(), err error) {
	fd, err := unix.Socket(unix.AF_UNIX, protocol, 0)
	if err != nil {
		return nil, fmt.Errorf("error creating null(%d) socket: %v", protocol, err)
	}

	if err := unix.Bind(fd, &unix.SockaddrUnix{Name: path}); err != nil {
		return nil, fmt.Errorf("error binding null(%d) socket: %v", protocol, err)
	}

	s, err := unet.NewSocket(fd)
	if err != nil {
		return nil, fmt.Errorf("error creating null(%d) unet socket: %v", protocol, err)
	}

	go func() {
		buf := make([]byte, 512)
		for {
			n, err := s.Read(buf)
			if err != nil {
				log.Warningf("failed to read: %d, %v", n, err)
				return
			}
		}
	}()

	cleanup = func() {
		if err := s.Close(); err != nil {
			log.Warningf("Failed to close null(%d) socket: %v", protocol, err)
		}
	}

	return cleanup, nil
}

type socketCreator func(path string, proto int) (cleanup func(), err error)

// CreateSocketTree creates a local tree of unix domain sockets for use in
// testing:
//   - /stream/echo
//   - /stream/nonlistening
//   - /seqpacket/echo
//   - /seqpacket/nonlistening
//   - /dgram/null
//
// Additionally, it will attempt to connect to sockets at the following
// locations, and turn into an echo server once connected:
//   - /stream/created-in-sandbox
//   - /seqpacket/created-in-sandbox
func CreateSocketTree(baseDir string) (dir string, cleanup func(), err error) {
	dir, err = ioutil.TempDir(baseDir, "sockets")
	if err != nil {
		return "", nil, fmt.Errorf("error creating temp dir: %v", err)
	}

	var protocols = []struct {
		protocol int
		name     string
		sockets  map[string]socketCreator
	}{
		{
			protocol: unix.SOCK_STREAM,
			name:     "stream",
			sockets: map[string]socketCreator{
				"echo":               createEchoSocket,
				"nonlistening":       createNonListeningSocket,
				"created-in-sandbox": connectAndBecomeEcho,
			},
		},
		{
			protocol: unix.SOCK_SEQPACKET,
			name:     "seqpacket",
			sockets: map[string]socketCreator{
				"echo":               createEchoSocket,
				"nonlistening":       createNonListeningSocket,
				"created-in-sandbox": connectAndBecomeEcho,
			},
		},
		{
			protocol: unix.SOCK_DGRAM,
			name:     "dgram",
			sockets: map[string]socketCreator{
				"null": createNullSocket,
			},
		},
	}

	var cleanups []func()
	for _, proto := range protocols {
		protoDir := filepath.Join(dir, proto.name)
		if err := os.Mkdir(protoDir, 0755); err != nil {
			return "", nil, fmt.Errorf("error creating %s dir: %v", proto.name, err)
		}

		for name, fn := range proto.sockets {
			path := filepath.Join(protoDir, name)
			cleanup, err := fn(path, proto.protocol)
			if err != nil {
				return "", nil, fmt.Errorf("error creating %s %s socket: %v", proto.name, name, err)
			}

			cleanups = append(cleanups, cleanup)
		}
	}

	cleanup = func() {
		for _, c := range cleanups {
			c()
		}

		os.RemoveAll(dir)
	}

	return dir, cleanup, nil
}
