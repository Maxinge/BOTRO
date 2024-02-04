package main

import(
    // "fmt"
    "time"
    "sort"
)

type States struct {
    IsNotMoving bool
    InLockMap bool
    HasDest bool
    OnTheRoad bool
    ReadyToTp bool
    InSaveMap bool
    HasTargetMob bool
    HasTargetItem bool
}

var(
    botStates = States{}

    curCoord = Coord{X:0, Y:0}
    nextPoint = Coord{X:0, Y:0}
	curMap = ""
    lockMap = ""
    saveMap = ""

	curPath = []Coord{}
    targetMobPath = []Coord{}

    pathIndex = 0
    countNotMoving = 0
    checkNotMoving = Coord{X:0, Y:0}
    targetMob = -1
    targetItem = -1

    useTpSearch = -1
    useGreed = -1

)

// var states = []string{
//     "standing",
//     "moving",
// }

func botLoop() {

    curPath = nil

    startTime := time.Now()
    elapsed := time.Now()

    for { time.Sleep(200 * time.Millisecond)
        if curCoord == (Coord{X:0, Y:0}){ continue }

        if countNotMoving > 30 { botStates.IsNotMoving = true }
        if checkNotMoving != curCoord{ countNotMoving = 0; botStates.IsNotMoving = false }else{ countNotMoving++; }
        checkNotMoving = curCoord

        botStates.ReadyToTp = false
        if elapsed.Sub(startTime).Milliseconds() > int64(3000) {
            botStates.ReadyToTp = true
            startTime = time.Now()
        }
        elapsed = time.Now()

        if curMap == lockMap { botStates.InLockMap = true }                 else{ botStates.InLockMap = false }
        if curMap == saveMap { botStates.InSaveMap = true }                 else{ botStates.InSaveMap = false }
        if _, exist := route[curMap]; exist { botStates.OnTheRoad = true }  else{ botStates.OnTheRoad = false }
        if curPath != nil { botStates.HasDest = true }                      else{ botStates.HasDest = false }
        if targetMob >= 0 { botStates.HasTargetMob = true }                 else{ botStates.HasTargetMob = false }
        if targetItem >= 0 { botStates.HasTargetItem = true }                else{ botStates.HasTargetItem = false }


        // #################################
        // #################################

        if botStates == (States{InLockMap:true, HasDest:true, HasTargetMob:true}) {
            mob := mobList[targetMob]
            // line := linearInterpolation(curCoord, mob.Coords)
        	// for _,vv := range line {
        	// 	gatcell := lgatMaps[curMap].cells[vv.X][vv.Y]
        	// 	if !isValidCell(gatcell) {
                    // delete(mobList, target) ; target = -1

            //         continue
            //     }
        	// }
            targetMobPath = pathfind(curCoord, mob.Coords, lgatMaps[curMap])
            time.Sleep(1000 * time.Millisecond)
            delete(mobList, targetMob) ; targetMob = -1
        }

        if botStates == (States{InLockMap:true, HasDest:true, IsNotMoving:true}) ||
           botStates == (States{OnTheRoad:true, HasDest:true, IsNotMoving:true})  {
               curPath = nil ; countNotMoving = 0; botStates.IsNotMoving = false
        }

        if botStates == (States{ReadyToTp:true, InSaveMap:false}) {
            if useTpSearch >= 0 { }
            // sendToServer("0x3804",[]byte{1,0,26,0,74,188,30,0})
            // time.Sleep(1000 * time.Millisecond)
            // curPath = nil ; nextPoint = (Coord{X:0, Y:0}); countNotMoving = 0;
        }

        if botStates == (States{OnTheRoad:true}) {
            nextPoint = (Coord{X:route[curMap][0], Y:route[curMap][1]})
            pathIndex = 0 ; curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
        }

        if botStates == (States{InLockMap:true}) {
            nextPoint = randomPoint(lgatMaps[curMap],curCoord, 100)
            pathIndex = 0 ; curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
        }

        if botStates == (States{InLockMap:true, HasDest:true}) ||
           botStates == (States{OnTheRoad:true, HasDest:true})  {
            for kk,vv := range mobList {
                if getDist(vv.Coords, curCoord) < 25 {
                if intInArray(vv.MobID, targetMobs){
                    targetMob = kk ; break
                }}
            }
            if getDist(nextPoint, curCoord) < 10 {
                sendToServer("0x5f03",coordsTo24Bits(nextPoint.X, nextPoint.Y))
                sendToServer("0x5f03",coordsTo24Bits(nextPoint.X, nextPoint.Y))
                time.Sleep(1000 * time.Millisecond)
                curPath = nil ; continue
            }
            if pathIndex > len(curPath)-1 {
                curPath = nil ; continue
            }
            if getDist(curPath[pathIndex], curCoord) < 7 {
                pathIndex += 9
            }else{
                sendToServer("0x5f03",coordsTo24Bits(curPath[pathIndex].X,curPath[pathIndex].Y))
                time.Sleep(100 * time.Millisecond)
            }
        }



        // if targetItem >= 0 {
        //     item := groundItems[targetItem]
        //     itempath := pathfind(curCoord, item.Coords, lgatMaps[curMap])
        //
        //     if len(itempath) > 0 {
        //         ii := 0
        //         for {
        //             sendToServer("0x5f03",coordsTo24Bits(itempath[ii].X,itempath[ii].Y))
        //             fmt.Printf("itempath -- %v / %v -- \n", itempath[ii].X,itempath[ii].Y )
        //             time.Sleep(250 * time.Millisecond)
        //             ii++
        //             if getDist(curCoord, item.Coords) <= 1 { break }
        //
        //         }
        //     }
        //
        //     lootBin := make([]byte, 4) ; lel := make([]byte, 4)
        //     binary.BigEndian.PutUint32(lootBin, uint32(targetItem))
        //
        //     for ii := 0; ii < len(lootBin); ii++ {
        //         lel[len(lootBin) -1 - ii ] = lootBin[ii]
        //     }
        //     fmt.Printf("lootBin -- %v -- \n", lel)
        //     time.Sleep(300 * time.Millisecond)
        //     sendToServer("0x6203",lel)
        //     time.Sleep(300 * time.Millisecond)
        //     curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
        //     pathIndex = 0
        //
        //     targetItem = -1; continue
        // }
        //
        // if target >= 0 {
        //     mob := mobList[target]
        //     line := linearInterpolation(curCoord, mob.Coords)
        // 	for _,vv := range line {
        // 		gatcell := lgatMaps[curMap].cells[vv.X][vv.Y]
        // 		if !isValidCell(gatcell) {
        //             delete(mobList, target) ; target = -1
        //             continue
        //         }
        // 	}
        //
        //     mobpath := pathfind(curCoord, mob.Coords, lgatMaps[curMap])
        //     if len(mobpath) < 20 {
        //         ii := 0
        //         for {
        //             if getDist(curCoord, mob.Coords) <= 8 { break }
        //             sendToServer("0x5f03",coordsTo24Bits(mobpath[ii].X,mobpath[ii].Y))
        //             fmt.Printf("mobpath -- %v / %v -- \n", mobpath[ii].X,mobpath[ii].Y )
        //             time.Sleep(200 * time.Millisecond)
        //
        //             ii ++
        //         }
        //     }
        //
        //     fmt.Printf("kill -- %v -- \n", uint8(target) )
        //
        //     sendToServer("0x3804",[]byte{10,0,15,0,uint8(target),109,222,6})
        //     sendToServer("0x3804",[]byte{10,0,15,0,uint8(target),110,222,6})
        //     time.Sleep(1600 * time.Millisecond)
        //     delete(mobList, target)
        //     curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
        //     pathIndex = 0
        //
        //     target = -1; continue
        // }
        //
        // if curMap == lockMap {
        //     for kk,vv := range groundItems {
        //         if getDist(vv.Coords, curCoord) < 20 {
        //             targetItem = kk
        //             break
        //         }
        //     }

        // }



    }
}

func infoUILoop() {
    for { time.Sleep(200 * time.Millisecond)
        strMobs = ""
        keys := []int{}
        mu.Lock()
        for kk,_ := range mobList { keys = append(keys,kk) }
        sort.Ints(keys)
        for _,kkk := range keys {
            mm := mobList[kkk]
            strMobs += "("+Itos(mm.Coords.X)+" / "+Itos(mm.Coords.Y)+") "+Itos(mm.MobID) +" "+mm.Name +"\n"
        }
        mu.Unlock()

        strGroundItems = ""
        for _,vv := range groundItems {
            strGroundItems += "("+Itos(vv.Coords.X)+" / "+Itos(vv.Coords.Y)+") "+Itos(vv.ItemID) +" "+Itos(vv.Amount) +"\n"
        }

        mu.Lock()
        for kk,vv := range mobList {
            if getDist(vv.Coords, curCoord) > 35 { delete(mobList, kk) }
        }
        mu.Unlock()

        for kk,vv := range groundItems {
            if getDist(vv.Coords, curCoord) > 35 { delete(groundItems, kk) }
        }
    }
}
