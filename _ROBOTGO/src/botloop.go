package main

import(
    "fmt"
    "time"
    "encoding/binary"
    "math"
    // "math/rand"
)

func initConf(){
    if exist := getConf(conf["General"],"Key","accountID"); exist != nil {
        accountID = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","lockMap"); exist != nil {
        lockMap = exist.(struct{Key string; Val string}).Val
    }
    if exist := getConf(conf["General"],"Key","saveMap"); exist != nil {
        saveMap = exist.(struct{Key string; Val string}).Val
    }
    if exist := getConf(conf["General"],"Key","useGreed"); exist != nil {
        useGreed = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPUnderHP"); exist != nil {
        useTPUnderHP = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPNbAggro"); exist != nil {
        useTPNbAggro = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPLockMap"); exist != nil {
        useTPLockMap = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPOnRoad"); exist != nil {
        useTPOnRoad = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPDelay"); exist != nil {
        useTPDelay = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useSitUnderSP"); exist != nil {
        useSitUnderSP = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useSitAboveSP"); exist != nil {
        useSitAboveSP = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","storageWeight"); exist != nil {
        storageWeight = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","storageX"); exist != nil {
        storageX = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","storageY"); exist != nil {
        storageY = exist.(struct{Key string; Val int}).Val
    }
}

func infoLoop() {

    loopTimeEnd := time.Now()
    loopTimeStart := time.Now()

    for {

        loopTimeEnd = time.Now()
        waitfor := 33 - int(loopTimeEnd.Sub(loopTimeStart).Milliseconds())
        time.Sleep(time.Duration(waitfor) * time.Millisecond)
        loopTimeStart = time.Now()

        if ccFrom != ccTo  {
            ff := float64(lastMoveTime) / (float64(MOVESPEED)*float64(0.9))
            ii := int(math.Round(ff))
            if ii < 0 { ii = 0 }
            if ii >= len(pathTo)-1{ ii = len(pathTo)-1 }
            ccFrom = pathTo[ii]
            lastMoveTime += 25
        }else{
            lastMoveTime = 0
        }
        charCoord = ccFrom


        MUbuffList.Lock()
        cleanBuffList()
        MUbuffList.Unlock()

        MUmobList.Lock()
        for kk,vv := range mobList {
            mm := vv

            if mm.CoordsFrom != mm.CoordsTo {
                ff := float64(mm.LastMoveTime) / (float64(mm.MoveSpeed)*float64(0.9))
                ii := int(math.Round(ff))
                if ii < 0 { ii = 0 }
                if ii >= len(mm.PathMoveTo)-1 { ii = len(mm.PathMoveTo)-1 }
                mm.CoordsFrom = mm.PathMoveTo[ii]
                mm.LastMoveTime += 25
            }else{
                mm.LastMoveTime = 0
            }

            if getDist(mm.CoordsFrom, charCoord) > 30 { mm.AtSight = false }else{ mm.AtSight = true }
            mobList[kk] = mm

        }
        cleanMobDeath()
        MUmobList.Unlock()

        cleans := map[int]bool{}
        for kk,vv := range mobDeadList {
            tDeath := time.Unix(vv.DeathTime, 0)
            if time.Now().Sub(tDeath).Seconds() > 300{
                cleans[kk] = true
            }
        }
        tmp := []Mob{}
        for kk,vv := range mobDeadList {
            if !cleans[kk]{ tmp = append(tmp,vv) }
        }
        mobDeadList = tmp

        MUgroundItems.Lock()
        MUmobList.Lock()
        flagGoodItems()
        MUmobList.Unlock()
        MUgroundItems.Unlock()

    }
}

