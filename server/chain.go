package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jacobious52/blockchainserver/blockchain"
)

type ChainResponse struct {
	chain []*blockchain.Block
	len   int
}

type ChainHandler struct {
	blockChainChan chan *blockchain.BlockChain
}

func (h ChainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bc := <-h.blockChainChan
	chain := bc.Chain()
	bytes, err := json.Marshal(ChainResponse{chain, len(chain)})
	h.blockChainChan <- bc

	if err != nil {
		log.Println(err)
	}
	w.Write(bytes)
}
