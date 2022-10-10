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

package pipe

import (
	"bytes"
	"testing"

	"github.com/utam0k/gvisor/pkg/errors/linuxerr"
	"github.com/utam0k/gvisor/pkg/sentry/contexttest"
	"github.com/utam0k/gvisor/pkg/usermem"
	"github.com/utam0k/gvisor/pkg/waiter"
)

func TestPipeRW(t *testing.T) {
	ctx := contexttest.Context(t)
	r, w := NewConnectedPipe(ctx, 65536)
	defer r.DecRef(ctx)
	defer w.DecRef(ctx)

	msg := []byte("here's some bytes")
	wantN := int64(len(msg))
	n, err := w.Writev(ctx, usermem.BytesIOSequence(msg))
	if n != wantN || err != nil {
		t.Fatalf("Writev: got (%d, %v), wanted (%d, nil)", n, err, wantN)
	}

	buf := make([]byte, len(msg))
	n, err = r.Readv(ctx, usermem.BytesIOSequence(buf))
	if n != wantN || err != nil || !bytes.Equal(buf, msg) {
		t.Fatalf("Readv: got (%d, %v) %q, wanted (%d, nil) %q", n, err, buf, wantN, msg)
	}
}

func TestPipeReadBlock(t *testing.T) {
	ctx := contexttest.Context(t)
	r, w := NewConnectedPipe(ctx, 65536)
	defer r.DecRef(ctx)
	defer w.DecRef(ctx)

	n, err := r.Readv(ctx, usermem.BytesIOSequence(make([]byte, 1)))
	if n != 0 || err != linuxerr.ErrWouldBlock {
		t.Fatalf("Readv: got (%d, %v), wanted (0, %v)", n, err, linuxerr.ErrWouldBlock)
	}
}

func TestPipeWriteBlock(t *testing.T) {
	const atomicIOBytes = 2
	const capacity = MinimumPipeSize

	ctx := contexttest.Context(t)
	r, w := NewConnectedPipe(ctx, capacity)
	defer r.DecRef(ctx)
	defer w.DecRef(ctx)

	msg := make([]byte, capacity+1)
	n, err := w.Writev(ctx, usermem.BytesIOSequence(msg))
	if wantN, wantErr := int64(capacity), linuxerr.ErrWouldBlock; n != wantN || err != wantErr {
		t.Fatalf("Writev: got (%d, %v), wanted (%d, %v)", n, err, wantN, wantErr)
	}
}

func TestPipeWriteUntilEnd(t *testing.T) {
	const atomicIOBytes = 2

	ctx := contexttest.Context(t)
	r, w := NewConnectedPipe(ctx, atomicIOBytes)
	defer r.DecRef(ctx)
	defer w.DecRef(ctx)

	msg := []byte("here's some bytes")

	wDone := make(chan struct{}, 0)
	rDone := make(chan struct{}, 0)
	defer func() {
		// Signal the reader to stop and wait until it does so.
		close(wDone)
		<-rDone
	}()

	go func() {
		defer close(rDone)
		// Read from r until done is closed.
		ctx := contexttest.Context(t)
		buf := make([]byte, len(msg)+1)
		dst := usermem.BytesIOSequence(buf)
		e, ch := waiter.NewChannelEntry(waiter.ReadableEvents)
		r.EventRegister(&e)
		defer r.EventUnregister(&e)
		for {
			n, err := r.Readv(ctx, dst)
			dst = dst.DropFirst64(n)
			if err == linuxerr.ErrWouldBlock {
				select {
				case <-ch:
					continue
				case <-wDone:
					// We expect to have 1 byte left in dst since len(buf) ==
					// len(msg)+1.
					if dst.NumBytes() != 1 || !bytes.Equal(buf[:len(msg)], msg) {
						t.Errorf("Reader: got %q (%d bytes remaining), wanted %q", buf, dst.NumBytes(), msg)
					}
					return
				}
			}
			if err != nil {
				t.Errorf("Readv: got unexpected error %v", err)
				return
			}
		}
	}()

	src := usermem.BytesIOSequence(msg)
	e, ch := waiter.NewChannelEntry(waiter.WritableEvents)
	w.EventRegister(&e)
	defer w.EventUnregister(&e)
	for src.NumBytes() != 0 {
		n, err := w.Writev(ctx, src)
		src = src.DropFirst64(n)
		if err == linuxerr.ErrWouldBlock {
			<-ch
			continue
		}
		if err != nil {
			t.Fatalf("Writev: got (%d, %v)", n, err)
		}
	}
}
