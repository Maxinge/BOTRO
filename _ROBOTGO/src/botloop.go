package main

import(
    // "fmt"
    "time"
    "sort"
)

var(
    curCoord = Coord{X:0, Y:0}
    nextPoint = Coord{X:0, Y:0}
	curMap = ""
    lockMap = ""
	curPath = []Coord{}
    pathIndex = 0
    countStuck = 0
    checkStuck = Coord{X:0, Y:0}
    target = -1
    targetItem = -1
)

var states = []string{
    "standing",
    "moving",
}

func botLoop() {

    curPath = nil
    // startTime := time.Now()
    // elapsed := time.Now()
    for { time.Sleep(200 * time.Millisecond)
        if curCoord == (Coord{X:0, Y:0}){ continue }
        
        // if elapsed.Sub(startTime).Milliseconds() > int64(3000) {
        //     sendToServer("0x3804",[]byte{1,0,26,0,74,188,30,0})
        //     time.Sleep(1000 * time.Millisecond)
        //     curPath = nil ; nextPoint = (Coord{X:0, Y:0}); countStuck = 0;
        //     startTime = time.Now()
        // }

        // elapsed = time.Now()

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
        //     for kk,vv := range mobList {
        //         if getDist(vv.Coords, curCoord) < 30 {
        //         if intInArray(vv.MobID, targetMobs){
        //             target = kk
        //             break
        //         }}
        //     }
        // }
        //
        // if curMap == lockMap {
        // // fmt.Println("#--- in lock map ---#")
        // if nextPoint == (Coord{X:0, Y:0}) {
        //     nextPoint = randomPoint(lgatMaps[curMap],curCoord, 100)
        //     fmt.Printf("nextPoint -- %v -- \n", nextPoint )
        //     // maskTexture = nil
        //     continue
        // }}
        //
        // if curMap != lockMap {
        // if nextPoint == (Coord{X:0, Y:0}){
        // if _, exist := route[curMap]; exist {
        //     nextPoint = (Coord{X:route[curMap][0], Y:route[curMap][1]})
        //     fmt.Printf("curCoord -- %v -- nextPoint -- %v -- \n", curCoord, nextPoint )
        // }}}
        //
        // if nextPoint != (Coord{X:0, Y:0}) {
        // if curPath == nil {
        //     curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
        //     pathIndex = 0
        // }}
        //
        // if curPath != nil  {
        //
        //     // fmt.Printf("countStuck -- %v - - \n", countStuck )
        //     if countStuck > 30 {
        //         curPath = nil ; nextPoint = (Coord{X:0, Y:0}); countStuck = 0;
        //     }
        //     if checkStuck != curCoord{ countStuck = 0; }else{ countStuck++; }
        //     checkStuck = curCoord
        //
        //     if getDist(nextPoint, curCoord) < 10 {
        //         sendToServer("0x5f03",coordsTo24Bits(nextPoint.X, nextPoint.Y))
        //         sendToServer("0x5f03",coordsTo24Bits(nextPoint.X, nextPoint.Y))
        //         sendToServer("0x5f03",coordsTo24Bits(nextPoint.X, nextPoint.Y))
        //         time.Sleep(1000 * time.Millisecond)
        //         curPath = nil ; nextPoint = (Coord{X:0, Y:0}) ; continue
        //     }
        //     if pathIndex > len(curPath)-1 {
        //         curPath = nil ; nextPoint = (Coord{X:0, Y:0}) ; continue
        //     }
        //     if getDist(curPath[pathIndex], curCoord) < 7 {
        //         pathIndex += 9
        //     }else{
        //         sendToServer("0x5f03",coordsTo24Bits(curPath[pathIndex].X,curPath[pathIndex].Y))
        //         time.Sleep(100 * time.Millisecond)
        //     }
        // }


    }
}

func infoUILoop() {
    for { time.Sleep(200 * time.Millisecond)
        strMobs = ""
        keys := []int{}
        for kk,_ := range mobList { keys = append(keys,kk) }
        sort.Ints(keys)
        for _,kkk := range keys {
            mm := mobList[kkk]
            strMobs += "("+Itos(mm.Coords.X)+" / "+Itos(mm.Coords.Y)+") "+Itos(mm.MobID) +" "+mm.Name +"\n"
        }

        strGroundItems = ""
        for _,vv := range groundItems {
            strGroundItems += "("+Itos(vv.Coords.X)+" / "+Itos(vv.Coords.Y)+") "+Itos(vv.ItemID) +" "+Itos(vv.Amount) +"\n"
        }

        for kk,vv := range mobList {
            if getDist(vv.Coords, curCoord) > 35 { delete(mobList, kk) }
        }

        for kk,vv := range groundItems {
            if getDist(vv.Coords, curCoord) > 35 { delete(groundItems, kk) }
        }
    }
}
