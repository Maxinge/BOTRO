package main

import(
    "fmt"
    "strconv"
    // "reflect"
    "encoding/binary"
    "strings"
    "time"
)

func parsePacket(bb []byte){
    hexID := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(bb[0:2]))
    bb = bb[2:]

    switch hexID {
    default:
        // fmt.Printf("### no_fct ### [%v][%v] -> [%v] \n", hexID, len(bb),bb)

    case "0091","0092":  //map_change
        resetInventoryList()
        resetMobItemList()
        resetPlayerList()
        resetNpcList()
        resetTrapList()
        resetTargets()

        timers.TuseSkillSelf = 400
        timers.TuseSkill = 400
        timers.TuseItem = 400
        timers.TclickMove = 400
        timers.TloadTP = 400

        SSphere = 0
        pauseLoop(300)
        ccFrom = Coord{
            X:int(binary.LittleEndian.Uint16(bb[16:16+2])),
            Y:int(binary.LittleEndian.Uint16(bb[18:18+2])),
        }
        strMAP := strings.Split(string(bb[0:0+16]), "\x00")[0]
        strMAP = strings.Split(strMAP, ".gat")[0]
        MAP = strMAP
        movePath = []Coord{Coord{X:0,Y:0},Coord{X:0,Y:0}}
        ccTo = Coord{X:ccFrom.X,Y:ccFrom.Y}
        pathTo = []Coord{ccFrom,ccTo}
        lastMoveTime = 0

        go func() {
            time.Sleep(time.Duration(300) * time.Millisecond)
            resetPath()
            // if exist := getConf(conf["StorageRoute"],"Map", MAP); exist != nil {  }
            // if MAP := saveMap{ resetPath() }
        }()


    case "0087":  //recv_self_move_to
        fromto := bits48ToCoords(bb[4:4+6])
        ccFrom = Coord{X:fromto[0],Y:fromto[1]}
        ccTo = Coord{X:fromto[2],Y:fromto[3]}
        lastMoveTime = 0
        pathTo = pathfind(ccFrom, ccTo, lgatMaps[MAP],[]Coord{})

    case "1414":  //mem_data
        ii := 0

        MOVESPEED = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));      ii += 4
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
        CARTMIN = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));        ii += 4
        CARTMAX = int(binary.LittleEndian.Uint32(bb[ii:ii+4]));        ii += 4

    case "0B09", "0B0A":  //inventory_info

        inventoryType := bb[2]
        pad := 34
        if hexID == "0B09" { pad = 34 }
        if hexID == "0B0A" { pad = 67 }

        for ii := 3; ii < len(bb); ii+=pad {
            inventoryID := int(binary.LittleEndian.Uint16(bb[ii:ii+2]))
            itemID := int(binary.LittleEndian.Uint32(bb[ii+2:ii+2+4]))
            amount := 1
            eqSlot := 0
            if hexID == "0B09" {
                amount = int(binary.LittleEndian.Uint16(bb[ii+7:ii+7+2]))

            }
            if hexID == "0B0A" {
                eqSlot = int(binary.LittleEndian.Uint16(bb[ii+11:ii+11+2]))
            }
            if inventoryType == 0 {
                MUinventoryItems.Lock()
                inventoryItems[inventoryID] = Item{ ItemID:itemID, Amount:amount, EqSlot:eqSlot}
                MUinventoryItems.Unlock()
            }
            if inventoryType == 1 {
                MUcartItems.Lock()
                cartItems[inventoryID] = Item{ ItemID:itemID, Amount:amount, EqSlot:eqSlot}
                MUcartItems.Unlock()
            }
            if inventoryType == 2 {
                MUstorageItems.Lock()
                storageItems[inventoryID] = Item{ ItemID:itemID, Amount:amount, EqSlot:eqSlot}
                MUstorageItems.Unlock()
            }
        }


    case "01C8":  //item_use
        inventoryID := int(binary.LittleEndian.Uint16(bb[0:0+2]))
        // itemID := int(binary.LittleEndian.Uint32(bb[2:2+4]))
        amountLeft := int(binary.LittleEndian.Uint16(bb[10:10+2]))
        MUinventoryItems.Lock()
        if ii, exist := inventoryItems[inventoryID]; exist {
            if amountLeft == 0{
                ii.Amount -= 1
            }else{
                ii.Amount = amountLeft
            }
            inventoryItems[inventoryID] = ii
        }
        MUinventoryItems.Unlock()

    case "0ADD":  //item_appear
        tnow := time.Now()
        mapID := int(binary.LittleEndian.Uint32(bb[0:0+4]))
        itemID := int(binary.LittleEndian.Uint16(bb[4:4+2]))
        x := int(binary.LittleEndian.Uint16(bb[11:13]))
        y := int(binary.LittleEndian.Uint16(bb[13:15]))
        amount := int(binary.LittleEndian.Uint16(bb[17:19]))
        prio := -1
        if exist := getConf(conf["ItemLoot"],"Id",itemID); exist != nil {
            prio = exist.(CItemLoot).Priority
        }

        MUgroundItems.Lock()
        groundItems[mapID] = Item{
            ItemID:itemID, Coords:Coord{X:x,Y:y}, Amount:amount,
            DropTime:tnow.Unix(), Priority:prio,
        }
        MUgroundItems.Unlock()

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
        MUinventoryItems.Lock()
        if ii, exist := inventoryItems[inventoryID]; exist {
            ii.Amount -= amount
            if ii.Amount <= 0 {
                delete(inventoryItems, inventoryID)
            }else{
                inventoryItems[inventoryID] = ii
            }
        }
        MUinventoryItems.Unlock()

    case "0B1A":  //skill_use_confirm
        sourceID := int(binary.LittleEndian.Uint32(bb[0:0+4]))
        // targetID := int(binary.LittleEndian.Uint32(bb[4:4+4]))
        // skillId := int(binary.LittleEndian.Uint16(bb[12:12+2]))
        castTime := int(binary.LittleEndian.Uint32(bb[18:18+4]))
        if sourceID == accountID {
            if castTime > 0 {
                pauseLoop(castTime)
                timers.TclickLoot =  100
                timers.TuseSkillSelf =  100
                timers.TclickMove =  100
            }
        }

    case "01D0":  //spirit_sphere
        targetID := int(binary.LittleEndian.Uint32(bb[0:0+4]))
        number := int(binary.LittleEndian.Uint16(bb[4:4+2]))
        if targetID == accountID { SSphere = number }

    case "0983","0196","043F":  //buff_active_time //buff_active_off // buff_active_on
        buffID := int(binary.LittleEndian.Uint16(bb[0:0+2]))
        target := int(binary.LittleEndian.Uint32(bb[2:2+4]))
        flag := bb[6]
        timeLeft := int(binary.LittleEndian.Uint32(bb[11:11+4]))
        if target == accountID {
        if buffID == 622{
            if flag == 1 { SIT = true } else { SIT = false }
            return
        }}

        if hexID == "0196" {
        if target == accountID {
        if flag == 0 {
            MUbuffList.Lock()
            if _, exist := buffList[buffID]; exist {
                delete(buffList, buffID)
            }
            MUbuffList.Unlock()
        }}}

        if target == accountID && hexID == "0983"{
            if buffID == 46 {
                if timeLeft == 0 { timeLeft = 200 }
                timers.TuseSkill = timeLeft
                return
            }
            if timeLeft == 9999 { timeLeft = 9999999999999999 };
            if timeLeft > 0 {
                tnow := time.Now()
                MUbuffList.Lock()
                buffList[buffID] = []int64{int64(timeLeft), tnow.Unix()}
                MUbuffList.Unlock()
                return
            }
        }

    case "0086":  //actor_moving
        // fmt.Printf("### actor_moving ### [%v][%v] -> [%v] \n", hexID, len(bb),bb)
        mapID := int(binary.LittleEndian.Uint32(bb[0:4]))
        fromto := bits48ToCoords(bb[4:4+6])
        // tick := int(binary.LittleEndian.Uint32(bb[10:10+4]))
        ccFrom := Coord{X:fromto[0],Y:fromto[1]}
        ccTo := Coord{X:fromto[2],Y:fromto[3]}
        pathto := pathfind(ccFrom, ccTo, lgatMaps[MAP],[]Coord{})
        MUmobList.Lock()
        if mm, exist := mobList[mapID]; exist {
            mm.CoordsFrom = ccFrom; mm.CoordsTo = ccTo
            mm.LastMoveTime = 0
            mm.PathMoveTo = pathto
            mobList[mapID] = mm
        }
        MUmobList.Unlock()

    case "0088":  //actor_moving_interrupt
        mapID := int(binary.LittleEndian.Uint32(bb[0:4]))
        x := int(binary.LittleEndian.Uint16(bb[4:4+2]))
        y := int(binary.LittleEndian.Uint16(bb[6:6+2]))
        if mapID == accountID {
        cc := Coord{X:x,Y:y}
           pathTo = []Coord{cc,cc}
           lastMoveTime = 0
        }
        MUmobList.Lock()
        if mm, exist := mobList[mapID]; exist {
            cc := Coord{X:x,Y:y}
            mm.CoordsFrom = cc; mm.CoordsTo = cc
            mm.LastMoveTime = 0
            mm.PathMoveTo = []Coord{cc,cc}
            mobList[mapID] = mm
        }
        MUmobList.Unlock()

    case "0080":  //actor_dead_disappear
        tnow := time.Now()
        mapID := int(binary.LittleEndian.Uint32(bb[0:0+4]))
        // if bb[4] == 1 { targetMobDead = mapID }
        MUmobList.Lock()
        if mm, exist := mobList[mapID]; exist {
            if bb[4] == 1 { // killed
                mm.DeathTime = tnow.Unix()
                if mapID != targetMobID {
                    mm.IsNotValid = true
                }else{
                    mobDeadList = append(mobDeadList,mm)
                    timers.TclickMove = 200
                    timers.TclickLoot = 180
                    timers.TuseSkillSelf = 180
                    timers.TnoMob = 180
                    pauseLoop(200)
                    targetMobID = -1
                }
                mobList[mapID] = mm
            }else{
                delete(mobList, mapID)
            }
        }
        MUmobList.Unlock()

        MUnpcList.Lock()
        if _, exist := npcList[mapID]; exist {
            delete(npcList, mapID)
        }
        MUnpcList.Unlock()

        MUplayerList.Lock()
        if _, exist := playerList[mapID]; exist {
            delete(playerList, mapID)
        }
        MUplayerList.Unlock()

        MUtrapList.Lock()
        if _, exist := trapList[mapID]; exist {
            delete(trapList, mapID)
        }
        MUtrapList.Unlock()

    case "09CA":  //traps
        // SkillID := int(binary.LittleEndian.Uint32(bb[2:2+4]))
        mapID := int(binary.LittleEndian.Uint32(bb[6:6+4]))
        x := int(binary.LittleEndian.Uint16(bb[10:10+2]))
        y := int(binary.LittleEndian.Uint16(bb[12:12+2]))
        job := int(binary.LittleEndian.Uint32(bb[14:14+4]))
        radius := int(bb[18])
        if job == 148 {
            MUtrapList.Lock()
            trapList[mapID] = Trap{ TrapID:job, Coords:Coord{X:x,Y:y}, Radius:radius }
            MUtrapList.Unlock()
        }

    case "09FD", "09FF", "09FE":     //actor_appear_exist // actor_spawn/connected
        if len(bb) == 0 { return }
        mapID := int(binary.LittleEndian.Uint32(bb[3:3+4]))
        actorID := int(binary.LittleEndian.Uint16(bb[21:21+4]))
        moveSpeed := int(binary.LittleEndian.Uint16(bb[11:11+2]))
        actorType := byte(bb[2])
        if bb[17] == 4 || bb[17] == 2 { return } // hided
        mccFrom := Coord{X:0,Y:0}
        mccTo := Coord{X:0,Y:0}
        mpathTo := []Coord{}

        if  hexID == "09FD" {
            index := 65
            bc := bits48ToCoords(bb[index:index+6])
            mccFrom.X = bc[0]; mccFrom.Y = bc[1];
            mccTo.X = bc[2]; mccTo.Y = bc[3];
            mpathTo = pathfind(mccFrom, mccTo, lgatMaps[MAP],[]Coord{})
        }
        if  hexID == "09FF" || hexID == "09FE" {
            index := 61
            bc := bits24ToCoords(bb[index:index+3])
            mccFrom.X = bc[0]; mccFrom.Y = bc[1];
            mccTo.X = bc[0]; mccTo.Y = bc[1];
        }

        if actorType == 5  {
            MUmobList.Lock()
            prio := -1
            aggro := false
            name := ""
            looter := false
            mobdata := findMobInDb(actorID)
            bexp := 0
            jexp := 0
            tpdist := 0
            if mobdata != nil {
                aggro = mobdata["IsAggressive"].(bool)
                name = mobdata["Name"].(string)
                looter = mobdata["IsLooter"].(bool)
                bexp = int(mobdata["BaseExp"].(float64))
                jexp = int(mobdata["JobExp"].(float64))
            }
            if exist := getConf(conf["Mob"],"Id",actorID); exist != nil {
                prio = exist.(CMob).Priority
                tpdist = exist.(CMob).TPdist
            }
            mobList[mapID] = Mob{
                MobID:actorID, CoordsFrom:mccFrom, CoordsTo:mccTo, MoveSpeed:moveSpeed,
                Priority: prio, Aggro:aggro, Name:name, IsLooter:looter,
                Bexp:bexp, Jexp:jexp, TPdist:tpdist, PathMoveTo:mpathTo, LastMoveTime:0,
            }
            MUmobList.Unlock()
        }
        if actorType == 0 {
            MUplayerList.Lock()
            playerList[mapID] = Player{Coords:mccFrom}
            MUplayerList.Unlock()
        }

        if actorType == 6 {
            MUnpcList.Lock()
            npcList[mapID] = Npc{NpcID:actorID,Coords:mccFrom,Name:""}
            MUnpcList.Unlock()
        }



    case "01DE", "09CB", "08C8":  //skill_used_on_target //skill_no_dmg //actor_action

        sourceii := 0 ; targetii := 0
        if hexID == "01DE" { sourceii = 2 ; targetii = 6 }
        if hexID == "09CB" { sourceii = 10 ; targetii = 6 }
        if hexID == "08C8" { sourceii = 0 ; targetii = 4 }
        sourceID := int(binary.LittleEndian.Uint32(bb[sourceii:sourceii+4]))
        targetID := int(binary.LittleEndian.Uint32(bb[targetii:targetii+4]))
        dmg := int(binary.LittleEndian.Uint32(bb[20:20+4]))
        if hexID == "08C8" && bb[27] != 1 {
            if targetID == accountID && dmg > 0{
                timers.TclickMove = 30
            }
            return
        } //autoattack

        MUmobList.Lock()
        if mm, exist := mobList[targetID]; exist {
            if sourceID != accountID {
                mm.IsNotValid = true;
                mobList[targetID] = mm
            }
        }
        if targetID == accountID {
            if mm, exist := mobList[sourceID]; exist {
                    mm.IsNotValid = false;
                    mm.Aggro = true;
                    mobList[sourceID] = mm
            }else{
                mobList[sourceID] = Mob{ Aggro:true, MoveSpeed:200 }
            }
        }
        MUmobList.Unlock()

    // #######################

    case "6847":  //???
    // case "0A30":  //actor_info
    case "00C0":  //emote
    case "0438":  //skill_use_send
    case "0360":  //send_sync_serv
    case "007F":  //recv_sync_serv
    case "02C1":  //chat_main
    case "009A":  //serv_announc
    case "0439":  //item_use_send
    case "009D":  //item_exist
    case "0362":  //try_item_loot
    case "0A0A":  //storage_item_added
    case "035F":  //send_self_move_to
    // case "0087":  //recv_self_move_to
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
