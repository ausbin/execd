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

package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os/exec"

	"code.austinjadams.com/execd"
)

func handle(c net.Conn) {
	// if we bail out, close the connection
	defer c.Close()

	for {
		bodyWriter := execd.NewBodyWriter(c)
		bodyReader := execd.NewArgBodyReader(c)

		args, err := bodyReader.Args()

		// we're done
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			break
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = bodyReader
		cmd.Stdout = bodyWriter

		log.Println("running command", args, "...")
		err = cmd.Run()
		log.Println("done")

		if err != nil {
			log.Println(err)
			break
		}

		// write command output into c
		err = bodyWriter.Flush()

		if err != nil {
			log.Println(err)
			break
		}
	}
}

func main() {
	listen := flag.String("listen", "127.0.0.1:4000", "where to listen")
	timestamps := flag.Bool("timedlog", false, "show timestamps in logging")
	flag.Parse()

	if !*timestamps {
		log.SetFlags(0)
	}

	sock, err := net.Listen("tcp", *listen)

	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := sock.Accept()

		if err != nil {
			log.Fatalln(err)
		}

		go handle(conn)
	}
}
