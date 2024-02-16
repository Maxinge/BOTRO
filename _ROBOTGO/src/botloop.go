package main

import(
    "fmt"
    "time"
    "encoding/binary"
    // "math"
)


func resetAll()  {
    fmt.Printf("# reset # \n")
}


func botLoop() {

    movePath = nil

    if exist := getConf(conf["General"],"Key","accountID"); exist != nil {
        accountID = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","lockMap"); exist != nil {
        lockMap = exist.(struct{Key string; Val string}).Val
    }
    if exist := getConf(conf["General"],"Key","saveMap"); exist != nil {
        saveMap = exist.(struct{Key string; Val string}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPLockMap"); exist != nil {
        useTPLockMap = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useTPDelay"); exist != nil {
        useTPDelay = exist.(struct{Key string; Val int}).Val
    }



    for { time.Sleep(50 * time.Millisecond)

        now = time.Now()
        time.Sleep(2000 * time.Millisecond)
        // if now.Sub(stateTime).Milliseconds() > 15000 { resetAll(); stateTime = time.Now(); continue }

        bbb := []byte{109,111,114,111,99,99,46,103,97,116,0,0,0,0,0,0,102,0,144,0}

        sendToClient("0091", bbb)

        continue

        if needWait > 0  { time.Sleep(time.Duration(needWait) * time.Millisecond); needWait = 0}

        charCoord = Coord{X:int(XPOS),Y:int(YPOS)}

        MUbuffList.Lock()
        cleanBuffList()
        MUbuffList.Unlock()

        MUmobList.Lock()
        for kk,vv := range mobList { if getDist(vv.Coords, charCoord) > 40 { delete(mobList, kk) } }
        refreshMobsDist()
        MUmobList.Unlock()


        MUgroundItems.Lock()
        for kk,vv := range groundItems { if getDist(vv.Coords, charCoord) > 40 { delete(groundItems, kk) } }
        MUgroundItems.Unlock()

        // ############################
        // ############################

        MUbuffList.Lock()
        skID, lv := needSkillSelf()
        MUbuffList.Unlock()
        if skID > -1 {
            sendUseSkill(skID, lv, accountID);
        continue }


        MUinventoryItems.Lock()
        itID := needUseItem()
        MUinventoryItems.Unlock()
        if itID > -1 {
            sendUseItem(itID); needWait = 500
        continue }

        if movePath != nil {
            if pathIndex > len(movePath)-2 {
                nextStep = nextPoint
            }else{
                nextStep = Coord{movePath[pathIndex].X,movePath[pathIndex].Y}
            }
            if getDist(charCoord, nextStep) < 6{ pathIndex += 8 }
            sendToServer("035F",coordsTo24Bits(nextStep.X,nextStep.Y))
            time.Sleep(50 * time.Millisecond)
        continue }


        if exist := getConf(conf["Route"],"Map", MAP); exist != nil {
            nextPoint = Coord{X:exist.(CRoute).X, Y:exist.(CRoute).Y}
            movePath = pathfind(charCoord, nextPoint, lgatMaps[MAP])
            pathIndex = 0 ; minDist = 1;
        continue }


    }

    fmt.Printf("# lel # %v \n","lel")
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

func needSkillSelf() (int, int) {
    for _, vv := range conf["SKillSelf"] {
        sk := vv.(CSKillSelf)
        if sk.MinHP > 0 && sk.MinSP > 0{
        if (float32(HPLEFT)/float32(HPMAX)*100) < float32(sk.MinHP) {
        if (float32(SPLEFT)/float32(SPMAX)*100) > float32(sk.MinSP) {
            return sk.Id, sk.Lv
        }}}
        if sk.BuffId > 0 {
        if !isInArray(sk.BuffId, keyMap(buffList)){
            return sk.Id, sk.Lv
        }}
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


func refreshMobsDist(){
    for kk,vv  := range mobList {
        if vv.LastMoveTime > 0 {
            mm := vv
            tMove := time.Unix(vv.LastMoveTime, 0)
            tMoveElapsed := now.Sub(tMove)
            index := (int(tMoveElapsed.Milliseconds())/(vv.MoveSpeed*3))+1
            if index >= 0 && index < len(vv.PathTo){
                mm.Coords = mm.PathTo[index]
            }
            if index >= len(vv.PathTo) {
                mm.LastMoveTime = 0
                mm.Coords = vv.CoordsTo
            }
            mobList[kk] = mm
        }
    }
}






// func botLoop()  {

// startTime := time.Now()
// elapsed := time.Now()
//
// TPstartTime := time.Now()
// TPelapsed := time.Now()



// for { time.Sleep(100 * time.Millisecond)

    // if botStates == prevState {
    //     timeInState = elapsed.Sub(startTime).Seconds()
    //     if timeInState > float64(15) { resetStates() }
    //     elapsed = time.Now()
    // }else{
    //     startTime = time.Now()
    // }
    // prevState = botStates


    // MUgroundItems.Lock()
    // if targetItem == targetItemLooted { targetItem = -1; targetItemLooted = -2; nextPoint = Coord{X:0, Y:0} }
    // if !isInArray(targetItem, keyMap(groundItems)){ targetItem = -1; }
    // for kk,vv := range groundItems { if getDist(vv.Coords, charCoord) > 40 { delete(groundItems, kk) } }
    // for kk,vv := range groundItems {
    //     if curMap != lockMap { continue }
    //     if exist := getConf(conf["ItemLoot"],"Id",vv.ItemID); exist != nil {
    //         if exist.(CItemLoot).Priority == -1 { continue }
    //     }
    //     targetItem = kk
    // }
    // MUgroundItems.Unlock()

    // MUmobList.Lock()
    // if targetMob == targetMobDead {
    //     targetMob = -1; targetMobDead = -2; nextPoint = Coord{X:0, Y:0}
    //     time.Sleep(200 * time.Millisecond)
    // }
    // if !isInArray(targetMob, keyMap(mobList)){ targetMob = -1; }
    // for kk,vv := range mobList { if getDist(vv.Coords, charCoord) > 40 { delete(mobList, kk) } }
    // distMobList := map[float64]int{}
    // for kk,vv := range mobList { distMobList[getDist(vv.Coords, charCoord)] = kk }
    // keys := sortFloatKeys(keyMap(distMobList))
    // for i := len(keys)-1; i >= 0; i-- {
    //     if curMap != lockMap { continue }
    //     mob := mobList[distMobList[keys[i]]]
    //     if exist := getConf(conf["Mob"],"Id",mob.MobID); exist == nil { continue }
    //     // if getDist(charCoord, mob.Coords) > 25 { continue }
    //     mobPath := pathfind(charCoord, mob.Coords, lgatMaps[curMap])
    //     if len(mobPath) < 50 { targetMob = distMobList[keys[i]] ; continue }
    //
    //     isValidLine := true
    //     line := linearInterpolation(charCoord, mob.Coords)
    // 	for _,vv := range line {
    // 		gatcell := lgatMaps[curMap].cells[vv.X][vv.Y]
    // 		if !isValidCell(gatcell) { isValidLine = false; break}
    // 	}
    //     if !isValidLine { delete(mobList, distMobList[keys[i]]); continue }
    // }
    // MUmobList.Unlock()


    // if useTPLockMap > 0 {
    //     botStates.ReadyToTp = false
    //     if TPelapsed.Sub(TPstartTime).Seconds() > float64(useTPDelay) {
    //         botStates.ReadyToTp = true
    //     }
    //     TPelapsed = time.Now()
    // }

    // distFromDest := getDist(charCoord, nextPoint)
    // // #################################
    // // #################################
    // if nextPoint != (Coord{X:0, Y:0}) { botStates.HasDest = true }          else{ botStates.HasDest = false }
    // if curMap == lockMap { botStates.InLockMap = true }                     else{ botStates.InLockMap = false }
    // if curMap == saveMap { botStates.InSaveMap = true }                     else{ botStates.InSaveMap = false }
    // if distFromDest <= float64(minDist) { botStates.AtRange = true }        else{ botStates.AtRange = false }
    // if exist := getConf(conf["Route"],"Map", curMap); exist != nil {
    //     botStates.OnTheRoad = true
    // }else{ botStates.OnTheRoad = false }
    // if targetMob >= 0 { botStates.HasTargetMob = true }                     else{ botStates.HasTargetMob = false }
    // if targetItem >= 0 { botStates.HasTargetItem = true }                   else{ botStates.HasTargetItem = false }
    // #################################

    // #################################

// }

    // curPath = nil

    // if exist := getConf(conf["General"],"Key","accountID"); exist != nil {
    //     accountID = exist.(struct{Key string; Val int}).Val
    // }
    // if exist := getConf(conf["General"],"Key","lockMap"); exist != nil {
    //     lockMap = exist.(struct{Key string; Val string}).Val
    // }
    // if exist := getConf(conf["General"],"Key","saveMap"); exist != nil {
    //     saveMap = exist.(struct{Key string; Val string}).Val
    // }
    // if exist := getConf(conf["General"],"Key","useTPLockMap"); exist != nil {
    //     useTPLockMap = exist.(struct{Key string; Val int}).Val
    // }
    // if exist := getConf(conf["General"],"Key","useTPDelay"); exist != nil {
    //     useTPDelay = exist.(struct{Key string; Val int}).Val
    // }

    // for { time.Sleep(50 * time.Millisecond)
        // if charCoord == (Coord{X:0, Y:0}){ continue }
        // cleanBuffList()

        // MUbuffList.Lock()
        // skID, lv := needSkillSelf()
        // MUbuffList.Unlock()
        // if id > -1 {
        //     sendUseSkill(skID, lv, accountID)
        //     time.Sleep(200 * time.Millisecond)
        // }
        //

        // MUbuffList.Unlock()
        //
        //
        // if botStates.HasTargetItem == true {  botStates.HasTargetMob = false }
        // if botStates.HasTargetItem == true {  botStates.ReadyToTp = false }
        // if botStates.HasTargetMob == true {  botStates.ReadyToTp = false }
        // if botStates.OnTheRoad == true {  botStates.ReadyToTp = false }

        // if lastCastTime > 0  {
        //     time.Sleep(time.Duration(lastCastTime) * time.Millisecond)
        //     lastCastTime = 0
        // }

        // if botStates == (States{InLockMap:true, ReadyToTp:true}) ||
        //    botStates == (States{InLockMap:true, HasDest:true, ReadyToTp:true}) {
        //     if useTPLockMap == 1 {
        //         resetStates()
        //         time.Sleep(800 * time.Millisecond)
        //         tpId := int(binary.LittleEndian.Uint16([]byte{26,0}))
        //         tpLv := int(binary.LittleEndian.Uint16([]byte{1,0}))
        //         sendUseSkill(tpId, tpLv, accountID)
        //         time.Sleep(1300 * time.Millisecond)
        //         TPstartTime = time.Now()
        //     }
        //     if useTPLockMap == 2 {
        //         inventID := itemInInventory(601,1) // fly wing
        //         if inventID > -1  {
        //             resetStates()
        //             time.Sleep(800 * time.Millisecond)
        //             sendUseItem(inventID)
        //             time.Sleep(1300 * time.Millisecond)
        //         }
        //         TPstartTime = time.Now()
        //     }
        // }
        //
        // if botStates == (States{InLockMap:true, HasTargetMob: true}) ||
        //    botStates == (States{InLockMap:true, HasTargetMob: true, HasDest:true}) ||
        //    // botStates == (States{InLockMap:true, HasTargetMob: true, AtRange:true}) ||
        //    botStates == (States{InLockMap:true, HasTargetMob: true, HasDest:true, AtRange:true}) {
        //     MUmobList.Lock() ;  mob := mobList[targetMob] ; MUmobList.Unlock()
        //     nextPoint = mob.Coords
        //     curPath = pathfind(charCoord, nextPoint, lgatMaps[curMap])
        //     pathIndex = 1 ; minDist = (conf["SkillTarget"][attackIndex].(CSkillTarget)).MinDist
        // }
        //
        // if botStates == (States{InLockMap:true, HasTargetItem: true}) ||
        //    botStates == (States{InLockMap:true, HasTargetItem: true, HasDest:true}) {
        //     MUgroundItems.Lock() ; item := groundItems[targetItem] ; MUgroundItems.Unlock()
        //     nextPoint = item.Coords
        //     curPath = pathfind(charCoord, nextPoint, lgatMaps[curMap])
        //     pathIndex = 1 ; minDist = 2;
        // }
        //
        // if botStates == (States{InLockMap:true, HasTargetItem: true, HasDest:true, AtRange:true})  {
        //     itemBin := make([]byte, 4) ;
        //     binary.LittleEndian.PutUint32(itemBin, uint32(targetItem))
        //     sendToServer("0362", itemBin)
        //     time.Sleep(200 * time.Millisecond)
        // }
        //
        // if botStates == (States{InLockMap:true, HasTargetMob: true, HasDest:true, AtRange:true})  {
        //     skill := conf["SkillTarget"][attackIndex].(CSkillTarget)
        //     delay := 0
        //     if skill.Id != -1 {
        //         sendUseSkill(skill.Id, skill.Lv, targetMob)
        //     }else{
        //         arrayBin := []byte{}
        //         mobBin := make([]byte, 4)
        //         binary.LittleEndian.PutUint32(mobBin, uint32(targetMob))
        //         arrayBin = append(arrayBin,mobBin...)
        //         // 0 = unique autoattack / 7 = start autoattack
        //         arrayBin = append(arrayBin,byte(7))
        //         sendToServer("0437", arrayBin)
        //         delay = 1000
        //     }
        //     if attackIndex < len(conf["SkillTarget"])-1 { attackIndex++ }else{ attackIndex = 0 }
        //     time.Sleep(time.Duration(delay) * time.Millisecond)
        // }
        //
        //
        // if botStates == (States{InLockMap:true, HasDest:true}) ||
        //    botStates == (States{OnTheRoad:true, HasDest:true}) ||
        //    botStates == (States{OnTheRoad:true, HasDest:true, InSaveMap:true}) ||
        //    botStates == (States{InLockMap:true, HasDest:true, HasTargetMob:true}) ||
        //    botStates == (States{InLockMap:true, HasDest:true, HasTargetItem:true}) {
        //     if pathIndex > len(curPath)-2 {
        //         nextStep = nextPoint
        //     }else{
        //         nextStep = Coord{curPath[pathIndex].X,curPath[pathIndex].Y}
        //     }
        //
        //     if getDist(charCoord, nextStep) < 6{ pathIndex += 8 }
        //     sendToServer("035F",coordsTo24Bits(nextStep.X,nextStep.Y))
        //     time.Sleep(50 * time.Millisecond)
        // }
        //
        // if botStates == (States{InLockMap:true, HasDest:true, AtRange:true}) ||
        //    botStates == (States{OnTheRoad:true, HasDest:true, AtRange:true}){
        //     curPath = nil ; pathIndex = 0 ; nextPoint = Coord{X:0, Y:0}
        // }
        //
        // if botStates == (States{OnTheRoad:true}) ||
        //    botStates == (States{OnTheRoad:true, InSaveMap:true}) {
        //     if exist := getConf(conf["Route"],"Map",curMap); exist != nil {
        //         time.Sleep(1000 * time.Millisecond)
        //         if exist.(CRoute).WarpPortal != "" {
        //         if itemInInventory(717,1) > 0 { // bluegem
        //             time.Sleep(2000 * time.Millisecond)
        //             warpPoint := randomPoint(lgatMaps[curMap],charCoord, 3)
        //             sendWarpPortal(4,warpPoint.X,warpPoint.Y)
        //             time.Sleep(2000 * time.Millisecond)
        //             sendWarpPortalConfirm(exist.(CRoute).WarpPortal)
        //             time.Sleep(2000 * time.Millisecond)
        //             sendToServer("035F",coordsTo24Bits(warpPoint.X,warpPoint.Y))
        //             time.Sleep(2000 * time.Millisecond)
        //             continue // ##### !!!
        //         }}
        //         nextPoint = Coord{X:exist.(CRoute).X, Y:exist.(CRoute).Y}
        //         curPath = pathfind(charCoord, nextPoint, lgatMaps[curMap])
        //         pathIndex = 0 ; minDist = 1;
        //     }
        // }
        //
        // if botStates == (States{InLockMap:true}) {
        //     nextPoint = randomPoint(lgatMaps[curMap],charCoord, 80)
        //     curPath = pathfind(charCoord, nextPoint, lgatMaps[curMap])
        //     pathIndex = 0 ; minDist = 1;
        // }

    // }

// }
