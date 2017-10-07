package main

import "github.com/Jacobious52/blockchainserver/server"

func main() {
	server := server.NewServer()
	server.Run()
}
