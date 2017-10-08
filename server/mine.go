package server

import (
	"encoding/json"
	"net/http"

	"github.com/Jacobious52/blockchainserver/blockchain"
)

type mineResponse struct {
	Message string
	Block   *blockchain.Block
}

type mineHandler struct {
	blockChainChan chan *blockchain.BlockChain
	nodeID         string
}

func (h mineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bc := <-h.blockChainChan
	lastBlock := bc.LastBlock()
	lastHash := lastBlock.Hash()
	lastProof := lastBlock.Proof
	h.blockChainChan <- bc

	proof := blockchain.ProofOfWork(lastHash, lastProof)

	bc = <-h.blockChainChan
	bc.NewTransaction("0", h.nodeID, 1)
	block := bc.NewBlock(proof, lastHash)
	h.blockChainChan <- bc

	response := mineResponse{
		Message: "New Block Created",
		Block:   block,
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
