package main

import(
    "fmt"
    "time"
    "sort"
    "encoding/binary"
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
    InLootItem bool
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

	curPath = []Coord{}
    nextStep = Coord{}
    minDist = 1
    // targetMobPath = []Coord{}
    // targetItemPath = []Coord{}

    pathIndex = 0

    targetMobDead = -1
    targetMob = -1
    targetItem = -1
    targetItemLooted = -1

    attackDist = 1

    tpSearch = -1
    tpTime = 5
    useGreed = -1
    combatTime = 10
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
    // targetItemPath = []Coord{}
    // targetMobPath = []Coord{}
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
            if timeInState > float64(20) { resetStates() }
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
            botStates.HasTargetMob = true; botStates.HasTargetItem = false;
        }else{ botStates.HasTargetMob = false }
        if targetItem >= 0 {
            botStates.HasTargetItem = true; botStates.HasTargetMob = false;
        }else{ botStates.HasTargetItem = false }



        // if botStates == (States{InLockMap:true, HasTargetItem:true}) {
        //
        //     MUgroundItems.Lock()
        //     item := groundItems[targetItem]
        //     MUgroundItems.Unlock()
        //
        //     if targetItem == targetItemLooted { targetItem = -1; continue startBotLoop }
        //     targetItemPath = pathfind(curCoord, item.Coords, lgatMaps[curMap])
        //     if len(targetItemPath) <= 0 { continue startBotLoop }
        //     sendToServer("0x5f03",coordsTo24Bits(targetItemPath[0].X,targetItemPath[0].Y))
        //     time.Sleep(200 * time.Millisecond)
        //     if getDist(curCoord, item.Coords) > 3 { continue startBotLoop }
        //
        //     lootBin := make([]byte, 4) ; lel := make([]byte, 4)
        //     binary.BigEndian.PutUint32(lootBin, uint32(targetItem))
        //     for ii := 0; ii < len(lootBin); ii++ {
        //         lel[len(lootBin) -1 - ii ] = lootBin[ii]
        //     }
        //     fmt.Printf("lootBin -- %v -- \n", lel)
        //     sendToServer("0x6203",lel)
        //     sendToServer("0x6203",lel)
        //     time.Sleep(100 * time.Millisecond)
        //     resetStates()
        //     continue startBotLoop
        // }
        //
        // if botStates == (States{InLockMap:true, HasTargetMob:true, InCombat:true}) {
        //

        //

        // }

        if botStates == (States{InLockMap:true, HasTargetMob:true, InCombat:true}) {

            if targetMob == targetMobDead { targetMob = -1; continue startBotLoop }

            mobBin := make([]byte, 4) ;
            binary.BigEndian.PutUint32(mobBin, uint32(targetMob))
            // for ii := 0; ii < len(lootBin); ii++ {
            //     lel[len(lootBin) -1 - ii ] = lootBin[ii]
            // }
            // fmt.Printf("lootBin -- %v -- \n", lel)

            lel := []byte{}
            lel = append(lel,[]byte{10,0,15}...)
            lel = append(lel,mobBin...)
            lel = append(lel,byte(6))


            fmt.Printf("kill -- %v -- \n", mobBin)
            sendToServer("0x3804",lel)
            time.Sleep(1400 * time.Millisecond)
            continue startBotLoop
            // '0114' => ['skill_use', 'v a4 a4 V3 v3 C',
            // [qw(skillID sourceID targetID tick src_speed dst_speed damage level option type)]]
        }

        if botStates == (States{InLockMap:true, IsWalking:true, HasTargetMob:true}) ||
           botStates == (States{InLockMap:true, IsWalking:true, HasTargetMob:true, GoToMob:true}){
            MUmobList.Lock()
            mob := mobList[targetMob]
            MUmobList.Unlock()
            line := linearInterpolation(curCoord, mob.Coords)
        	for _,vv := range line {
        		gatcell := lgatMaps[curMap].cells[vv.X][vv.Y]
        		if !isValidCell(gatcell) {
                    MUmobList.Lock()
                    delete(mobList, targetMob);
                    MUmobList.Unlock()
                    continue startBotLoop
                }
        	}
            curPath = pathfind(curCoord, mob.Coords, lgatMaps[curMap])
            pathIndex = 0 ; minDist = attackDist;
            botStates.GoToMob = true
        }


        if botStates == (States{OnTheRoad:true}) {
            fmt.Printf("# newPath OnTheRoad # \n")
            nextPoint = (Coord{X:route[curMap][0], Y:route[curMap][1]})
            pathIndex = 0 ; minDist = 0;
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
        }

        if botStates == (States{InLockMap:true}) {
            fmt.Printf("# newPath InLockMap # \n")
            nextPoint = randomPoint(lgatMaps[curMap],curCoord, 80)
            pathIndex = 0 ; minDist = 0;
            curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
        }

        if botStates == (States{InLockMap:true, IsWalking:true}) ||
           botStates == (States{OnTheRoad:true, IsWalking:true}) ||
           botStates == (States{InLockMap:true, IsWalking:true, HasTargetMob:true, GoToMob:true}) ||
           botStates == (States{OnTheRoad:true, IsWalking:true, HasTargetItem:true, GoToItem:true}) {
            if getDist(nextPoint, curCoord) <= float64(minDist) {
                if botStates.GoToMob == true { botStates.InCombat = true; botStates.GoToMob = false }
                if botStates.GoToItem == true { botStates.InLootItem = true; botStates.GoToItem = false }
                curPath = nil; continue startBotLoop
            }
            if pathIndex > len(curPath)-2 {
                nextStep = nextPoint
            }else{
                nextStep = Coord{curPath[pathIndex].X,curPath[pathIndex].Y}
            }
            if getDist(curCoord, nextStep) < 6{ pathIndex += 8 }
            sendToServer("0x5f03",coordsTo24Bits(nextStep.X,nextStep.Y))
            time.Sleep(50 * time.Millisecond)
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
               targetItem = kk ; continue startBotLoop
           }
        }
        MUgroundItems.Unlock()

        MUmobList.Lock()
        targetMob = -1
        for kk,vv := range mobList { if getDist(vv.Coords, curCoord) > 35 { delete(mobList, kk) } }
        for kk,vv := range mobList {
            if getDist(vv.Coords, curCoord) < 25 {
            if intInArray(vv.MobID, targetMobs){
                MUmobList.Unlock()
                targetMob = kk ; continue startBotLoop
            }}
        }
        MUmobList.Unlock()

    }

}

func infoUILoop() {
    for { time.Sleep(200 * time.Millisecond)

        strMobs = ""
        keys := []int{}
        MUmobList.Lock()
        for kk,_ := range mobList { keys = append(keys,kk) }
        sort.Ints(keys)
        for _,kkk := range keys {
            mm := mobList[kkk]
            strMobs += "["+Itos(kkk)+"] ("+Itos(mm.Coords.X)+" / "+Itos(mm.Coords.Y)+") "+Itos(mm.MobID) +" "+mm.Name +"\n"
        }
        MUmobList.Unlock()

        strGroundItems = ""
        MUgroundItems.Lock()
        for _,vv := range groundItems {
            strGroundItems += "("+Itos(vv.Coords.X)+" / "+Itos(vv.Coords.Y)+") "+Itos(vv.ItemID) +" "+Itos(vv.Amount) +"\n"
        }
        MUgroundItems.Unlock()

    }
}
