package main

import(
    "fmt"
    "strconv"
    "strings"
)


func fctpackInit()  {
    fctpack = map[string]func([]byte, []byte){}

    // fctpack["recv_self_move_to"] = func (HexID []byte, bb []byte)  {
    //     cc := bits48ToCoords(bb[4:])
    //     curCoord = Coord{X:cc[2],Y:cc[3]}
    //     // fmt.Printf("recv_self_move_to -> [%v]\n",string(bb))
    //     fmt.Printf("curCoord -> \t[%v]\n",curCoord)
    // }


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

    // fctpack["actor_moved"] = func (HexID []byte, bb []byte)  {
    //     rr := splitBitsArray(bb,[]byte{135,0})
    //     for _,vv := range rr {
    //         mapID := int(vv[0])
    //         MUmobList.Lock()
    //         fromto := bits48ToCoords(vv[4:4+6])
    //         if mm, exist := mobList[mapID]; exist {
    //             mm.Coords.X = fromto[2]
    //             mm.Coords.Y = fromto[3]
    //             mobList[mapID] = mm
    //         }
    //         MUmobList.Unlock()
    //     }
    // }

    // fctpack["actor_appear"] = func (HexID []byte, bb []byte)  {
    //     parseMob(bb)
    // }
    //
    // fctpack["load_map_data"] = func (HexID []byte, bb []byte)  {
    //     rr := getActorsFromArray(bb)
    //     for _,vv := range rr {
    //         fmt.Printf("vv -> \t[%v]\n",vv)
    //         parseMob(bb[2:])
    //     }
    // }



    // fctpack["actor_something_happen"] = func (HexID []byte, bb []byte)  {
    //     // fmt.Printf("actor_something_happen [%v] \t",fmt.Sprintf("%#x", HexID))
    //     // fmt.Printf("bb -> \t[%v] \n",bb)
    //     mapID := int(bb[0])
    //     MUmobList.Lock()
    //     if bb[4] == 1{ // isDead
    //     if _, exist := mobList[mapID]; exist {
    //         targetMobDead = mapID
    //         delete(mobList, mapID)
    //         // fmt.Printf("mapID is ded-> \t[%v]\n",mapID)
    //     }}
    //     MUmobList.Unlock()
    // }
    //
    // fctpack["item_appear"] = func (HexID []byte, bb []byte)  {
    //     // fmt.Printf("item_appear [%v] \n",fmt.Sprintf("%#x", HexID))
    //     rrr := splitBitsArray(bb,[]byte{221,10})
    //     for _, vv := range rrr {
    //         parseItem(vv)
    //     }
    // }
    //
    // fctpack["mob_info"] = func (HexID []byte, bb []byte)  {
    //     rr := splitBitsArray(bb,[]byte{255,255})
    //     if len(rr) > 1{
    //         //mob dropped items
    //         rrr := splitBitsArray(rr[1],[]byte{221,10})
    //         for ii := 1; ii < len(rrr) ; ii++ {
    //             parseItem(rrr[ii])
    //         }
    //     }
    // }
    //
    // fctpack["use_skill"] = func (HexID []byte, bb []byte)  {
    //     fmt.Printf("use_skill [%v] \t",fmt.Sprintf("%#x", HexID))
    //     fmt.Printf("bb -> \t[%v]\n",bb)
    // }
    //
    // fctpack["item_disappear"] = func (HexID []byte, bb []byte)  {
    //     // fmt.Printf("item_disappear [%v] \t",fmt.Sprintf("%#x", HexID))
    //     // fmt.Printf("bb -> \t[%v]\n",bb)
    //     // 253 9
    //     // 255 9
    //     MUgroundItems.Lock()
    //     for ii := 0; ii < len(bb); ii+=4 {
    //         delete(groundItems, int(byteArrayToUInt32(bb[ii:ii+4])))
    //     }
    //     MUgroundItems.Unlock()
    // }
    // // fctpack["try_loot_item"] = func (HexID []byte, bb []byte)  {
    // //     fmt.Printf("try_loot_item [%v] \t",fmt.Sprintf("%#x", HexID))
    // //     fmt.Printf("bb -> \t[%v]\n",bb)
    // // }
    // fctpack["loot_item_confirm"] = func (HexID []byte, bb []byte)  {
    //
    //     rrr := splitBitsArray(bb,[]byte{74,188,30,0})
    //     if len(rrr) > 1{
    //         MUgroundItems.Lock()
    //         for ii := 1; ii < len(rrr) ; ii++ {
    //             targetItemLooted = int(byteArrayToUInt32(rrr[ii][0:4]))
    //             // fmt.Printf("loot -> \t[%v]\n",rrr[ii][0:4])
    //             delete(groundItems, int(byteArrayToUInt32(rrr[ii][0:4])))
    //         }
    //         MUgroundItems.Unlock()
    //     }
    // }
    //
    // fctpack["uknw_greed"] = func (HexID []byte, bb []byte)  {
    //     rr := splitBitsArray(bb,[]byte{74,188,30,0,74,188,30,0})
    //     if len(rr) > 2{
    //         rrr := splitBitsArray(rr[2],[]byte{74,188,30,0})
    //         if len(rrr) > 1{
    //             MUgroundItems.Lock()
    //             for ii := 1; ii < len(rrr) ; ii++ {
    //                 // fmt.Printf("loot -> \t[%v]\n",rrr[ii][0:4])
    //                 delete(groundItems, int(byteArrayToUInt32(rrr[ii][0:4])))
    //             }
    //             MUgroundItems.Unlock()
    //         }
    //     }
    // }

}

