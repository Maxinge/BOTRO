package main

import(
    "fmt"
    "time"
    "encoding/binary"
    // "math"
    // "math/rand"
)



func botLoop() {

    addWait(3000)

    for {

        addWait(50)
        for{
            if needWait <= 0 { break }
            time.Sleep(time.Duration(50) * time.Millisecond)
            needWait -= 50
        }

        now = time.Now()


        // fmt.Printf("# SIT # %v \n",SIT)

        // #####################################################################

        if (float32(SPLEFT)/float32(SPMAX)*100) <= float32(10) {
        if !SIT {
            sendToServer("0437", []byte{0,0,0,0,2})
        }}

        if (float32(SPLEFT)/float32(SPMAX)*100) >= float32(100) {
        if SIT {
            sendToServer("0437", []byte{0,0,0,0,3})
        }}

        if SIT { addWait(1000); continue }

        scroll := itemInInventory(7433,1)
        tail := itemInInventory(904,3) // [82 47 0 0] // fire
        snail := itemInInventory(946,3) // [83 47 0 0] // water
        horn := itemInInventory(947,3) // [84 47 0 0] // earth
        rainbow := itemInInventory(1013,3) // [85 47 0 0] // wind

        if scroll > -1 && tail > -1 {
            sendToServer("0438", []byte{1,0,239,3,85,188,30,0})
            time.Sleep(time.Duration(200) * time.Millisecond)
            sendToServer("01AE", []byte{82,47,0,0})
            time.Sleep(time.Duration(200) * time.Millisecond)
            refreshUI()
            time.Sleep(time.Duration(1000) * time.Millisecond)
        }
        if scroll > -1 && snail > -1 {
            sendToServer("0438", []byte{1,0,239,3,85,188,30,0})
            time.Sleep(time.Duration(200) * time.Millisecond)
            sendToServer("01AE", []byte{83,47,0,0})
            time.Sleep(time.Duration(200) * time.Millisecond)
            refreshUI()
            time.Sleep(time.Duration(1000) * time.Millisecond)
        }
        if scroll > -1 && horn > -1 {
            sendToServer("0438", []byte{1,0,239,3,85,188,30,0})
            time.Sleep(time.Duration(200) * time.Millisecond)
            sendToServer("01AE", []byte{84,47,0,0})
            time.Sleep(time.Duration(200) * time.Millisecond)
            refreshUI()
            time.Sleep(time.Duration(1000) * time.Millisecond)
        }
        if scroll > -1 && rainbow > -1 {
            sendToServer("0438", []byte{1,0,239,3,85,188,30,0})
            time.Sleep(time.Duration(200) * time.Millisecond)
            sendToServer("01AE", []byte{85,47,0,0})
            time.Sleep(time.Duration(200) * time.Millisecond)
            refreshUI()
            time.Sleep(time.Duration(1000) * time.Millisecond)
        }

    }

    fmt.Printf("# lel # %v \n","lel")
}

func refreshUI(){
    arrayBin := []byte{}
    XBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(XBin, uint16(XPOS))
    YBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(YBin, uint16(YPOS))
    arrayBin = append(arrayBin,[]byte{109,111,114,111,99,99,46,103,97,116,0,0,0,0,0,0}...)
    arrayBin = append(arrayBin,XBin...)
    arrayBin = append(arrayBin,YBin...)
    sendToClient("0091",arrayBin)
}
func addWait(nw int){ if nw > needWait { needWait = nw } }

func sendToClient(hexID string,data []byte){
    var ii uint16
	fmt.Sscanf(hexID, "%x", &ii)
    bb := []byte{byte(ii), byte(ii >> 8)}
    bb = append(bb,data...)
    proxyCoClient.Write(bb)
}

func sendToServer(hexID string,data []byte){
    var ii uint16
	fmt.Sscanf(hexID, "%x", &ii)
    bb := []byte{byte(ii), byte(ii >> 8)}
    bb = append(bb,data...)
    proxyCo.Write(bb)
}

func itemInInventory(id int, amount int)  int{
    for kk,ii := range inventoryItems {
        if ii.ItemID == id && ii.Amount >= amount   { return kk }
    }
    return -1
}

func sendUseItem(id int){
    arrayBin := []byte{}
    inventoryIDBin := make([]byte, 2);
    binary.LittleEndian.PutUint16(inventoryIDBin, uint16(id))
    accountIDBin := make([]byte, 4) ;
    binary.LittleEndian.PutUint32(accountIDBin, uint32(accountID))
    arrayBin = append(arrayBin,inventoryIDBin...)
    arrayBin = append(arrayBin,accountIDBin...)
    sendToServer("0439", arrayBin)
}

func sendUseSkill(id int, lv int, target int){
    arrayBin := []byte{}
    skillIDBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(skillIDBin, uint16(id))
    skillLVBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(skillLVBin, uint16(lv))
    targetBin := make([]byte, 4)
    binary.LittleEndian.PutUint32(targetBin, uint32(target))
    arrayBin = append(arrayBin,skillLVBin...)
    arrayBin = append(arrayBin,skillIDBin...)
    arrayBin = append(arrayBin,targetBin...)
    sendToServer("0438", arrayBin)
}

func sendWarpPortal(lv int, x int, y int){
    arrayBin := []byte{}
    skillLVBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(skillLVBin, uint16(lv))
    skillId := []byte{27, 0}
    XBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(XBin, uint16(x))
    YBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(YBin, uint16(y))
    arrayBin = append(arrayBin,skillLVBin...)
    arrayBin = append(arrayBin,skillId...)
    arrayBin = append(arrayBin,XBin...)
    arrayBin = append(arrayBin,YBin...)
    arrayBin = append(arrayBin,byte(0))
    sendToServer("0AF4", arrayBin)
}

func sendWarpPortalConfirm(choice string){
    arrayBin := []byte{}
    Id := []byte{27, 0}
    byteStr := []byte(choice)
    arrayBin = append(arrayBin, Id...)
    arrayBin = append(arrayBin, byteStr[0:16]...)
    sendToServer("011B", arrayBin)
}

func cleanBuffList()  {
    for kk,vv  := range buffList {
        tBuffTot := vv[0]
        tBuff := time.Unix(vv[1], 0)
        timeLeft := now.Sub(tBuff)
        if (tBuffTot - timeLeft.Milliseconds()) < 0 { delete(buffList, kk); }
    }
}
