package server

import (
	"io"
	"log"
	"net"
)

// reading input from the client
func readCommand(conn net.Conn) (string, error) {
	var buf []byte = make([]byte, 1024)
	n, err := conn.Read(buf[:]) // blocking call, server will wait until something is sent by the client
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

// repsonding back to client
func respond(cmd string, conn net.Conn) error {
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}

func RunSyncTCPServer() {
	// Listen for incoming connections on a given port
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	var con_client = 0

	for {
		// Accept incoming connection
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		con_client += 1
		log.Println("client connected with address:", conn.RemoteAddr(), "concurrent clients", con_client)

		for {
			// over the socket, continuously reading the command and responding back to client
			cmd, err := readCommand(conn)
			if err != nil {
				conn.Close()
				con_client -= 1
				log.Println("client disconnected", conn.RemoteAddr(), "concurrent clients", con_client)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}

			log.Println("command", cmd)
			if err = respond("Server: "+cmd, conn); err != nil {
				log.Println("err write:", err)
			}
		}

	}
}
