package main

import(
    "fmt"
    "net"
    "encoding/json"
    "reflect"
    "github.com/cimgui-go"
    // "github.com/go-gl/gl/v2.1/gl"
    // "github.com/go-gl/glfw/v3.3/glfw"
    "io/ioutil"
    "strings"
    "time"
    "sync"
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
    Coords Coord
    Amount int
}

var (
    err error
    proxyCo net.Conn
    exit = make(chan bool)
    mu sync.Mutex

    fctpack map[string]func([]byte, []byte)
    // gatMaps = map[string]ROGatMap{}
    lgatMaps = map[string]ROLGatMap{}
    packetsmap map[string]Packet

    profil map[string]interface{}
    route map[string][]int
    targetMobs []int

    mobList = map[int]Mob{}
    groundItems = map[int]Item{}
    strMobs = ""
    strGroundItems = ""
)

func main() {

    fmt.Println("#--- ROBOTGO START ---#")
    fmt.Printf("current dir -- %v -- \n", CurDir())

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
            fmt.Printf("name -- %v -- \n", name)
            loadLGatMap(name)
        }
    }


    // ########################
    go func() {
        buffer := make([]byte, 100000)
        for {
            n, _ := proxyCo.Read(buffer)
            // if err != nil { fmt.Printf("err localConn -- %v -- \n", err); return }
            HexID := fmt.Sprintf("%#x", buffer[0:2])
            if len(buffer) < 3 { continue }
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

    go botLoop()
    go infoUILoop()

    // ########################
    backend := imgui.CreateBackend(imgui.NewGLFWBackend())
    backend.SetAfterCreateContextHook(func () {  })

    backend.SetBeforeDestroyContextHook(func () {  })
    backend.SetBgColor(imgui.NewVec4(0.45, 0.55, 0.6, 1.0))
    backend.CreateWindow("ROBOTGO", 500, 800)

    // imgui.CreateWindow("ROBOTGO", 500, 500)

    targetFPS := 10
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

        imgui.SetNextWindowPosV(imgui.NewVec2(basePos.X, basePos.Y + 400), 0, imgui.NewVec2(0, 0))
        imgui.SetNextWindowSize(imgui.Vec2{X: baseSize.X, Y: baseSize.Y - 400})
        imgui.Begin("Info")

        imgui.Text(fmt.Sprintf("Coords = X : %v / Y : %v", curCoord.X, curCoord.Y ))
        imgui.Text(fmt.Sprintf("Map : %v - Next point : %v", curMap, nextPoint))
        imgui.Text(fmt.Sprintf("\n states --- \n%v", printStruct(botStates) ))
        imgui.Text(fmt.Sprintf("\n Mobs --- \n%v", strMobs ))
        imgui.Text(fmt.Sprintf("\n groundItems --- \n%v", strGroundItems ))
        imgui.End()

        // drawList := imgui.WindowDrawList()
        drawList := imgui.BackgroundDrawListNil()

        // imgui.PushStyleColor(imgui.StyleColorBorder, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1}) // Couleur rouge
    	// imgui.Button("Mon Carré Rouge")
    	// imgui.PopStyleColor()

        scale := float32(3)
        sightDist := float32(66)
        if _, exist := lgatMaps[curMap]; exist {
            lgatMap := lgatMaps[curMap]
            for x := 0; x < lgatMap.width; x++{
            for y := 0; y < lgatMap.height; y++{

                if getDist(curCoord,(Coord{X:x,Y:y})) > float64(sightDist) { continue }
                size := float32(1)
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

                for _,vv := range targetMobPath {
                    if vv.X == x && vv.Y == y{
                        bbcolor[0] = 22; bbcolor[1] = 180; bbcolor[2] = 17;
                    }
                }

                if curCoord.X == x && curCoord.Y == y{
                    bbcolor[0] = 150; bbcolor[1] = 100; bbcolor[2] = 50;
                    xpos -= float32(2)
                    ypos -= float32(2)
                    size = float32(5)
                }

                mu.Lock()
                for _,vv := range mobList {
                    if vv.Coords.X == x && vv.Coords.Y == y{
                        bbcolor[0] = 33; bbcolor[1] = 200; bbcolor[2] = 220;
                        xpos -= float32(2)
                        ypos -= float32(2)
                        size = float32(5)
                    }
                }
                mu.Unlock()
                drawList.AddRectFilled(imgui.Vec2{X: xpos, Y:  ypos}, imgui.Vec2{X: xpos+(size*scale),Y: ypos+(size*scale)}, byteArrayToUInt32(bbcolor))
            }}
        }
        imgui.Render()
        lastFrameTime = currentTime
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
