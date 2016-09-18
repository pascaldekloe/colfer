package rpc

import (
	"bytes"
	"net/rpc"
	"strings"
	"testing"
	"testing/iotest"
)

type connMock struct {
	buf bytes.Buffer
}

func (c *connMock) Read(buf []byte) (n int, err error) {
	return iotest.OneByteReader(&c.buf).Read(buf)
}

func (c *connMock) Write(buf []byte) (n int, err error) {
	return c.buf.Write(buf)
}

func (c *connMock) Close() error {
	return nil
}

func TestRequest(t *testing.T) {
	rwc := new(connMock)
	s := NewServerCodec(rwc)
	c := NewClientCodec(rwc)

	h := &rpc.Request{Seq: 42}
	// body can be any Colfer struct
	b := &Header{Error: "body " + strings.Repeat("A", 64*1024)}

	if err := c.WriteRequest(h, b); err != nil {
		t.Fatalf("write error: %s", err)
	}

	gotH := new(rpc.Request)
	if err := s.ReadRequestHeader(gotH); err != nil {
		t.Fatalf("read header error: %s", err)
	} else if gotH.Seq != h.Seq {
		t.Errorf("got sequence ID %d, want %d", gotH.Seq, h.Seq)
	}

	gotB := new(Header)
	if err := s.ReadRequestBody(gotB); err != nil {
		t.Fatalf("read body error: %s", err)
	} else if gotB.Error != b.Error {
		t.Errorf("got body %q, want %q", gotB.Error, b.Error)
	}

}

func TestResponse(t *testing.T) {
	rwc := new(connMock)
	s := NewServerCodec(rwc)
	c := NewClientCodec(rwc)

	h := &rpc.Response{Seq: 42}
	// body can be any Colfer struct
	b := &Header{Error: "body " + strings.Repeat("A", 64*1024)}

	if err := s.WriteResponse(h, b); err != nil {
		t.Fatalf("write error: %s", err)
	}

	gotH := new(rpc.Response)
	if err := c.ReadResponseHeader(gotH); err != nil {
		t.Fatalf("read header error: %s", err)
	} else if gotH.Seq != h.Seq {
		t.Errorf("got sequence ID %d, want %d", gotH.Seq, h.Seq)
	}

	gotB := new(Header)
	if err := c.ReadResponseBody(gotB); err != nil {
		t.Fatalf("read body error: %s", err)
	} else if gotB.Error != b.Error {
		t.Errorf("got body %q, want %q", gotB.Error, b.Error)
	}
}
