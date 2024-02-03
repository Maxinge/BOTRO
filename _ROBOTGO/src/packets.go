package main

import(
    "fmt"
    "strconv"
    "strings"
)


func fctpackInit()  {
    fctpack = map[string]func([]byte, []byte){}

    fctpack["mem_coord"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf("mem_coord [%v] \t",fmt.Sprintf("%#x", HexID))
        // fmt.Printf("bb -> \t[%v]\n",bb)
        cx := byteArrayToUInt32(bb[0:4])
        cy := byteArrayToUInt32(bb[4:8])
        curCoord = Coord{X:int(cx),Y:int(cy)}
    }
    fctpack["mem_map"] = func (HexID []byte, bb []byte)  {
        curMap = strings.Split(string(bb), ".rsw")[0]
    }

    fctpack["actor_moved"] = func (HexID []byte, bb []byte)  {
        rr := splitBitsArray(bb,[]byte{135,0})
        for _,vv := range rr {
            mapID := int(vv[0])
            fromto := bits48ToCoords(vv[4:4+6])
            if mm, exist := mobList[mapID]; exist {
                mm.Coords.X = fromto[2]
                mm.Coords.Y = fromto[3]
                mobList[mapID] = mm
            }
        }
    }

    fctpack["actor_appear"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf("actor_appear [%v] \t",fmt.Sprintf("%#x", HexID))
        // fmt.Printf("bb -> \t[%v] \n",bb)
        rr := getActorsFromArray(bb)
        for _,vv := range rr {
            // fmt.Printf("vv -> \t[%v]\n",vv)
            if isMob(vv[0:3]) {
                mapID := int(vv[3])
                mobID := int(byteArrayToUInt16(vv[21:21+2]))
                sss := splitBitsArray(vv,[]byte{255,255,255,255,255,255,255,255})
                mobName := strings.Replace(string(sss[1]),"\u0000", "", -1)
                cc := Coord{X:0,Y:0} ; index := 0
                if sliceEqual(vv[0:3], []byte{108,0,5}) { index = 61 }
                if sliceEqual(vv[0:3], []byte{114,0,5}) { index = 65 }
                bb := bits24ToCoords(vv[index:index+3])
                cc.X = bb[0]; cc.Y = bb[1];
                mobList[mapID] = Mob{ MobID:mobID, Name:mobName, Coords:cc }
            }
        }
    }

    fctpack["load_map_data"] = fctpack["actor_appear"]

    fctpack["actor_something_happen"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf("actor_something_happen [%v] \t",fmt.Sprintf("%#x", HexID))
        // fmt.Printf("bb -> \t[%v] \n",bb)
        mapID := int(bb[0])
        if bb[4] == 1{ // isDead
        if _, exist := mobList[mapID]; exist {
            delete(mobList, mapID)
            // fmt.Printf("mapID is ded-> \t[%v]\n",mapID)
        }}
    }

    fctpack["item_appear"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf("item_appear [%v] \n",fmt.Sprintf("%#x", HexID))
        rrr := splitBitsArray(bb,[]byte{221,10})
        for _, vv := range rrr {
            parseItem(vv)
        }
    }

    fctpack["mob_info"] = func (HexID []byte, bb []byte)  {
        // fmt.Printf("mob_info [%v] \n",fmt.Sprintf("%#x", HexID))
        rr := splitBitsArray(bb,[]byte{255,255})
        if len(rr) > 1{
            rrr := splitBitsArray(rr[1],[]byte{221,10})
            for ii := 1; ii < len(rrr) ; ii++ {
                parseItem(rrr[ii])
            }
        }
    }

    fctpack["use_skill"] = func (HexID []byte, bb []byte)  {
        fmt.Printf("use_skill [%v] \t",fmt.Sprintf("%#x", HexID))
        fmt.Printf("bb -> \t[%v]\n",bb)
    }

    fctpack["item_disappear"] = func (HexID []byte, bb []byte)  {
        fmt.Printf("item_disappear [%v] \t",fmt.Sprintf("%#x", HexID))
        fmt.Printf("bb -> \t[%v]\n",bb)
        for ii := 0; ii < len(bb); ii+=4 {
            delete(groundItems, int(byteArrayToUInt32(bb[ii:ii+4])))
        }
    }
    fctpack["try_loot_item"] = func (HexID []byte, bb []byte)  {
        fmt.Printf("try_loot_item [%v] \t",fmt.Sprintf("%#x", HexID))
        fmt.Printf("bb -> \t[%v]\n",bb)

    }
    fctpack["loot_item_confirm"] = func (HexID []byte, bb []byte)  {
        rrr := splitBitsArray(bb,[]byte{74,188,30,0})
        if len(rrr) > 1{
            for ii := 1; ii < len(rrr) ; ii++ {
                fmt.Printf("loot -> \t[%v]\n",rrr[ii][0:4])
                delete(groundItems, int(byteArrayToUInt32(rrr[ii][0:4])))
            }
        }
    }

    fctpack["uknw_greed"] = func (HexID []byte, bb []byte)  {
        rr := splitBitsArray(bb,[]byte{74,188,30,0,74,188,30,0})
        if len(rr) > 2{
            rrr := splitBitsArray(rr[2],[]byte{74,188,30,0})
            if len(rrr) > 1{
                for ii := 1; ii < len(rrr) ; ii++ {
                    fmt.Printf("loot -> \t[%v]\n",rrr[ii][0:4])
                    delete(groundItems, int(byteArrayToUInt32(rrr[ii][0:4])))
                }
            }
        }
    }

}

