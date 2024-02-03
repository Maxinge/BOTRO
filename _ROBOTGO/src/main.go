package main

import(
    "fmt"
    "net"
    "time"
    "encoding/json"
    "reflect"
    "github.com/cimgui-go"
    // "github.com/go-gl/gl/v2.1/gl"
    // "github.com/go-gl/glfw/v3.3/glfw"
    "io/ioutil"
    "strings"
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

    // gatMaps = map[string]ROGatMap{}
    lgatMaps = map[string]ROLGatMap{}

    profil map[string]interface{}
    packetsmap map[string]Packet

    curCoord = Coord{X:0, Y:0}
    nextPoint = Coord{X:0, Y:0}
	curMap = ""
    lockMap = ""

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

    proxyCo, _ = net.Dial("tcp", "127.0.0.1:6666")
    defer proxyCo.Close()

    json.Unmarshal([]byte(readFileString(CurDir()+"data/packets.json")), &packetsmap)
    // fmt.Printf("packetsmap -- %v -- \n", packetsmap)
    loadprofil()
    fctpackInit()

    maps, _ := ioutil.ReadDir(CurDir()+"data/lgats/")
    for _, m := range maps {
        if !m.IsDir() {
            name := strings.Split(m.Name(), ".lgat")[0]
            fmt.Printf("name -- %v -- \n", name)
            loadLGatMap(name)
        }
    }


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
    go func() {
        buffer := make([]byte, 100000)
        for {
            n, _ := proxyCo.Read(buffer)
            // if err != nil { fmt.Printf("err localConn -- %v -- \n", err); return }
            HexID := fmt.Sprintf("%#x", buffer[0:2])
            if _, exist := packetsmap[HexID]; !exist {
                fmt.Printf("## !! [%v] len [%v] \t -> [%v]\n", HexID, len(buffer[:n]), buffer[:n])
                // fmt.Printf(" !! [%v] len [%v] \t -> [%v]\n", HexID, len(buffer[:n]), string(buffer[:n]))
            }else{
                function := reflect.ValueOf(fctpack[packetsmap[HexID].Ident])
                if function.Kind() == reflect.Func && fctpack[packetsmap[HexID].Ident] != nil{
                    args := []reflect.Value{reflect.ValueOf(buffer[0:2]),reflect.ValueOf(buffer[2:n])}
                    function.Call(args)
                }
            }
        }
    }()


    // ########################
    backend := imgui.CreateBackend(imgui.NewGLFWBackend())
    backend.SetAfterCreateContextHook(func () {
        for kk,_ := range lgatMaps {
            loadGatTexture(kk)
        }
    })

    backend.SetBeforeDestroyContextHook(func () {  })
    backend.SetBgColor(imgui.NewVec4(0.45, 0.55, 0.6, 1.0))
    backend.CreateWindow("ROBOTGO", 500, 500)

    backend.Run(func () {
        basePos := imgui.MainViewport().Pos()
        baseSize := imgui.MainViewport().Size()
        imgui.SetNextWindowPosV(imgui.NewVec2(basePos.X, basePos.Y), 0, imgui.NewVec2(0, 0))
        imgui.SetNextWindowSize(imgui.Vec2{X: baseSize.X, Y: baseSize.Y})
        imgui.Begin("robot")
        imgui.Text(fmt.Sprintf("Coords = X : %v / Y : %v", curCoord.X, curCoord.Y ))
        imgui.Text(fmt.Sprintf("Map : %v", curMap ))
        if curMap != "" {
        if _, exist := mapTextures[curMap]; exist {
            imgui.ImageV(mapTextures[curMap].ID(), imgui.NewVec2(float32(baseSize.X/1.5),float32(baseSize.Y/1.5)), imgui.NewVec2(0, 0), imgui.NewVec2(1, 1), imgui.NewVec4(1, 1, 1, 1), imgui.NewVec4(0, 0, 0, 0))
        }}
        imgui.End()
        imgui.Render()
    })

    <-exit
}


func sendToServer(hexID string,data []byte){
    var ii int16
	fmt.Sscanf(hexID, "0x%x", &ii)
    bb := []byte{ byte(ii >> 8), byte(ii) }
    bb = append(bb,data...)
    proxyCo.Write(bb)
}
