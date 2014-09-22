package main

import (
    //"flag"
    //"bytes"
    "encoding/binary"
    "encoding/hex"
    "fmt"
    "strings"
    //"os"
    "math"
   // "io/ioutil"
    "io"
    "math/rand"
    "net"
    //"bufio"
    "bytes"
    "time"
    //"log"
    "strconv"
    "crypto/sha256"
    //secp256k1 "github.com/haltingstate/secp256k1-go"
)
//0xDAB5BFFA <-- magic for testnet or FABFB5DA or 0B110907 or 0709110b
//0xd9b4bef9 <-- magic for mainnet
//https://en.bitcoin.it/wiki/Protocol_specification#Message_types

//54.210.107.2:18333

var address string = "54.210.107.2"

func main() {
    //var transaction string
    //flag.StringVar(&transaction, "transaction", "", "")
    //flag.Parse()

    bufaddr := new(bytes.Buffer)
    bufaddr.WriteString(address)
    bufaddr.WriteString(":18333")

    conn, err := net.Dial("tcp", 'http://www.google.com:80')
    if err != nil {
        fmt.Println(err)
        // handle error
    }

    versionMessage := makeMessage("0709110b", "version", getVersionMessage())

    n, err := conn.Write(versionMessage)
    if err != nil {
        fmt.Println(err,n)
        // handle error
    }


   /* status, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println(err)
        // handle error
    }*/

    fmt.Println(n)

    fmt.Println(hex.EncodeToString(versionMessage))

   // fmt.Println(bufio.NewReader(conn).ReadByte())

    //fmt.Println(n)
    var buf bytes.Buffer
    io.Copy(&buf, conn)
    response := make([]byte, buf.Len())
    conn.Read(response)

    fmt.Println(buf.Len())
    fmt.Println(response)
}