func botLoop() {

    movePath = nil

    loopTimeEnd := time.Now()
    loopTimeStart := time.Now()

    for {

        for{
            if needWait <= 0 { break }
            time.Sleep(time.Duration(50) * time.Millisecond)
            needWait -= 50
        }

        if timers.ThpTeleport >= -10000000000 { timers.ThpTeleport -= 50 }
        if timers.TnoMob >= -10000000000 { timers.TnoMob -= 50 }
        if timers.TuseItem >= -10000000000 { timers.TuseItem -= 50 }
        if timers.TuseSkill >= -10000000000 { timers.TuseSkill -= 50 }
        if timers.TuseSkillSelf >= -10000000000 { timers.TuseSkillSelf -= 50 }
        if timers.TclickMove >= -10000000000 { timers.TclickMove -= 50 }
        if timers.TsameCoord >= -10000000000 { timers.TsameCoord -= 50 }
        if timers.TsameMob >= -10000000000 { timers.TsameMob -= 50 }
        if timers.TsameItem >= -10000000000 { timers.TsameItem -= 50 }


        loopTimeEnd = time.Now()
        waitfor := 50 - int(loopTimeEnd.Sub(loopTimeStart).Milliseconds())
        if waitfor < 0 { waitfor = 0 }
        time.Sleep(time.Duration(waitfor) * time.Millisecond)
        loopTimeStart = time.Now()

        // #####################################################################
        // #####################################################################
        if HPLEFT <= 0 {
            sendToServer("00B2", []byte{0})
            pauseLoop(1000)
            continue
        }

        MUmobList.Lock()
        countAggro = 0
        for _,vv := range mobList {
            if vv.Priority <= -5 && vv.AtSight { countAggro = 999; break }
            if vv.Priority <= -4 && int(getDist(charCoord,vv.CoordsFrom)) <= vv.TPdist { countAggro = 999; break }
            if getDist(charCoord,vv.CoordsFrom) <= 4 && vv.Aggro && vv.DeathTime <= 0{ countAggro++ }
        }
        MUmobList.Unlock()

        MUgroundItems.Lock()
        if targetItemID < 0 { targetItemID = pickItemTarget() }
        MUgroundItems.Unlock()

        if SIT && countAggro > 0{
            sendToServer("0437", []byte{0,0,0,0,3})
            useTeleport(); continue
        }

        if (float32(SPLEFT)/float32(SPMAX)*100) <= float32(useSitUnderSP) {
        if targetItemID < 0{
            if countAggro > 0 { useTeleport(); continue }
        if !SIT {
            sendToServer("0437", []byte{0,0,0,0,2})
        }}}

        if (float32(SPLEFT)/float32(SPMAX)*100) >= float32(useSitAboveSP) {
        if SIT {
            sendToServer("0437", []byte{0,0,0,0,3})
        }}

        if SIT { pauseLoop(1000); continue }

        MUmobList.Lock()
        if targetMobID < 0 { targetMobID = pickMobTarget() }
        if targetItemID > 0 && targetMobID > 0 {
            mob := mobList[targetMobID]
            if int(getDist(charCoord, mob.CoordsFrom)) <= 2{
                targetItemID = -1
            }else{
                targetMobID = -1
            }
        }
        MUmobList.Unlock()

        bannedCells := []Coord{}
        MUtrapList.Lock()
        for _,vv := range trapList {
            bannedCells = append(bannedCells,vv.Coords)
            bannedCells = append(bannedCells, firstCircle(vv.Coords)...)
            bannedCells = append(bannedCells, secondCircle(vv.Coords)...)
        }
        MUtrapList.Unlock()
        // #####################################################################
        // #####################################################################
        if countAggro >= useTPNbAggro && targetItemID < 0{ useTeleport(); continue}

        if (float32(HPLEFT)/float32(HPMAX)*100) <= float32(useTPUnderHP) {
        if MAP != saveMap {
            if timers.ThpTeleport <= 0 {
                useTeleport(); timers.ThpTeleport = 10000
            }
        }}
        // #####################################################################
        // #####################################################################
        if sameCoord != charCoord || targetMobID > 0 || targetItemID > 0{
            timers.TsameCoord = 10000
        }
        sameCoord = charCoord
        if timers.TsameCoord <= 0 {
            timers.TsameCoord = 10000
            resetTargets(); resetPath(); useTeleport()
        }
        if sameMob != targetMobID || targetMobID < 0{
            timers.TsameMob = 15000
        }
        sameMob = targetMobID
        if timers.TsameMob <= 0 {
            timers.TsameMob = 15000
            resetTargets();resetPath();
        }
        if sameItem != targetItemID || targetItemID < 0{
            timers.TsameItem = 10000
        }
        sameItem = targetItemID
        if timers.TsameItem <= 0 {
            timers.TsameItem = 10000
            resetTargets();resetPath();
        }
        // #####################################################################
        // #####################################################################
        for _,vv := range conf["CartTransfert"] {
            it := vv.(CCartTransfert)
            if it.From {
                am, itId := itemInInventory(it.Id, it.Am)
                if itId >= 0 {
                    putItemIn("inventory","cart", itId, am)
                    pauseLoop(250)
                    MUinventoryItems.Lock()
                    ii := inventoryItems[itId]
                    ii.Amount -= am
                    if ii.Amount <= 0 {
                        delete(inventoryItems, itId)
                    }else{
                        inventoryItems[itId] = ii
                    }
                    MUinventoryItems.Unlock()
                    continue
                }
            }else{
                am, _ := itemInInventory(it.Id, it.Am)
                if am <= it.Am {
                    amCart, itIdCart := itemInCart(it.Id, it.Am)
                    if amCart >= it.Am {
                        putItemIn("cart","inventory", itIdCart, it.Am)
                        pauseLoop(250)
                        ii := cartItems[itIdCart]
                        ii.Amount -= am
                        if ii.Amount <= 0 {
                            MUcartItems.Lock()
                            delete(cartItems, itIdCart)
                            MUcartItems.Unlock()
                        }else{
                            cartItems[itIdCart] = ii
                        }
                        continue
                    }
                }
            }
        }
        // #####################################################################
        // #####################################################################
        if MAP == saveMap && !townRun{
            if exist := getConf(conf["General"],"Key","WarpPortal"); exist != nil {
            portalChoice := exist.(struct{Key string; Val string}).Val
            _,inventID := itemInInventory(717,1) // bluegem
            if inventID > 0 && portalChoice != "" {
                time.Sleep(1000 * time.Millisecond)
                warpPoint := randomPoint(lgatMaps[MAP],charCoord, 3)
                sendWarpPortal(4,warpPoint.X,warpPoint.Y)
                time.Sleep(2000 * time.Millisecond)
                sendWarpPortalConfirm(portalChoice)
                time.Sleep(2000 * time.Millisecond)
                sendToServer("035F",coordsTo24Bits(warpPoint.X,warpPoint.Y))
                time.Sleep(2000 * time.Millisecond)
                continue
            }}
        }
        // #####################################################################
        // #####################################################################
        townRun = false
        if (float32(WEIGHT)/float32(WEIGHTMAX)*100) >= float32(storageWeight) {
            townRun = true;
        }

        for _,vv := range conf["Storage"] {
            it := vv.(CStorage)
            _,itId := itemInInventory(it.Id, it.Min)
            if itId == -1 {
                townRun = true;
            }
        }

        // #####################################################################
        // #####################################################################
        if MAP != saveMap && townRun{
        if targetMobID < 0 && targetItemID < 0{
            goTown()
        }}
        if MAP == saveMap && townRun{
            if movePath == nil {
                npcCoord := Coord{X:storageX, Y:storageY}
                movePath = pathfind(charCoord, npcCoord, lgatMaps[MAP], []Coord{})
            }else{

                if timers.TclickMove <= 0 {
                    ii := getClosestPoint(charCoord,movePath) + 5
                    if ii >= len(movePath)-1{ ii = len(movePath)-1 }; if ii < 0 { ii = 0 }
                    sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
                    timers.TclickMove = 250
                }

                if int(getDist(movePath[len(movePath)-1],charCoord)) <= 8 {
                    ActorID := 0
                    MUnpcList.Lock()
                    for kk,vv := range npcList {
                        if vv.Coords.X == storageX && vv.Coords.Y == storageY  {
                            ActorID = kk
                        }
                    }
                    MUnpcList.Unlock()
                    if ActorID != 0 {
                        talkNpc(ActorID)
                        time.Sleep(time.Duration(500) * time.Millisecond)
                        talkNpcNext(ActorID)
                        time.Sleep(time.Duration(500) * time.Millisecond)
                        talkNpcChoice(ActorID, 2)
                        time.Sleep(time.Duration(500) * time.Millisecond)
                        talkNpcClose(ActorID)
                        time.Sleep(time.Duration(500) * time.Millisecond)

                        MUinventoryItems.Lock()
                        for kk,vv := range inventoryItems {
                            if vv.EqSlot <= 0 {
                                putItemIn("inventory","storage", kk, vv.Amount)
                                time.Sleep(time.Duration(500) * time.Millisecond)
                            }
                        }
                        inventoryItems = map[int]Item{}
                        MUinventoryItems.Unlock()

                        MUcartItems.Lock()
                        for kk,vv := range cartItems {
                            if vv.EqSlot <= 0 {
                                putItemIn("cart","storage", kk, vv.Amount)
                                time.Sleep(time.Duration(500) * time.Millisecond)
                            }
                        }
                        cartItems = map[int]Item{}
                        MUcartItems.Unlock()

                        closeStorage()
                        time.Sleep(time.Duration(500) * time.Millisecond)
                        // #######################
                        storageItems = map[int]Item{}
                        talkNpc(ActorID)
                        time.Sleep(time.Duration(500) * time.Millisecond)
                        talkNpcNext(ActorID)
                        time.Sleep(time.Duration(500) * time.Millisecond)
                        talkNpcChoice(ActorID, 2)
                        time.Sleep(time.Duration(500) * time.Millisecond)
                        talkNpcClose(ActorID)
                        time.Sleep(time.Duration(500) * time.Millisecond)

                        for _,vv := range conf["Storage"] {
                            it := vv.(CStorage)
                            am, itId := itemInStorage(it.Id, 1)
                            if itId >= 0 {
                                if am <= it.Max{
                                    putItemIn("storage","inventory", itId, am)
                                }else{
                                    putItemIn("storage","inventory", itId, it.Max)
                                }
                                time.Sleep(time.Duration(500) * time.Millisecond)
                            }
                        }
                        for _,vv := range conf["StorageCart"] {
                            it := vv.(CStorageCart)
                            am, itId := itemInStorage(it.Id, 1)
                            if itId >= 0 {
                                if am <= it.Max{
                                    putItemIn("storage","cart", itId, am)
                                }else{
                                    putItemIn("storage","cart", itId, it.Max)
                                }
                                time.Sleep(time.Duration(500) * time.Millisecond)
                            }
                        }
                        closeStorage()
                        resetPath()
                        townRun = false
                        timers.TsameCoord = 10000
                    }
                }
            }
        }
        // #####################################################################
        // #####################################################################
        if exist := getConf(conf["Route"],"Map", MAP); exist != nil || MAP == lockMap {
        if !townRun && MAP != saveMap{
            MUbuffList.Lock()
            skID, lv := needSkillSelf()
            itID := needUseItem()
            MUbuffList.Unlock()
            if itID > 0 {
                if timers.TuseItem <= 0 {
                    sendUseItem(itID)
                    timers.TuseItem = 300
                    MUinventoryItems.Lock()
                    it := inventoryItems[itID]
                    MUinventoryItems.Unlock()
                    if it.ItemID == 12114 || it.ItemID == 12115 ||
                        it.ItemID == 12116 || it.ItemID == 12117 ||
                        it.ItemID == 655 || it.ItemID == 656 || it.ItemID == 657 {
                        timers.TuseSkillSelf = 300
                        timers.TuseSkill = 300
                    }
                }
            }
            if skID > 0 {
                if timers.TuseSkillSelf <= 0 {
                    sendUseSkill(skID, lv, accountID)
                    timers.TuseSkillSelf = 300
                    timers.TuseSkill = 300
                    timers.TuseItem = 300
                    timers.TclickMove = 300
                }
            }
            if exist := getConf(conf["SKillSelf"],"Id",666666); exist != nil {
            if SSphere < exist.(CSKillSelf).Lv {
            if timers.TuseSkillSelf <= 0 {
                sendUseSkill(261, 5, accountID)
                timers.TuseSkillSelf = 300
                timers.TuseSkill = 300
                timers.TuseItem = 300
                timers.TclickMove = 300
            }}}
        }}
        // #####################################################################
        // #####################################################################
        if exist := getConf(conf["Route"],"Map", MAP); exist != nil {
        if !townRun {
            if movePath != nil {
                if exist.(CRoute).UseTPdist > 0 {
                    if len(movePath) > exist.(CRoute).UseTPdist || len(movePath) <= 2{
                        if useTPOnRoad == 1 {
                            tpId := int(binary.LittleEndian.Uint16([]byte{26,0}))
                            tpLv := int(binary.LittleEndian.Uint16([]byte{1,0}))
                            sendUseSkill(tpId, tpLv, accountID)
                        }
                        if useTPOnRoad == 2 {
                            _,inventID := itemInInventory(601,1) // fly wing
                            if inventID > -1  {  sendUseItem(inventID)  }
                        }
                        continue
                    }
                }
                if timers.TclickMove <= 0 {
                    ii := getClosestPoint(charCoord,movePath) + 5
                    if ii >= len(movePath)-1{ ii = len(movePath)-1 }; if ii < 0 { ii = 0 }
                    sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
                    timers.TclickMove = 250
                    if int(getDist(movePath[len(movePath)-1],charCoord)) <= 1{
                        resetPath(); pauseLoop(500); continue
                    }
                }
            }else{
                warpCoord := Coord{X:exist.(CRoute).X, Y:exist.(CRoute).Y}
                movePath = pathfind(charCoord, warpCoord, lgatMaps[MAP], []Coord{})
            }
        continue }}
        // #####################################################################
        // #####################################################################
        if MAP == lockMap && targetMobID < 0 && targetItemID < 0 {
            if timers.TnoMob <= 0 {
                useTeleport()
                timers.TnoMob = useTPDelay
            }
            if movePath != nil {
                if timers.TclickMove <= 0 {
                    ii := getClosestPoint(charCoord,movePath) + 5
                    if ii >= len(movePath)-1{ ii = len(movePath)-1 }; if ii < 0 { ii = 0 }
                    sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
                    timers.TclickMove = 250
                    if int(getDist(movePath[len(movePath)-1],charCoord)) <= 2{
                        resetPath(); continue
                    }
                }
            }else{
                rndCoord := randomPoint(lgatMaps[MAP], charCoord, 100)
                movePath = pathfind(charCoord, rndCoord, lgatMaps[MAP], bannedCells)
            }
        continue }

        // #####################################################################
        // #####################################################################
        if MAP == lockMap && targetMobID > 0 {
            MUmobList.Lock();
            mob := mobList[targetMobID];
            if !isInArray(targetMobID, keyMap(mobList)) || mob.DeathTime > 0 {
                targetMobID = -1; resetPath(); MUmobList.Unlock(); continue  // ## !! ##
            }
            MUmobList.Unlock()
            resetPath()
            movePath = pathfind(charCoord, mob.CoordsFrom, lgatMaps[MAP], bannedCells)

            atkDist := 1
            if exist := getConf(conf["Mob"],"Id",mob.MobID); exist != nil {
                atkDist = exist.(CMob).MinDist
            }
            movePath = movePath[:len(movePath)-1]

            if int(math.Round(getDist(charCoord,mob.CoordsFrom))) <= atkDist{ // ### int !!
                AtkId := 0; AtkLv := 0
                if exist := getConf(conf["Mob"],"Id",mob.MobID); exist != nil {
                    AtkId = exist.(CMob).AtkId; AtkLv = exist.(CMob).AtkLv
                }
                if AtkId != 0 {
                    if timers.TuseSkill <= 0 {
                        sendUseSkill(AtkId, AtkLv, targetMobID);
                        timers.TuseSkill = 100
                        timers.TuseItem = 200
                        timers.TuseSkillSelf = 300
                    }
                }else{
                    arrayBin := []byte{}
                    mobBin := make([]byte, 4)
                    binary.LittleEndian.PutUint32(mobBin, uint32(targetMobID))
                    arrayBin = append(arrayBin,mobBin...)
                    // 0 = unique autoattack / 7 = start autoattack
                    arrayBin = append(arrayBin,byte(7))
                    sendToServer("0437", arrayBin)
                    pauseLoop(200)
                }
            }else{
                if timers.TclickMove <= 0 {
                    ii := getClosestPoint(charCoord,movePath) + 5
                    if ii >= len(movePath)-1{ ii = len(movePath)-1 }; if ii < 0 { ii = 0 }
                    sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
                    timers.TclickMove = 250
                }
            }
        }
        // #####################################################################
        // #####################################################################
        if MAP == lockMap && targetItemID > 0 {
            MUgroundItems.Lock();
            it := groundItems[targetItemID];
            if !isInArray(targetItemID, keyMap(groundItems)){
                targetItemID = -1; resetPath(); MUgroundItems.Unlock(); continue  // ## !! ##
            }
            MUgroundItems.Unlock()

            movePath = pathfind(charCoord, it.Coords, lgatMaps[MAP], []Coord{})

            if timers.TclickMove <= 0 {
                ii := getClosestPoint(charCoord,movePath) + 5
                if ii >= len(movePath)-1{ ii = len(movePath)-1 }; if ii < 0 { ii = 0 }
                sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
                timers.TclickMove = 250
            }

            if int(math.Round(getDist(charCoord,it.Coords))) <= 3{
                itemBin := make([]byte, 4) ;
                binary.LittleEndian.PutUint32(itemBin, uint32(targetItemID))
                sendToServer("0362", itemBin)
                pauseLoop(300)
            }
        }
        // #####################################################################
        // #####################################################################

    }

    fmt.Printf("# lel # %v \n","lel")
}

