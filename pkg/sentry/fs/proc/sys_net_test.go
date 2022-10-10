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
	"testing"

	"github.com/utam0k/gvisor/pkg/context"
	"github.com/utam0k/gvisor/pkg/sentry/inet"
	"github.com/utam0k/gvisor/pkg/usermem"
)

func TestQuerySendBufferSize(t *testing.T) {
	ctx := context.Background()
	s := inet.NewTestStack()
	s.TCPSendBufSize = inet.TCPBufferSize{100, 200, 300}
	tmi := &tcpMemInode{s: s, dir: tcpWMem}
	tmf := &tcpMemFile{tcpMemInode: tmi}

	buf := make([]byte, 100)
	dst := usermem.BytesIOSequence(buf)
	n, err := tmf.Read(ctx, nil, dst, 0)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if got, want := string(buf[:n]), "100\t200\t300\n"; got != want {
		t.Fatalf("Bad string: got %v, want %v", got, want)
	}
}

func TestQueryRecvBufferSize(t *testing.T) {
	ctx := context.Background()
	s := inet.NewTestStack()
	s.TCPRecvBufSize = inet.TCPBufferSize{100, 200, 300}
	tmi := &tcpMemInode{s: s, dir: tcpRMem}
	tmf := &tcpMemFile{tcpMemInode: tmi}

	buf := make([]byte, 100)
	dst := usermem.BytesIOSequence(buf)
	n, err := tmf.Read(ctx, nil, dst, 0)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if got, want := string(buf[:n]), "100\t200\t300\n"; got != want {
		t.Fatalf("Bad string: got %v, want %v", got, want)
	}
}

var cases = []struct {
	str     string
	initial inet.TCPBufferSize
	final   inet.TCPBufferSize
}{
	{
		str:     "",
		initial: inet.TCPBufferSize{1, 2, 3},
		final:   inet.TCPBufferSize{1, 2, 3},
	},
	{
		str:     "100\n",
		initial: inet.TCPBufferSize{1, 100, 200},
		final:   inet.TCPBufferSize{100, 100, 200},
	},
	{
		str:     "100 200 300\n",
		initial: inet.TCPBufferSize{1, 2, 3},
		final:   inet.TCPBufferSize{100, 200, 300},
	},
}

func TestConfigureSendBufferSize(t *testing.T) {
	ctx := context.Background()
	s := inet.NewTestStack()
	for _, c := range cases {
		s.TCPSendBufSize = c.initial
		tmi := &tcpMemInode{s: s, dir: tcpWMem}
		tmf := &tcpMemFile{tcpMemInode: tmi}

		// Write the values.
		src := usermem.BytesIOSequence([]byte(c.str))
		if n, err := tmf.Write(ctx, nil, src, 0); n != int64(len(c.str)) || err != nil {
			t.Errorf("Write, case = %q: got (%d, %v), wanted (%d, nil)", c.str, n, err, len(c.str))
		}

		// Read the values from the stack and check them.
		if s.TCPSendBufSize != c.final {
			t.Errorf("TCPSendBufferSize, case = %q: got %v, wanted %v", c.str, s.TCPSendBufSize, c.final)
		}
	}
}

func TestConfigureRecvBufferSize(t *testing.T) {
	ctx := context.Background()
	s := inet.NewTestStack()
	for _, c := range cases {
		s.TCPRecvBufSize = c.initial
		tmi := &tcpMemInode{s: s, dir: tcpRMem}
		tmf := &tcpMemFile{tcpMemInode: tmi}

		// Write the values.
		src := usermem.BytesIOSequence([]byte(c.str))
		if n, err := tmf.Write(ctx, nil, src, 0); n != int64(len(c.str)) || err != nil {
			t.Errorf("Write, case = %q: got (%d, %v), wanted (%d, nil)", c.str, n, err, len(c.str))
		}

		// Read the values from the stack and check them.
		if s.TCPRecvBufSize != c.final {
			t.Errorf("TCPRecvBufferSize, case = %q: got %v, wanted %v", c.str, s.TCPRecvBufSize, c.final)
		}
	}
}

// TestIPForwarding tests the implementation of
// /proc/sys/net/ipv4/ip_forwarding
func TestIPForwarding(t *testing.T) {
	ctx := context.Background()
	s := inet.NewTestStack()

	var cases = []struct {
		comment string
		initial bool
		str     string
		final   bool
	}{
		{
			comment: `Forwarding is disabled; write 1 and enable forwarding`,
			initial: false,
			str:     "1",
			final:   true,
		},
		{
			comment: `Forwarding is disabled; write 0 and disable forwarding`,
			initial: false,
			str:     "0",
			final:   false,
		},
		{
			comment: `Forwarding is enabled; write 1 and enable forwarding`,
			initial: true,
			str:     "1",
			final:   true,
		},
		{
			comment: `Forwarding is enabled; write 0 and disable forwarding`,
			initial: true,
			str:     "0",
			final:   false,
		},
		{
			comment: `Forwarding is disabled; write 2404 and enable forwarding`,
			initial: false,
			str:     "2404",
			final:   true,
		},
		{
			comment: `Forwarding is enabled; write 2404 and enable forwarding`,
			initial: true,
			str:     "2404",
			final:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.comment, func(t *testing.T) {
			s.IPForwarding = c.initial
			ipf := &ipForwarding{stack: s}
			file := &ipForwardingFile{
				stack: s,
				ipf:   ipf,
			}

			// Write the values.
			src := usermem.BytesIOSequence([]byte(c.str))
			if n, err := file.Write(ctx, nil, src, 0); n != int64(len(c.str)) || err != nil {
				t.Errorf("file.Write(ctx, nil, %q, 0) = (%d, %v); want (%d, nil)", c.str, n, err, len(c.str))
			}

			// Read the values from the stack and check them.
			if got, want := s.IPForwarding, c.final; got != want {
				t.Errorf("s.IPForwarding incorrect; got: %v, want: %v", got, want)
			}

		})
	}
}
