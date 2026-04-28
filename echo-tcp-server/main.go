package main

import (
	"echotcpserver/server"
	"log"
)

func main() {
	log.Println("Running Echo Server on localhost:8080")
	server.RunSyncTCPServer()
}