func pauseLoop(nw int){ if nw > needWait { needWait = nw } }

func getClosestPoint(cc Coord, path []Coord) int{
    bestcc := len(path)-1
    bdist := getDist(cc, path[len(path)-1])
    for kk, vv := range path {
        dist := getDist(cc, vv)
        if dist <= bdist{
            bdist = dist
            bestcc = kk
        }
    }
    return bestcc
}

func putItemIn(from string, to string, inventoryID int, amount int)  {

    arrayBin := []byte{}
    inventoryIDBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(inventoryIDBin, uint16(inventoryID))
    amountBin := make([]byte, 4)
    binary.LittleEndian.PutUint32(amountBin, uint32(amount))
    arrayBin = append(arrayBin,inventoryIDBin...)
    arrayBin = append(arrayBin,amountBin...)

    packet := "0000"

    if from == "inventory" && to == "cart"      { packet = "0126" }
    if from == "cart" && to == "inventory"      { packet = "0127" }

    if from == "inventory" && to == "storage"   { packet = "0364" }
    if from == "storage" && to == "inventory"   { packet = "0365" }

    if from == "storage" && to == "cart"        { packet = "0128" }
    if from == "cart" && to == "storage"        { packet = "0129" }

    sendToServer(packet, arrayBin)
}

func talkNpc(Id int){
    arrayBin := []byte{}
    IdBin := make([]byte, 4)
    binary.LittleEndian.PutUint32(IdBin, uint32(Id))
    arrayBin = append(arrayBin,IdBin...)
    arrayBin = append(arrayBin,1)
    sendToServer("0090", arrayBin)
}

