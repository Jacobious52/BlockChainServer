package server

import (
	"net/http"

	"github.com/Jacobious52/blockchainserver/blockchain"
)

type mineResponse struct {
	message      string
	index        int64
	transactions []*blockchain.Transaction
	proof        int64
	previousHash string
}

type mineHandler struct {
	blockChainChan chan *blockchain.BlockChain
}

func (h mineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bc := <-h.blockChainChan

	lastBlock := bc.LastBlock()
	lastProof := lastBlock.Proof
	proof := bc.ProofOfWork(lastProof)

	bc.NewTransaction("0", "id", 1)
	bc.NewBlock(proof, lastBlock.Hash())

	h.blockChainChan <- bc

	w.Write([]byte("TODO: return mine result"))
}
