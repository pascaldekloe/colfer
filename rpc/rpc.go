// Package rpc implements the net/rpc codecs.
package rpc

import (
	"errors"
	"io"
	"net/rpc"
)

var errBodyMismatch = errors.New("colfer/rpc: body not a Colfer type")

type colferer interface {
	MarshalTo([]byte) int
	MarshalLen() (int, error)
	Unmarshal([]byte) (int, error)
}

type conn struct {
	rwc io.ReadWriteCloser

	// buf is the read buffer.
	buf []byte

	// offset is the index of the first data byte in buf.
	offset int

	// i is the index of the data end (exclusive) in buf.
	i int
}

func NewServerCodec(rwc io.ReadWriteCloser) rpc.ServerCodec {
	return &conn{
		rwc: rwc,
		buf: make([]byte, 2048),
	}
}

func NewClientCodec(rwc io.ReadWriteCloser) rpc.ClientCodec {
	return &conn{
		rwc: rwc,
		buf: make([]byte, 2048),
	}
}

func (c *conn) ReadResponseHeader(r *rpc.Response) error {
	// escapes to heap
	h := new(Header)
	if err := c.deserialize(h); err != nil {
		return err
	}

	r.ServiceMethod = h.Method
	r.Seq = h.SeqID
	r.Error = h.Error
	return nil
}

func (c *conn) ReadRequestHeader(r *rpc.Request) error {
	// escapes to heap
	h := new(Header)
	if err := c.deserialize(h); err != nil {
		return err
	}

	r.ServiceMethod = h.Method
	r.Seq = h.SeqID
	return nil
}

func (c *conn) ReadResponseBody(r interface{}) error {
	b, ok := r.(colferer)
	if !ok {
		return errBodyMismatch
	}
	return c.deserialize(b)
}

func (c *conn) ReadRequestBody(r interface{}) error {
	b, ok := r.(colferer)
	if !ok {
		return errBodyMismatch
	}
	return c.deserialize(b)
}

func (c *conn) WriteRequest(header *rpc.Request, body interface{}) error {
	// escapes to heap
	h := &Header{
		Method: header.ServiceMethod,
		SeqID:  header.Seq,
	}
	b, ok := body.(colferer)
	if !ok {
		return errBodyMismatch
	}
	return c.serialize(h, b)
}

func (c *conn) WriteResponse(header *rpc.Response, body interface{}) error {
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
	return c.serialize(h, b)
}

func (c *conn) Close() error {
	return c.rwc.Close()
}

func (c *conn) serialize(h *Header, body colferer) error {
	l, err := h.MarshalLen()
	if err != nil {
		return err
	}

	if bl, err := body.MarshalLen(); err != nil {
		return err
	} else if bl > l {
		l = bl
	}

	buf := make([]byte, l)
	_, err = c.rwc.Write(buf[:h.MarshalTo(buf)])
	if err != nil {
		return err
	}
	_, err = c.rwc.Write(buf[:body.MarshalTo(buf)])
	return err
}

func (c *conn) deserialize(v colferer) error {
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

		n, err := c.rwc.Read(c.buf[c.i:])
		c.i += n
		if err != nil {
			return err
		}
	}
}