// #######################
func parseItem(item []byte){
    // fmt.Printf("mapID -> \t[%v]\n",item[0:4])
    mapID := int(byteArrayToUInt32(item[0:4]))
    itemID := int(byteArrayToUInt16(item[4:6]))
    x := int(byteArrayToUInt16(item[11:13]))
    y := int(byteArrayToUInt16(item[13:15]))
    amount := int(byteArrayToUInt16(item[17:19]))
    groundItems[mapID] = Item{ ItemID:itemID, Coords:Coord{X:x,Y:y}, Amount:amount}
}

func isMob(bb []byte) bool{
    flags := [][]byte{
        []byte{108,0,5},
        []byte{114,0,5},
    }
    for _,vv := range flags {
        if sliceEqual(vv, bb) { return true }
    }
    return false
}

func getActorsFromArray(bb []byte) [][]byte{
    flags := [][]byte{
        []byte{108,0,5}, // mobs
        []byte{114,0,5}, // mobs
        []byte{63,4,184}, // ???
        []byte{0,49,1}, // player
        []byte{164,1,3}, // ???
        []byte{108,0,0}, // player vending ?
    }
    var res [][]byte
    for i := 0; i <= len(bb)-3 ; i++ {
        if inArrayByte(flags,bb[i:i+3]) {
            j := 0
            for {
                if i+j >= len(bb){ break }
                if sliceEqual(bb[i+j:i+j+8], []byte{255,255,255,255,255,255,255,255}) {
                    res = append(res,bb[i:i+j+32])
                    i = i+j+32 ; break
                }
                j++
            }
        }
    }
    return res
}


func int16ToBitString(ii int) string {
    ss := ""
	for i := 16 - 1; i >= 0; i-- {
		bit := (ii >> uint(i)) & 1
        ss += fmt.Sprintf("%d", bit)
	}
	return ss
}

func byteArrayToUInt64(b []byte) uint64 {
	if len(b) < 4 { return 0 }
	result := uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
              uint64(b[3])<<32 | uint64(b[3])<<40 | uint64(b[3])<<48 | uint64(b[3])<<56
	return result
}

func byteArrayToUInt32(b []byte) uint32 {
	if len(b) < 4 { return 0 }
	result := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	return result
}

func byteArrayToUInt16(b []byte) uint16 {
	if len(b) < 2 { return 0 }
	result := uint16(b[0]) | uint16(b[1])<<8
	return result
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
func bitArrayToBitString(bb []byte) string{
    ss := ""
    for _, b := range bb {
        ss += fmt.Sprintf("%08b", b)
    }
    return ss
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