func talkNpcChoice(Id int, choice int){
    arrayBin := []byte{}
    IdBin := make([]byte, 4)
    binary.LittleEndian.PutUint32(IdBin, uint32(Id))
    arrayBin = append(arrayBin,IdBin...)
    arrayBin = append(arrayBin,byte(choice))
    sendToServer("00B8", arrayBin)
}

func talkNpcNext(Id int){
    IdBin := make([]byte, 4)
    binary.LittleEndian.PutUint32(IdBin, uint32(Id))
    sendToServer("00B9", IdBin)
}

func talkNpcClose(Id int){
    IdBin := make([]byte, 4)
    binary.LittleEndian.PutUint32(IdBin, uint32(Id))
    sendToServer("0146", IdBin)
}

func closeStorage(){
    sendToServer("0193", []byte{})
}

func closeShop(){
    sendToServer("09D4", []byte{})
}


func getIndex(cc Coord, path []Coord) int{
    for kk,vv := range path {
        if vv == cc { return kk }
    }
    return -1
}

func refreshGame(){
    arrayBin := []byte{}
    XBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(XBin, uint16(ccFrom.X))
    YBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(YBin, uint16(ccFrom.Y))
    arrayBin = append(arrayBin,[]byte{109,111,114,111,99,99,46,103,97,116,0,0,0,0,0,0}...)
    arrayBin = append(arrayBin,XBin...)
    arrayBin = append(arrayBin,YBin...)
    sendToClient("0091",arrayBin)
}

