package server

import (
	"net/http"

	"github.com/Jacobious52/blockchainserver/blockchain"
)

type NewTransactionHandler struct {
	blockChainChan chan *blockchain.BlockChain
}

func (h NewTransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bc := <-h.blockChainChan
	w.Write([]byte("TODO: new transaction"))
	h.blockChainChan <- bc
}