// #######################




// func parseMob(mob []byte){
//     if mob[0] == 108 || mob[0] == 114 {
//         mapID := int(byteArrayToUInt32(mob[0:0+4]))
//         mobID := int(byteArrayToUInt16(mob[21:21+2]))
//         sss := splitBitsArray(mob,[]byte{255,255,255,255,255,255,255,255})
//         mobName := strings.Replace(string(sss[1]),"\u0000", "", -1)
//         cc := Coord{X:0,Y:0} ; index := 0
//         if mob[0] == 108 { index = 61 }
//         if mob[0] == 114 { index = 65 }
//         bc := bits24ToCoords(mob[index:index+3])
//         cc.X = bc[0]; cc.Y = bc[1];
//         MUmobList.Lock()
//         mobList[mapID] = Mob{ MobID:mobID, Name:mobName, Coords:cc }
//         MUmobList.Unlock()
//     }
// }
// func parseItem(item []byte){
//     // fmt.Printf("mapID -> \t[%v]\n",item[0:4])
//     mapID := int(byteArrayToUInt32(item[0:4]))
//     itemID := int(byteArrayToUInt16(item[4:6]))
//     x := int(byteArrayToUInt16(item[11:13]))
//     y := int(byteArrayToUInt16(item[13:15]))
//     amount := int(byteArrayToUInt16(item[17:19]))
//     MUgroundItems.Lock()
//     groundItems[mapID] = Item{ ItemID:itemID, Coords:Coord{X:x,Y:y}, Amount:amount}
//     MUgroundItems.Unlock()
// }

// func isMob(bb []byte) bool{
//     flags := [][]byte{
//         []byte{0,108},
//         []byte{0,114},
//     }
//     for _,vv := range flags {
//         if sliceEqual(vv, bb) { return true }
//     }
//     return false
// }

// func getActorsFromArray(bb []byte) [][]byte{
//     flags := [][]byte{
//         []byte{255,9,108}, // mobs
//         []byte{255,9,114}, // mobs
//         []byte{255,9,63}, // ???
//         []byte{255,9,49}, // player
//         []byte{255,9,164}, // ???
//         []byte{255,9,108}, // player vending ?
//     }
//     var res [][]byte
//     for i := 0; i <= len(bb)-3 ; i++ {
//         if inArrayByte(flags,bb[i:i+3]) {
//             j := 0
//             for {
//                 if i+j >= len(bb){ break }
//                 if sliceEqual(bb[i+j:i+j+8], []byte{255,255,255,255,255,255,255,255}) {
//                     res = append(res,bb[i:i+j+32])
//                     i = i+j+32 ; break
//                 }
//                 j++
//             }
//         }
//     }
//     return res
// }


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
