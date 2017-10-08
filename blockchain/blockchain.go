package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

type nodeSet map[string]interface{}

func (n nodeSet) insert(uri string) {
	n[uri] = nil
}

// BlockChain is the main class containing the Chain of blocks and transactions
type BlockChain struct {
	Chain               []*Block
	CurrentTransactions []*Transaction
	Nodes               nodeSet
}

// NewBlockChain creates and returns a new block Chain
func NewBlockChain() *BlockChain {
	bc := &BlockChain{
		Chain:               make([]*Block, 0),
		CurrentTransactions: make([]*Transaction, 0),
		Nodes:               make(nodeSet, 0),
	}
	// Create genisis block
	bc.NewBlock(100, "1")
	return bc
}

// RegisterNode registers a url for another node
func (bc *BlockChain) RegisterNode(uri string) {
	uri = url.PathEscape(uri)
	if uri == "" {
		return
	}
	bc.Nodes.insert(uri)
}

// ValidChain checks if the chain is valid
func ValidChain(chain []*Block) bool {
	if len(chain) == 0 {
		return false
	}

	lastBlock := chain[0]
	currentIndex := 1

	for currentIndex < len(chain) {
		currentBlock := chain[currentIndex]
		if currentBlock.PreviousHash != lastBlock.Hash() {
			log.Println("Invalid Chain at index ", currentBlock.Index)
			return false
		}
		lastBlock = currentBlock
		currentIndex++
	}

	return true
}

// ResolveConflicts replaces our chain with the longest chain in the known network
func (bc *BlockChain) ResolveConflicts() bool {
	var newChain []*Block
	maxLength := len(bc.Chain)

	for node := range bc.Nodes {
		response, err := http.Get(fmt.Sprint("http://", node, "/chain"))
		if err != nil {
			log.Println("Count not reach node", node, ". Error:", err)
			continue
		}

		if response.StatusCode != http.StatusOK {
			log.Println("Count not reach node", node, ". Status:", response.StatusCode)
			continue
		}

		var chainResponse struct {
			Chain []*Block
			Len   int
		}

		err = json.NewDecoder(response.Body).Decode(&chainResponse)
		if err != nil {
			log.Println("Failed to decode node", node, ". Error:", err)
			continue
		}

		if chainResponse.Len > maxLength && ValidChain(chainResponse.Chain) {
			maxLength = chainResponse.Len
			newChain = chainResponse.Chain
		}
	}

	if newChain != nil {
		bc.Chain = newChain
		return true
	}

	return false
}

/*
ProofOfWork Algorithm:
Ref: https://hackernoon.com/learn-blockchains-by-building-one-117428612f46
 - Find a number p' such that hash(pp') contains leading 4 zeroes, where p is the previous p'
 - p is the previous proof, and p' is the new proof
*/
func ProofOfWork(lastProof int64) int64 {
	var proof int64
	for !ValidProof(lastProof, proof) {
		proof++
	}
	return proof
}

// ValidProof Validates the Proof: Does hash(last_proof, proof) contain 4 leading zeroes?
func ValidProof(lastProof, proof int64) bool {
	guess := fmt.Sprint(lastProof, proof)
	guessHash := fmt.Sprintf("%x", sha256.Sum256([]byte(guess)))
	return guessHash[:4] == "0000"
}

// NewBlock creates a new block and adds it to the block Chain
func (bc *BlockChain) NewBlock(proof int64, previousHash string) *Block {
	block := &Block{
		Index:        int64(len(bc.Chain)) + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: bc.CurrentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}
	bc.CurrentTransactions = make([]*Transaction, 0)
	bc.Chain = append(bc.Chain, block)
	return block
}

// NewTransaction creates a new transaction ready to be put on the next block of the Chain
func (bc *BlockChain) NewTransaction(sender, recipient string, amount int64) int64 {
	bc.CurrentTransactions = append(bc.CurrentTransactions, &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	})
	return bc.LastBlock().Index + 1
}

// LastBlock returns the lastest block in the Chain
func (bc *BlockChain) LastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}
