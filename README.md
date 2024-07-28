# merkle-patricia-trie

# Description

Implementation of the merkle-patricia trie data structure described in the ehtereum yellow paper.

https://ethereum.github.io/yellowpaper/paper.pdf

# Installation

`go run .`

or

`go build`

# Test

To run the tests

`go test`

# Notes

- Due to lake of time RLP serialization is not implemented
- All data are kept as simple bytes
- Two storage mode are available, in memory and using badgerDB
