package main

import(
    "fmt"
    "strconv"
    "reflect"
    "encoding/binary"
    "strings"
    "time"
)

func parsePacket(fname string, args []reflect.Value){
    function := reflect.ValueOf(fctpack[fname])
    if function.Kind() == reflect.Func && fctpack[fname] != nil{
        function.Call(args)
    }
}


func fctpackInit()  {
    fctpack = map[string]func([]byte, []byte){}

    // for k := range packetsMap {
    //     fctpack[packetsMap[k].Ident] = func (HexID []byte, bb []byte)  {
    //         HexID1 := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(HexID))
    //         name := packetsMap[HexID1].Ident
    //         fmt.Printf("[%v][%v][%v]\t ",name, HexID1, len(bb)+2)
    //         fmt.Printf("-> [%v]\n", bb)
    //     }
    // }

    // fctpack["account_id"] = func (HexID []byte, bb []byte)  {
    //     myActorID = int(binary.LittleEndian.Uint32(bb[0:0+4]))
    //     fmt.Printf("#### myActorID-> [%v]\n", myActorID)
    // }

    fctpack["uknw_pck"] = func (HexID []byte, bb []byte)  {
        HexID1 := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(HexID))
        HexID2 := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(bb[0:2]))
        fmt.Printf("uknw_pck ####### [%v][%v] -> [%v] \t", HexID1, HexID2, len(bb)+2)
        fmt.Printf("-> [%v]\n", bb)
    }

    fctpack["chat_main"] = func (HexID []byte, bb []byte)  {}
    fctpack["send_sync_serv"] = func (HexID []byte, bb []byte)  {}
    fctpack["recv_sync_serv"] = func (HexID []byte, bb []byte)  {}
    fctpack["recv_self_move_to"] = func (HexID []byte, bb []byte)  {}


    fctpack["skill_use"] = func (HexID []byte, bb []byte)  {
        fmt.Printf(" ####### [%v][%v] -> [%v]\n","skill_use", len(bb)+2, bb)
    }

    fctpack["warp_portal_choice_recv"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf(" ####### [%v][%v] -> [%v]\n","warp_portal_choice_recv", len(bb)+2, bb)
        fmt.Printf("####### warp_portal_choice_recv -> [%s]\n", bb[4:])
    }


    fctpack["mem_data"] = func (HexID []byte, bb []byte)  {
        ii := 0
        MAP := strings.Split(string(bb[ii:ii+40]), ".rsw")[0]; ii += 40
        XPOS := binary.LittleEndian.Uint32(bb[ii:ii+4]); ii += 4
        YPOS := binary.LittleEndian.Uint32(bb[ii:ii+4]); ii += 4
        HPLEFT := binary.LittleEndian.Uint32(bb[ii:ii+4]); ii += 4
        HPMAX := binary.LittleEndian.Uint32(bb[ii:ii+4]); ii += 4
        WEIGHTMAX := binary.LittleEndian.Uint32(bb[ii:ii+4]); ii += 4
        WEIGHT := binary.LittleEndian.Uint32(bb[ii:ii+4]); ii += 4
        SPLEFT := binary.LittleEndian.Uint32(bb[ii:ii+4]); ii += 4
        SPMAX := binary.LittleEndian.Uint32(bb[ii:ii+4]); ii += 4

        curMap = MAP
        curCoord = Coord{X:int(XPOS),Y:int(YPOS)}
        HPLeft = int(HPLEFT)
        HPMax = int(HPMAX)
        maxWeight = int(WEIGHTMAX)
        weight = int(WEIGHT)
        SPLeft = int(SPLEFT)
        SPMax = int(SPMAX)
    }

    // fctpack["warp_portal_send"] = func (HexID []byte, bb []byte)  {
    //     fmt.Printf(" ####### [%v][%v] -> [%v]\n","warp_portal_send", len(bb)+2, bb)
    // }
    // fctpack["warp_portal_choice_send"] = func (HexID []byte, bb []byte)  {
    //     fmt.Printf(" ####### [%v][%v] -> [%v]\n","warp_portal_choice_send", len(bb)+2, bb)
    //     fmt.Printf("-> [%s]\n", bb[2:])
    // }



    fctpack["actor_status_active"] = func (HexID []byte, bb []byte)  {
        // HexID1 := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(HexID))
        // fmt.Printf("HexID1 [%v] \t", HexID1)
        // fmt.Printf(" ####### [%v][%v] -> [%v]\n","actor_status_active", len(bb)+2, bb)
        buffID := int(binary.LittleEndian.Uint16(bb[0:0+2]))
        target := int(binary.LittleEndian.Uint32(bb[2:2+4]))
        // timeLeft := int(binary.LittleEndian.Uint32(bb[7:7+4]))
        timeLeft := int(binary.LittleEndian.Uint32(bb[11:11+4]))
        if target == accountId {
            MUbuffList.Lock()
            buffList[buffID] = []int64{int64(timeLeft), time.Now().Unix()}
            MUbuffList.Unlock()
        }
    }

    fctpack["inventory_info"] = func (HexID []byte, bb []byte)  {
        HexID1 := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(HexID))
        // inventoryType := bb[2]
        if HexID1 == "0B09" {
            for ii := 3; ii < len(bb); ii+=34 {
                inventoryID := int(binary.LittleEndian.Uint16(bb[ii:ii+2]))
                itemID := int(binary.LittleEndian.Uint32(bb[ii+2:ii+2+4]))
                amount := int(binary.LittleEndian.Uint16(bb[ii+7:ii+7+2]))
                MUinventoryItems.Lock()
                inventoryItems[inventoryID] = Item{ ItemID:itemID, Coords:Coord{X:0,Y:0}, Amount:amount}
                MUinventoryItems.Unlock()
            }
        }
        if HexID1 == "0B0A" {
            for ii := 3; ii < len(bb); ii+=67 {
                inventoryID := int(binary.LittleEndian.Uint16(bb[ii:ii+2]))
                itemID := int(binary.LittleEndian.Uint32(bb[ii+2:ii+2+4]))
                MUinventoryItems.Lock()
                inventoryItems[inventoryID] = Item{ ItemID:itemID, Coords:Coord{X:0,Y:0}, Amount:1}
                MUinventoryItems.Unlock()
            }
        }
    }



    fctpack["inventory_item_added"] = func (HexID []byte, bb []byte)  {
        inventoryID := int(binary.LittleEndian.Uint16(bb[0:0+2]))
        amount := int(binary.LittleEndian.Uint16(bb[2:2+2]))
        itemID := int(binary.LittleEndian.Uint32(bb[4:4+4]))
        MUinventoryItems.Lock()
        if ii, exist := inventoryItems[inventoryID]; exist {
            ii.Amount += amount
            inventoryItems[inventoryID] = ii
        }else{
            inventoryItems[inventoryID] = Item{ ItemID:itemID, Coords:Coord{X:0,Y:0}, Amount:amount}
        }
        MUinventoryItems.Unlock()
    }

    fctpack["inventory_item_removed"] = func (HexID []byte, bb []byte)  {
        inventoryID := int(binary.LittleEndian.Uint16(bb[0:0+2]))
        amount := int(binary.LittleEndian.Uint16(bb[2:2+2]))
        if ii, exist := inventoryItems[inventoryID]; exist {
            ii.Amount -= amount
            if ii.Amount <= 0 {
                MUinventoryItems.Lock()
                delete(inventoryItems, inventoryID)
                MUinventoryItems.Unlock()
            }else{
                inventoryItems[inventoryID] = ii
            }
        }
    }

    fctpack["item_used"] = func (HexID []byte, bb []byte)  {
        inventoryID := int(binary.LittleEndian.Uint16(bb[0:0+2]))
        // itemID := int(binary.LittleEndian.Uint32(bb[2:2+4]))
        amountLeft := int(binary.LittleEndian.Uint16(bb[10:10+2]))
        if ii, exist := inventoryItems[inventoryID]; exist {
            ii.Amount = amountLeft
            inventoryItems[inventoryID] = ii
        }
    }

    fctpack["item_appear"] = func (HexID []byte, bb []byte)  {
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
        targetItemLooted = int(byteArrayToUInt32(bb[0:4]))
        MUgroundItems.Lock()
        delete(groundItems, targetItemLooted)
        MUgroundItems.Unlock()
    }

    fctpack["actor_moving_interrupt"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf(" ####### [%v][%v] -> [%v]\n","actor_moved", len(bb)+2, bb)
        mapID := int(binary.LittleEndian.Uint32(bb[0:4]))
        x := int(binary.LittleEndian.Uint16(bb[4:4+2]))
        y := int(binary.LittleEndian.Uint16(bb[6:6+2]))
        MUmobList.Lock()
        if mm, exist := mobList[mapID]; exist {
            mm.Coords.X = x; mm.Coords.Y = y
            mobList[mapID] = mm
        }
        MUmobList.Unlock()
    }

    fctpack["actor_moving"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf(" ####### [%v][%v] -> [%v]\n","actor_moved", len(bb)+2, bb)
        mapID := int(binary.LittleEndian.Uint32(bb[0:4]))
        MUmobList.Lock()
        fromto := bits48ToCoords(bb[4:4+6])
        if mm, exist := mobList[mapID]; exist {
            mm.Coords.X = fromto[2]; mm.Coords.Y = fromto[3]
            mobList[mapID] = mm
        }
        MUmobList.Unlock()
    }

    // ## type : 5 = mob / 6 = npc
    fctpack["actor_appear"] = func (HexID []byte, bb []byte)  {
        mapID := int(binary.LittleEndian.Uint32(bb[3:3+4]))
        mobID := int(binary.LittleEndian.Uint16(bb[21:21+4]))
        actorType := byte(bb[2])
        cc := Coord{X:0,Y:0}
        index := 0
        if sliceEqual(bb[0:2],[]byte{114,0}){ index = 65 }
        if sliceEqual(bb[0:2],[]byte{108,0}){ index = 61 }
        bc := bits24ToCoords(bb[index:index+3])
        cc.X = bc[0]; cc.Y = bc[1];
        if actorType == 5  {
            MUmobList.Lock()
            mobList[mapID] = Mob{ MobID:mobID, Coords:cc }
            MUmobList.Unlock()
        }
    }

    fctpack["actor_dead_disapear"] = func (HexID []byte, bb []byte)  {
        mapID := int(binary.LittleEndian.Uint32(bb[0:0+4]))
        if bb[4] == 1 { targetMobDead = mapID }
        MUmobList.Lock()
        delete(mobList, targetMobDead)
        MUmobList.Unlock()
    }



}

// #######################
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
