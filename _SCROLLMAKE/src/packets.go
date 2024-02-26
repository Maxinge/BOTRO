package main

import(
    "fmt"
    "strconv"
    // "reflect"
    "encoding/binary"
    "strings"
    // "time"
)

func parsePacket(bb []byte){
    hexID := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(bb[0:2]))
    bb = bb[2:]

    switch hexID {
    default:
        // fmt.Printf("### no_fct ### [%v][%v] -> [%v] \n", hexID, len(bb),bb)

    case "1414":  //mem_data
        ii := 0
        XPOS = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));           ii += 4
        YPOS = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));           ii += 4
        BASEXPMAX = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));      ii += 4
        BASEEXP = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));        ii += 4
        JOBXPMAX = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));       ii += 4
        JOBEXP = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));         ii += 4
        CHARNAME = strings.Split(string(bb[ii:ii+24]), "\x00")[0];     ii += 24
        BASELV = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));         ii += 4
        JOBLV = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));          ii += 4
        ZENY = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));           ii += 4
        MAP = strings.Split(string(bb[ii:ii+24]), ".rsw")[0];          ii += 24
        HPLEFT = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));         ii += 4
        HPMAX = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));          ii += 4
        WEIGHTMAX = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));      ii += 4
        WEIGHT = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));         ii += 4
        SPLEFT = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));         ii += 4
        SPMAX = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));          ii += 4

    // case "0438":  //skill_use_send
    //     fmt.Printf("### skill_use_send ### [%v][%v] -> [%v] \n", hexID, len(bb),bb)
    // case "01AE":  //scrollmakesend
    //     fmt.Printf("### scroll_make_send ### [%v][%v] -> [%v] \n", hexID, len(bb),bb)


    case "0B09", "0B0A":  //inventory_info
        // inventoryType := bb[2]
        if hexID == "0B09" {
            for ii := 3; ii < len(bb); ii+=34 {
                inventoryID := int(binary.LittleEndian.Uint16(bb[ii:ii+2]))
                itemID := int(binary.LittleEndian.Uint32(bb[ii+2:ii+2+4]))
                amount := int(binary.LittleEndian.Uint16(bb[ii+7:ii+7+2]))
                MUinventoryItems.Lock()
                inventoryItems[inventoryID] = Item{ ItemID:itemID, Coords:Coord{X:0,Y:0}, Amount:amount}
                MUinventoryItems.Unlock()
            }
        }
        if hexID == "0B0A" {
            for ii := 3; ii < len(bb); ii+=67 {
                inventoryID := int(binary.LittleEndian.Uint16(bb[ii:ii+2]))
                itemID := int(binary.LittleEndian.Uint32(bb[ii+2:ii+2+4]))
                MUinventoryItems.Lock()
                inventoryItems[inventoryID] = Item{ ItemID:itemID, Coords:Coord{X:0,Y:0}, Amount:1}
                MUinventoryItems.Unlock()
            }
        }

    case "01C8":  //item_use
        inventoryID := int(binary.LittleEndian.Uint16(bb[0:0+2]))
        // itemID := int(binary.LittleEndian.Uint32(bb[2:2+4]))
        amountLeft := int(binary.LittleEndian.Uint16(bb[10:10+2]))
        if ii, exist := inventoryItems[inventoryID]; exist {
            ii.Amount = amountLeft
            inventoryItems[inventoryID] = ii
        }


    case "00A1":  //item_disappear
        MUgroundItems.Lock()
        delete(groundItems, int(byteArrayToUInt32(bb[0:4])))
        MUgroundItems.Unlock()

    case "0A37":  //inventory_item_added
        inventoryID := int(binary.LittleEndian.Uint16(bb[0:0+2]))
        amount := int(binary.LittleEndian.Uint16(bb[2:2+2]))
        itemID := int(binary.LittleEndian.Uint32(bb[4:4+4]))
        MUinventoryItems.Lock()
        if ii, exist := inventoryItems[inventoryID]; exist {
            ii.Amount += amount; inventoryItems[inventoryID] = ii
        }else{
            inventoryItems[inventoryID] = Item{ ItemID:itemID, Coords:Coord{X:0,Y:0}, Amount:amount}
        }
        MUinventoryItems.Unlock()

    case "00AF":  //inventory_item_removed
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

    case "0B1A":  //skill_use_confirm
        sourceID := int(binary.LittleEndian.Uint32(bb[0:0+4]))
        // targetID := int(binary.LittleEndian.Uint32(bb[4:4+4]))
        // skillId := int(binary.LittleEndian.Uint16(bb[12:12+2]))
        castTime := int(binary.LittleEndian.Uint32(bb[18:18+4]))
        if sourceID == accountID {
            nw := castTime + 0
            // if skillId == 267{ nw += 300 } // TSS
            addWait(nw)
        }

    case "0983","0196","043F":  //buff_active_time //buff_active_off // buff_active_on
        buffID := int(binary.LittleEndian.Uint16(bb[0:0+2]))
        target := int(binary.LittleEndian.Uint32(bb[2:2+4]))
        flag := bb[6]

        if target == accountID {
        if buffID == 622{
            if flag == 1 { SIT = true } else { SIT = false }
            return
        }}



    // #######################
    case "0A30":  //actor_info
    case "00C0":  //emote
    case "0360":  //send_sync_serv
    case "007F":  //recv_sync_serv
    case "02C1":  //chat_main
    case "009A":  //serv_announc
    case "0439":  //item_use_send
    case "009D":  //item_exist
    case "0362":  //try_item_loot
    case "0A0A":  //storage_item_added
    case "035F":  //send_self_move_to
    case "0087":  //recv_self_move_to
    case "011B":  //warp_portal_choice_send
    case "0ABE":  //warp_portal_choice_recv
    case "0AF4":  //skill_use_aoe_recv
    case "0110":  //skill_use_failed
    case "0118":  //stop_attack
    case "0437":  //player_action_send
    case "0ACC":  //get_exp
    case "0368":  //actor_info_request
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
