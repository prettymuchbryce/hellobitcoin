package main

import (
	"code.google.com/p/go.crypto/ripemd160"
	"crypto/sha256"
	"fmt"
	secp256k1 "github.com/toxeus/go-secp256k1"
	"log"
	"math"
	"math/rand"
	"time"
)

func main() {
	//BTC private key is 256 bits of "random" data.
	privateKey := generatePrivateKey()

	//Multiple steps follow. I've encapsulated this functionality into
	//the base58CheckEncode method because a similar process is used to generate
	//a readable public key as well. Here are the steps for the private key.

	//First generate "extended" private key from private key
	//The difference between a private key and an extended
	//private key is this prefix, which determines the
	//network the key belongs to (real btc network, or test network)

	//EF is the testnet prefix
	//80 is the mainnet prefix

	//Perform SHA-256 on the extended key twice
	//First 4 bytes if this double-sha'd byte array are the checksum
	//Append this checksum to the extended private key
	//Convert the extended private key to a big Int
	//Encoded the big int extended private key into a Base58Checked string

	privateKeyWif := base58CheckEncode("80", privateKey)

	publicKey := generatePublicKey(privateKey)

	//There is also a prefix on the public key
	//This is known as the Network ID Byte, or the version byte
	//6f is the testnet prefix
	//00 is the mainnet prefix
	publicKeyEncoded := base58CheckEncode("00", publicKey)

	//Print the keys
	fmt.Println("Your private key is")
	fmt.Println(privateKeyWif)

	fmt.Println("Your public key is")
	fmt.Println(publicKeyEncoded)
}

func generatePublicKey(privateKeyBytes []byte) []byte {
	//Generate the public key from the private key.
	//Unfortunately golang ecdsa package does not include a
	//secp256k1 curve as this is fairly specific to bitcoin
	//as I understand it, so I have used this one by toxeus which wraps the official bitcoin/c-secp256k1 with cgo.
	var privateKeyBytes32 [32]byte
	for i := 0; i < 32; i++ {
		privateKeyBytes32[i] = privateKeyBytes[i]
	}
	secp256k1.Start()
	publicKeyBytes, success := secp256k1.Pubkey_create(privateKeyBytes32, false)
	if !success {
		log.Fatal("Failed to create public key.")
	}

	secp256k1.Stop()

	//Next we get a sha256 hash of the public key generated
	//via ECDSA, and then get a ripemd160 hash of the sha256 hash.
	shaHash := sha256.New()
	shaHash.Write(publicKeyBytes)
	shadPublicKeyBytes := shaHash.Sum(nil)

	ripeHash := ripemd160.New()
	ripeHash.Write(shadPublicKeyBytes)
	ripeHashedBytes := ripeHash.Sum(nil)

	return ripeHashedBytes
}

func generatePrivateKey() []byte {
	bytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		//This is not "cryptographically random"
		bytes[i] = byte(randInt(0, math.MaxUint8))
	}
	return bytes
}

func randInt(min int, max int) uint8 {
	rand.Seed(time.Now().UTC().UnixNano())
	return uint8(min + rand.Intn(max-min))
}
