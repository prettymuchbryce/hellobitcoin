/*
    Generating a BTC address using Go.
    Inspired by Ken Sherriff's blog post, "Bitcoins the hard way".
    
    Obviously do not actually use this to generate your BTC address as the random method
    is not cryptographically strong.
*/

package main

import (
    "fmt"
    "time"
    "strings"
    "math"
    "encoding/hex"
    //"math/big"
    //"github.com/tv42/base58"
    "math/rand"
)

func main() {
    //So far this is only a raw private key, not a WIF private key
    rand.Seed(time.Now().UTC().UnixNano())
    privateKey := randomBytes(32)
    fmt.Println(strings.ToUpper(hex.EncodeToString(privateKey)))
    //privateKeyInt := big.NewInt(0)
    //privateKeyInt.SetBytes(privateKey)
   // buf := base58.EncodeBig(nil, privateKeyInt)
}

func randomBytes(l int ) []byte {
    bytes := make([]byte, l)
    for i:=0 ; i<l ; i++ {
        bytes[i] = byte(randInt(0,math.MaxInt8))
    }
    return bytes
}

func randInt(min int, max int) int {
    return min + rand.Intn(max-min)
}