func sendHardRefresh(){
    arrayBin := []byte{}
    arrayBin = append(arrayBin,[]byte{22,0}...)
    arrayBin = append(arrayBin,[]byte(CHARNAME)...)
    arrayBin = append(arrayBin,[]byte{32,58,32,64,114,101,102,114,101,115,104}...)
    sendToServer("00F3", arrayBin)
}

func findMobInDb(id int) map[string]interface{}{
    for _,vv := range mobDB {
        if int(vv["Id"].(float64)) == id { return vv }
    }
    return nil
}
func goTown(){
    _,inventID := itemInInventory(602,1) // butt fly  wing
    if inventID > -1  {  sendUseItem(inventID) ; return }

    // tpId := int(binary.LittleEndian.Uint16([]byte{26,0}))
    // tpLv := int(binary.LittleEndian.Uint16([]byte{3,0}))
    // sendUseSkill(tpId, tpLv, accountID) ; return
}

func useTeleport()  {
    if useTPLockMap == 1 {
        tpId := int(binary.LittleEndian.Uint16([]byte{26,0}))
        tpLv := int(binary.LittleEndian.Uint16([]byte{1,0}))
        sendUseSkill(tpId, tpLv, accountID)
    }
    if useTPLockMap == 2 {
        _,inventID := itemInInventory(601,1) // fly wing
        if inventID > -1  {
            sendUseItem(inventID)
        }
    }

    if useTPLockMap == 3 {
        _,inventID := itemInInventory(601,1) // fly wing
        if inventID > -1 && (float32(SPLEFT)/float32(SPMAX)*100) <= float32(50){
            sendUseItem(inventID)
        }else{
            tpId := int(binary.LittleEndian.Uint16([]byte{26,0}))
            tpLv := int(binary.LittleEndian.Uint16([]byte{1,0}))
            sendUseSkill(tpId, tpLv, accountID)
        }
    }

}


