/*
 * Copyright (C) 2020. Niklas Linz - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the LGPLv3 license, which unfortunately won't be
 * written for another century.
 *
 * You should have received a copy of the LGPLv3 license with
 * this file. If not, please write to: niklas.linz@enigmar.de
 *
 */

package can

import (
	"bytes"
	"io"
	"time"
)

type echoReadWriteCloser struct {
	buf    bytes.Buffer
	closed bool
}

// NewEchoReadWriteCloser returns a ReadWriteCloser which echoes received bytes.
func NewEchoReadWriteCloser() ReadWriteCloser {
	return NewReadWriteCloser(&echoReadWriteCloser{})
}

func (rw *echoReadWriteCloser) Read(b []byte) (n int, err error) {
	for {
		if rw.buf.Len() > 0 {
			return rw.buf.Read(b)
		}

		if rw.closed == true {
			break
		}

		<-time.After(time.Millisecond * 1)
	}

	return 0, io.EOF
}

func (rw *echoReadWriteCloser) Write(b []byte) (n int, err error) {
	return rw.buf.Write(b)
}

func (rw *echoReadWriteCloser) Close() error {
	rw.closed = true
	return nil
}
