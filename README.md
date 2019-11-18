# gochain

Modelling a dead-simple blockchain in Go (mostly to learn Go itself).

## Usage

```
// Init chain
chain := NewChain("We ❤️ blockchains")

// Init communication channels
unverified := make(chan Block)
updates := make(chan Chain)
verified := make(chan Block)

// Init node(s)
node := &Node{
    Updates:    updates,
    Unverified: unverified,
    Verified:   verified,
}
go node.Run()

// Send initial version of the chain to all nodes
updates <- *chain

// Run appender
go append(chain, verified)
```