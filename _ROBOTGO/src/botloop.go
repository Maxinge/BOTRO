package main

import(
    "fmt"
    "time"
    "encoding/binary"
    "math"
    // "math/rand"
)

func infoLoop() {

    for {time.Sleep(50 * time.Millisecond)

        refreshPlayerCoords()

        charCoord = Coord{X:int(XPOS),Y:int(YPOS)}

        MUbuffList.Lock()
        cleanBuffList()
        MUbuffList.Unlock()

        MUmobList.Lock()
        for kk,vv := range mobList {
            mm := vv
            if getDist(vv.Coords, charCoord) > 30 { mm.AtSight = false }else{ mm.AtSight = true }
            mobList[kk] = mm
        }
        refreshMobsCoords()
        cleanMobDeath()
        MUmobList.Unlock()

        MUgroundItems.Lock()
        MUmobList.Lock()
        flagGoodItems()
        MUmobList.Unlock()
        MUgroundItems.Unlock()

    }
}

func botLoop() {

    if exist := getConf(conf["General"],"Key","accountID"); exist != nil {
        accountID = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","lockMap"); exist != nil {
        lockMap = exist.(struct{Key string; Val string}).Val
    }
    if exist := getConf(conf["General"],"Key","saveMap"); exist != nil {
        saveMap = exist.(struct{Key string; Val string}).Val
    }
    if exist := getConf(conf["General"],"Key","killBeforeLoot"); exist != nil {
        killBeforeLoot = exist.(struct{Key string; Val bool}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPNbAggro"); exist != nil {
        useTPNbAggro = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPNbAggroLoot"); exist != nil {
        useTPNbAggroLoot = exist.(struct{Key string; Val int}).Val
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

    movePath = nil


    for {
        waitime := time.Now()

        addWait(50)
        for{
            if needWait <= 0 { break }
            time.Sleep(time.Duration(50) * time.Millisecond)
            needWait -= 50
        }

        nowLoop = time.Now()
        looptime := int(nowLoop.Sub(waitime).Milliseconds())

        // #####################################################################
        // #####################################################################
        if !SIT {
            if chkcharCoord == charCoord { chkTimecharCoord += looptime }else{ chkTimecharCoord = 0 }
            if chktargetMobID == targetMobID && targetMobID > 0 { chkTimetargetMobID += looptime }else{ chkTimetargetMobID = 0; chkTimecharCoord = 0 }
            if chktargetItemID == targetItemID && targetItemID > 0 { chkTimetargetItemID += looptime }else{ chkTimetargetItemID = 0 }
            if chkTimecharCoord > 15000 { resetPath(); resetMobItemList(); resetTargets(); useTeleport() }
            if chkTimetargetMobID > 20000 && targetItemID < 0{ resetPath(); resetMobItemList(); resetTargets() }
            if chkTimetargetItemID > 5000 { resetPath(); resetMobItemList(); resetTargets() }
            chkcharCoord = charCoord
            chktargetMobID = targetMobID
            chktargetItemID = targetItemID
        }



        MUgroundItems.Lock()
        MUmobList.Lock()

        if targetItemID < 0 { targetItemID = pickItemTarget() }

        if targetMobID < 0 {
            targetMobID = pickMobTarget()
        }

        if targetItemID > 0 && targetMobID > 0 {
            mob := mobList[targetMobID]
            if mob.Priority > 1{
                targetItemID = -1
            }else{
                targetMobID = -1
            }
        }

        countAggro := 0
        for _,vv := range mobList {
            if vv.Priority <= -5 && vv.AtSight { countAggro = 999; break }
            if getDist(charCoord,vv.Coords) <= 10 && vv.Aggro && vv.DeathTime <= 0{ countAggro++ }
        }

        MUmobList.Unlock()
        MUgroundItems.Unlock()



        if len(playerList) > 0 { targetMobID = -1; }

        MUbuffList.Lock()
        skID, lv := needSkillSelf()
        itID := needUseItem()
        MUbuffList.Unlock()


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



        if HPLEFT <= 0 {
            sendToServer("00B2", []byte{0})
            resetMobItemList()
            resetPath()
            resetBuffList()
            resetTargets()
            addWait(1000)
            continue
        }

        townRun = false
        if (float32(WEIGHT)/float32(WEIGHTMAX)*100) >= float32(storageWeight) {
            townRun = true
        }

        for _,vv := range conf["Storage"] {
            it := vv.(CStorage)
            itId := itemInInventory(it.Id, it.Min)
            if itId == -1 {
                townRun = true
            }
        }

        if MAP != saveMap && townRun{
            goTown()
        }

        if countAggro < useTPNbAggroLoot && targetItemID > 0{ countAggro = 0 }
        if countAggro >= useTPNbAggro && MAP == lockMap{ useTeleport(); continue}

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

        if SIT { addWait(1000); continue }

        if MAP == saveMap && !townRun{
            if exist := getConf(conf["General"],"Key","WarpPortal"); exist != nil {
            portalChoice := exist.(struct{Key string; Val string}).Val
            if itemInInventory(717,1) > 0 && portalChoice != "" { // bluegem
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

        if MAP == saveMap && townRun{
            if movePath != nil {

                ii := getClosestPoint(charCoord,movePath) + 6
                if ii >= len(movePath)-6 { ii = len(movePath)-6 }

                if ii == len(movePath)-6 {
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
                        for _,vv := range conf["Storage"] {
                            it := vv.(CStorage)
                            itId := itemInStorage(it.Id, 1)
                            if itId >= 0 {
                                putItemIn("storage","inventory", itId, it.Max)
                                time.Sleep(time.Duration(500) * time.Millisecond)
                            }

                        }
                        closeStorage()
                        resetPath()
                        townRun = false
                    }

                }else{
                    sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
                    addWait(300)
                }

            }else{
                nextPoint = Coord{X:storageX, Y:storageY}
                movePath = pathfind(charCoord, nextPoint, lgatMaps[MAP], []Coord{})
            }
        continue }


        if exist := getConf(conf["Route"],"Map", MAP); exist != nil && !townRun{
            if skID > 0 {
                sendUseSkill(skID, lv, accountID);
            continue }

            if movePath != nil {
                if exist.(CRoute).UseTPdist > 0 {
                    if len(movePath) > exist.(CRoute).UseTPdist || len(movePath) <= 1{
                        if useTPOnRoad == 1 {
                            tpId := int(binary.LittleEndian.Uint16([]byte{26,0}))
                            tpLv := int(binary.LittleEndian.Uint16([]byte{1,0}))
                            sendUseSkill(tpId, tpLv, accountID)
                        }
                        if useTPOnRoad == 2 {
                            inventID := itemInInventory(601,1) // fly wing
                            if inventID > -1  {  sendUseItem(inventID)  }
                        }
                        continue
                    }
                }

                ii := getClosestPoint(charCoord,movePath) + 6
                if ii >= len(movePath)-1{ ii = len(movePath)-1 }
                sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
                addWait(300)
                if ii == len(movePath)-1 { resetPath() }
            }else{
                nextPoint = Coord{X:exist.(CRoute).X, Y:exist.(CRoute).Y}
                movePath = pathfind(charCoord, nextPoint, lgatMaps[MAP], []Coord{})
            }

        continue }



        if MAP == lockMap {

            if skID > 0 {
                sendUseSkill(skID, lv, accountID); addWait(100)
            continue }

            if itID > 0 {
                sendUseItem(itID); addWait(300)
            continue }

            if exist := getConf(conf["SKillSelf"],"Id",666666); exist != nil {
            if targetItemID < 0 {
            if SSphere < exist.(CSKillSelf).Lv {
                sendUseSkill(261, 5, accountID);
            continue }}}


            if targetMobID > 0 { timerNoMob = 0 }else{ timerNoMob += looptime }
            if timerNoMob > useTPDelay && targetItemID < 0 { useTeleport() }

            if targetMobID < 0 && targetItemID < 0 {

                if movePath != nil {
                    ii := getClosestPoint(charCoord,movePath) + 6
                    if ii >= len(movePath)-1{ ii = len(movePath)-1 }
                    sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
                    addWait(300)
                    if ii == len(movePath)-1 { resetPath() }
                }else{
                    nextPoint = randomPoint(lgatMaps[MAP], charCoord, 100)
                    movePath = pathfind(charCoord, nextPoint, lgatMaps[MAP], bannedCells)
                }
            continue }


            if targetMobID > 0 {

                MUmobList.Lock();
                mob := mobList[targetMobID];
                if !isInArray(targetMobID, keyMap(mobList)) || mob.DeathTime > 0 {
                    targetMobID = -1; resetPath(); MUmobList.Unlock(); continue  // ## !! ##
                }
                MUmobList.Unlock()


                resetPath()
                nextPoint = mob.Coords
                movePath = pathfind(charCoord, nextPoint, lgatMaps[MAP], bannedCells)
                atkDist := 2
                if exist := getConf(conf["Mob"],"Id",mob.MobID); exist != nil {
                    atkDist = exist.(CMob).MinDist
                }
                if len(movePath) <= atkDist{
                    AtkId := 0; AtkLv := 0
                    if exist := getConf(conf["Mob"],"Id",mob.MobID); exist != nil {
                        AtkId = exist.(CMob).AtkId; AtkLv = exist.(CMob).AtkLv
                    }
                    if AtkId != 0 {
                        sendUseSkill(AtkId, AtkLv, targetMobID);
                    }else{
                        arrayBin := []byte{}
                        mobBin := make([]byte, 4)
                        binary.LittleEndian.PutUint32(mobBin, uint32(targetMobID))
                        arrayBin = append(arrayBin,mobBin...)
                        // 0 = unique autoattack / 7 = start autoattack
                        arrayBin = append(arrayBin,byte(7))
                        sendToServer("0437", arrayBin)
                        addWait(100)
                    }
                }else{
                    ii := getClosestPoint(charCoord,movePath) + 6
                    if ii >= len(movePath)-1{ ii = len(movePath)-1 }
                    sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
                    addWait(300)
                }
            }

            if targetItemID > 0 {
                MUgroundItems.Lock();
                it := groundItems[targetItemID];
                if !isInArray(targetItemID, keyMap(groundItems)){
                    targetItemID = -1; resetPath(); MUgroundItems.Unlock(); continue  // ## !! ##
                }
                MUgroundItems.Unlock()

                nextPoint = it.Coords
                movePath = pathfind(charCoord, nextPoint, lgatMaps[MAP], []Coord{})

                ii := getClosestPoint(charCoord,movePath) + 6
                if ii >= len(movePath)-1{ ii = len(movePath)-1 }
                sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
                addWait(300)

                if len(movePath) <= 4{
                    itemBin := make([]byte, 4) ;
                    binary.LittleEndian.PutUint32(itemBin, uint32(targetItemID))
                    sendToServer("0362", itemBin)
                    addWait(300)
                }

            }

        }

    }

    fmt.Printf("# lel # %v \n","lel")
}

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
    binary.LittleEndian.PutUint16(XBin, uint16(XPOS))
    YBin := make([]byte, 2)
    binary.LittleEndian.PutUint16(YBin, uint16(YPOS))
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
    inventID := itemInInventory(602,1) // butt fly  wing
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
        inventID := itemInInventory(601,1) // fly wing
        if inventID > -1  {  sendUseItem(inventID)  }
    }

    if useTPLockMap == 3 {
        inventID := itemInInventory(601,1) // fly wing
        if inventID > -1 && (float32(SPLEFT)/float32(SPMAX)*100) <= float32(50){
            sendUseItem(inventID)
        }else{
            tpId := int(binary.LittleEndian.Uint16([]byte{26,0}))
            tpLv := int(binary.LittleEndian.Uint16([]byte{1,0}))
            sendUseSkill(tpId, tpLv, accountID)
        }
    }

}

func addWait(nw int){ if nw > needWait { needWait = nw } }

func resetTargets(){
    targetItemID = -1
    targetMobID = -1
    timerNoMob = 0
    chkTimecharCoord = 0
    chkTimetargetMobID = 0
    chkTimetargetItemID = 0
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

func resetPath() {
    movePath = nil; nextPoint = Coord{}
}

func pickMobTarget() int{
    distMobList := map[float64]int{}

    for kk,vv := range mobList {
        if exist := getConf(conf["Mob"],"Id",vv.MobID); exist == nil { continue }
        if !vv.AtSight { continue }
        if vv.IsNotValid { continue }
        if vv.DeathTime > 0 { continue }
        // distMobList[getDist(vv.Coords, charCoord)] = kk
        dist := float64(len(pathfind(charCoord, vv.Coords, lgatMaps[MAP], []Coord{})))
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
        if mob.IsLooter == true && keys[i] <= 3 { return mapID }
        if exist := getConf(conf["Mob"],"Id",mob.MobID); exist == nil { continue }
        if keys[i] > 50 { mob.IsNotValid = true; mobList[mapID] = mob; continue }
        if !mob.AtSight { continue }
        if mob.Priority < 0 { continue }
        return mapID
    }

    return -1
}

func pickItemTarget() int{
    for kk,vv := range groundItems {
        if !vv.IsValid { continue }
        if exist := getConf(conf["ItemLoot"],"Id",vv.ItemID); exist != nil {
            if exist.(CItemLoot).Priority == -1 { continue }
        }
        return kk
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

func itemInInventory(id int, amount int)  int{
    MUinventoryItems.Lock()
    for kk,ii := range inventoryItems {
        if ii.ItemID == id && ii.Amount >= amount { MUinventoryItems.Unlock(); return kk  }
    }
    MUinventoryItems.Unlock()
    return -1
}
func itemInStorage(id int, amount int)  int{
    MUstorageItems.Lock()
    for kk,ii := range storageItems {
        if ii.ItemID == id && ii.Amount >= amount { MUstorageItems.Unlock(); return kk  }
    }
    MUstorageItems.Unlock()
    return -1
}
func itemInCart(id int, amount int)  int{
    MUcartItems.Lock()
    for kk,ii := range cartItems {
        if ii.ItemID == id && ii.Amount >= amount { MUcartItems.Unlock(); return kk  }
    }
    MUcartItems.Unlock()
    return -1
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
            inventID := itemInInventory(it.Id, 1)
            if inventID > -1  { return inventID }
        }}
        if it.MinSP > 0  {
        if (float32(SPLEFT)/float32(SPMAX)*100) < float32(it.MinSP) {
            inventID := itemInInventory(it.Id, 1)
            if inventID > -1  { return inventID }
        }}
        if it.BuffId > 0 {
        if !isInArray(it.BuffId, keyMap(buffList)){
            inventID := itemInInventory(it.Id, 1)
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
            if getDist(vv.Coords,vvv.Coords) <= 5 {
            if vvv.DeathTime > 0 {
            if vvv.IsNotValid == false{
                tDeath := time.Unix(vvv.DeathTime, 0)
                tItem := time.Unix(vv.DropTime, 0)
                if tItem.Sub(tDeath).Milliseconds() < 1200{
                    ii.IsValid = true  ; groundItems[kk] = ii
                }
            }}}
        }
    }
}
func refreshPlayerCoords(){
    if lastMoveTime > 0 {
        tMove := time.Unix(lastMoveTime, 0)
        tMoveElapsed := time.Now().Sub(tMove)
        ff := float64(tMoveElapsed.Milliseconds()) / float64(MOVESPEED)*float64(0.33)
        index := int(math.Round(ff)-1)
        if index >= 0 && index < len(pathTo){
            XPOS = pathTo[index].X
            YPOS = pathTo[index].Y
        }
        if index >= len(pathTo) {
            lastMoveTime = 0
            XPOS = ccTo.X
            YPOS = ccTo.Y
        }
    }
}
func refreshMobsCoords(){
    for kk,vv  := range mobList {
        if vv.LastMoveTime > 0 {
            mm := vv
            tMove := time.Unix(vv.LastMoveTime, 0)
            tMoveElapsed := time.Now().Sub(tMove)
            ff := float64(tMoveElapsed.Milliseconds()) / float64(vv.MoveSpeed)*float64(1)
            index := int(math.Round(ff)-1)
            if index >= 0 && index < len(vv.PathMoveTo){
                mm.Coords = mm.PathMoveTo[index]
            }
            if index >= len(vv.PathMoveTo) {
                mm.LastMoveTime = 0
                mm.Coords = vv.CoordsTo
            }
            mobList[kk] = mm
        }
    }
}
