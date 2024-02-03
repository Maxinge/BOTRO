package main

import(
    "fmt"
    "net"
    "time"
    "encoding/json"
    "reflect"
    "github.com/cimgui-go"
    "io/ioutil"
    "strings"
    "os"
    "strconv"
)

type Packet struct {
    Ident  string
    Desc  string
}

type Mob struct {
    MobID int
    Name string
    Coords Coord
    HP int
    HPLeft int
}

type Item struct {
    ItemID int
    Name string
    Type int
    Coords Coord
}

var (
    proxyCo net.Conn
    exit = make(chan bool)


    mapTextures = map[string]*imgui.Texture{}
    // maskTexture  *imgui.Texture

    gatMaps = map[string]ROGatMap{}

    lgatMaps = map[string]ROLGatMap{}

    profil map[string]interface{}
    packetsmap map[string]Packet

    curCoord = Coord{X:0, Y:0}
    nextPoint Coord
	curMap string
    lockMap string

	curPath = []Coord{}
    pathIndex = 0
    countStuck = 0
    checkStuck = curCoord

    route map[string][]int
    mobList = map[int]Mob{}

    // groundItems = map[int]

    targetMobs []int
    target = -1

    fctpack map[string]func([]byte, []byte)
)

func main() {

    fmt.Println("#--- ROBOTGO START ---#")
    fmt.Printf("current dir -- %v -- \n", CurDir())


    maps, _ := ioutil.ReadDir(CurDir()+"data/gats/")
    for _, m := range maps {
        if !m.IsDir() {
            name := strings.Split(m.Name(), ".gat")[0]
            fmt.Printf("name -- %v -- \n", name)
            loadGatMap(name)
        }
    }


    for kk, ggg := range gatMaps {

        fichier, _ := os.Create(CurDir()+"data/lgats/"+kk+".lgat")
        defer fichier.Close()

        ccc := []byte{}
        ww := int16ToBitString(ggg.width)
        ww1,_ := strconv.ParseInt(ww[0:8], 2, 8)
        ww2,_ := strconv.ParseInt(ww[8:16], 2, 8)
        bb := []byte{byte(ww1),byte(ww2)}
        ccc = append(ccc,bb...)

        hh := int16ToBitString(ggg.height)
        hh1,_ := strconv.ParseInt(hh[0:8], 2, 8)
        hh2,_ := strconv.ParseInt(hh[8:16], 2, 8)
        bbb := []byte{byte(hh1),byte(hh2)}
        ccc = append(ccc,bbb...)

        for x := 0; x < ggg.width; x++{
        for y := 0; y < ggg.height; y++{
             ccc = append(ccc,byte(ggg.cells[x][y].cell_type))
        }}

        fichier.Write(ccc)
    }



    proxyCo, _ = net.Dial("tcp", "127.0.0.1:6666")
    defer proxyCo.Close()

    json.Unmarshal([]byte(readFileString(CurDir()+"data/packets.json")), &packetsmap)
    // fmt.Printf("packetsmap -- %v -- \n", packetsmap)
    loadprofil()
    fctpackInit()

    loadLGatMap("morocc")
    loadLGatMap("payon")
    loadLGatMap("moc_fild10")
    loadLGatMap("moc_fild09")
    loadLGatMap("moc_fild15")
    loadLGatMap("moc_fild12")
    loadLGatMap("moc_fild18")
    loadLGatMap("in_sphinx1")


    go func() {

        curPath = nil
        // startTime := time.Now()
        // elapsed := time.Now()
        for {
            if curCoord == (Coord{X:0, Y:0}){ continue }

            // if elapsed.Sub(startTime).Milliseconds() > int64(3000) {
            //     sendToServer("0x3804",[]byte{1,0,26,0,74,188,30,0})
            //     time.Sleep(1000 * time.Millisecond)
            //     curPath = nil ; nextPoint = (Coord{X:0, Y:0}); countStuck = 0;
            //     startTime = time.Now()
            // }

            // elapsed = time.Now()


            for kk,vv := range mobList {
                if getDist(vv.Coords, curCoord) > 35 { delete(mobList, kk) }
            }

            // if target >= 0 {
            //     mob := mobList[target]
            //     // line := linearInterpolation(curCoord, mob.Coords)
			// 	// for _,vv := range line {
			// 	// 	gatcell := lgatMaps[curMap].cells[vv.X][vv.Y]
			// 	// 	if !isValidCell(gatcell) {
            //     //         delete(mobList, target) ; target = -1
            //     //         continue
            //     //     }
			// 	// }
            //
            //     mobpath := pathfind(curCoord, mob.Coords, lgatMaps[curMap])
            //     if len(mobpath) < 20 {
            //
            //         for ii := 0; ii < len(mobpath); ii++ {
            //             sendToServer("0x5f03",coordsTo24Bits(mobpath[ii].X,mobpath[ii].Y))
            //             fmt.Printf("mobpath -- %v / %v -- \n", mobpath[ii].X,mobpath[ii].Y )
            //             time.Sleep(200 * time.Millisecond)
            //             if getDist(curCoord, mob.Coords) <= 8 { break }
            //         }
            //
            //         // [1 0 14 0 238 105 211 6
            //         sendToServer("0x3804",[]byte{2,0,14,0,uint8(target),105,211,6})
            //         sendToServer("0x3804",[]byte{2,0,14,0,uint8(target),106,211,6})
            //         time.Sleep(1600 * time.Millisecond)
            //         delete(mobList, target)
            //         curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
            //         pathIndex = 0
            //     }
            //
            //     target = -1; continue
            // }

            if curMap == lockMap {
                for kk,vv := range mobList {
                    if getDist(vv.Coords, curCoord) < 30 {
                    if intInArray(vv.MobID, targetMobs){
                        target = kk
                        continue
                    }}
                }
            }

            if curMap == lockMap {
            // fmt.Println("#--- in lock map ---#")
            if nextPoint == (Coord{X:0, Y:0}) {
                nextPoint = randomPoint(lgatMaps[curMap],curCoord, 100)
                fmt.Printf("nextPoint -- %v -- \n", nextPoint )
                continue
            }}

            if curMap != lockMap {
            if nextPoint == (Coord{X:0, Y:0}){
            if _, exist := route[curMap]; exist {
                nextPoint = (Coord{X:route[curMap][0], Y:route[curMap][1]})
                fmt.Printf("curCoord -- %v -- nextPoint -- %v -- \n", curCoord, nextPoint )
            }}}

            if nextPoint != (Coord{X:0, Y:0}) {
            if curPath == nil {
                curPath = pathfind(curCoord, nextPoint, lgatMaps[curMap])
                pathIndex = 0
            }}

            if curPath != nil  {
                // fmt.Printf("countStuck -- %v - - \n", countStuck )
                if countStuck > 30 {
                    curPath = nil ; nextPoint = (Coord{X:0, Y:0}); countStuck = 0;
                }
                if checkStuck != curCoord{ countStuck = 0; }else{ countStuck++; }
                checkStuck = curCoord

                if getDist(nextPoint, curCoord) < 10 {
                    sendToServer("0x5f03",coordsTo24Bits(nextPoint.X, nextPoint.Y))
                    sendToServer("0x5f03",coordsTo24Bits(nextPoint.X, nextPoint.Y))
                    sendToServer("0x5f03",coordsTo24Bits(nextPoint.X, nextPoint.Y))
                    time.Sleep(1000 * time.Millisecond)
                    curPath = nil ; nextPoint = (Coord{X:0, Y:0}) ; continue
                }
                if pathIndex > len(curPath)-1 {
                    curPath = nil ; nextPoint = (Coord{X:0, Y:0}) ; continue
                }
                if getDist(curPath[pathIndex], curCoord) < 7 {
                    pathIndex += 9
                }else{
                    sendToServer("0x5f03",coordsTo24Bits(curPath[pathIndex].X,curPath[pathIndex].Y))
                    time.Sleep(100 * time.Millisecond)
                }
            }
            time.Sleep(300 * time.Millisecond)

        }
    }()

    // ########################
    backend := imgui.CreateBackend(imgui.NewGLFWBackend())
    backend.SetAfterCreateContextHook(func () {
        loadGatTexture("morocc")
    })

    backend.SetBeforeDestroyContextHook(func () {  })
    backend.SetBgColor(imgui.NewVec4(0.45, 0.55, 0.6, 1.0))
    backend.CreateWindow("ROBOTGO", 800, 800)

    backend.Run(func () {
        basePos := imgui.MainViewport().Pos()
        baseSize := imgui.MainViewport().Size()
        imgui.SetNextWindowPosV(imgui.NewVec2(basePos.X, basePos.Y), 0, imgui.NewVec2(0, 0))
        imgui.SetNextWindowSize(imgui.Vec2{X: baseSize.X, Y: baseSize.Y})
        imgui.Begin("robot")
        imgui.Text(fmt.Sprintf("Coords = X : %v / Y : %v", curCoord.X, curCoord.Y ))
        imgui.Text(fmt.Sprintf("Map : %v", curMap ))

        imgui.ImageV(mapTextures["morocc"].ID(), imgui.NewVec2(float32(baseSize.X/2),float32(baseSize.Y/2)), imgui.NewVec2(0, 0), imgui.NewVec2(1, 1), imgui.NewVec4(1, 1, 1, 1), imgui.NewVec4(0, 0, 0, 0))

        imgui.End()
        imgui.Render()
    })


    // ########################
    go func() {
        buffer := make([]byte, 100000)
        for {
            n, _ := proxyCo.Read(buffer)
            // if err != nil { fmt.Printf("err localConn -- %v -- \n", err); return }
            HexID := fmt.Sprintf("%#x", buffer[0:2])
            if _, exist := packetsmap[HexID]; !exist {
                fmt.Printf("[%v] len [%v] \t -> [%v]\n", HexID, len(buffer[:n]), buffer[:n])
                fmt.Printf("[%v] len [%v] \t -> [%v]\n", HexID, len(buffer[:n]), string(buffer[:n]))
            }else{
                function := reflect.ValueOf(fctpack[packetsmap[HexID].Ident])
                if function.Kind() == reflect.Func && fctpack[packetsmap[HexID].Ident] != nil{
                    args := []reflect.Value{reflect.ValueOf(buffer[0:2]),reflect.ValueOf(buffer[2:n])}
                    function.Call(args)
                }
            }
        }
    }()

    <-exit
}


func sendToServer(hexID string,data []byte){
    var ii int16
	fmt.Sscanf(hexID, "0x%x", &ii)
    bb := []byte{ byte(ii >> 8), byte(ii) }
    bb = append(bb,data...)
    proxyCo.Write(bb)
}
