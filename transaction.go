package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/prettymuchbryce/hellobitcoin/base58check"
	secp256k1 "github.com/toxeus/go-secp256k1"
)

var flagPrivateKey string
var flagPublicKey string
var flagDestination string
var flagInputTransaction string
var flagInputIndex int
var flagSatoshis int

func main() {
	//This transaction code is not completely robust.
	//It expects that you have exactly 1 input transaction, and 1 output address.
	//It also expects that your transaction is a standard Pay To Public Key Hash (P2PKH) transaction.
	//This is the most common form used to send a transaction to one or multiple Bitcoin addresses.

	//Parse flags
	flag.StringVar(&flagPrivateKey, "private-key", "", "The private key of the bitcoin wallet which contains the bitcoins you wish to send.")
	flag.StringVar(&flagPublicKey, "public-key", "", "The public address of the bitcoin wallet which contains the bitcoins you wish to send.")
	flag.StringVar(&flagDestination, "destination", "", "The public address of the bitcoin wallet to which you wish to send the bitcoins.")
	flag.StringVar(&flagInputTransaction, "input-transaction", "", "An unspent input transaction hash which contains the bitcoins you wish to send. (Note: HelloBitcoin assumes a single input transaction, and a single output transaction for simplicity.)")
	flag.IntVar(&flagInputIndex, "input-index", 0, "The output index of the unspent input transaction which contains the bitcoins you wish to send. Defaults to 0 (first index).")
	flag.IntVar(&flagSatoshis, "satoshis", 0, "The number of bitcoins you wish to send as represented in satoshis (100,000,000 satoshis = 1 bitcoin). (Important note: the number of satoshis left unspent in your input transaction will be spent as the transaction fee.)")
	flag.Parse()

	//First we create the raw transaction.
	//In order to construct the raw transaction we need the input transaction hash,
	//the destination address, the number of satoshis to send, and the scriptSig
	//which is temporarily (prior to signing) the ScriptPubKey of the input transaction.
	tempScriptSig := createScriptPubKey(flagPublicKey)

	rawTransaction := createRawTransaction(flagInputTransaction, flagInputIndex, flagDestination, flagSatoshis, tempScriptSig)

	//After completing the raw transaction, we append
	//SIGHASH_ALL in little-endian format to the end of the raw transaction.
	hashCodeType, err := hex.DecodeString("01000000")
	if err != nil {
		log.Fatal(err)
	}

	var rawTransactionBuffer bytes.Buffer
	rawTransactionBuffer.Write(rawTransaction)
	rawTransactionBuffer.Write(hashCodeType)
	rawTransactionWithHashCodeType := rawTransactionBuffer.Bytes()

	//Sign the raw transaction, and output it to the console.
	finalTransaction := signRawTransaction(rawTransactionWithHashCodeType, flagPrivateKey)
	finalTransactionHex := hex.EncodeToString(finalTransaction)

	fmt.Println("Your final transaction is")
	fmt.Println(finalTransactionHex)
}

func createScriptPubKey(publicKeyBase58 string) []byte {
	publicKeyBytes := base58check.Decode(publicKeyBase58)

	var scriptPubKey bytes.Buffer
	scriptPubKey.WriteByte(byte(118))                 //OP_DUP
	scriptPubKey.WriteByte(byte(169))                 //OP_HASH160
	scriptPubKey.WriteByte(byte(len(publicKeyBytes))) //PUSH
	scriptPubKey.Write(publicKeyBytes)
	scriptPubKey.WriteByte(byte(136)) //OP_EQUALVERIFY
	scriptPubKey.WriteByte(byte(172)) //OP_CHECKSIG
	return scriptPubKey.Bytes()
}

