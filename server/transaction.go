package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jacobious52/blockchainserver/blockchain"
)

type newTransactionHandler struct {
	blockChainChan chan *blockchain.BlockChain
}

func (h newTransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", 400)
		return
	}

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var transaction blockchain.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	bc := <-h.blockChainChan
	index := bc.NewTransaction(transaction.Sender, transaction.Recipient, transaction.Amount)
	h.blockChainChan <- bc

	w.Write([]byte(fmt.Sprintln("Transaction will be added to Block", index)))
}

type transactionsHandler struct {
	blockChainChan chan *blockchain.BlockChain
}

func (h transactionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bc := <-h.blockChainChan
	bytes, err := json.Marshal(bc.CurrentTransactions)
	h.blockChainChan <- bc
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(bytes)
}
