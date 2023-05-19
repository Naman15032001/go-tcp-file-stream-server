package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct {
}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	//buffer := make([]byte, 2048)

	buffer := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(buffer, conn, size)

		if err != nil {
			log.Fatal(err)
		}

		// file := buffer[:n]
		file := buffer.Bytes()

		fmt.Println(file)
		fmt.Printf("recieved %d bytes over the network \n ", n)
	}
}

func main() {
	go func() {
		time.Sleep(4 * time.Second)
		sendFile(4000)
	}()
	server := &FileServer{}
	server.start()
}

func sendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)

	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	// n, err := conn.Write(file)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//n, err := io.Copy(conn, bytes.NewReader(file) , )

	binary.Write(conn, binary.LittleEndian, int64(size))

	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("written %d bytes over the network \n", n)
	return nil
}
