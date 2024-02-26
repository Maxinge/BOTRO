package main

import(
    "fmt"
    "net"
    // "encoding/json"
    // "github.com/cimgui-go"
    // "io/ioutil"
    "strings"
    // "time"
    "encoding/binary"
    // "reflect"
)


var (
    err error
    proxyCo net.Conn
    proxyCoClient net.Conn
    exit = make(chan bool)
    // gatMaps = map[string]ROGatMap{}
    packetsMap = map[string]int{}
    profil map[string][]interface{}
    conf = map[string][]interface{}{}
    mobDB []map[string]interface{}
)


func main() {

    fmt.Println("#--- ROBOTGO START ---#")
    fmt.Printf("current dir -- %v -- \n", CurDir())



    fff := readFileString(CurDir()+"data/recvpackets.txt")
    ss := strings.Split(fff, "\n")
    for _,vv := range ss[:len(ss)-1] {
        sss := strings.Split(vv, " ")
        packetsMap[sss[0]] = Stoi(sss[1][:len(sss[1])-1])
    }


    // ########################
    // ########################

    proxyCo, err = net.Dial("tcp", "127.0.0.1:6666")
    if err != nil { fmt.Printf("err -- %v -- \n", err); return }
    defer proxyCo.Close()

    proxyCoClient, err = net.Dial("tcp", "127.0.0.1:6667")
    if err != nil { fmt.Printf("err -- %v -- \n", err); return }
    defer proxyCoClient.Close()


    go func() {
        buffer := make([]byte, 100000)
        for {
            n ,_ := proxyCo.Read(buffer)
            // if err != nil { fmt.Printf("err localConn -- %v -- \n", err); return }
            if n < 2 { continue }
            ii := -1
            for {
                ii += 1; if ii >= n { break }
                HexID := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(buffer[ii:ii+2]));
                plen := n-ii
                if _, exist := packetsMap[HexID]; exist {
                    plen = packetsMap[HexID]
                    if plen < 0 { plen = int(binary.LittleEndian.Uint16(buffer[ii+2:ii+4])) }
                    if plen <= 2 { plen = 2 }
                    parsePacket(buffer[ii:ii+plen])
                    ii += plen-1;
                    continue
                }
                fmt.Printf("####### uknw_pck [%v][%v] -> [%v] \n", HexID, plen, buffer[ii:ii+plen])
                break;
            }
        }
    }()

    go botLoop()

    <-exit
}
