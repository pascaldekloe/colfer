package rpc

import (
	"bytes"
	"io"
	"net"
	"net/rpc"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/pascaldekloe/colfer/testdata"
)

type mockConn struct {
	buf bytes.Buffer
}

func (c *mockConn) Read(buf []byte) (n int, err error) {
	return iotest.OneByteReader(&c.buf).Read(buf)
}

func (c *mockConn) Write(buf []byte) (n int, err error) {
	return c.buf.Write(buf)
}

func (c *mockConn) Close() error {
	return nil
}

func TestRequest(t *testing.T) {
	conn := new(mockConn)
	s := NewServerCodec(conn)
	c := NewClientCodec(conn)

	h := &rpc.Request{Seq: 42}
	b := &testdata.O{S: "body " + strings.Repeat("A", 64*1024)}
	if err := c.WriteRequest(h, b); err != nil {
		t.Fatalf("write error: %s", err)
	}

	gotH := new(rpc.Request)
	if err := s.ReadRequestHeader(gotH); err != nil {
		t.Fatalf("read header error: %s", err)
	} else if gotH.Seq != h.Seq {
		t.Errorf("got sequence ID %d, want %d", gotH.Seq, h.Seq)
	}

	gotB := new(testdata.O)
	if err := s.ReadRequestBody(gotB); err != nil {
		t.Fatalf("read body error: %s", err)
	} else if gotB.S != b.S {
		t.Errorf("got body %q, want %q", gotB.S, b.S)
	}

}

func TestResponse(t *testing.T) {
	conn := new(mockConn)
	s := NewServerCodec(conn)
	c := NewClientCodec(conn)

	h := &rpc.Response{Seq: 42}
	b := &testdata.O{S: "body " + strings.Repeat("A", 64*1024)}
	if err := s.WriteResponse(h, b); err != nil {
		t.Fatalf("write error: %s", err)
	}

	gotH := new(rpc.Response)
	if err := c.ReadResponseHeader(gotH); err != nil {
		t.Fatalf("read header error: %s", err)
	} else if gotH.Seq != h.Seq {
		t.Errorf("got sequence ID %d, want %d", gotH.Seq, h.Seq)
	}

	gotB := new(testdata.O)
	if err := c.ReadResponseBody(gotB); err != nil {
		t.Fatalf("read body error: %s", err)
	} else if gotB.S != b.S {
		t.Errorf("got body %q, want %q", gotB.S, b.S)
	}
}

// TestRequestBodySkip calls ReadRequestBody with nil.
func TestRequestBodySkip(t *testing.T) {
	conn := new(mockConn)
	s := NewServerCodec(conn)
	c := NewClientCodec(conn)

	if err := c.WriteRequest(&rpc.Request{Seq: 1}, &testdata.O{S: "body 1"}); err != nil {
		t.Fatalf("write error: %s", err)
	}
	if err := c.WriteRequest(&rpc.Request{Seq: 2}, &testdata.O{S: "body 2"}); err != nil {
		t.Fatalf("write error: %s", err)
	}

	// skip first body
	if err := s.ReadRequestHeader(new(rpc.Request)); err != nil {
		t.Fatalf("read header error: %s", err)
	}
	if err := s.ReadRequestBody(nil); err != nil {
		t.Fatalf("read body error: %s", err)
	}

	gotH := new(rpc.Request)
	if err := s.ReadRequestHeader(gotH); err != nil {
		t.Fatalf("read header error: %s", err)
	} else if want := uint64(2); gotH.Seq != want {
		t.Errorf("got sequence ID %d, want %d", gotH.Seq, want)
	}

	gotB := new(testdata.O)
	if err := s.ReadRequestBody(gotB); err != nil {
		t.Fatalf("read body error: %s", err)
	} else if want := "body 2"; gotB.S != want {
		t.Errorf("got body %q, want %q", gotB.S, want)
	}
}

// TestResponseBodySkip calls ReadResponseBody with nil.
func TestResponseBodySkip(t *testing.T) {
	conn := new(mockConn)
	s := NewServerCodec(conn)
	c := NewClientCodec(conn)

	if err := s.WriteResponse(&rpc.Response{Seq: 1}, &testdata.O{S: "body 1"}); err != nil {
		t.Fatalf("write error: %s", err)
	}
	if err := s.WriteResponse(&rpc.Response{Seq: 2}, &testdata.O{S: "body 2"}); err != nil {
		t.Fatalf("write error: %s", err)
	}

	// skip first body
	if err := c.ReadResponseHeader(new(rpc.Response)); err != nil {
		t.Fatalf("read header error: %s", err)
	}
	if err := c.ReadResponseBody(nil); err != nil {
		t.Fatalf("read body error: %s", err)
	}

	gotH := new(rpc.Response)
	if err := c.ReadResponseHeader(gotH); err != nil {
		t.Fatalf("read header error: %s", err)
	} else if want := uint64(2); gotH.Seq != want {
		t.Errorf("got sequence ID %d, want %d", gotH.Seq, want)
	}

	gotB := new(testdata.O)
	if err := c.ReadResponseBody(gotB); err != nil {
		t.Fatalf("read body error: %s", err)
	} else if want := "body 2"; gotB.S != want {
		t.Errorf("got body %q, want %q", gotB.S, want)
	}
}

type pipeConn struct {
	r *io.PipeReader
	w *io.PipeWriter
}

func (c *pipeConn) Read(buf []byte) (n int, err error) {
	return c.r.Read(buf)
}

func (c *pipeConn) Write(buf []byte) (n int, err error) {
	return c.w.Write(buf)
}

func (c *pipeConn) Close() error {
	c.r.Close()
	return c.w.Close()
}

func BenchmarkCodec(b *testing.B) {
	cc, sc := net.Pipe()
	c := NewClientCodec(cc)
	s := NewServerCodec(sc)

	b.ReportAllocs()
	b.ResetTimer()

	// client requests
	go func() {
		header := new(rpc.Request)
		body := new(testdata.O)
		for i := 0; i < b.N; i++ {
			id := uint64(i)
			header.Seq = id
			body.U64 = id
			if err := c.WriteRequest(header, body); err != nil {
				b.Fatal(err)
			}
		}
	}()

	// server response
	go func() {
		req := new(rpc.Request)
		res := new(rpc.Response)
		body := new(testdata.O)
		for i := 0; i < b.N; i++ {
			if err := s.ReadRequestHeader(req); err != nil {
				b.Fatal(err)
			}

			if err := s.ReadRequestBody(body); err != nil {
				b.Fatal(err)
			}

			res.Seq = req.Seq
			if err := s.WriteResponse(res, body); err != nil {
				b.Fatal(err)
			}
		}
	}()

	header := new(rpc.Response)
	body := new(testdata.O)
	for i := 0; i < b.N; i++ {
		id := uint64(i)

		if err := c.ReadResponseHeader(header); err != nil {
			b.Fatal(err)
		}
		if header.Seq != id {
			b.Fatalf("got response sequence ID %d, want %d", header.Seq, id)
		}

		if err := c.ReadResponseBody(body); err != nil {
			b.Fatal(err)
		}
		if body.U64 != id {
			b.Errorf("got body %d, want %d", body.U64, id)
		}
	}
}
