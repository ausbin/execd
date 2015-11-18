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
	"bufio"
	"bytes"
	"io"
	"strconv"
)

type bodyReader struct {
	read   int64
	length int64
	reader *bufio.Reader
}

func NewBodyReader(in io.Reader) *bodyReader {
	return &bodyReader{0, -1, bufio.NewReader(in)}
}

func (br *bodyReader) Read(p []byte) (n int, err error) {
	if br.length == -1 {
		var line string
		line, err = br.reader.ReadString('\n')

		if err != nil {
			return
		}

		// kill newline
		line = line[:len(line)-1]

		br.length, err = strconv.ParseInt(line, 10, 64)

		if err != nil {
			return
		}
	}

	if br.read < br.length {
		n, err = br.reader.Read(p)
		br.read += int64(n)
	}

	if err == nil && br.read >= br.length {
		err = io.EOF
	}

	return
}

type argBodyReader struct {
	*bodyReader
}

func NewArgBodyReader(in io.Reader) *argBodyReader {
	return &argBodyReader{NewBodyReader(in)}
}

func (abr *argBodyReader) Args() (args []string, err error) {
	// read arguments
	for {
		var line string
		line, err = abr.bodyReader.reader.ReadString('\n')

		if err != nil {
			args = nil
			break
		}

		// remove trailing \n
		line = line[:len(line)-1]

		if line == "" {
			break
		} else {
			args = append(args, line)
		}
	}
	return
}

type bodyWriter struct {
	*bytes.Buffer
	out io.Writer
}

func NewBodyWriter(out io.Writer) *bodyWriter {
	return &bodyWriter{&bytes.Buffer{}, out}
}

func (bw *bodyWriter) Flush() (err error) {
	_, err = bw.out.Write([]byte(strconv.Itoa(bw.Buffer.Len()) + "\n"))

	if err == nil {
		_, err = bw.Buffer.WriteTo(bw.out)
	}

	return
}

type argBodyWriter struct {
	*bodyWriter
}

func NewArgBodyWriter(out io.Writer) *argBodyWriter {
	return &argBodyWriter{NewBodyWriter(out)}
}

func (abw *argBodyWriter) WriteArgs(args []string) (err error) {
	for i := 0; i <= len(args); i++ {
		var line string

		if i < len(args) {
			line = args[i]
		}

		_, err = abw.bodyWriter.out.Write([]byte(line + "\n"))
		if err != nil {
			break
		}
	}
	return
}
