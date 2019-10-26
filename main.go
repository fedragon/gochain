package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// ChainHandler provides a chain to HTTP handlers
type ChainHandler struct {
	Chain *Chain
}

// AddBlockHandler manages requests to add blocks to the chain
type AddBlockHandler struct {
	ChainHandler
	Unverified chan Block
}

// GetBlockHandler manages requests to get blocks
type GetBlockHandler struct {
	ChainHandler
}

// GetChainHandler manages requests to get the whole chain
type GetChainHandler struct {
	ChainHandler
}

// BlockRequest describes the POST body of a request to create a new block
type BlockRequest struct {
	Data Data `json:"data"`
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (adh *AddBlockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	block, err := Create(adh.Chain, blockReq.Data)
	if err != nil {
		log.Println("Unable to create block")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	adh.Unverified <- *block
}

func (gbh *GetBlockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := Hash(vars["hash"])

	block, err := gbh.Chain.Get(hash)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(block)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (gch *GetChainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(gch.Chain)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func serve(chain *Chain, unverified chan Block) {
	chainHandler := ChainHandler{chain}

	router := mux.NewRouter()
	router.HandleFunc("/healthy", healthCheckHandler).Methods("GET")

	router.Handle("/chain", &GetChainHandler{chainHandler}).Methods("GET")
	sub := router.PathPrefix("/block").Subrouter()
	sub.Handle("/",
		&AddBlockHandler{chainHandler, unverified}).Methods("POST")
	sub.Handle("/{hash}",
		&GetBlockHandler{chainHandler}).Methods("GET")

	http.ListenAndServe(":9999", router)
}

func append(chain *Chain, verified <-chan Block) {
	for b := range verified {
		block := b
		chain.Append(&block)
		fmt.Println("Appended block", block.Hash, "to the chain")
	}
}

func main() {
	chain := NewChain("We ❤️ blockchains")
	unverified := make(chan Block)

	updates := make(chan Chain)
	verified := make(chan Block)
	node := &Node{
		Updates:    updates,
		Unverified: unverified,
		Verified:   verified,
	}

	go node.Run()
	updates <- *chain

	go append(chain, verified)

	serve(chain, unverified)
}
