package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jacobious52/blockchainserver/blockchain"
)

type chainResponse struct {
	Chain []*blockchain.Block
	Len   int
}

type chainHandler struct {
	blockChainChan chan *blockchain.BlockChain
}

func (h chainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bc := <-h.blockChainChan
	chain := bc.Chain
	bytes, err := json.Marshal(chainResponse{chain, len(chain)})
	h.blockChainChan <- bc

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
