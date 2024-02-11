package main

import(
    "fmt"
    "time"
    "encoding/binary"
    "sync"
)

type Coord struct {
	X,Y int
}

type Mob struct {
    MobID int
    Coords Coord
    HPMax int
    HPLeft int
}

type Item struct {
    ItemID int
    Coords Coord
    Amount int
}

type States struct {
    InSaveMap bool
    OnTheRoad bool
    InLockMap bool
    HasDest bool
    HasTargetMob bool
    HasTargetItem bool
    AtRange bool
    ReadyToTp bool
}

var(
    botStates = States{}
    prevState = botStates
    timeInState = float64(0)

    accountId = 0

    curCoord = Coord{X:0, Y:0}
    nextPoint = Coord{X:0, Y:0}
	curMap = ""
    HPLeft = 0
    HPMax = 0
    maxWeight = 0
    weight = 0
    SPLeft = 0
    SPMax = 0

    pathIndex = 0
	curPath = []Coord{}
    nextStep = Coord{}
    minDist = 1

    MUmobList sync.Mutex
    mobList = map[int]Mob{}
    MUgroundItems sync.Mutex
    groundItems = map[int]Item{}
    MUinventoryItems sync.Mutex
    inventoryItems = map[int]Item{}
    MUbuffList sync.Mutex
    buffList = map[int][]int64{}

    targetMob = -1
    targetMobDead = -2
    targetItem = -1
    targetItemLooted = -2
    attackDist = 1
    attackIndex = 0

    lockMap = ""
    saveMap = ""
    useTPLockMap = 0
    useTPDelay = 10
)

