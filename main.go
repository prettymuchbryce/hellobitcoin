package main

import (
    "fmt"
    "crypto/sha256"
    "time"
    "math"
    "log"
    "encoding/hex"
    "math/big"
    "github.com/tv42/base58"
    "math/rand"
)

func main() {
    //BTC private key is 256 bits of random data
    privateKey := generateExtendedPrivateKey()

    //Generate "extended" private key from private key
    //The difference between an private key and an extended
    //private key is this prefix, which determines which
    //network the key is of (real btc network, or test network)
    //EF is the testnet prefix
    //80 is the mainnet prefix
    prefixBytes, err := hex.DecodeString("EF")
    if err != nil {
        log.Fatal(err)
    }

    extendedPrivateKey := make([]byte, 33)
    extendedPrivateKey[0] = prefixBytes[0]
    for i:=1 ; i<33 ; i++ {
        //This is not "cryptographically random"
        extendedPrivateKey[i] = privateKey[i-1]
    }

    //Perform SHA-256 on the extended key twice
    var hash [32]byte = sha256.Sum256(extendedPrivateKey)
    var hashSlice []byte = make([]byte, 33)
    for i:=0 ; i<32 ; i++ {
        hashSlice[i] = hash[i]
    }

    hash2 := sha256.Sum256(hashSlice)

    //First 4 bytes if this double-sha'd byte array is the checksum
    checksum := hash2[0:4]

    //Append this checksum to the extended private key
    extendedPrivateKey = append(extendedPrivateKey, checksum[0], checksum[1], checksum[2], checksum[3])

    //Convert the extended private key to a big Int
    bigIntPrivateKey := big.NewInt(0)
    bigIntPrivateKey.SetBytes(extendedPrivateKey)

    //Encoded the big int extended private key into a Base58Checked string
    buf := base58.EncodeBig(nil, bigIntPrivateKey)

    //Print the private key in Wallet Import Format
    fmt.Println(string(buf))
}

func generateExtendedPrivateKey() []byte {
    bytes := make([]byte, 32)
    for i:=0 ; i<32 ; i++ {
        //This is not "cryptographically random"
        bytes[i] = byte(randInt(0,math.MaxUint8))
    }
    return bytes
}

func randInt(min int, max int) uint8 {
    rand.Seed(time.Now().UTC().UnixNano())
    return uint8(min + rand.Intn(max-min))
}