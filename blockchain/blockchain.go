package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Transaction is the transaction in the blockchain
type Transaction struct {
	Sender    string
	Recipient string
	Amount    int64
}

// Block is a block in the blockchain
type Block struct {
	Index        int64
	Timestamp    int64
	Transactions []*Transaction
	Proof        int64
	PreviousHash string
}

// Hash this block
func (b *Block) Hash() string {
	jsonData := b.JSON()
	hash := sha256.Sum256(jsonData)
	return fmt.Sprintf("%x", hash)
}

// JSON returns the block in json format
func (b *Block) JSON() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		log.Fatalln(err)
		return []byte{'{', '}'}
	}
	return data
}

// BlockChain is the main class containing the chain of blocks and transactions
type BlockChain struct {
	chain               []*Block
	currentTransactions []*Transaction
}

// NewBlockChain creates and returns a new block chain
func NewBlockChain() *BlockChain {
	bc := &BlockChain{
		chain:               make([]*Block, 0),
		currentTransactions: make([]*Transaction, 0),
	}
	// Create genisis block
	bc.NewBlock(100, "1")
	return bc
}

/*
ProofOfWork Algorithm:
Ref: https://hackernoon.com/learn-blockchains-by-building-one-117428612f46
 - Find a number p' such that hash(pp') contains leading 4 zeroes, where p is the previous p'
 - p is the previous proof, and p' is the new proof
*/
func (bc *BlockChain) ProofOfWork(lastProof int64) int64 {
	var proof int64
	for !bc.ValidProof(lastProof, proof) {
		proof++
	}
	return proof
}

// ValidProof Validates the Proof: Does hash(last_proof, proof) contain 4 leading zeroes?
func (bc *BlockChain) ValidProof(lastProof, proof int64) bool {
	guess := fmt.Sprint(lastProof, proof)
	guessHash := fmt.Sprintf("%x", sha256.Sum256([]byte(guess)))
	return guessHash[:4] == "0000"
}

// NewBlock creates a new block and adds it to the block chain
func (bc *BlockChain) NewBlock(proof int64, previousHash string) *Block {
	block := &Block{
		Index:        int64(len(bc.chain)) + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: bc.currentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}
	bc.currentTransactions = make([]*Transaction, 0)
	bc.chain = append(bc.chain, block)
	return block
}

// NewTransaction creates a new transaction ready to be put on the next block of the chain
func (bc *BlockChain) NewTransaction(sender, recipient string, amount int64) int64 {
	bc.currentTransactions = append(bc.currentTransactions, &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	})
	return bc.LastBlock().Index + 1
}

func (bc *BlockChain) Chain() []*Block {
	return bc.chain
}

// LastBlock returns the lastest block in the chain
func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}
