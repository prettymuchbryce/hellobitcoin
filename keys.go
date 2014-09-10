package main

import (
    "fmt"
    "crypto/sha256"
    "code.google.com/p/go.crypto/ripemd160"
    "time"
    "math"
    "log"
    "bytes"
    "encoding/hex"
    "math/big"
    secp256k1 "github.com/haltingstate/secp256k1-go"
    "hellobitcoin/base58"
    "math/rand"
)

func main() {
    //BTC private key is 256 bits of random data
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

    privateKeyWif := base58CheckEncode("EF", privateKey)

    publicKey := generatePublicKey(privateKey)

    fmt.Println(hex.EncodeToString(publicKey))

    //There is also a prefix on the public key
    //This is known as the Network ID Byte, or the version byte
    //6f is the testnet prefix
    //00 is the mainnet prefix
    publicKeyEncoded := base58CheckEncode("6f", publicKey)

    //Print the keys
    fmt.Println("---")
    fmt.Println("Your private key is")
    fmt.Println(privateKeyWif)

    fmt.Println("Your public key is")
    fmt.Println(publicKeyEncoded)
}

func base58CheckEncode(prefix string, byteData []byte) string {
    prefixBytes, err := hex.DecodeString(prefix)
    if err != nil {
        log.Fatal(err)
    }

    length := len(byteData)+1
    encoded := make([]byte, length)
    encoded[0] = prefixBytes[0]
    for i:=1; i<length; i++ {
        encoded[i] = byteData[i-1]
    }

    //Perform SHA-256 twice
    shaHash := sha256.New()
    shaHash.Write(encoded)
    var hash []byte = shaHash.Sum(nil)

    shaHash2 := sha256.New()
    shaHash2.Write(hash)
    hash2 := shaHash2.Sum(nil)

    //First 4 bytes if this double-sha'd byte array is the checksum
    checksum := hash2[0:4]

    //Append this checksum to the input bytes
    encodedChecksum := append(encoded, checksum[0], checksum[1], checksum[2], checksum[3])

    //base58 alone is not enough. We need to first count each of the zero bytes
    //which are at the beginning of the encodedCheckSum
    zeroBytes := 0
    for i:=0; i < len(encodedChecksum); i++ {
        if encodedChecksum[i] == 0 {
            zeroBytes += 1
        } else {
            break
        }
    }

    //Convert this checksum'd version to a big Int
    bigIntEncodedChecksum := big.NewInt(0)
    bigIntEncodedChecksum.SetBytes(encodedChecksum)

    //Encode the big int checksum'd version into a Base58Checked string
    base58EncodedChecksum := string(base58.EncodeBig(nil, bigIntEncodedChecksum))

    //Now for each zero byte we counted above we need to prepend a 1 to our
    //base58 encoded string. The rational behind this is that base58 removes 0's (0x00).
    //So bitcoin demands we add leading 0s back on as 1s.
    var buffer bytes.Buffer
    for i := 0; i < zeroBytes; i++ {
        buffer.WriteString("1")
    }

    buffer.WriteString(base58EncodedChecksum)

    return buffer.String()
}

func generatePublicKey(privateKey []byte) []byte {
    //Generate the public key from the private key.
    //Unfortunately golang ecdsa package does not include a
    //secp256k1 curve as this is fairly specific to bitcoin
    //as I understand it, so I have used this one by haltingstate.
    publicKeyBytes := secp256k1.UncompressedPubkeyFromSeckey(privateKey)

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