func resetTargets(){
    targetItemID = -1
    targetMobID = -1
}
func resetBuffList() {
    MUbuffList.Lock()
    buffList = map[int][]int64{}
    MUbuffList.Unlock()
}
func resetPlayerList(){
    MUplayerList.Lock()
    playerList = map[int]Player{}
    MUplayerList.Unlock()
}
func resetTrapList(){
    MUtrapList.Lock()
    trapList = map[int]Trap{}
    MUtrapList.Unlock()
}
func resetMobItemList() {
    MUmobList.Lock()
    mobList = map[int]Mob{}
    MUmobList.Unlock()
    MUgroundItems.Lock()
    groundItems = map[int]Item{}
    MUgroundItems.Unlock()
}

func resetInventoryList(){
    MUinventoryItems.Lock()
    inventoryItems = map[int]Item{}
    MUinventoryItems.Unlock()
    MUstorageItems.Lock()
    storageItems = map[int]Item{}
    MUstorageItems.Unlock()
    MUcartItems.Lock()
    cartItems = map[int]Item{}
    MUcartItems.Unlock()
}

func resetPath() {
    movePath = nil;
}

func pickMobTarget() int{
    distMobList := map[float64]int{}

    for kk,vv := range mobList {
        if exist := getConf(conf["Mob"],"Id",vv.MobID); exist == nil { continue }
        if !vv.AtSight { continue }
        if vv.IsNotValid { continue }
        if vv.DeathTime > 0 { continue }
        // distMobList[getDist(vv.CoordsFrom, charCoord)] = kk
        dist := float64(len(pathfind(charCoord, vv.CoordsFrom, lgatMaps[MAP], []Coord{})))
        distMobList[dist] = kk
    }

    keys := sortFloatKeys(keyMap(distMobList))
    // for i := len(keys)-1; i >= 0; i-- {
    for i := 0; i < len(keys); i++ {
        mapID := distMobList[keys[i]]
        mob := mobList[mapID]
        if mob.Priority > 3 && keys[i] <= 25{ return mapID }
        if mob.Priority > 2 && keys[i] <= 25{ return mapID }
        if mob.Priority > 1 && keys[i] <= 25{ return mapID }
        if keys[i] <= 3 && mob.Aggro && mob.Priority >= 0 { return mapID }
        if mob.IsLooter == true && keys[i] <= 4 { return mapID }
        if exist := getConf(conf["Mob"],"Id",mob.MobID); exist == nil { continue }
        if keys[i] > 40 { mob.IsNotValid = true; mobList[mapID] = mob; continue }
        if !mob.AtSight { continue }
        if mob.Priority < 0 { continue }
        return mapID
    }

    return -1
}

