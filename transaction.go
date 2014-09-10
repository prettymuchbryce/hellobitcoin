package main

import (
    /*"fmt"
    "crypto/sha256"
    "code.google.com/p/go.crypto/ripemd160"
    "time"
    "math"
    "log"
    "bytes"
    "math/big"
    secp256k1 "github.com/haltingstate/secp256k1-go"*/
    "hellobitcoin/base58"
    //"flag"
    "log"
    "encoding/hex"
    "fmt"
)

func main() {
    /*var inputTransactionHash string
    var publicKeyBase58 string
    var publicKeyBase58Destination string
    var inputTransactionOutputIndex int
    var satoshis int

    flag.StringVar(&inputTransactionHash, "input", "", "The hexidecimal value of the hashed transaction input")
    flag.StringVar(&publicKeyBase58, "pubkey", "", "Your public key in base58 format")
    flag.StringVar(&publicKeyBase58Destination, "destination", "", "The destination public key in base58 format")

    flag.IntVar(&inputTransactionOutputIndex, "outputIndex", 0, "The index of the source transaction output you wish to use as input for this transaction")
    flag.IntVar(&satoshis, "value", 0, "The number of satoshis you wish to send (the remainder will be given as a mining fee)")

    flag.Parse()

    createRawTransaction(inputTransactionHash, publicKeyBase58, publicKeyBase58Destination, inputTransactionOutputIndex, satoshis)*/

    createRawTransaction("7756bf4ed3b495adb63e05c02398e799c471b885d10523028b6e1b10f0ae181f", "mujf6HNVrAFUX2gjNgirTWyxaT7XzeKUrj", "msj42CCGruhRsFrGATiUuh25dtxYtnpbTx", 1, 25000000)
}

func createRawTransaction(inputTransactionHash string, publicKeyBase58 string, publicKeyBase58Destination string, inputTransactionOutputIndex int, satoshis int) {
    //112

    publicKeyInt, err := base58.DecodeToBig([]byte(publicKeyBase58))
    if err != nil {
        log.Fatal(err)
    }

    publicKeyBytes := publicKeyInt.Bytes()

    version, err := hex.DecodeString("010000")
    if err != nil {
        log.Fatal(err)
    }

    inputs, err := hex.DecodeString("01")
    if err != nil {
        log.Fatal(err)
    }

    input := &inputTransactionHash //this should be reversed

    outputIndex, err := hex.DecodeString("01000000")
    if err != nil {
        log.Fatal(err)
    }

    scriptSigLength := 4 + len(publicKeyBytes)
    scriptSig := make([]byte, scriptSigLength)

    scriptSig[0] = 118 //OP_DUP
    scriptSig[1] = 169 //OP_HASH160
    for i:=2; i<scriptSigLength-2; i++ {
        scriptSig[i] = publicKeyBytes[i-2]
    }
    scriptSig[scriptSigLength-2] = 136 //OP_EQUALVERIFY
    scriptSig[scriptSigLength-1] = 172 //OP_CHECKSIG

    fmt.Println(publicKeyBytes)
    fmt.Println(scriptSigLength, version, inputs, input, outputIndex, []byte(satoshis))

    //scriptSig 0x76a9 + 9bf8cee4ce4532eab13454490dbdfb346d5e37f8 + 0x88ac
    //scriptSigLength := make([]byte, 1)
    //scriptSigLength[0] = 

    //four byte version field
    //one byte for # of inputs
    //32 byte hash of thee transaction for which want to redeem an output 7756bf4ed3b495adb63e05c02398e799c471b885d10523028b6e1b10f0ae181f
    //four byte field denoting the output index (01000000)
    //script sig lengh (1 byte)
    //script sig (24 bytes) OP_DUP OP_HASH160 9bf8cee4ce4532eab13454490dbdfb346d5e37f8 OP_EQUALVERIFY OP_CHECKSIG
    //four byte field which is always 0xfffffff (lol?)
    //1 byte varint containg number of outputs (01)
    //8 byte field containing the amount we want to redeem (left over is for miners)
    //1 byte for output script
    //24 bytes for actual script OP_DUP OP_HASH160 9bf8cee4ce4532eab13454490dbdfb346d5e37f8 OP_EQUALVERIFY OP_CHECKSIG (NOT BASE 58.. need to decode FUUUUU)
    //four byte "lock time" field ? 0x00000000
    //four byte hashcode type. 0x00000001
}