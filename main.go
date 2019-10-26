package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var chain *Chain
var unverified chan Block

// BlockRequest describes the POST body of a request to create a new block
type BlockRequest struct {
	Data Data `json:"data"`
}

func blockHandler(w http.ResponseWriter, r *http.Request) {
	blockReq := &BlockRequest{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, blockReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	block, err := Create(chain, blockReq.Data)
	if err != nil {
		log.Println("Unable to create block")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	unverified <- *block
}

func chainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(chain)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func serve() {
	router := mux.NewRouter()
	router.HandleFunc("/chain", chainHandler).Methods("GET")
	router.HandleFunc("/block", blockHandler).Methods("POST")

	http.ListenAndServe(":9999", router)
}

func append(verified <-chan Block) {
	for b := range verified {
		block := b
		chain.Append(&block)
		fmt.Println("Appended block", block.Hash, "to the chain")
	}
}

func main() {
	chain = NewChain("We ❤️ blockchains")
	unverified = make(chan Block)

	updates := make(chan Chain)
	verified := make(chan Block)
	node := &Node{
		Updates:    updates,
		Unverified: unverified,
		Verified:   verified,
	}

	go node.Run()
	updates <- *chain

	go append(verified)

	serve()
}
