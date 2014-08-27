Ultimately I would like to be able to generate a BTC wallet, and send a transaction over the bitcoin test network (TestNet) using golang.

#### Dependencies

* https://github.com/haltingstate/secp256k1-go

#### Other libraries I've used

* https://github.com/tv42/base58

#### Notes

##### Generating a private key

via https://en.bitcoin.it/wiki/Wallet_import_format

* Take a private key
`0C28FCA386C7A227600B2FE50B7CAE11EC86D3BF1FBE471BE89827E19D72AA1D`
* Add a 0x80 byte in front of it for mainnet addresses or 0xef for testnet addresses. Also add a 0x01 byte at the end if the private key will correspond to a compressed public key
`800C28FCA386C7A227600B2FE50B7CAE11EC86D3BF1FBE471BE89827E19D72AA1D`
* Perform SHA-256 hash on the extended key
`8147786C4D15106333BF278D71DADAF1079EF2D2440A4DDE37D747DED5403592`
* Perform SHA-256 hash on result of SHA-256 hash
`507A5B8DFED0FC6FE8801743720CEDEC06AA5C6FCA72B07C49964492FB98A714`
* Take the first 4 bytes of the second SHA-256 hash, this is the checksum
`507A5B8D`
* Add the 4 checksum bytes from point 5 at the end of the extended key from point 2
`800C28FCA386C7A227600B2FE50B7CAE11EC86D3BF1FBE471BE89827E19D72AA1D507A5B8D`
* Convert the result from a byte string into a base58 string using Base58Check encoding. This is the Wallet Import Format
`5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ`

##### Generating a public key from a private key

![](http://i.stack.imgur.com/N93Nn.png)

##### Resources

- Ken Shirriff's blog post "Bitcoins the hard way": http://www.righto.com/2014/02/bitcoins-hard-way-using-raw-bitcoin.html
- The Bitcoin wiki: https://en.bitcoin.it/
- Bitcoin developer guide: https://bitcoin.org/en/developer-guide
- Lots of googling