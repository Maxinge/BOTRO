package main

import(
    "fmt"
    "strconv"
    "reflect"
    "encoding/binary"
    "strings"
)

func parsePacket(fname string, args []reflect.Value){
    function := reflect.ValueOf(fctpack[fname])
    if function.Kind() == reflect.Func && fctpack[fname] != nil{
        function.Call(args)
    }
}


func fctpackInit()  {
    fctpack = map[string]func([]byte, []byte){}

    for k := range packetsmap {
        fctpack[packetsmap[k].Ident] = func (HexID []byte, bb []byte)  {
            HexID1 := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(HexID))
            HexID2 := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(bb[0:2]))
            fmt.Printf("[%v][%v]   [%v]\t ", HexID1, HexID2, len(bb)+2)
            fmt.Printf("-> [%v]\n", bb)
        }
    }

    fctpack["uknw"] = func (HexID []byte, bb []byte)  {}

    fctpack["recv_self_move_to"] = func (HexID []byte, bb []byte)  {
        // cc := bits48ToCoords(bb[4:])
    }

    fctpack["chat_main"] = func (HexID []byte, bb []byte)  {}

    fctpack["send_self_move_to"] = func (HexID []byte, bb []byte)  {}
    fctpack["send_sync_serv"] = func (HexID []byte, bb []byte)  {}
    fctpack["recv_sync_serv"] = func (HexID []byte, bb []byte)  {}
    fctpack["actor_info"] = func (HexID []byte, bb []byte)  {}
    fctpack["stat_info"] = func (HexID []byte, bb []byte)  {}
    fctpack["actor_info_request"] = func (HexID []byte, bb []byte)  {}

    fctpack["item_appear"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf("[%v][%v] -> [%v]\n","item_appear", len(bb)+2, bb)
        mapID := int(binary.LittleEndian.Uint32(bb[0:0+4]))
        itemID := int(binary.LittleEndian.Uint16(bb[4:4+2]))
        x := int(binary.LittleEndian.Uint16(bb[11:13]))
        y := int(binary.LittleEndian.Uint16(bb[13:15]))
        amount := int(binary.LittleEndian.Uint16(bb[17:19]))
        MUgroundItems.Lock()
        groundItems[mapID] = Item{ ItemID:itemID, Coords:Coord{X:x,Y:y}, Amount:amount}
        MUgroundItems.Unlock()
    }

    fctpack["item_disappear"] = func (HexID []byte, bb []byte)  {
        MUgroundItems.Lock()
        for ii := 0; ii < len(bb); ii+=4 {
            delete(groundItems, int(binary.LittleEndian.Uint32(bb[ii:ii+4])))
        }
        MUgroundItems.Unlock()
    }

    fctpack["actor_action"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf("[%v][%v] -> [%v]\n","item_disappear", len(bb)+2, bb)
    }


    fctpack["mem_map"] = func (HexID []byte, bb []byte)  {
        curMap = strings.Split(string(bb), ".rsw")[0]
    }

    fctpack["mem_coord"] = func (HexID []byte, bb []byte)  {
        cx := binary.LittleEndian.Uint32(bb[0:4])
        cy := binary.LittleEndian.Uint32(bb[4:8])
        curCoord = Coord{X:int(cx),Y:int(cy)}
    }

    fctpack["actor_moved"] = func (HexID []byte, bb []byte)  {
        mapID := int(binary.LittleEndian.Uint32(bb[0:4]))
        MUmobList.Lock()
        fromto := bits48ToCoords(bb[4:4+6])
        if mm, exist := mobList[mapID]; exist {
            mm.Coords.X = fromto[2]
            mm.Coords.Y = fromto[3]
            mobList[mapID] = mm
        }
        MUmobList.Unlock()
    }

    fctpack["actor_appear"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf("[%v][%v] -> [%v]\n","actor_appear", len(bb)+2, bb)
        mapID := int(binary.LittleEndian.Uint32(bb[3:3+4]))
        mobID := int(binary.LittleEndian.Uint16(bb[21:21+4]))
        sss := splitBitsArray(bb,[]byte{255,255,255,255,255,255,255,255})
        mobName := ""
        if len(sss) > 1 { mobName = strings.Replace(string(sss[1]),"\u0000", "", -1) }
        cc := Coord{X:0,Y:0};
        index := 0
        if bb[0] == 108 { index = 61 }
        if bb[0] == 114 { index = 65 }
        bc := bits24ToCoords(bb[index:index+3])
        cc.X = bc[0]; cc.Y = bc[1];
        MUmobList.Lock()
        mobList[mapID] = Mob{ MobID:mobID, Name:mobName, Coords:cc }
        MUmobList.Unlock()
    }

    fctpack["actor_dead_disapear"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf("[%v][%v] -> [%v]\n","actor_dead_disapear", len(bb)+2, bb)
        mapID := int(binary.LittleEndian.Uint32(bb[0:0+4]))
        MUmobList.Lock()
        if bb[4] == 1{ // isDead
        if _, exist := mobList[mapID]; exist {
            targetMobDead = mapID
            delete(mobList, mapID)
            // fmt.Printf("mapID is ded-> \t[%v]\n",mapID)
        }}
        MUmobList.Unlock()
    }


    fctpack["loot_item_confirm"] = func (HexID []byte, bb []byte)  {
        rrr := splitBitsArray(bb,[]byte{74,188,30,0})
        if len(rrr) > 1{
            MUgroundItems.Lock()
            for ii := 1; ii < len(rrr) ; ii++ {
                targetItemLooted = int(byteArrayToUInt32(rrr[ii][0:4]))
                // fmt.Printf("loot -> \t[%v]\n",rrr[ii][0:4])
                delete(groundItems, int(byteArrayToUInt32(rrr[ii][0:4])))
            }
            MUgroundItems.Unlock()
        }
    }


}

