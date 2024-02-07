package main

import(
    "fmt"
    "net"
    "encoding/json"
    "github.com/cimgui-go"
    "io/ioutil"
    "strings"
    "time"
    "sync"
    "encoding/binary"
    "reflect"
)

type Packet struct {
    Ident  string
    Desc  string
    Size  int
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
    Coords Coord
    Amount int
}

var (
    err error
    proxyCo net.Conn
    exit = make(chan bool)

    fctpack map[string]func([]byte, []byte)
    // gatMaps = map[string]ROGatMap{}
    lgatMaps = map[string]ROLGatMap{}
    packetsmap map[string]Packet

    profil map[string]interface{}
    route map[string][]int
    targetMobs []int

    MUmobList sync.Mutex
    mobList = map[int]Mob{}

    MUgroundItems sync.Mutex
    groundItems = map[int]Item{}

    strMobs = ""
    strGroundItems = ""
)


func main() {

    fmt.Println("#--- ROBOTGO START ---#")
    fmt.Printf("current dir -- %v -- \n", CurDir())

    // iiiii := fmt.Sprintf("%04X",binary.LittleEndian.Uint16([]byte{150,1}))
    // fmt.Printf("kek -- [%v] -- \n", iiiii); return

    proxyCo, err = net.Dial("tcp", "127.0.0.1:6666")
    if err != nil { fmt.Printf("err -- %v -- \n", err); return }
    defer proxyCo.Close()

    json.Unmarshal([]byte(readFileString(CurDir()+"data/packets.json")), &packetsmap)
    // fmt.Printf("packetsmap -- %v -- \n", packetsmap)
    loadprofil()
    fctpackInit()

    maps, _ := ioutil.ReadDir(CurDir()+"data/lgats/")
    for _, m := range maps {
        if !m.IsDir() {
            name := strings.Split(m.Name(), ".lgat")[0]
            fmt.Printf("read -- %v -- \n", name)
            loadLGatMap(name)
        }
    }

    // ########################
    go func() {
        buffer := make([]byte, 100000)
        for {
            n ,_ := proxyCo.Read(buffer)
            // if err != nil { fmt.Printf("err localConn -- %v -- \n", err); return }
            if n < 3 { continue }
            ii := -1
            for {
                ii += 1; if ii >= n { break }
                HexID := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(buffer[ii:ii+2]));
                plen := n - ii
                if _, exist := packetsmap[HexID]; exist {
                    plen = packetsmap[HexID].Size
                    if plen < 0 { plen = int(binary.LittleEndian.Uint16(buffer[ii+2:ii+4])) }
                    if plen < 2 { plen = 2 }
                    args := []reflect.Value{ reflect.ValueOf(buffer[ii:ii+2]), reflect.ValueOf(buffer[ii+2:ii+2+plen-2]) }
                    parsePacket(packetsmap[HexID].Ident, args)
                    ii += plen -1 ;
                    continue
                }
                args := []reflect.Value{ reflect.ValueOf([]byte{255,255}), reflect.ValueOf(buffer[ii:ii+plen]) }
                parsePacket("uknw_pck", args)
                break;
            }
        }
    }()

    // go botLoop()
    go infoUILoop()

    // ########################
    backend := imgui.CreateBackend(imgui.NewGLFWBackend())
    backend.SetAfterCreateContextHook(func () {  })

    backend.SetBeforeDestroyContextHook(func () {  })
    backend.SetBgColor(imgui.NewVec4(0.45, 0.55, 0.6, 1.0))
    backend.CreateWindow("ROBOTGO", 500, 800)

    targetFPS := 15
	frameTime := time.Second / time.Duration(targetFPS)

	lastFrameTime := time.Now()


    backend.Run(func () {
        currentTime := time.Now()
		elapsedTime := currentTime.Sub(lastFrameTime)
		if elapsedTime < frameTime {
            time.Sleep(frameTime - elapsedTime);
        }

        basePos := imgui.MainViewport().Pos()
        baseSize := imgui.MainViewport().Size()

        imgui.SetNextWindowPosV(imgui.NewVec2(basePos.X+1, basePos.Y + 400+1), 0, imgui.NewVec2(0, 0))
        imgui.SetNextWindowSize(imgui.Vec2{X: baseSize.X-2, Y: baseSize.Y - 400-2})
        imgui.Begin("Info")

        imgui.Text(fmt.Sprintf("Coords = X : %v / Y : %v", curCoord.X, curCoord.Y ))
        imgui.Text(fmt.Sprintf("Map : %v - Next point : %v", curMap, nextPoint))

        imgui.Text(fmt.Sprintf("\n timeInState --- %v", timeInState ))
        imgui.Text(fmt.Sprintf("\n targetMob [%v] ---  targetItem[%v]\n", targetMob, targetItem ))

        imgui.Text(fmt.Sprintf("\n states --- \n%v", printStruct(botStates) ))
        imgui.Text(fmt.Sprintf("\n Mobs --- \n%v", strMobs ))
        imgui.Text(fmt.Sprintf("\n groundItems --- \n%v", strGroundItems ))
        imgui.End()

        drawList := imgui.BackgroundDrawListNil()

        scale := float32(3)
        sightDist := float32(66)
        if _, exist := lgatMaps[curMap]; exist {
            lgatMap := lgatMaps[curMap]
            for x := 0; x < lgatMap.width; x++{
            for y := 0; y < lgatMap.height; y++{

                if getDist(curCoord,(Coord{X:x,Y:y})) > float64(sightDist) { continue }
                size := float32(1) * scale
                bbcolor := []byte{111,111,111,255}
                xpos := float32(x) - float32(curCoord.X)
                ypos := float32(lgatMap.height - 1 - y) - float32(curCoord.Y*-1)

                xpos = (xpos*scale) + basePos.X
                ypos = ((ypos - float32(lgatMap.height))*scale) + basePos.Y

                xpos += (sightDist * scale)
                ypos += (sightDist * scale)

                if lgatMap.cells[x][y] == 0 || lgatMap.cells[x][y] == 3 {
                    bbcolor[0] = 255; bbcolor[1] = 255; bbcolor[2] = 255;
                }

                for _,vv := range curPath {
                    if vv.X == x && vv.Y == y{
                        bbcolor[0] = 50; bbcolor[1] = 100; bbcolor[2] = 150;
                    }
                }

                if curCoord.X == x && curCoord.Y == y{
                    bbcolor[0] = 150; bbcolor[1] = 100; bbcolor[2] = 50;
                    size = float32(3) * scale
                }

                MUmobList.Lock()
                for _,vv := range mobList {
                    if vv.Coords.X == x && vv.Coords.Y == y{
                        bbcolor[0] = 33; bbcolor[1] = 200; bbcolor[2] = 220;
                        size = float32(3) * scale
                    }
                }
                MUmobList.Unlock()
                drawList.AddRectFilled(imgui.Vec2{X: xpos, Y:  ypos}, imgui.Vec2{X: xpos+size,Y: ypos+size}, byteArrayToUInt32(bbcolor))
            }}
        }
        imgui.Render()
        lastFrameTime = currentTime
    })

    <-exit
}


func sendToServer(hexID string,data []byte){
    var ii uint16
	fmt.Sscanf(hexID, "%x", &ii)
    bb := []byte{byte(ii), byte(ii >> 8)}
    bb = append(bb,data...)
    proxyCo.Write(bb)
}
