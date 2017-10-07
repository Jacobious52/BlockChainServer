package server

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"github.com/Jacobious52/blockchainserver/blockchain"
)

type Server struct {
	nodeIdentifier string
	blockChain     *blockchain.BlockChain
	blockChainChan chan *blockchain.BlockChain
}

func newUUID() string {

	b := make([]byte, 16)
	n, err := rand.Read(b)
	if err != nil || n != 16 {
		log.Println(err)
		return "0000-00-00-00-000000"
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func NewServer() Server {
	return Server{
		newUUID(),
		blockchain.NewBlockChain(),
		make(chan *blockchain.BlockChain, 1),
	}
}

// Run the server
func (s Server) Run() {
	s.blockChainChan <- s.blockChain

	http.Handle("/", welcomeHandler{s.nodeIdentifier})
	http.Handle("/chain", chainHandler{s.blockChainChan})
	http.Handle("/mine", mineHandler{s.blockChainChan})
	http.Handle("/transaction/new", newTransactionHandler{s.blockChainChan})
	http.Handle("/transactions", transactionsHandler{s.blockChainChan})
	log.Fatalln(http.ListenAndServe("0.0.0.0:3000", nil))
}