func pickItemTarget() int{
    for kk,vv := range groundItems {
        if exist := getConf(conf["ItemLoot"],"Id",vv.ItemID); exist != nil {
            if exist.(CItemLoot).Priority == -1 { continue }
        }
        if len(playerList) <= 0 {
            return kk
        }else{
            if vv.IsValid { return kk }
        }
    }
    return -1
}


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

func itemInInventory(id int, amount int)  (int,int){
    MUinventoryItems.Lock()
    for kk,ii := range inventoryItems {
        if ii.ItemID == id && ii.Amount >= amount { MUinventoryItems.Unlock(); return ii.Amount, kk  }
    }
    MUinventoryItems.Unlock()
    return -1,-1
}
func itemInStorage(id int, amount int)  (int,int){
    MUstorageItems.Lock()
    for kk,ii := range storageItems {
        if ii.ItemID == id && ii.Amount >= amount { MUstorageItems.Unlock(); return ii.Amount, kk  }
    }
    MUstorageItems.Unlock()
    return -1,-1
}
func itemInCart(id int, amount int)  (int,int){
    MUcartItems.Lock()
    for kk,ii := range cartItems {
        if ii.ItemID == id && ii.Amount >= amount { MUcartItems.Unlock(); return ii.Amount, kk  }
    }
    MUcartItems.Unlock()
    return -1,-1
}

