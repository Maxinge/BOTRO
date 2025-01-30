package main

import(
    "fmt"
    "time"
    "encoding/binary"
    "math"
    // "math/rand"
    // "strings"
)

func initConf(){
    if exist := getConf(conf["General"],"Key","accountID"); exist != nil {
        accountID = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","autoSkillAspd"); exist != nil {
        autoSkillAspd = exist.(struct{Key string; Val int}).Val
    }
    
    if exist := getConf(conf["General"],"Key","useHpItem"); exist != nil {
        useHpItem = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useHpId"); exist != nil {
        useHpId = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useSpItem"); exist != nil {
        useSpItem = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useSpId"); exist != nil {
        useSpId = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","autoGreenPot"); exist != nil {
        autoGreenPot = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","autolootkey"); exist != nil {
        autolootkey = exist.(struct{Key string; Val int}).Val
    }
    if exist := getConf(conf["General"],"Key","useStealAfter"); exist != nil {
        useStealAfter = exist.(struct{Key string; Val int}).Val
    }  
    if exist := getConf(conf["General"],"Key","useStealSP"); exist != nil {
        useStealSP = exist.(struct{Key string; Val int}).Val
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

            if getDist(mm.CoordsFrom, charCoord) > 40 { mm.AtSight = false }else{ mm.AtSight = true }
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
        if timers.TclickLoot >= -10000000000 { timers.TclickLoot -= 50 }
        if timers.TloadTP >= -10000000000 { timers.TloadTP -= 50 }

        loopTimeEnd = time.Now()
        waitfor := 50 - int(loopTimeEnd.Sub(loopTimeStart).Milliseconds())
        if waitfor < 0 { waitfor = 0 }
        time.Sleep(time.Duration(waitfor) * time.Millisecond)
        loopTimeStart = time.Now()

        if useHpItem > 0 && useHpId > 0 {
        if (float32(HPLEFT)/float32(HPMAX)*100) < float32(useHpItem) {
            _,inventID := itemInInventory(useHpId, 1)
            if inventID > -1  { 
            if timers.TuseItem <= 0 {
                sendUseItem(inventID)
                timers.TuseItem = 300
            }}
        }}


        if targetMobID < 0 { targetStealed = -1; nbAuto = 0}

        if targetMobID > 0 {
            MUmobList.Lock();
            mob := mobList[targetMobID];
            if !isInArray(targetMobID, keyMap(mobList)) || mob.DeathTime > 0 || mob.IsNotValid{
                targetMobID = -1; MUmobList.Unlock();  // ## !! ##
                continue 
            }
            MUmobList.Unlock()

            if useStealSP > 0 {
            if (float32(SPLEFT)/float32(SPMAX)*100) >= float32(useStealSP) {
            if targetStealed < 0 {
                if nbAuto >= useStealAfter {
                    sendUseSkill(50, 10, targetMobID) // steal

                    arrayBin := []byte{}
                    mobBin := make([]byte, 4)
                    binary.LittleEndian.PutUint32(mobBin, uint32(targetMobID))
                    arrayBin = append(arrayBin,mobBin...)
                    // 0 = unique autoattack / 7 = start autoattack
                    arrayBin = append(arrayBin,byte(7))
                    sendToServer("0437", arrayBin)
                }
            }}}


        }

     
        // if (float32(HPLEFT)/float32(HPMAX)*100) < float32(it.MinHP) {
        //     _,inventID := itemInInventory(it.Id, 1)
        //     if inventID > -1  { return inventID }
        // }}
        // if it.MinSP > 0  {
        // if (float32(SPLEFT)/float32(SPMAX)*100) < float32(it.MinSP) {
        //     _,inventID := itemInInventory(it.Id, 1)
        //     if inventID > -1  { return inventID }
        // }}

        // if timers.TuseItem <= 0 {
        //     sendUseItem(itID)
        //     timers.TuseItem = 300
        // }
   
      
      
        // #####################################################################
        // #####################################################################
        // if MAP == lockMap && targetItemID > 0 {
        //     MUgroundItems.Lock();
        //     it := groundItems[targetItemID];
        //     if !isInArray(targetItemID, keyMap(groundItems)){
        //         targetItemID = -1; resetPath(); MUgroundItems.Unlock(); continue  // ## !! ##
        //     }
        //     MUgroundItems.Unlock()

        //     movePath = pathfind(charCoord, it.Coords, lgatMaps[MAP], []Coord{})

        //     if timers.TclickMove <= 0 {
        //         ii := getClosestPoint(charCoord,movePath) + 5
        //         if ii >= len(movePath)-1{ ii = len(movePath)-1 }
        //         sendToServer("035F",coordsTo24Bits(movePath[ii].X,movePath[ii].Y))
        //         timers.TclickMove = 250
        //     }

        //     if int(math.Round(getDist(charCoord,it.Coords))) <= 3{
        //         if timers.TclickLoot <= 0 {
        //             itemBin := make([]byte, 4) ;
        //             binary.LittleEndian.PutUint32(itemBin, uint32(targetItemID))
        //             sendToServer("0362", itemBin)
        //             timers.TclickLoot = 300
        //             timers.TclickMove = 200
        //             timers.TuseSkill = 300
        //             timers.TuseSkillSelf = 300
        //         }
        //     }
        // }
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
    if inventID > -1  { sendUseItem(inventID); resetPath(); return }

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
func resetNpcList(){
    MUnpcList.Lock()
    npcList = map[int]Npc{}
    MUnpcList.Unlock()
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
        pf := pathfind(charCoord, vv.CoordsFrom, lgatMaps[MAP], []Coord{})
        if pf[0] == pf[1] { continue }
        dist := float64(len(pf))
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
        if keys[i] > 45 { mob.IsNotValid = true; mobList[mapID] = mob; continue }
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
        if vv.IsValid { return kk }
        // if len(playerList) <= 0 {
        //     return kk
        // }else{
        //
        // }
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
    arrayBin = append(arrayBin, byteStr...)
    arrayBin = append(arrayBin, []byte{'\x00','\x00','\x00','\x00'}...)
    arrayBin = append(arrayBin, []byte{'\x00','\x00','\x00','\x00'}...)
    arrayBin = append(arrayBin, []byte{'\x00','\x00','\x00','\x00'}...)
    arrayBin = append(arrayBin, []byte{'\x00','\x00','\x00','\x00'}...)
    sendToServer("011B", arrayBin[0:18])
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
        if sk.DeBuffId > 0 && sk.MinHP < 0{
        if isInArray(sk.DeBuffId, keyMap(buffList)){
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
        if it.DeBuffId > 0 {
        if isInArray(it.DeBuffId, keyMap(buffList)){
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
            if time.Now().Sub(tDeath).Milliseconds() > 3000{ delete(mobList,kk) }
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
                    ii.IsValid = true; groundItems[kk] = ii
                }
            }}}
        }
    }
}
