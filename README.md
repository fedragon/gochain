# gochain

Modelling a dead-simple blockchain in Go (mostly to learn Go itself).

## Test

    go test

## Run

    go build
    ./gochain

## API Endpoints

### Healthcheck

    curl http://localhost:9999/healthy

### Get chain

    curl http://localhost:9999/chain/

### Add block to chain

    curl -X POST http://localhost:9999/block/ -H 'Content-Type: application/json' -d '{ "data": "Really!"}'

### Get block by hash

    curl http://localhost:9999/block/<hash>