func sendUseItem(id int){
    arrayBin := []byte{}
    inventoryIDBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(inventoryIDBin, uint16(id))
    accountIDBin := make([]byte, 4)
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
        if tBuffTot == 9999{ continue }
        tBuff := time.Unix(vv[1], 0)
        timeLeft := time.Now().Sub(tBuff)
        if (tBuffTot - timeLeft.Milliseconds()) < 0 { delete(buffList, kk); }
    }
}

func needSkillSelf() (int, int) {
    for _, vv := range conf["SKillSelf"] {
        sk := vv.(CSKillSelf)
        if sk.MinHP > 0 && sk.MinSP > 0 && sk.BuffId < 0{
        if (float32(HPLEFT)/float32(HPMAX)*100) < float32(sk.MinHP) {
        if (float32(SPLEFT)/float32(SPMAX)*100) > float32(sk.MinSP) {
            return sk.Id, sk.Lv
        }}}
        if sk.BuffId > 0 && sk.MinHP < 0{
        if !isInArray(sk.BuffId, keyMap(buffList)){
        if (float32(SPLEFT)/float32(SPMAX)*100) > float32(sk.MinSP) {
            return sk.Id, sk.Lv
        }}}
    }
    return -1, -1
}

func needUseItem() int {
    for _, vv := range conf["ItemUse"] {
        it := vv.(CItemUse)
        if it.MinHP > 0  {
        if (float32(HPLEFT)/float32(HPMAX)*100) < float32(it.MinHP) {
            _,inventID := itemInInventory(it.Id, 1)
            if inventID > -1  { return inventID }
        }}
        if it.MinSP > 0  {
        if (float32(SPLEFT)/float32(SPMAX)*100) < float32(it.MinSP) {
            _,inventID := itemInInventory(it.Id, 1)
            if inventID > -1  { return inventID }
        }}
        if it.BuffId > 0 {
        if !isInArray(it.BuffId, keyMap(buffList)){
            _,inventID := itemInInventory(it.Id, 1)
            if inventID > -1  { return inventID }
        }}
    }

    return -1
}

func cleanMobDeath(){
    for kk,vv  := range mobList {
        if vv.DeathTime > 0 {
            tDeath := time.Unix(vv.DeathTime, 0)
            if time.Now().Sub(tDeath).Milliseconds() > 3000{
                delete(mobList,kk)
            }
        }
    }
}

func flagGoodItems(){
    for kk, vv  := range groundItems {
        ii := vv
        if vv.IsValid { continue }
        for _,vvv  := range mobList {
            if getDist(vv.Coords,vvv.CoordsFrom) <= 5 {
            if vvv.DeathTime > 0 {
            if vvv.IsNotValid == false{
                tDeath := time.Unix(vvv.DeathTime, 0)
                tItem := time.Unix(vv.DropTime, 0)
                if tItem.Sub(tDeath).Milliseconds() < 1500{
                    ii.IsValid = true  ; groundItems[kk] = ii
                }
            }}}
        }
    }
}
