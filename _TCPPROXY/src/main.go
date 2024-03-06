package main

import(
    "fmt"
    "net"
    "time"
    "os/exec"
    "strings"
    "unsafe"
    "golang.org/x/sys/windows"
    "encoding/binary"
)

var(
    servAddr = "51.81.56.97"
    exit = make(chan bool)
    ports = []int{6900, 5121, 6121, 6666, 6667}

    sendConn net.Conn
    botConn net.Conn
    clientConn net.Conn
)

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
            time.Sleep(200 * time.Millisecond)
            if botConn == nil { continue }

            bb := []byte{20,20}

            //
            // // XPOS ## 0x00F2EA98
            // windows.ReadProcessMemory(processHandle, 0x00F2EA98, &buffer2[0], uintptr(4), nil)
            // bb = append(bb,buffer2[0:4]...)
            //
            // // YPOS ## 0x00F2EA9C
            // windows.ReadProcessMemory(processHandle, 0x00F2EA9C, &buffer2[0], uintptr(4), nil)
            // bb = append(bb,buffer2[0:4]...)

            // MOVESPEED ## 0x00F4235C
            windows.ReadProcessMemory(processHandle, 0x00F4235C, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // BASEXPMAX ## 0x00F422A0
            windows.ReadProcessMemory(processHandle, 0x00F422A0, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // BASEEXP ## 0x00F42298
            windows.ReadProcessMemory(processHandle, 0x00F42298, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // JOBXPMAX ## 0x00F422A8
            windows.ReadProcessMemory(processHandle, 0x00F422A8, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // JOBEXP ## 0x00F422B0
            windows.ReadProcessMemory(processHandle, 0x00F422B0, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // CHARNAME ## 0x00F48798
            windows.ReadProcessMemory(processHandle, 0x00F48798, &buffer2[0], uintptr(24), nil)
            bb = append(bb,buffer2[0:24]...)

            // BASELV ## 0x00F422B8
            windows.ReadProcessMemory(processHandle, 0x00F422B8, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // JOBLV ## 0x00F422C0
            windows.ReadProcessMemory(processHandle, 0x00F422C0, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // ZENY ## 0x00F42358
            windows.ReadProcessMemory(processHandle, 0x00F42358, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // MAP ## 0x00CBACF0
            windows.ReadProcessMemory(processHandle, 0x00CBACF0, &buffer2[0], uintptr(24), nil)
            bb = append(bb,buffer2[0:24]...)

            // HPLEFT ## 0x00F45E54
            windows.ReadProcessMemory(processHandle, 0x00F45E54, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // HPMAX ## 0x00F45E58
            windows.ReadProcessMemory(processHandle, 0x00F45E58, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // WEIGHTMAX ## 0x00F42364
            windows.ReadProcessMemory(processHandle, 0x00F42364, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // WEIGHT ## 0x00F42368
            windows.ReadProcessMemory(processHandle, 0x00F42368, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // SP ## 0x00F45E5C
            windows.ReadProcessMemory(processHandle, 0x00F45E5C, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            // MAXSP ## 0x00F45E60
            windows.ReadProcessMemory(processHandle, 0x00F45E60, &buffer2[0], uintptr(4), nil)
            bb = append(bb,buffer2[0:4]...)

            if botConn != nil { botConn.Write(bb) }
        }
    }()

    for _, port := range ports {
        go handleFromClient(port)
    }

    <-exit
}


func routeFromClient(tcpConn net.Conn, port int){
    serverConn, _ := net.Dial("tcp", servAddr+":"+Itos(port))
    defer serverConn.Close()

    if port == 5121 { sendConn = serverConn; clientConn = tcpConn}

    go func() {
        recvbuffer := make([]byte, 100000)
        for {
            n, err := serverConn.Read(recvbuffer)
            if err != nil { fmt.Printf("err serverConn -- %v -- \n", err); return }
            HexID := binary.LittleEndian.Uint16(recvbuffer[0:2]);
            fmt.Printf("recv : [%04X] len [%v] \n", HexID, len(recvbuffer[:n]))
            tcpConn.Write(recvbuffer[:n])
            if botConn != nil { botConn.Write(recvbuffer[:n]) }
    	}
    }()


    sendbuffer := make([]byte, 100000)
	for {
        n, err := tcpConn.Read(sendbuffer)
        if err != nil { fmt.Printf("err tcpConn -- %v -- \n", err); return }
        HexID := binary.LittleEndian.Uint16(sendbuffer[0:2]);
        fmt.Printf("send : [%04X] len [%v] \n", HexID, len(sendbuffer[:n]))
        if botConn != nil { botConn.Write(sendbuffer[:n]) }
        serverConn.Write(sendbuffer[:n])
	}
}

func routeToServ(tcpConn net.Conn){
    botConn = tcpConn
    buffer := make([]byte, 100000)
    for {
        n, err := botConn.Read(buffer)
        if err != nil { fmt.Printf("err tcpConn -- %v -- \n", err); return }
        sendConn.Write(buffer[:n])
	}
}

func routeToClient(tcpConn net.Conn){
    buffer := make([]byte, 100000)
    for {
        n, err := tcpConn.Read(buffer)
        if err != nil { fmt.Printf("err tcpConn -- %v -- \n", err); return }
        clientConn.Write(buffer[:n])
	}
}

func handleFromClient(port int){
    listener, err := net.Listen("tcp", "127.0.0.1:"+Itos(port))
    if err != nil { fmt.Printf("err : %v\n", err); return}
    defer listener.Close()
    for {
		tcpConn, err := listener.Accept()
        if err != nil { fmt.Printf("err : %v\n", err) }
		fmt.Printf("new co port %d de %s\n", port, tcpConn.RemoteAddr())
        if port == 6666 {
            go routeToServ(tcpConn); continue
        }
        if port == 6667 {
            go routeToClient(tcpConn); continue
        }
        go routeFromClient(tcpConn,port)

	}
    exit<-true
}
