package main

import(
    // "fmt"
    "time"
    "sort"
    "encoding/binary"
    "sync"
)


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

    curCoord = Coord{X:0, Y:0}
    nextPoint = Coord{X:0, Y:0}
	curMap = ""
    lockMap = ""
    saveMap = ""

    pathIndex = 0
	curPath = []Coord{}
    nextStep = Coord{}
    minDist = 1

    MUmobList sync.Mutex
    mobList = map[int]Mob{}
    MUgroundItems sync.Mutex
    groundItems = map[int]Item{}

    targetMob = -1
    targetMobDead = -2
    targetItem = -1
    targetItemLooted = -2
    ignoreItem = []int{}

    attackDist = 1

    tpUse = 0
    tpTime = 0
    useAttacks = [][]string{}
    attackIndex = 0

    hp = 0
    maxHP = 0
    maxWeight = 0
    weight = 0
    sp = 0
    spMAx = 0

)

// var states = []string{
//     "standing",
//     "moving",
// }
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

func botLoop() {

    curPath = nil

    startTime := time.Now()
    elapsed := time.Now()

    TPstartTime := time.Now()
    TPelapsed := time.Now()

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

        if targetMob == targetMobDead { targetMob = -1; targetMobDead = -2; nextPoint = Coord{X:0, Y:0} }
        if targetItem == targetItemLooted { targetItem = -1; targetItemLooted = -2; nextPoint = Coord{X:0, Y:0} }

        MUgroundItems.Lock()
        targetItem = -1
        for kk,vv := range groundItems { if getDist(vv.Coords, curCoord) > 40 { delete(groundItems, kk) } }
        for kk,vv := range groundItems {
            if curMap == lockMap {
            if getDist(vv.Coords, curCoord) < 25 {
                targetItem = kk
            }}
        }
        MUgroundItems.Unlock()

        MUmobList.Lock()
        targetMob = -1
        for kk,vv := range mobList { if getDist(vv.Coords, curCoord) > 40 { delete(mobList, kk) } }
        distMobList := map[float64]int{}
        for kk,vv := range mobList { distMobList[getDist(vv.Coords, curCoord)] = kk }
        keys := []float64{}
        for kk,_ := range distMobList { keys = append(keys,kk) }
        sort.Sort(sort.Float64Slice(keys))
        for i := len(keys)-1; i >= 0; i-- {
            if curMap != lockMap { continue }
            mob := mobList[distMobList[keys[i]]]
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

        distFromDest := getDist(curCoord, nextPoint)

        if tpUse > 0 {
            botStates.ReadyToTp = false
            if TPelapsed.Sub(TPstartTime).Seconds() > float64(tpTime) {
                botStates.ReadyToTp = true
            }
            TPelapsed = time.Now()
        }

        // #################################
        // #################################
        if nextPoint != (Coord{X:0, Y:0}) { botStates.HasDest = true }          else{ botStates.HasDest = false }
        if curMap == lockMap { botStates.InLockMap = true }                     else{ botStates.InLockMap = false }
        if curMap == saveMap { botStates.InSaveMap = true }                     else{ botStates.InSaveMap = false }
        if distFromDest <= float64(minDist) { botStates.AtRange = true }        else{ botStates.AtRange = false }
        if _, exist := route[curMap]; exist { botStates.OnTheRoad = true }      else{ botStates.OnTheRoad = false }
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
            if tpUse == 1 {
                resetStates()
                time.Sleep(800 * time.Millisecond)
                sendToServer("0438",[]byte{1,0,26,0,74,188,30,0})
                time.Sleep(1300 * time.Millisecond)
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
                // fmt.Printf("# loot loot # \n")
                sendToServer("0362", itemBin)
                time.Sleep(200 * time.Millisecond)
        }

        if botStates == (States{InLockMap:true, HasTargetMob: true, HasDest:true, AtRange:true})  {
            arrayBin := []byte{}
            att := useAttacks[attackIndex]
            skillIDBin := make([]byte, 2) ;
            binary.LittleEndian.PutUint16(skillIDBin, uint16(Stoi(att[0])))
            skillLVBin := make([]byte, 2) ;
            binary.LittleEndian.PutUint16(skillLVBin, uint16(Stoi(att[1])))
            mobBin := make([]byte, 4) ;
            binary.LittleEndian.PutUint32(mobBin, uint32(targetMob))
            delay := Stoi(att[2])

            if Stoi(att[0]) != 0 {
                arrayBin = append(arrayBin,skillLVBin...)
                arrayBin = append(arrayBin,skillIDBin...)
                arrayBin = append(arrayBin,mobBin...)
                sendToServer("0438", arrayBin)
            }else{
                arrayBin = append(arrayBin,mobBin...)
                // 0 = unique autoattack / 7 = start autoattack
                arrayBin = append(arrayBin,byte(0))
                sendToServer("0437", arrayBin)
            }
            if attackIndex < len(useAttacks)-1 { attackIndex++ }else{ attackIndex = 0 }
            time.Sleep(time.Duration(delay) * time.Millisecond)
        }


        if botStates == (States{InLockMap:true, HasDest:true}) ||
           botStates == (States{OnTheRoad:true, HasDest:true}) ||
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

        if botStates == (States{OnTheRoad:true}) {
            nextPoint = (Coord{X:route[curMap][0], Y:route[curMap][1]})
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
            pathIndex = 0 ; minDist = 1;
        }

        if botStates == (States{InLockMap:true}) {
            nextPoint = randomPoint(lgatMaps[curMap],curCoord, 80)
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
            pathIndex = 0 ; minDist = 1;
        }

    }

}

var strInfo = ""
var strMobs = ""
var strGroundItems = ""

func infoUILoop() {
    for { time.Sleep(200 * time.Millisecond)

        strInfo = Itos(hp)+" ## "+Itos(maxHP)+" ## "+Itos(maxWeight)+" ## "+Itos(weight)+" ## "+Itos(sp)+" ## "+Itos(spMAx)+" ## "

        keys := []int{}
        strMobs = ""
        MUmobList.Lock()
        for kk,_ := range mobList { keys = append(keys,kk) }
        sort.Ints(keys)
        for _,kkk := range keys {
            mm := mobList[kkk]
            strMobs += "["+Itos(kkk)+"] ("+Itos(mm.Coords.X)+" / "+Itos(mm.Coords.Y)+") "+Itos(mm.MobID) +" "+mm.Name +"\n"
        }
        MUmobList.Unlock()

        keys = []int{}
        strGroundItems = ""
        MUgroundItems.Lock()
        for kk,_ := range groundItems { keys = append(keys,kk) }
        sort.Ints(keys)
        for _,kkk := range keys {
            ii := groundItems[kkk]
            strGroundItems += "["+Itos(kkk)+"] ("+Itos(ii.Coords.X)+" / "+Itos(ii.Coords.Y)+") "+Itos(ii.ItemID) +" "+Itos(ii.Amount) +"\n"
        }
        MUgroundItems.Unlock()

    }
}
