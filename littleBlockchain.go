package main

// Add opportunity to show and add new blocks with html forms
import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64  // date when block is created
	Data          []byte // transactions or smartcontract
	PrevBlockhash []byte
	Hash          []byte
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockhash, b.Data, timestamp}, []byte{})
	// sha256.Sum appends only byte slice
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// Golang Constructor Type
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	// redefinition Hash
	block.SetHash()
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis String", []byte{})
}

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	// choose the last block in chain
	prevBlock := bc.blocks[len(bc.blocks)-1]
	// create new block with data from parametris
	newBlock := NewBlock(data, prevBlock.Hash)
	// append new block to slice
	bc.blocks = append(bc.blocks, newBlock)
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

// Web part
func Renfer(w http.ResponseWriter, r *http.Request) {
	bc := NewBlockchain()
	bc.AddBlock("Send 1 BTC to Dinc")
	bc.AddBlock("Send 3 EXP to Dyrk")
	temp, _ := json.Marshal(bc.blocks)
	fmt.Fprintf(w, string(temp))
}

func main() {
	http.HandleFunc("/", Renfer)
	http.ListenAndServe(":8080", nil)
}
