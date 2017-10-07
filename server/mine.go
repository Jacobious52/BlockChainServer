package server

import (
	"net/http"

	"github.com/Jacobious52/blockchainserver/blockchain"
)

type MineHandler struct {
	blockChainChan chan *blockchain.BlockChain
}

func (h MineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bc := <-h.blockChainChan
	w.Write([]byte("TODO: mine"))
	h.blockChainChan <- bc
}
