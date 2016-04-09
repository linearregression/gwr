package protocol

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

type chanBuf struct {
	sync.Mutex
	bytes.Buffer
	ready   chan<- *chanBuf
	closed  bool
	pending bool
	p       []byte
}

func (cb *chanBuf) close() {
	cb.Lock()
	cb.closed = true
	cb.Unlock()
}

var errBufClosed = errors.New("buffer closed")

func (cb *chanBuf) Reset() {
	cb.pending = false
	cb.Buffer.Reset()
}

func (cb *chanBuf) Write(p []byte) (int, error) {
	cb.Lock()

	if cb.closed {
		cb.Unlock()
		return 0, errBufClosed
	}

	send := false
	n, err := cb.Buffer.Write(p)
	if n > 0 && !cb.pending {
		cb.pending = true
		send = true
	}
	cb.Unlock()

	if send {
		cb.ready <- cb
	}
	return n, err
}

func (cb *chanBuf) writeTo(w io.Writer) (int, error) {
	return w.Write(cb.drain())
}

func (cb *chanBuf) drain() []byte {
	cb.Lock()
	if cap(cb.p) < cb.Len() {
		cb.p = make([]byte, cb.Len())
	}

	n := copy(cb.p[:cap(cb.p)], cb.Bytes())
	cb.p = cb.p[:n]
	cb.Reset()
	cb.Unlock()
	return cb.p
}