func signRawTransaction(rawTransaction []byte, privateKeyBase58 string) []byte {
	//Here we start the process of signing the raw transaction.

	secp256k1.Start()
	privateKeyBytes := base58check.Decode(privateKeyBase58)
	var privateKeyBytes32 [32]byte

	for i := 0; i < 32; i++ {
		privateKeyBytes32[i] = privateKeyBytes[i]
	}

	//Get the raw public key
	publicKeyBytes, success := secp256k1.Pubkey_create(privateKeyBytes32, false)
	if !success {
		log.Fatal("Failed to convert private key to public key")
	}

	//Hash the raw transaction twice before the signing
	shaHash := sha256.New()
	shaHash.Write(rawTransaction)
	var hash []byte = shaHash.Sum(nil)

	shaHash2 := sha256.New()
	shaHash2.Write(hash)
	rawTransactionHashed := shaHash2.Sum(nil)

	//Sign the raw transaction
	signedTransaction, success := secp256k1.Sign(rawTransactionHashed, privateKeyBytes32, generateNonce())
	if !success {
		log.Fatal("Failed to sign transaction")
	}

	//Verify that it worked.
	verified := secp256k1.Verify(rawTransactionHashed, signedTransaction, publicKeyBytes)
	if !verified {
		log.Fatal("Failed to sign transaction")
	}

	secp256k1.Stop()

	hashCodeType, err := hex.DecodeString("01")
	if err != nil {
		log.Fatal(err)
	}

	//+1 for hashCodeType
	signedTransactionLength := byte(len(signedTransaction) + 1)

	var publicKeyBuffer bytes.Buffer
	publicKeyBuffer.Write(publicKeyBytes)
	pubKeyLength := byte(len(publicKeyBuffer.Bytes()))

	var buffer bytes.Buffer
	buffer.WriteByte(signedTransactionLength)
	buffer.Write(signedTransaction)
	buffer.WriteByte(hashCodeType[0])
	buffer.WriteByte(pubKeyLength)
	buffer.Write(publicKeyBuffer.Bytes())

	scriptSig := buffer.Bytes()

	//Return the final transaction
	return createRawTransaction(flagInputTransaction, flagInputIndex, flagDestination, flagSatoshis, scriptSig)
}

func generateNonce() [32]byte {
	var bytes [32]byte
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

func createRawTransaction(inputTransactionHash string, inputTransactionIndex int, publicKeyBase58Destination string, satoshis int, scriptSig []byte) []byte {
	//Create the raw transaction.

	//Version field
	version, err := hex.DecodeString("01000000")
	if err != nil {
		log.Fatal(err)
	}

	//# of inputs (always 1 in our case)
	inputs, err := hex.DecodeString("01")
	if err != nil {
		log.Fatal(err)
	}

	//Input transaction hash
	inputTransactionBytes, err := hex.DecodeString(inputTransactionHash)
	if err != nil {
		log.Fatal(err)
	}

	//Convert input transaction hash to little-endian form
	inputTransactionBytesReversed := make([]byte, len(inputTransactionBytes))
	for i := 0; i < len(inputTransactionBytes); i++ {
		inputTransactionBytesReversed[i] = inputTransactionBytes[len(inputTransactionBytes)-i-1]
	}

	//Output index of input transaction
	outputIndexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(outputIndexBytes, uint32(inputTransactionIndex))

	//Script sig length
	scriptSigLength := len(scriptSig)

	//sequence_no. Normally 0xFFFFFFFF. Always in this case.
	sequence, err := hex.DecodeString("ffffffff")
	if err != nil {
		log.Fatal(err)
	}

	//Numbers of outputs for the transaction being created. Always one in this example.
	numOutputs, err := hex.DecodeString("01")
	if err != nil {
		log.Fatal(err)
	}

	//Satoshis to send.
	satoshiBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(satoshiBytes, uint64(satoshis))

	//Script pub key
	scriptPubKey := createScriptPubKey(publicKeyBase58Destination)
	scriptPubKeyLength := len(scriptPubKey)

	//Lock time field
	lockTimeField, err := hex.DecodeString("00000000")
	if err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer
	buffer.Write(version)
	buffer.Write(inputs)
	buffer.Write(inputTransactionBytesReversed)
	buffer.Write(outputIndexBytes)
	buffer.WriteByte(byte(scriptSigLength))
	buffer.Write(scriptSig)
	buffer.Write(sequence)
	buffer.Write(numOutputs)
	buffer.Write(satoshiBytes)
	buffer.WriteByte(byte(scriptPubKeyLength))
	buffer.Write(scriptPubKey)
	buffer.Write(lockTimeField)

	return buffer.Bytes()
}
