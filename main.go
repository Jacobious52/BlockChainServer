package main

import (
	"log"
	"os"

	"github.com/Jacobious52/blockchainserver/server"
)

func main() {

	if len(os.Args) != 2 {
		log.Println("usage: ./blockchainserver <port>")
		return
	}
	port := os.Args[1]

	server := server.NewServer()
	server.Run(port)
}