func resetStates(){
    botStates = States{}
    pathIndex = 0 ; curPath = nil ;
    nextPoint = Coord{X:0, Y:0}
    targetMob = -1; targetMobDead = -2
    targetItem = -1; targetItemLooted = -2
    MUgroundItems.Lock()
    mobList = map[int]Mob{}
    MUgroundItems.Unlock()
    MUmobList.Lock()
    groundItems = map[int]Item{}
    MUmobList.Unlock()
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
    accountIdBin := make([]byte, 4) ;
    binary.LittleEndian.PutUint32(accountIdBin, uint32(accountId))
    arrayBin = append(arrayBin,inventoryIDBin...)
    arrayBin = append(arrayBin,accountIdBin...)
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

func botLoop() {

    curPath = nil

    startTime := time.Now()
    elapsed := time.Now()

    TPstartTime := time.Now()
    TPelapsed := time.Now()

    if exist := getConf(conf["General"],"Key","accountId"); exist != nil {
        accountId = exist.(struct{Key string; Val int}).Val
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

    for { time.Sleep(100 * time.Millisecond)
        if curCoord == (Coord{X:0, Y:0}){ continue }

        if botStates == prevState {
            timeInState = elapsed.Sub(startTime).Seconds()
            if timeInState > float64(15) { resetStates() }
            elapsed = time.Now()
        }else{
            startTime = time.Now()
        }
        prevState = botStates
        // #################################
        // #################################

        MUbuffList.Lock()
        for kk,vv  := range buffList {
            tBuffTot := vv[0]
            tBuffAt := time.Unix(vv[1], 0)
            timeLeft := time.Now().Sub(tBuffAt)
            if (tBuffTot - timeLeft.Milliseconds()) < 0 { delete(buffList, kk); }
        }
        MUbuffList.Unlock()

        for _, vv := range conf["SKillSelf"] {
            sk := vv.(CSKillSelf)
            if sk.MinHP > 0 && sk.MinSP > 0{
            if (float32(HPLeft)/float32(HPMax)*100) < float32(sk.MinHP) {
            if (float32(SPLeft)/float32(SPMax)*100) > float32(sk.MinSP) {
                sendUseSkill(sk.Id, sk.Lv, accountId)
            }}}
            if sk.BuffId > 0 {
                MUbuffList.Lock()
                if !isInArray(sk.BuffId, keyMap(buffList)){
                    sendUseSkill(sk.Id, sk.Lv, accountId)
                }
                MUbuffList.Unlock()
            }
        }

        // #################################
        // #################################

        for _, vv := range conf["ItemUse"] {
            it := vv.(CItemUse)
            if it.MinHP > 0  {
            if (float32(HPLeft)/float32(HPMax)*100) < float32(it.MinHP) {
                MUinventoryItems.Lock()
                inventID := itemInInventory(it.Id, 1)
                if inventID > -1  { sendUseItem(inventID) }
                MUinventoryItems.Unlock()
            }}
            if it.MinSP > 0  {
            if (float32(SPLeft)/float32(SPMax)*100) < float32(it.MinSP) {
                MUinventoryItems.Lock()
                inventID := itemInInventory(it.Id, 1)
                if inventID > -1  { sendUseItem(inventID) }
                MUinventoryItems.Unlock()
            }}
            if it.BuffId > 0 {
                MUbuffList.Lock()
                if !isInArray(it.BuffId, keyMap(buffList)){
                    inventID := itemInInventory(it.Id, 1)
                    if inventID > -1  { sendUseItem(inventID) }
                }
                MUbuffList.Unlock()
            }
        }

        // #################################
        // #################################

        if targetMob == targetMobDead {
            targetMob = -1; targetMobDead = -2; nextPoint = Coord{X:0, Y:0}
             time.Sleep(200 * time.Millisecond)
        }
        if targetItem == targetItemLooted { targetItem = -1; targetItemLooted = -2; nextPoint = Coord{X:0, Y:0} }

        MUgroundItems.Lock()
        targetItem = -1
        for kk,vv := range groundItems { if getDist(vv.Coords, curCoord) > 40 { delete(groundItems, kk) } }
        for kk,vv := range groundItems {
            if curMap != lockMap { continue }
            if exist := getConf(conf["ItemLoot"],"Id",vv.ItemID); exist != nil {
                if exist.(CItemLoot).Priority == -1 { continue }
            }
            targetItem = kk
        }
        MUgroundItems.Unlock()

        MUmobList.Lock()
        targetMob = -1
        for kk,vv := range mobList { if getDist(vv.Coords, curCoord) > 40 { delete(mobList, kk) } }
        distMobList := map[float64]int{}
        for kk,vv := range mobList { distMobList[getDist(vv.Coords, curCoord)] = kk }
        keys := sortFloatKeys(keyMap(distMobList))
        for i := len(keys)-1; i >= 0; i-- {
            if curMap != lockMap { continue }
            mob := mobList[distMobList[keys[i]]]
            if exist := getConf(conf["Mob"],"Id",mob.MobID); exist == nil { continue }
            // if getDist(curCoord, mob.Coords) > 25 { continue }
            mobPath := pathfind(curCoord, mob.Coords, lgatMaps[curMap])
            if len(mobPath) < 50 { targetMob = distMobList[keys[i]] ; continue }

            isValidLine := true
            line := linearInterpolation(curCoord, mob.Coords)
        	for _,vv := range line {
        		gatcell := lgatMaps[curMap].cells[vv.X][vv.Y]
        		if !isValidCell(gatcell) { isValidLine = false; break}
        	}
            if !isValidLine { delete(mobList, distMobList[keys[i]]); continue }
        }
        MUmobList.Unlock()


        if useTPLockMap > 0 {
            botStates.ReadyToTp = false
            if TPelapsed.Sub(TPstartTime).Seconds() > float64(useTPDelay) {
                botStates.ReadyToTp = true
            }
            TPelapsed = time.Now()
        }

        distFromDest := getDist(curCoord, nextPoint)
        // #################################
        // #################################
        if nextPoint != (Coord{X:0, Y:0}) { botStates.HasDest = true }          else{ botStates.HasDest = false }
        if curMap == lockMap { botStates.InLockMap = true }                     else{ botStates.InLockMap = false }
        if curMap == saveMap { botStates.InSaveMap = true }                     else{ botStates.InSaveMap = false }
        if distFromDest <= float64(minDist) { botStates.AtRange = true }        else{ botStates.AtRange = false }
        if exist := getConf(conf["Route"],"Map", curMap); exist != nil {
            botStates.OnTheRoad = true
        }else{ botStates.OnTheRoad = false }
        if targetMob >= 0 { botStates.HasTargetMob = true }                     else{ botStates.HasTargetMob = false }
        if targetItem >= 0 { botStates.HasTargetItem = true }                   else{ botStates.HasTargetItem = false }
        // #################################
        if botStates.HasTargetItem == true {  botStates.HasTargetMob = false }
        if botStates.HasTargetItem == true {  botStates.ReadyToTp = false }
        if botStates.HasTargetMob == true {  botStates.ReadyToTp = false }
        if botStates.OnTheRoad == true {  botStates.ReadyToTp = false }
        // #################################


        if botStates == (States{InLockMap:true, ReadyToTp:true}) ||
           botStates == (States{InLockMap:true, HasDest:true, ReadyToTp:true}) {
            if useTPLockMap == 1 {
                resetStates()
                time.Sleep(800 * time.Millisecond)
                tpId := int(binary.LittleEndian.Uint16([]byte{26,0}))
                tpLv := int(binary.LittleEndian.Uint16([]byte{1,0}))
                sendUseSkill(tpId, tpLv, accountId)
                time.Sleep(1300 * time.Millisecond)
                TPstartTime = time.Now()
            }
            if useTPLockMap == 2 {
                inventID := itemInInventory(601,1) // fly wing
                if inventID > -1  {
                    resetStates()
                    time.Sleep(800 * time.Millisecond)
                    sendUseItem(inventID)
                    time.Sleep(1300 * time.Millisecond)
                }
                TPstartTime = time.Now()
            }
        }

        if botStates == (States{InLockMap:true, HasTargetMob: true}) ||
           botStates == (States{InLockMap:true, HasTargetMob: true, HasDest:true}) ||
           botStates == (States{InLockMap:true, HasTargetMob: true, HasDest:true, AtRange:true}) {
            MUmobList.Lock() ;  mob := mobList[targetMob] ; MUmobList.Unlock()
            nextPoint = mob.Coords
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
            pathIndex = 1 ; minDist = attackDist
        }

        if botStates == (States{InLockMap:true, HasTargetItem: true}) ||
           botStates == (States{InLockMap:true, HasTargetItem: true, HasDest:true}) {
            MUgroundItems.Lock() ; item := groundItems[targetItem] ; MUgroundItems.Unlock()
            nextPoint = item.Coords
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
            pathIndex = 1 ; minDist = 2;
        }

        if botStates == (States{InLockMap:true, HasTargetItem: true, HasDest:true, AtRange:true})  {
            itemBin := make([]byte, 4) ;
            binary.LittleEndian.PutUint32(itemBin, uint32(targetItem))
            sendToServer("0362", itemBin)
            time.Sleep(200 * time.Millisecond)
        }

        if botStates == (States{InLockMap:true, HasTargetMob: true, HasDest:true, AtRange:true})  {
            skill := conf["SkillTarget"][attackIndex]
            delay := 0
            if skill.(CSkillTarget).Id != -1 {
                sendUseSkill(skill.(CSkillTarget).Id, skill.(CSkillTarget).Lv, targetMob)
            }else{
                arrayBin := []byte{}
                mobBin := make([]byte, 4)
                binary.LittleEndian.PutUint32(mobBin, uint32(targetMob))
                arrayBin = append(arrayBin,mobBin...)
                // 0 = unique autoattack / 7 = start autoattack
                arrayBin = append(arrayBin,byte(7))
                sendToServer("0437", arrayBin)
                delay = 1000
            }
            if attackIndex < len(conf["SkillTarget"])-1 { attackIndex++ }else{ attackIndex = 0 }
            time.Sleep(time.Duration(delay) * time.Millisecond)
        }


        if botStates == (States{InLockMap:true, HasDest:true}) ||
           botStates == (States{OnTheRoad:true, HasDest:true}) ||
           botStates == (States{OnTheRoad:true, HasDest:true, InSaveMap:true}) ||
           botStates == (States{InLockMap:true, HasDest:true, HasTargetMob:true}) ||
           botStates == (States{InLockMap:true, HasDest:true, HasTargetItem:true}) {
            if pathIndex > len(curPath)-2 {
                nextStep = nextPoint
            }else{
                nextStep = Coord{curPath[pathIndex].X,curPath[pathIndex].Y}
            }

            if getDist(curCoord, nextStep) < 6{ pathIndex += 8 }
            sendToServer("035F",coordsTo24Bits(nextStep.X,nextStep.Y))
            time.Sleep(50 * time.Millisecond)
        }

        if botStates == (States{InLockMap:true, HasDest:true, AtRange:true}) ||
           botStates == (States{OnTheRoad:true, HasDest:true, AtRange:true}){
            curPath = nil ; pathIndex = 0 ; nextPoint = Coord{X:0, Y:0}
        }

        if botStates == (States{OnTheRoad:true}) ||
           botStates == (States{OnTheRoad:true, InSaveMap:true}) {
            if exist := getConf(conf["Route"],"Map",curMap); exist != nil {
                nextPoint = Coord{X:exist.(CRoute).X, Y:exist.(CRoute).Y}
                curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
                pathIndex = 0 ; minDist = 1;
            }
        }

        if botStates == (States{InLockMap:true}) {
            nextPoint = randomPoint(lgatMaps[curMap],curCoord, 80)
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
            pathIndex = 0 ; minDist = 1;
        }

    }

    fmt.Printf("# lel # \n")
}

var (
    strInfo = ""
    strMobs = ""
    strBuffs = ""
    strInventoryItems = ""
    strGroundItems = ""
)
func infoUILoop() {
    for { time.Sleep(200 * time.Millisecond)

        HPpc := int(float32(HPLeft)/float32(HPMax)*100)
        SPpc := int(float32(SPLeft)/float32(SPMax)*100)

        strInfo = "HP : "+Itos(HPLeft)+"/"+Itos(HPMax)+"("+Itos(HPpc)+"%) "
        strInfo += "| SP : "+Itos(SPLeft)+"/"+Itos(SPMax)+" ("+Itos(SPpc)+"%)"
        strInfo += "| W: "+Itos(maxWeight)+"/"+Itos(weight)




        strBuffs = ""
        MUbuffList.Lock()
        for _, kkk := range sortIntKeys(keyMap(buffList)) {
            iii := buffList[kkk]
            timeLeft := iii[0]-( (time.Now().Sub(time.Unix(iii[1], 0))).Milliseconds() )
            strBuffs += "["+Itos(kkk)+"] " + fmt.Sprintf("%v",timeLeft)+"\n"
        }
        MUbuffList.Unlock()

        strInventoryItems = ""
        MUinventoryItems.Lock()
        for _, kkk := range sortIntKeys(keyMap(inventoryItems)) {
            iii := inventoryItems[kkk]
            strInventoryItems += "["+Itos(iii.ItemID)+"] "+Itos(iii.Amount)  +" ea \n"
        }
        MUinventoryItems.Unlock()

        strMobs = ""
        MUmobList.Lock()
        for _, kkk := range sortIntKeys(keyMap(mobList)) {
            mm := mobList[kkk]
            strMobs += "["+Itos(kkk)+"] ("+Itos(mm.Coords.X)+" / "+Itos(mm.Coords.Y)+") "+Itos(mm.MobID) +"\n"
        }
        MUmobList.Unlock()

        strGroundItems = ""
        MUgroundItems.Lock()
        for _, kkk := range sortIntKeys(keyMap(groundItems)) {
            ii := groundItems[kkk]
            strGroundItems += "["+Itos(kkk)+"] ("+Itos(ii.Coords.X)+" / "+Itos(ii.Coords.Y)+")"
            strGroundItems += Itos(ii.ItemID) +" "+Itos(ii.Amount) +"\n"
        }
        MUgroundItems.Unlock()

    }
}
