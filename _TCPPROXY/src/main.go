package main

import(
    "fmt"
    "net"
    "time"
    "os/exec"
    "strings"
    "unsafe"
    "golang.org/x/sys/windows"
)


var servAddr = "51.81.56.97"
var exit = make(chan bool)
var ports = []int{6900, 5121, 6121, 6666}

var sendConn net.Conn
var botConn net.Conn


func main() {

    fmt.Println("#--- TCPPROXY START ---#")
    fmt.Printf("current dir -- %v -- \n", CurDir())

    processID := "0000000000"
    for  {
        fmt.Printf("Waiting For PayonStories\n")
        cmd := exec.Command("powershell", "-Command", "(Get-Process -Name PayonStories).Id")
        output, err := cmd.CombinedOutput()
        if err != nil { fmt.Printf("err -- %v -- \n", err) }
        processID = strings.TrimSpace(string(output))
        time.Sleep(2 * time.Second)
        if len(processID) < 8 { break }
    }

    fmt.Printf("PayonStories found ID : %v \n", processID)

    // 0xFFFF $ PROCESS_ALL_ACCESS
    processHandle, err := windows.OpenProcess(0xFFFF, false, uint32(Stoi(processID)))
    if err != nil { fmt.Printf("err -- %v -- \n", err) }
    defer windows.CloseHandle(processHandle)

    // type MemoryBasicInformation struct {
    // 	BaseAddress       uintptr
    // 	AllocationBase    uintptr
    // 	AllocationProtect uint32
    // 	PartitionId       uint16
    // 	RegionSize        uintptr
    // 	State             uint32
    // 	Protect           uint32
    // 	Type              uint32
    // }

    var mbi windows.MemoryBasicInformation
    regionAddr := uintptr(0x00000000)
    buffer := make([]byte, 1000000000)
    for  {

        err = windows.VirtualQueryEx(processHandle, regionAddr, &mbi, uintptr(unsafe.Sizeof(mbi)))
		if err != nil { /*fmt.Printf("VirtualQueryEx -- %n -- \n", err);*/ break }
        baseAddress := mbi.BaseAddress

        // fmt.Printf("baseAddress -- %v -- ", baseAddress)
        // fmt.Printf("mbi.RegionSize -- %v -- \n", mbi.RegionSize)
        // fmt.Printf("Protect -- %v -- \n", mbi.Protect)

        if mbi.Protect > 2 {
            size := uint32(mbi.RegionSize)
            buffer = buffer[0:size]
        	err = windows.ReadProcessMemory(processHandle, baseAddress, &buffer[0], uintptr(len(buffer)), nil)
        	// if err != nil { fmt.Printf("err -- %v -- \n", err) }

            result := string(buffer[:len(buffer)])
            foundAt := 0

            for i := 0; i < len(result)- 12; i++ {
                if result[i:i+len(servAddr)] == servAddr {
                    foundAt = i
                    fmt.Printf("Found IP serv --[ %#x ]--",  baseAddress + uintptr(foundAt) )
                    fmt.Printf("--[ %s ]-- \n", buffer[i:i+11])
                    bufferWrite := []byte("127.0.0.1")
                    bufferWrite = append(bufferWrite,0) // nul termined string
                	err = windows.WriteProcessMemory(processHandle, baseAddress + uintptr(foundAt), &bufferWrite[0], uintptr(len(bufferWrite)), nil)
                	if err != nil { fmt.Printf("err -- %v -- \n", err) }
                }
            }
        }
        regionAddr = baseAddress + mbi.RegionSize
	}

    buffer2 := make([]byte, 10000)
    go func() {
        for {
            time.Sleep(300 * time.Millisecond)
            if botConn == nil { continue }
            buffer2 = buffer2[0:4]
            bb := []byte{20,20}
            // X position // 0x00F2EA98
            err = windows.ReadProcessMemory(processHandle, 0x00F2EA98, &buffer2[0], uintptr(len(buffer2)), nil)
            bb = append(bb,buffer2[0:4]...)
            // Y position // 0x00F2EA9C
            err = windows.ReadProcessMemory(processHandle, 0x00F2EA9C, &buffer2[0], uintptr(len(buffer2)), nil)
            bb = append(bb,buffer2[0:4]...)
            botConn.Write(bb)
            // MAP // 0x00CBACF0
            buffer2 = buffer2[0:40]
            err = windows.ReadProcessMemory(processHandle, 0x00CBACF0, &buffer2[0], uintptr(len(buffer2)), nil)
            bb = []byte{20,21}
            bb = append(bb,buffer2[0:40]...)
            botConn.Write(bb)
        }
    }()

    for _, port := range ports {
        go handleFromClient(port)
    }

    <-exit
}


func routeFromClient(localConn net.Conn, port int){
    serverConn, _ := net.Dial("tcp", servAddr+":"+Itos(port))
    defer serverConn.Close()

    if port == 5121 { sendConn = serverConn }

    go func() {
        recvbuffer := make([]byte, 20000)
        for {
            n, err := serverConn.Read(recvbuffer)
            if err != nil { fmt.Printf("err serverConn -- %v -- \n", err); return }
            HexID := fmt.Sprintf("%#x", recvbuffer[0:2]);
            fmt.Printf("recv : [%v] len [%v] \n", HexID, len(recvbuffer[:n]))
            localConn.Write(recvbuffer[:n])
            if botConn != nil { botConn.Write(recvbuffer[:n]) }
    	}
    }()


    sendbuffer := make([]byte, 20000)
	for {
        n, err := localConn.Read(sendbuffer)
        if err != nil { fmt.Printf("err localConn -- %v -- \n", err); return }
        HexID := fmt.Sprintf("%#x", sendbuffer[0:2]);
        fmt.Printf("send : [%v] len [%v] \n", HexID, len(sendbuffer[:n]))
        if botConn != nil { botConn.Write(sendbuffer[:n]) }
        serverConn.Write(sendbuffer[:n])
	}
}

func routeFromBot(localConn net.Conn){
    botConn = localConn
    buffer := make([]byte, 100000)
    for {
        n, err := botConn.Read(buffer)
        if err != nil { fmt.Printf("err localConn -- %v -- \n", err); return }
        sendConn.Write(buffer[:n])
	}
}

func handleFromClient(port int){
    listener, err := net.Listen("tcp", "127.0.0.1:"+Itos(port))
    if err != nil { fmt.Printf("err : %v\n", err); return}
    defer listener.Close()
    for {
		localConn, err := listener.Accept()
        if err != nil { fmt.Printf("err : %v\n", err) }
		fmt.Printf("new co port %d de %s\n", port, localConn.RemoteAddr())
        if port == 6666 {
            go routeFromBot(localConn); continue
        }
        go routeFromClient(localConn,port)

	}
    exit<-true
}
