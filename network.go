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
    //"io"
    "math/rand"
    "net"
    "bufio"
    "bytes"
    "time"
    //"log"
    "strconv"
    "crypto/sha256"
    //secp256k1 "github.com/haltingstate/secp256k1-go"
)

//0xDAB5BFFA <-- magic for testnet
//0xd9b4bef9 <-- magic for mainnet
//https://en.bitcoin.it/wiki/Protocol_specification#Message_types

//54.210.107.2:18333

func main() {
    //var transaction string
    //flag.StringVar(&transaction, "transaction", "", "")
    //flag.Parse()

    conn, err := net.Dial("tcp", "95.85.39.28:18333")
    if err != nil {
        fmt.Println(err)
        // handle error
    }

    versionMessage := makeMessage("DAB5BFFA", "version", getVersionMessage())
    n, err := conn.Write(versionMessage)
    if err != nil {
        fmt.Println(err,n)
        // handle error
    }

    fmt.Println(hex.EncodeToString(versionMessage))

    fmt.Println(bufio.NewReader(conn).ReadByte())

    /*fmt.Println(n)
    var buf bytes.Buffer
    io.Copy(&buf, conn)
    response := make([]byte, buf.Len())
    conn.Read(response)*/

//    fmt.Println(buf.Len())
  //  fmt.Println(response)
}

func makeMessage(magic string, command string, payload []byte) []byte {
    magicBytes, err := hex.DecodeString(magic)
    if err != nil {
        fmt.Println(err)
    }

    shaHash := sha256.New()
    shaHash.Write(payload)
    shaHashFirst := shaHash.Sum(nil)

    shaHash2 := sha256.New()
    shaHash2.Write(shaHashFirst)
    hashedPayload := shaHash2.Sum(nil)

    checksum := hashedPayload[0:4]

    length := uint32(len(payload))
    lengthBytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(lengthBytes, length)

    commandBytes := make([]byte, 12)
    for i := 0; i < 12; i++ {
        if i >= len(command) {
            commandBytes[i] = 0
        } else {
            commandBytes[i] = command[i]
        }
    }

    buffer := new(bytes.Buffer)
    //buffer.Write(magicBytes)
    binary.Write(buffer, binary.LittleEndian, magicBytes)
    binary.Write(buffer, binary.LittleEndian, magicBytes)
    buffer.Write(commandBytes)
    buffer.Write(lengthBytes)
    buffer.Write(checksum)
    buffer.Write(payload)
    return buffer.Bytes()
}

func getNetworkAddress(ip string) []byte {
    //timestamp := make([]byte, 4)
    //binary.LittleEndian.PutUint32(timestamp, uint32(time.Now().Unix()))

    services, err := hex.DecodeString("0100000000000000")
    if err != nil {
        fmt.Println(err)
    }

    ipv4Strings := strings.Split(ip, ".")
    ipv4Bytes := make([]byte, 4)

    for i := 3; i >= 0; i-- {
        ipByte, err := strconv.Atoi(ipv4Strings[i])
        if err != nil {
            fmt.Println(err)
        }

        ipv4Bytes[i] = byte(ipByte)
    }

    var ipv64 bytes.Buffer
    prefix, err := hex.DecodeString("00000000000000000000FFFF")
    if err != nil {
        fmt.Println(err)
    }
    ipv64.Write(prefix)
    ipv64.Write(ipv4Bytes)

    port := make([]byte, 2)
    binary.BigEndian.PutUint16(port, uint16(18333))

    var networkAddressBuffer bytes.Buffer

    //networkAddressBuffer.Write(timestamp)
    networkAddressBuffer.Write(services)
    networkAddressBuffer.Write(ipv64.Bytes())
    networkAddressBuffer.Write(port)

    return networkAddressBuffer.Bytes()
}

func getVersionMessage() []byte {
    version, err := hex.DecodeString("62EA0000")
    if err != nil {
        fmt.Println(err)
    }

    services, err := hex.DecodeString("0100000000000000")
    if err != nil {
        fmt.Println(err)
    }

    timestamp := make([]byte, 8)
    binary.PutVarint(timestamp, int64(time.Now().Unix()))

    fmt.Println(time.Now().Unix())

    /*
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        fmt.Println(err)
    }

    var ip string
 
    for _, a := range addrs {
        if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                ip = ipnet.IP.String()
            }
        }
    }*/

    addrRecv := getNetworkAddress("95.85.39.28")
    addrFrom := getNetworkAddress("76.102.229.234") //me

    nonce := make([]byte, 8)
    for i := 0; i < 8; i++ {
        nonce[i] = byte(randInt(0, math.MaxUint8))
    }

    userAgent, err := hex.DecodeString("00")
    if err != nil {
        fmt.Println(err)
    }

    startHeight, err := hex.DecodeString("00000000")
    if err != nil {
        fmt.Println(err)
    }


    var buffer bytes.Buffer

    buffer.Write(version)
    buffer.Write(services)
    buffer.Write(timestamp)
    buffer.Write(addrRecv)
    buffer.Write(addrFrom)
    buffer.Write(nonce)
    buffer.Write(userAgent)
    buffer.Write(startHeight)

    return buffer.Bytes()

}

func randInt(min int, max int) uint8 {
    rand.Seed(time.Now().UTC().UnixNano())
    return uint8(min + rand.Intn(max-min))
}