// #######################
func int16ToBitString(ii int) string {
    ss := ""
	for i := 16 - 1; i >= 0; i-- {
		bit := (ii >> uint(i)) & 1
        ss += fmt.Sprintf("%d", bit)
	}
	return ss
}

func bitArrayToBitString(bb []byte) string{
    ss := ""
    for _, b := range bb {
        ss += fmt.Sprintf("%08b", b)
    }
    return ss
}

func coordsTo24Bits(x int, y int /*, direction int*/) []byte {
    // coords to (x,y), so in 3 bytes
    // packed in 2 x 10 bits trunks  + 00
    // those are not "bytes-aligned" in the packet
    ss := ""
    ss += int16ToBitString(x)[6:16] // 10
    ss += int16ToBitString(y)[6:16] // 10
    ss += "0000"
    r1, _ := strconv.ParseInt(ss[0:8], 2, 16)
    r2, _ := strconv.ParseInt(ss[8:16], 2, 16)
    r3, _ := strconv.ParseInt(ss[16:], 2, 16)
    return([]byte{ byte(r1), byte(r2), byte(r3) })
}

func bits48ToCoords(bb []byte) []int  {
    // coords from (x,y) -> to (x,y), so 4 bytes
    // packed in 4 x 10 bits trunks
    // those are not "bytes-aligned" in the packet
    ss := bitArrayToBitString(bb)
    r1, _ := strconv.ParseInt(ss[0:10], 2, 10)
    r2, _ := strconv.ParseInt(ss[10:20], 2, 10)
    r3, _ := strconv.ParseInt(ss[20:30], 2, 10)
    r4, _ := strconv.ParseInt(ss[30:40], 2, 10)
     return([]int{ int(r1), int(r2), int(r3), int(r4) })
}

func bits24ToCoords(bb []byte) []int  {
    // same as coordsTo24Bits be reversed
    ss := bitArrayToBitString(bb)
    r1, _ := strconv.ParseInt(ss[0:10], 2, 10)
    r2, _ := strconv.ParseInt(ss[10:20], 2, 10)
    return([]int{ int(r1), int(r2) })
}
