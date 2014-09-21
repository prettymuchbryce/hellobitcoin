package main

import (
	"flag"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"crypto/sha256"
	secp256k1 "github.com/haltingstate/secp256k1-go"
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

	var privateKey string
	flag.StringVar(&privateKey, "private-key", "", "")
	flag.Parse()


	publicKeyBytes := base58CheckDecode("mujf6HNVrAFUX2gjNgirTWyxaT7XzeKUrj")

	//SCRIPT SIG (TEMPORARILY SCRIPTPUBKEY OF INPUT) + LENGTH
	scriptSigLength := 4 + len(publicKeyBytes)
	scriptSig := make([]byte, scriptSigLength)

	scriptSig[0] = 118 //OP_DUP
	scriptSig[1] = 169 //OP_HASH160
	for i := 2; i < scriptSigLength-2; i++ {
		scriptSig[i] = publicKeyBytes[i-2]
	}
	scriptSig[scriptSigLength-2] = 136 //OP_EQUALVERIFY
	scriptSig[scriptSigLength-1] = 172 //OP_CHECKSIG


	rawTransaction := createRawTransaction("7756bf4ed3b495adb63e05c02398e799c471b885d10523028b6e1b10f0ae181f", "msj42CCGruhRsFrGATiUuh25dtxYtnpbTx", 1, 25000000, scriptSig)

	//HASHCODETYPE
	hashCodeType, err := hex.DecodeString("01000000")
	if err != nil {
		log.Fatal(err)
	}

	rawTransactionWithHashCodeType := append(rawTransaction, hashCodeType[0])

	finalTransaction := signRawTransaction(rawTransactionWithHashCodeType, privateKey, "mujf6HNVrAFUX2gjNgirTWyxaT7XzeKUrj")
	
	//HASHCODETYPE
	finalTransactionHex := hex.EncodeToString(finalTransaction)

	fmt.Println("Your final transaction is")
	fmt.Println(finalTransactionHex)
}

func signRawTransaction(rawTransaction []byte, privateKeyBase58 string, publicKeyBase58 string) []byte {
	publicKeyBytes := base58CheckDecode(publicKeyBase58)
	privateKeyBytes := base58CheckDecode(privateKeyBase58)

	//Perform SHA-256 twice
	shaHash := sha256.New()
	shaHash.Write(rawTransaction)
	var hash []byte = shaHash.Sum(nil)

	shaHash2 := sha256.New()
	shaHash2.Write(hash)
	rawTransactionHashed := shaHash2.Sum(nil)

	signedTransaction := secp256k1.Sign(rawTransactionHashed, privateKeyBytes)

	hashCodeType, err := hex.DecodeString("01")
	if err != nil {
		log.Fatal(err)
	}

	signedTransactionWithHashCodeType := append(signedTransaction, hashCodeType[0])

	signedTransactionWithHashCodeTypeLength := byte(len(signedTransactionWithHashCodeType))

	pubKeyLength := byte(len(publicKeyBytes))

	var buffer bytes.Buffer
	buffer.WriteByte(signedTransactionWithHashCodeTypeLength)
	buffer.Write(signedTransactionWithHashCodeType)
	buffer.WriteByte(pubKeyLength)
	buffer.Write(publicKeyBytes)

	scriptSig := buffer.Bytes()

	return createRawTransaction("7756bf4ed3b495adb63e05c02398e799c471b885d10523028b6e1b10f0ae181f", "msj42CCGruhRsFrGATiUuh25dtxYtnpbTx", 1, 25000000, scriptSig)

	//add 01 byte
}

func createRawTransaction(inputTransactionHash string, publicKeyBase58Destination string, inputTransactionOutputIndex int, satoshis int, scriptSig []byte) []byte {
	//VERSION FIELD
	version, err := hex.DecodeString("01000000")
	if err != nil {
		log.Fatal(err)
	}

	//# of INPUTS
	inputs, err := hex.DecodeString("01")
	if err != nil {
		log.Fatal(err)
	}

	//Reversed Transaction Input
	inputTransactionBytes, err := hex.DecodeString(inputTransactionHash) //this should be reversed ?
	if err != nil {
		log.Fatal(err)
	}

	inputTransactionBytesReversed := make([]byte, len(inputTransactionBytes))
	j := 0
	for i := len(inputTransactionBytes) - 1; i > 0; i-- {
		inputTransactionBytesReversed[j] = inputTransactionBytes[i]
		j++
	}

	//OUTPUT INDEX OF REFERENCED TRANSACTION
	outputIndex, err := hex.DecodeString("01000000")
	if err != nil {
		log.Fatal(err)
	}

	//SCRIPT SIG (TEMPORARILY SCRIPTPUBKEY OF INPUT) + LENGTH
	scriptSigLength := len(scriptSig)

	//SEQUENCE
	sequence, err := hex.DecodeString("ffffff")
	if err != nil {
		log.Fatal(err)
	}

	//NUMBER OF OUTPUTS
	numOutputs, err := hex.DecodeString("01")
	if err != nil {
		log.Fatal(err)
	}

	//SATOSHIS TO SEND
	satoshiBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(satoshiBytes, uint64(satoshis))

	publicKeyDestinationBytes := base58CheckDecode(publicKeyBase58Destination)

	//SCRIPTPUBKEY + LENGTH
	scriptPubKeyLength := 4 + len(publicKeyDestinationBytes)
	scriptPubKey := make([]byte, scriptPubKeyLength)

	scriptPubKey[0] = 118 //OP_DUP
	scriptPubKey[1] = 169 //OP_HASH160
	for i := 2; i < scriptPubKeyLength-2; i++ {
		scriptPubKey[i] = publicKeyDestinationBytes[i-2]
	}
	scriptPubKey[scriptPubKeyLength-2] = 136 //OP_EQUALVERIFY
	scriptPubKey[scriptPubKeyLength-1] = 172 //OP_CHECKSIG

	//LOCKTIMEFIELD
	lockTimeField, err := hex.DecodeString("00000000")
	if err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer
	buffer.Write(version)
	buffer.Write(inputs)
	buffer.Write(inputTransactionBytesReversed)
	buffer.Write(outputIndex)
	buffer.WriteByte(byte(scriptSigLength))
	buffer.Write(scriptSig)
	buffer.Write(sequence)
	buffer.Write(numOutputs)
	buffer.Write(satoshiBytes)
	buffer.WriteByte(byte(scriptPubKeyLength))
	buffer.Write(scriptPubKey)
	buffer.Write(lockTimeField)

	return buffer.Bytes()

	//four byte version field
	//one byte for # of inputs
	//32 byte hash of thee transaction for which want to redeem an output 7756bf4ed3b495adb63e05c02398e799c471b885d10523028b6e1b10f0ae181f
	//four byte field denoting the output index (01000000)
	//script sig lengh (1 byte)
	//script sig (24 bytes) OP_DUP OP_HASH160 9bf8cee4ce4532eab13454490dbdfb346d5e37f8 OP_EQUALVERIFY OP_CHECKSIG
	//four byte field which is always 0xfffffff (lol?)
	//1 byte varint containg number of outputs (01)
	//8 byte field containing the amount we want to redeem (left over is for miners)
	//1 byte for output script size
	//24 bytes for actual script OP_DUP OP_HASH160 9bf8cee4ce4532eab13454490dbdfb346d5e37f8 OP_EQUALVERIFY OP_CHECKSIG (NOT BASE 58.. need to decode FUUUUU)
	//four byte "lock time" field ? 0x00000000
	//four byte hashcode type. 0x00000001
}
