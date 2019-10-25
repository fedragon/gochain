package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var ledger *Ledger

func ledgerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(ledger)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func serve() {
	router := mux.NewRouter()
	router.HandleFunc("/ledger", ledgerHandler).Methods("GET")

	http.ListenAndServe(":9999", router)
}

func main() {
	ledger = NewLedger("We ❤️ blockchains")

	updates := make(chan Ledger)
	unverified := make(chan Block)
	verified := make(chan Block)
	node := &Node{
		Updates:    updates,
		Unverified: unverified,
		Verified:   verified,
	}

	go node.Run()

	updates <- *ledger

	go func() {
		t := time.NewTicker(time.Second * 5)

		for {
			select {
			case <-t.C:
				block, err := Create(ledger, "foobar")
				if err != nil {
					log.Println("Unable to create block")
				}

				unverified <- *block
			case block := <-verified:
				ledger.Append(&block)
				fmt.Println("Appended block", block.Hash, "to the ledger")
			}
		}
	}()

	serve()
}
