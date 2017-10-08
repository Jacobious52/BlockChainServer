package server

import (
	"encoding/json"
	"net/http"

	"github.com/Jacobious52/blockchainserver/blockchain"
)

type registerNodeHandler struct {
	blockChainChan chan *blockchain.BlockChain
}

func (h registerNodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", 400)
		return
	}

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	var nodeRequest struct {
		Nodes []string
	}
	err := json.NewDecoder(r.Body).Decode(&nodeRequest)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if nodeRequest.Nodes == nil {
		http.Error(w, "Please send some Nodes", 400)
		return
	}

	bc := <-h.blockChainChan
	for _, node := range nodeRequest.Nodes {
		bc.RegisterNode(node)
	}

	nodesList := make([]string, 0, len(bc.Nodes))
	for node := range bc.Nodes {
		nodesList = append(nodesList, node)
	}
	h.blockChainChan <- bc

	nodeResponse := struct {
		Message    string
		TotalNodes []string
	}{
		Message:    "New nodes have been added",
		TotalNodes: nodesList,
	}

	bytes, err := json.Marshal(nodeResponse)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

type resolveNodeHandler struct {
	blockChainChan chan *blockchain.BlockChain
}

func (h resolveNodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bc := <-h.blockChainChan
	changes := bc.ResolveConflicts()

	message := "This chain is authoritative"
	if changes {
		message = "This chain was replaced"
	}

	response := struct {
		Message string
		Chain   []*blockchain.Block
	}{message, bc.Chain}

	bytes, err := json.Marshal(response)
	h.blockChainChan <- bc

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
