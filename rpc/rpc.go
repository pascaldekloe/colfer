// Package rpc implements the net/rpc codecs.
package rpc

import (
	"errors"
	"io"
	"net/rpc"
)

var errBodyMismatch = errors.New("colfer/rpc: body not a Colfer type")

// colferer covers the encoding methods.
type colferer interface {
	MarshalTo([]byte) int
	MarshalLen() (int, error)
	Unmarshal([]byte) (int, error)
}

type codec struct {
	conn io.ReadWriteCloser

	// buf is the read buffer.
	buf []byte

	// offset is the index of the first data byte in buf.
	offset int

	// i is the index of the data end (exclusive) in buf.
	i int

	// header is a reusable read struct
	header Header
}

// NewClientCodec returns a new RPC codec.
func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &codec{
		conn: conn,
		buf: make([]byte, 2048),
	}
}

// NewServerCodec returns a new RPC codec.
func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	return &codec{
		conn: conn,
		buf: make([]byte, 2048),
	}
}

func (c *codec) ReadRequestHeader(r *rpc.Request) error {
	c.header = Header{}	// reset
	if err := c.decode(&c.header); err != nil {
		return err
	}

	r.ServiceMethod = c.header.Method
	r.Seq = c.header.SeqID
	return nil
}

func (c *codec) ReadResponseHeader(r *rpc.Response) error {
	c.header = Header{}	// reset
	if err := c.decode(&c.header); err != nil {
		return err
	}

	r.ServiceMethod = c.header.Method
	r.Seq = c.header.SeqID
	r.Error = c.header.Error
	return nil
}

func (c *codec) ReadRequestBody(r interface{}) error {
	b, ok := r.(colferer)
	if !ok {
		return errBodyMismatch
	}
	return c.decode(b)
}

func (c *codec) ReadResponseBody(r interface{}) error {
	b, ok := r.(colferer)
	if !ok {
		return errBodyMismatch
	}
	return c.decode(b)
}

func (c *codec) WriteRequest(header *rpc.Request, body interface{}) error {
	// escapes to heap
	h := &Header{
		Method: header.ServiceMethod,
		SeqID:  header.Seq,
	}
	b, ok := body.(colferer)
	if !ok {
		return errBodyMismatch
	}
	return c.encode(h, b)
}

func (c *codec) WriteResponse(header *rpc.Response, body interface{}) error {
	// escapes to heap
	h := &Header{
		Method: header.ServiceMethod,
		SeqID:  header.Seq,
		Error:  header.Error,
	}
	b, ok := body.(colferer)
	if !ok {
		return errBodyMismatch
	}
	return c.encode(h, b)
}

func (c *codec) Close() error {
	return c.conn.Close()
}

func (c *codec) encode(h *Header, body colferer) error {
	l, err := body.MarshalLen()
	if err != nil {
		return err
	}

	if hl, err := h.MarshalLen(); err != nil {
		return err
	} else if hl > l {
		l = hl
	}

	buf := make([]byte, l)
	_, err = c.conn.Write(buf[:h.MarshalTo(buf)])
	if err != nil {
		return err
	}
	_, err = c.conn.Write(buf[:body.MarshalTo(buf)])
	return err
}

func (c *codec) decode(v colferer) error {
	for {
		if c.offset < c.i {
			n, err := v.Unmarshal(c.buf[c.offset:c.i])
			switch err {
			case nil:
				c.offset += n
				return nil

			default:
				return err

			case io.EOF:
			}
		}
		// not enough data

		if c.i >= len(c.buf) {
			if c.offset == 0 {
				// grow
				bigger := make([]byte, len(c.buf)*4)
				copy(bigger, c.buf)
				c.buf = bigger
			} else {
				// move data to start of buffer
				copy(c.buf, c.buf[c.offset:])
				c.i -= c.offset
				c.offset = 0
			}
		}

		n, err := c.conn.Read(c.buf[c.i:])
		c.i += n
		if err != nil {
			return err
		}
	}
}
