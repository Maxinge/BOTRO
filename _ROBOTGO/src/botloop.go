package main

import(
    "fmt"
    "time"
    "sort"
    "encoding/binary"
    "sync"
)

type States struct {
    InLockMap bool
    IsWalking bool
    GoToMob bool
    GoToItem bool
    OnTheRoad bool
    ReadyToTp bool
    InSaveMap bool
    HasTargetMob bool
    HasTargetItem bool
    InCombat bool
    IsLooting bool
}

var(
    botStates = States{}
    checkState = botStates
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
    targetMobDead = -1
    targetItem = -1
    targetItemLooted = -1
    ignoreItem = []int{}

    attackDist = 1

    tpSearch = -1
    tpTime = 5
    useAttacks = [][]string{}
    attackIndex = 0
)

// var states = []string{
//     "standing",
//     "moving",
// }
func resetStates(){
    fmt.Printf("## resetStates ## \n")
    botStates = States{};
    pathIndex = 0 ; curPath = nil ;
    targetItem = -1; targetMob = -1;
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

    // TPstartTime := time.Now()
    // TPelapsed := time.Now()

    startBotLoop:
    for { time.Sleep(200 * time.Millisecond)
        if curCoord == (Coord{X:0, Y:0}){ continue }

        if checkState == botStates {
            timeInState = elapsed.Sub(startTime).Seconds()
            if timeInState > float64(15) { resetStates() }
            elapsed = time.Now()
        }else{
            startTime = time.Now()
        }
        checkState = botStates

        // #################################
        // #################################

        if curMap == lockMap { botStates.InLockMap = true }                 else{ botStates.InLockMap = false }
        if curMap == saveMap { botStates.InSaveMap = true }                 else{ botStates.InSaveMap = false }
        if _, exist := route[curMap]; exist { botStates.OnTheRoad = true }  else{ botStates.OnTheRoad = false }
        if curPath != nil { botStates.IsWalking = true }                    else{ botStates.IsWalking = false }
        if targetMob >= 0 {
            botStates.HasTargetMob = true;
        }else{ botStates.HasTargetMob = false ; botStates.GoToMob = false; botStates.InCombat = false}
        if targetItem >= 0 {
            botStates.HasTargetItem = true; botStates.HasTargetMob = false; botStates.InCombat = false;
        }else{ botStates.HasTargetItem = false ; botStates.GoToItem = false; botStates.IsLooting = false;}

        // #################################

        if botStates == (States{InLockMap:true, HasTargetItem:true, IsLooting:true}) {
            if targetItem == targetItemLooted { targetItem = -1; botStates.IsLooting = false; continue startBotLoop }
            itemBin := make([]byte, 4) ;
            binary.LittleEndian.PutUint32(itemBin, uint32(targetItem))
            // fmt.Printf("# loot loot # \n")
            sendToServer("0362", itemBin)
            time.Sleep(200 * time.Millisecond)
        }

        if botStates == (States{InLockMap:true, HasTargetItem:true}) ||
           botStates == (States{InLockMap:true, IsWalking:true, HasTargetItem:true}) ||
           botStates == (States{InLockMap:true, IsWalking:true, HasTargetItem:true, GoToItem:true}) {
            MUgroundItems.Lock()
            item := groundItems[targetItem]
            MUgroundItems.Unlock()
            nextPoint = item.Coords
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
            pathIndex = 1 ; minDist = 2;
            botStates.GoToItem = true
        }

        if botStates == (States{InLockMap:true, HasTargetMob:true, InCombat:true}) {
            if targetMob == targetMobDead { targetMob = -1; botStates.InCombat = false; continue startBotLoop }
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
            continue startBotLoop
        }

        if botStates == (States{InLockMap:true, HasTargetMob:true}) ||
           botStates == (States{InLockMap:true, IsWalking:true, HasTargetMob:true}) ||
           botStates == (States{InLockMap:true, IsWalking:true, HasTargetMob:true, GoToMob:true}) {
            botStates.InCombat = false
            MUmobList.Lock() ;  mob := mobList[targetMob] ; MUmobList.Unlock()
            line := linearInterpolation(curCoord, mob.Coords)
        	for _,vv := range line {
        		gatcell := lgatMaps[curMap].cells[vv.X][vv.Y]
        		if !isValidCell(gatcell) {
                    MUmobList.Lock()
                    delete(mobList, targetMob);
                    MUmobList.Unlock()
                    targetMob = -1;
                    continue startBotLoop
                }
        	}
            nextPoint = mob.Coords
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
            pathIndex = 1 ; minDist = attackDist;
            botStates.GoToMob = true
        }

        if botStates == (States{OnTheRoad:true}) {
            // fmt.Printf("# newPath OnTheRoad # \n")
            nextPoint = (Coord{X:route[curMap][0], Y:route[curMap][1]})
            pathIndex = 0 ; minDist = 1;
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
        }

        if botStates == (States{InLockMap:true}) {
            // fmt.Printf("# newPath InLockMap # \n")
            nextPoint = randomPoint(lgatMaps[curMap],curCoord, 80)
            pathIndex = 0 ; minDist = 1;
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
        }

        if botStates == (States{InLockMap:true, IsWalking:true}) ||
           botStates == (States{OnTheRoad:true, IsWalking:true}) ||
           botStates == (States{InLockMap:true, IsWalking:true, HasTargetMob:true, GoToMob:true}) ||
           botStates == (States{InLockMap:true, IsWalking:true, HasTargetItem:true, GoToItem:true}) {
            if getDist(nextPoint, curCoord) <= float64(minDist) {
                if botStates.GoToMob == true { botStates.InCombat = true; botStates.GoToMob = false }
                if botStates.GoToItem == true { botStates.IsLooting = true; botStates.GoToItem = false }
                curPath = nil; pathIndex = 0 ; continue startBotLoop
            }
            if pathIndex > len(curPath)-2 {
                nextStep = nextPoint
            }else{
                nextStep = Coord{curPath[pathIndex].X,curPath[pathIndex].Y}
            }
            if getDist(curCoord, nextStep) < 6{ pathIndex += 8 }
            sendToServer("035F",coordsTo24Bits(nextStep.X,nextStep.Y))
            time.Sleep(50 * time.Millisecond)

            if botStates.GoToMob == true { continue startBotLoop }
            if botStates.GoToItem == true { continue startBotLoop }
        }

        // botStates.ReadyToTp = false
        // if TPelapsed.Sub(TPstartTime).Seconds() > float64(tpTime) {
        //     botStates.ReadyToTp = true
        //     TPstartTime = time.Now()
        //     continue startBotLoop
        // }
        // TPelapsed = time.Now()
        //
        // if botStates == (States{InLockMap:true, ReadyToTp:true}) ||
        //    botStates == (States{InLockMap:true, IsWalking:true, ReadyToTp:true}) {
        //     if tpSearch >= 0 {
        //         fmt.Printf("use tp -- \n")
        //         sendToServer("0x3804",[]byte{1,0,26,0,74,188,30,0})
        //         time.Sleep(1000 * time.Millisecond)
        //         resetStates()
        //         continue startBotLoop
        //     }
        // }

        // #################################
        // #################################

        MUgroundItems.Lock()
        targetItem = -1
        for kk,vv := range groundItems { if getDist(vv.Coords, curCoord) > 35 { delete(groundItems, kk) } }
        for kk,vv := range groundItems {
           if getDist(vv.Coords, curCoord) < 25 {
               MUgroundItems.Unlock()
               targetItem = kk; continue startBotLoop
           }
        }
        MUgroundItems.Unlock()

        MUmobList.Lock()
        targetMob = -1
        for kk,vv := range mobList { if getDist(vv.Coords, curCoord) > 35 { delete(mobList, kk) } }

        distMobList := map[float64]int{}
        for kk,vv := range mobList { distMobList[getDist(vv.Coords, curCoord)] = kk }
        keys := []float64{}
        for kk,_ := range distMobList { keys = append(keys,kk) }
        sort.Sort(sort.Float64Slice(keys))

        for _,kkk := range keys {
           mob := mobList[distMobList[kkk]]
            if getDist(mob.Coords, curCoord) < 25 {
                MUmobList.Unlock()
                targetMob = distMobList[kkk] ; continue startBotLoop
            }
        }
        MUmobList.Unlock()

    }

}

// func infoUILoop() {
//     for { time.Sleep(200 * time.Millisecond)
//
//         keys := []int{}
//         strMobs = ""
//         MUmobList.Lock()
//         for kk,_ := range mobList { keys = append(keys,kk) }
//         sort.Ints(keys)
//         for _,kkk := range keys {
//             mm := mobList[kkk]
//             strMobs += "["+Itos(kkk)+"] ("+Itos(mm.Coords.X)+" / "+Itos(mm.Coords.Y)+") "+Itos(mm.MobID) +" "+mm.Name +"\n"
//         }
//         MUmobList.Unlock()
//
//         keys = []int{}
//         strGroundItems = ""
//         MUgroundItems.Lock()
//         for kk,_ := range groundItems { keys = append(keys,kk) }
//         sort.Ints(keys)
//         for _,kkk := range keys {
//             ii := groundItems[kkk]
//             strGroundItems += "["+Itos(kkk)+"] ("+Itos(ii.Coords.X)+" / "+Itos(ii.Coords.Y)+") "+Itos(ii.ItemID) +" "+Itos(ii.Amount) +"\n"
//         }
//         MUgroundItems.Unlock()
//
//     }
// }
