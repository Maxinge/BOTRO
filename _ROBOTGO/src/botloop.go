package main

import(
    "fmt"
    "time"
    "encoding/binary"
    // "math"
    "math/rand"
)



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
    if exist := getConf(conf["General"],"Key","useTPNbAggro"); exist != nil {
        useTPNbAggro = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPNbAggroLoot"); exist != nil {
        useTPNbAggroLoot = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPLockMap"); exist != nil {
        useTPLockMap = exist.(struct{Key string; Val int}).Val
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

    movePath = nil

    go func() {
        for {time.Sleep(50 * time.Millisecond)

            charCoord = Coord{X:int(XPOS),Y:int(YPOS)}

            MUbuffList.Lock()
            cleanBuffList()
            MUbuffList.Unlock()

            MUmobList.Lock()
            for kk,vv := range mobList {
                mm := vv
                if getDist(vv.Coords, charCoord) > 25 { mm.AtSight = false }else{ mm.AtSight = true }
                mobList[kk] = mm
            }
            refreshMobsCoords()
            cleanMobDeath()
            MUmobList.Unlock()

            MUgroundItems.Lock()
            for kk,vv := range groundItems { if getDist(vv.Coords, charCoord) > 40 { delete(groundItems, kk) } }
            MUmobList.Lock()
            flagGoodItems()
            MUmobList.Unlock()
            MUgroundItems.Unlock()

        }
    }()


    for {


        waitime := time.Now()

        addWait(50)
        for{
            if needWait <= 0 { break }
            time.Sleep(time.Duration(50) * time.Millisecond)
            needWait -= 50
        }

        now = time.Now()
        looptime := int(now.Sub(waitime).Milliseconds())

        // #####################################################################
        // #####################################################################
        if !SIT {
            if chkcharCoord == charCoord { chkTimecharCoord += looptime }else{ chkTimecharCoord = 0 }
            if chktargetMobID == targetMobID && targetMobID > 0 { chkTimetargetMobID += looptime }else{ chkTimetargetMobID = 0 }
            if chktargetItemID == targetItemID && targetItemID > 0 { chkTimetargetItemID += looptime }else{ chkTimetargetItemID = 0 }
            if chkTimecharCoord > 20000 { resetPath(); resetMobItemList(); resetTargets() }
            if chkTimetargetMobID > 15000 { resetPath(); resetMobItemList(); resetTargets() }
            if chkTimetargetItemID > 5000 { resetPath(); resetMobItemList(); resetTargets() }
            chkcharCoord = charCoord
            chktargetMobID = targetMobID
            chktargetItemID = targetItemID
        }


        MUgroundItems.Lock()
        MUmobList.Lock()

        if targetItemID < 0 { targetItemID = pickItemTarget() }
        if targetMobID < 0 { targetMobID = pickMobTarget() }

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
            if vv.Priority >= 5 && vv.AtSight { countAggro = 999; break }
            if getDist(charCoord,vv.Coords) <= 3 && vv.Aggro{ countAggro++ }
        }

        MUmobList.Unlock()
        MUgroundItems.Unlock()

        MUbuffList.Lock()
        MUinventoryItems.Lock()
        skID, lv := needSkillSelf()
        itID := needUseItem()
        MUinventoryItems.Unlock()
        MUbuffList.Unlock()

        distFromDest = getDist(charCoord, nextPoint)

        // #####################################################################
        // #####################################################################

        if HPLEFT <= 0 {
            sendToServer("00B2", []byte{0})
            addWait(1000)
            resetMobItemList()
            resetPath()
            resetBuffList()
            resetTargets()
            continue
        }


        if countAggro < useTPNbAggroLoot && targetItemID > 0{ countAggro = 0 }
        if countAggro >= useTPNbAggro{ useTeleport(); continue}


        if (float32(SPLEFT)/float32(SPMAX)*100) <= float32(useSitUnderSP) {
        if targetItemID < 0 {
            if countAggro > 0 { useTeleport(); continue }
        if !SIT {
            sendToServer("0437", []byte{0,0,0,0,2})
        }}}

        if (float32(SPLEFT)/float32(SPMAX)*100) >= float32(useSitAboveSP) {
        if SIT {
            sendToServer("0437", []byte{0,0,0,0,3})
        }}

        if SIT { addWait(1000); continue }


        if distFromDest <= float64(minDist){ resetPath() }

        if movePath != nil {
            if pathIndex > len(movePath)-1 {
                nextStep = nextPoint
            }else{
                nextStep = Coord{movePath[pathIndex].X,movePath[pathIndex].Y}
            }
            if getDist(charCoord, nextStep) <= 6{ pathIndex += 2 }
            sendToServer("035F",coordsTo24Bits(nextStep.X,nextStep.Y))
            addWait(150)
        }

        if MAP == lockMap {

            if skID > 0 {
                sendUseSkill(skID, lv, accountID); addWait(100)
            continue }

            if itID > 0 {
                sendUseItem(itID); addWait(500)
            continue }

            if exist := getConf(conf["SKillSelf"],"Id",666666); exist != nil {
            if targetItemID < 0 {
            if SSphere < exist.(CSKillSelf).Lv {
                sendUseSkill(261, 5, accountID);
            continue }}}

            if targetMobID < 0 && targetItemID < 0 && movePath == nil{
                nextPoint = randomPoint(lgatMaps[MAP], charCoord, 50)
                movePath = pathfind(charCoord, nextPoint, lgatMaps[MAP])
                pathIndex = 1 ; minDist = 3;
            }

            if targetMobID > 0{
                MUmobList.Lock();
                mob := mobList[targetMobID];
                if !isInArray(targetMobID, keyMap(mobList)) || mob.DeathTime > 0 {
                    targetMobID = -1; resetPath(); MUmobList.Unlock(); continue  // ## !! ##
                }
                MUmobList.Unlock()

                resetPath()
                nextPoint = mob.Coords
                movePath = pathfind(charCoord, nextPoint, lgatMaps[MAP])
                if exist := getConf(conf["Mob"],"Id",mob.MobID); exist != nil {
                    minDist = exist.(CMob).MinDist
                }
                pathIndex = 2

                if distFromDest <= float64(minDist){
                    AtkId := 0; AtkLv := 0
                    if exist := getConf(conf["Mob"],"Id",mob.MobID); exist != nil {
                        AtkId = exist.(CMob).AtkId; AtkLv = exist.(CMob).AtkLv
                    }

                    if AtkId != 0 {
                        sendUseSkill(AtkId, AtkLv, targetMobID)
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
                }

            }

            if targetItemID > 0 {
                MUgroundItems.Lock();
                it := groundItems[targetItemID];
                if !isInArray(targetItemID, keyMap(groundItems)){
                    targetItemID = -1; resetPath(); MUgroundItems.Unlock(); continue  // ## !! ##
                }
                MUgroundItems.Unlock()
                if distFromDest >= float64(minDist){
                    resetPath()
                    allCells := firstCircle(it.Coords)
                    allCells = append(allCells,it.Coords)
                    rand.Seed(time.Now().UnixNano())
            		rnd := rand.Intn(len(allCells))
                    nextPoint = allCells[rnd]
                    movePath = pathfind(charCoord, nextPoint, lgatMaps[MAP])
                    pathIndex = 2 ; minDist = 2;
                }else{
                    itemBin := make([]byte, 4) ;
                    binary.LittleEndian.PutUint32(itemBin, uint32(targetItemID))
                    sendToServer("0362", itemBin)
                    addWait(150)
                }
            }

        }

        if MAP == lockMap && targetMobID < 0 && targetItemID < 0{ timerNoMob += looptime }else{ timerNoMob = 0 }
        if timerNoMob > useTPDelay { useTeleport(); continue }


        if exist := getConf(conf["Route"],"Map",MAP); exist != nil {
            if skID > 0 {
                sendUseSkill(skID, lv, accountID);
            continue }
            if movePath == nil {
                nextPoint = Coord{X:exist.(CRoute).X, Y:exist.(CRoute).Y}
                movePath = pathfind(charCoord, nextPoint, lgatMaps[MAP])
                pathIndex = 0 ; minDist = 1;
            }
        }

        if MAP == saveMap {
            if exist := getConf(conf["General"],"Key","WarpPortal"); exist != nil {
            if itemInInventory(717,1) > 0 { // bluegem
                portalChoice := exist.(struct{Key string; Val string}).Val
                time.Sleep(1000 * time.Millisecond)
                warpPoint := randomPoint(lgatMaps[MAP],charCoord, 3)
                sendWarpPortal(4,warpPoint.X,warpPoint.Y)
                time.Sleep(2000 * time.Millisecond)
                sendWarpPortalConfirm(portalChoice)
                time.Sleep(2000 * time.Millisecond)
                sendToServer("035F",coordsTo24Bits(warpPoint.X,warpPoint.Y))
                time.Sleep(2000 * time.Millisecond)
            }}
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

func findMobInDb(id int) map[string]interface{}{
    for _,vv := range mobDB {
        if int(vv["Id"].(float64)) == id { return vv }
    }
    return nil
}

func useTeleport()  {
    if useTPLockMap == 1 {
        tpId := int(binary.LittleEndian.Uint16([]byte{26,0}))
        tpLv := int(binary.LittleEndian.Uint16([]byte{1,0}))
        sendUseSkill(tpId, tpLv, accountID)
    }
    if useTPLockMap == 2 {
        inventID := itemInInventory(601,1) // fly wing
        if inventID > -1  { sendUseItem(inventID) }
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
func resetMobItemList() {
    MUmobList.Lock()
    mobList = map[int]Mob{}
    MUmobList.Unlock()
    MUgroundItems.Lock()
    groundItems = map[int]Item{}
    MUgroundItems.Unlock()
}
func resetPath() {
    movePath = nil; pathIndex = 0; nextPoint = Coord{};
}

func isValidLine(start Coord, dest Coord) bool{
    line := linearInterpolation(start, dest)
    for _,vv := range line {
    	gatcell := lgatMaps[MAP].cells[vv.X][vv.Y]
    	if !isValidCell(gatcell) { return false }
    }
    return true
}

func pickMobTarget() int{
    distMobList := map[float64]int{}

    for kk,vv := range mobList {
        if exist := getConf(conf["Mob"],"Id",vv.MobID); exist == nil { continue }

        if vv.Aggro {
        if vv.Priority >= 0 {
        if vv.DeathTime <= 0 {
        if getDist(charCoord, vv.Coords) <= 5 {
            return kk
        }}}}

        if !vv.AtSight { continue }
        if vv.IsNotValid { continue }
        if vv.DeathTime > 0 { continue }
        if vv.Priority > 3 { return kk }
        if vv.Priority > 2 { return kk }
        if vv.Priority > 1 { return kk }
        distMobList[getDist(vv.Coords, charCoord)] = kk
    }

    keys := sortFloatKeys(keyMap(distMobList))
    // for i := len(keys)-1; i >= 0; i-- {
    for i := 0; i < len(keys); i++ {
        mapID := distMobList[keys[i]]
        mob := mobList[mapID]
        if mob.IsNotValid { continue }
        if mob.DeathTime > 0 { continue }
        if exist := getConf(conf["Mob"],"Id",mob.MobID); exist == nil { continue }
        if !mob.AtSight { continue }
        if mob.Priority < 0 { continue }
        mobPath := pathfind(charCoord, mob.Coords, lgatMaps[MAP])
        if len(mobPath) > 50 { mob.IsNotValid = true; mobList[mapID] = mob; continue }
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
    for kk,ii := range inventoryItems {
        if ii.ItemID == id && ii.Amount >= amount   { return kk }
    }
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
        timeLeft := now.Sub(tBuff)
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
            if now.Sub(tDeath).Milliseconds() > 3000{
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

func refreshMobsCoords(){
    for kk,vv  := range mobList {
        if vv.LastMoveTime > 0 {
            mm := vv
            tMove := time.Unix(vv.LastMoveTime, 0)
            tMoveElapsed := now.Sub(tMove)
            index := (int(tMoveElapsed.Milliseconds())/(vv.MoveSpeed*2)) -1
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
