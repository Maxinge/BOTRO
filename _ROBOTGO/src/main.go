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
    backend.SetAfterCreateContextHook(func () {
        // for kk,_ := range lgatMaps { loadGatTexture(kk)  }
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
        imgui.Text(fmt.Sprintf("Mobs --- \n\n%v", strMobs ))
        imgui.Text(fmt.Sprintf("groundItems --- \n\n%v", strGroundItems ))

        
        // if curMap != "" {
        // if _, exist := mapTextures[curMap]; exist {
        //     imgui.ImageV(mapTextures[curMap].ID(), imgui.NewVec2(float32(baseSize.X/1.5),float32(baseSize.Y/1.5)), imgui.NewVec2(0, 0), imgui.NewVec2(1, 1), imgui.NewVec4(1, 1, 1, 1), imgui.NewVec4(0, 0, 0, 0))
        // }}
        // if maskTexture != nil {
        //     imgui.ImageV(maskTexture.ID(), imgui.NewVec2(float32(baseSize.X/2),float32(baseSize.Y/2)), imgui.NewVec2(0, 0), imgui.NewVec2(1, 1), imgui.NewVec4(1, 1, 1, 1), imgui.NewVec4(0, 0, 0, 0))
        // }
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
