### Overview

This is a collection of simple programs which can generate bitcoin wallets, create and sign transactions, and send transactions over the bitcoin network.

It was a learning project for me to learn more about both golang, and the bitcoin protocol.

This project consists of three programs which each contain the most basic usecase.

* keys.go - Generates a public/private key pair
* transaction.go - Creates, and signs a bitcoin transaction
* network.go - Connects to a peer, and sends a transaction over the network

### Disclaimer

These programs are not "crytographically" random, and should not be used for any purpose other than educational use.

### Installation

1. Install [go](http://golang.org/)
2. run `go get` to install dependencies
3. Follow the instructions at [go-secp256k1](https://github.com/toxeus/go-secp256k1) to compile bitcoin/c-secp256k1
4. Run one of the programs using the syntax below

### Usage

##### Creating a key pair

	go run keys.go

	options (optional)
	--testnet

##### Generating a transaction

	go run transaction.go
	
	options (required)
	--private-key yourPrivateKey
	--public-key yourPublicKey
	--destination destinationPublicKey
	--input-transaction inputTransactionHash
	--satoshis satoshisToSend

	options (optional)
	--input-index inputTransactionIndex


##### Sending a transaction over the bitcoin network

	go run network.go
	
	options (required)
	--transaction yourTransaction
	--node-address 255.255.255.255 (IPv4 address of the bitcoin node to connect to)
	--network-address 255.255.255.255 (IPv4 address of your public IP address)

	options (optional)
	--testnet

### Dependencies

##### https://github.com/toxeus/go-secp256k1
This library is used for the creation of public keys from private keys, as well as signing transactions. It is a project which wraps the official bitcoin/c-secp256k1 bitcoin library.

##### https://github.com/tv42/base58
This library does the base58 conversion. I have included the base58 project in this codebase rather than importing it from the aforementioned github, because I needed to change the dictionary that was used.

### Resources

- Bitpay's insight for testnet: https://test-insight.bitpay.com/
- TP's TestNet Faucet: http://tpfaucet.appspot.com/
- Ken Shirriff's blog post "Bitcoins the hard way": http://www.righto.com/2014/02/bitcoins-hard-way-using-raw-bitcoin.html
- The Bitcoin wiki: https://en.bitcoin.it/
- Bitcoin developer guide: https://bitcoin.org/en/developer-guide
- http://blockexplorer.com/ for general transction/address searching
- http://blockchain.info/ for general transction/address searching
- https://blockchain.info/connected-nodes to see a list of connected nodes
- bitcoin.stackexchange.com (http://bitcoin.stackexchange.com/questions/3374/how-to-redeem-a-basic-tx) Information on redeeming a raw transaction, and explanation of fields.

### License

MIT
