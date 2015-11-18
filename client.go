// Copyright (c) 2015 Austin Adams
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.

package execd

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
)

type Client struct {
	net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}

func DialClient(network, addr string) (*Client, error) {
	conn, err := net.Dial(network, addr)

	if err != nil {
		return nil, err
	}

	return NewClient(conn), nil
}

func (c *Client) Exec(in io.Reader, out io.Writer, args ...string) error {
	bodyWriter := NewArgBodyWriter(c.Conn)
	bodyReader := NewBodyReader(c.Conn)

	if err := bodyWriter.WriteArgs(args); err != nil {
		return err
	}

	if _, err := io.Copy(bodyWriter, in); err != nil {
		return err
	}

	// now that we know the size of the data to send, send the length
	// and then the data
	if err := bodyWriter.Flush(); err != nil {
		return err
	}

	if _, err := io.Copy(out, bodyReader); err != nil {
		return err
	}

	return nil
}

func (c *Client) ExecString(input string, args ...string) (output string, err error) {
	in := bytes.NewBufferString(input)
	out := bytes.NewBuffer(nil)

	err = c.Exec(in, out, args...)

	if err == nil {
		output = out.String()
	}

	return
}

type devNull struct{}

func (dn devNull) Read(_ []byte) (n int, err error) {
	return 0, io.EOF
}

func (dn devNull) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (dn devNull) ReadFrom(in io.Reader) (n int64, err error) {
	// ioutil has a nice implementation of this, so let's not reinvent the wheel
	// XXX this type assertion is a hack, but it's safe with the current
	//     implementation of the go stdlib
	return ioutil.Discard.(io.ReaderFrom).ReadFrom(in)
}

// useful for black-holing input or output from Exec()
var DevNull devNull
