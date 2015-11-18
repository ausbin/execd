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
	"log"
	"net"
	"os"

	"code.austinjadams.com/execd"
)

const (
	argProgram = iota
	argWhere
	argProg
	argCount
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < argCount {
		log.Fatalln("usage:", os.Args[0], "<where> <prog> [args...]")
	}

	conn, err := net.Dial("tcp", os.Args[argWhere])

	if err != nil {
		log.Fatalln(err)
	}

	client := execd.NewClient(conn)

	args := os.Args[argProg:]

	if err = client.Exec(os.Stdin, os.Stdout, args...); err != nil {
		log.Fatalln(err)
	}

	client.Close()
}
