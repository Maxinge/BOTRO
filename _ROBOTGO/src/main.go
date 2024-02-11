package main

import(
    "fmt"
    "net"
    "encoding/json"
    "github.com/cimgui-go"
    "io/ioutil"
    "strings"
    "time"
    "encoding/binary"
    "reflect"
)

type Packet struct {
    Ident  string
    Desc  string
    Size  int
}


var (
    err error
    proxyCo net.Conn
    exit = make(chan bool)
    fctpack map[string]func([]byte, []byte)
    // gatMaps = map[string]ROGatMap{}
    lgatMaps = map[string]ROLGatMap{}
    packetsMap = map[string]Packet{}
    profil map[string][]interface{}
    conf = map[string][]interface{}{}
)


func main() {

    fmt.Println("#--- ROBOTGO START ---#")
    fmt.Printf("current dir -- %v -- \n", CurDir())

    loadprofil()
    // fmt.Printf("conf -- %v -- \n", prettyPrint(conf)); return

    proxyCo, err = net.Dial("tcp", "127.0.0.1:6666")
    if err != nil { fmt.Printf("err -- %v -- \n", err); return }
    defer proxyCo.Close()

    fff := readFileString(CurDir()+"data/recvpackets.txt")
    ss := strings.Split(fff, "\n")
    for _,vv := range ss[:len(ss)-1] {
        sss := strings.Split(vv, " ")
        packetsMap[sss[0]] = Packet{Ident:"", Desc:"", Size:Stoi(sss[1][:len(sss[1])-1])}
    }

    var tempp map[string][]string
    err := json.Unmarshal([]byte(readFileString(CurDir()+"data/packets.json")), &tempp)
    if err != nil { fmt.Printf("err -- %v -- \n", err); return }

    for kk,vv := range tempp {
        tt := Packet{Ident:vv[0], Desc:vv[1], Size:0}
        if _, exist := packetsMap[kk]; exist {
            tt.Size = packetsMap[kk].Size
        }
        packetsMap[kk] = tt
    }

    fctpackInit()

    maps, _ := ioutil.ReadDir(CurDir()+"data/lgats/")
    for _, m := range maps {
        if !m.IsDir() {
            name := strings.Split(m.Name(), ".lgat")[0]
            // fmt.Printf("read -- %v -- \n", name)
            loadLGatMap(name)
        }
    }

    // ########################
    go func() {
        buffer := make([]byte, 100000)
        for {
            n ,_ := proxyCo.Read(buffer)
            // if err != nil { fmt.Printf("err localConn -- %v -- \n", err); return }
            if n < 2 { continue }
            ii := -1
            for {
                ii += 1; if ii >= n { break }
                HexID := fmt.Sprintf("%04X",binary.LittleEndian.Uint16(buffer[ii:ii+2]));
                plen := n-ii
                if _, exist := packetsMap[HexID]; exist {
                    plen = packetsMap[HexID].Size
                    if plen < 0 { plen = int(binary.LittleEndian.Uint16(buffer[ii+2:ii+4])) }
                    if plen <= 2 { plen = 2 }
                    args := []reflect.Value{ reflect.ValueOf(buffer[ii:ii+2]), reflect.ValueOf(buffer[ii+2:ii+2+plen-2]) }
                    parsePacket(packetsMap[HexID].Ident, args)
                    ii += plen-1;
                    continue
                }
                args := []reflect.Value{ reflect.ValueOf([]byte{255,255}), reflect.ValueOf(buffer[ii:ii+plen]) }
                parsePacket("uknw_pck", args)
                break;
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

        imgui.Text(fmt.Sprintf("Coords = X : %v / Y : %v | ID: %v", curCoord.X, curCoord.Y, accountId ))
        imgui.Text(fmt.Sprintf("Map : %v - Next point : %v (dist:%v)", curMap, nextPoint, getDist(curCoord, nextPoint)))

        imgui.Text(fmt.Sprintf("\n timeInState --- %v", timeInState ))
        imgui.Text(fmt.Sprintf("\n targetMob [%v] ---  targetItem[%v]\n", targetMob, targetItem ))
        imgui.Text(fmt.Sprintf("\n strInfo --- \n%v", strInfo ))

        imgui.Text(fmt.Sprintf("\n states --- \n%v", printStruct(botStates) ))
        imgui.Text(fmt.Sprintf("\n strBuffs --- \n%v", strBuffs ))
        imgui.Text(fmt.Sprintf("\n strMobs --- \n%v", strMobs ))
        imgui.Text(fmt.Sprintf("\n strInventoryItems --- \n%v", strInventoryItems ))
        imgui.Text(fmt.Sprintf("\n strGroundItems --- \n%v", strGroundItems ))

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

                if curPath != nil {
                    for _,vv := range curPath {
                        if vv.X == x && vv.Y == y{
                            bbcolor[0] = 50; bbcolor[1] = 100; bbcolor[2] = 150;
                        }
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


type CRoute struct{ Map string; X int; Y int; }
type CMob struct{ Priority int; Id int; Name string; }
type CItemLoot struct{ Priority int; Id int; Name string; }
type CItemUse struct{ Id int; Name string; MinHP int; MinSP int; BuffId int; }
type CSKillSelf struct{ Id int; Lv int; Name string; MinHP int; MinSP int; BuffId int; }
type CSkillTarget struct{ Id int; Lv int; Name string; MinDist int; MinHP int; }

func loadprofil(){
    err := json.Unmarshal([]byte(readFileString(CurDir()+"profils/_profil.json")), &profil)
    if err != nil { fmt.Printf("err json conf -- %v -- \n", err) }


    for _,vv := range profil["General"] {
        tt := vv.([]interface{})
        if (reflect.TypeOf(tt[1]).Kind()) == reflect.Float64{
            stru := struct{ Key string; Val int }{ Key: tt[0].(string), Val: int(tt[1].(float64)) }
            conf["General"] = append(conf["General"], stru)
        }
        if (reflect.TypeOf(tt[1]).Kind()) == reflect.String{
            stru := struct{ Key string; Val string }{ Key: tt[0].(string), Val: tt[1].(string) }
            conf["General"] = append(conf["General"], stru)
        }
    }
    for _,vv := range profil["Route"] {
        stru := CRoute{}
        for kkk,vvv := range vv.(map[string]interface{}) {
            fld := reflect.ValueOf(&stru).Elem().FieldByName(kkk); convertField(vvv, fld)
        }
        conf["Route"] = append(conf["Route"], stru)
    }
    for _,vv := range profil["Mob"] {
        stru := CMob{ Priority:1 }
        for kkk,vvv := range vv.(map[string]interface{}) {
            fld := reflect.ValueOf(&stru).Elem().FieldByName(kkk); convertField(vvv, fld)
        }
        conf["Mob"] = append(conf["Mob"], stru)
    }
    for _,vv := range profil["ItemLoot"] {
        stru := CItemLoot{ Priority:-1 }
        for kkk,vvv := range vv.(map[string]interface{}) {
            fld := reflect.ValueOf(&stru).Elem().FieldByName(kkk); convertField(vvv, fld)
        }
        conf["ItemLoot"] = append(conf["ItemLoot"], stru)
    }
    for _,vv := range profil["ItemUse"] {
        stru := CItemUse{ MinHP:-1, MinSP:-1, BuffId:-1 }
        for kkk,vvv := range vv.(map[string]interface{}) {
            fld := reflect.ValueOf(&stru).Elem().FieldByName(kkk); convertField(vvv, fld)
        }
        conf["ItemUse"] = append(conf["ItemUse"], stru)
    }
    for _,vv := range profil["SKillSelf"] {
        stru := CSKillSelf{ MinHP:-1, MinSP:-1, BuffId:-1 }
        for kkk,vvv := range vv.(map[string]interface{}) {
            fld := reflect.ValueOf(&stru).Elem().FieldByName(kkk); convertField(vvv, fld)
        }
        conf["SKillSelf"] = append(conf["SKillSelf"], stru)
    }
    for _,vv := range profil["SkillTarget"] {
        stru := CSkillTarget{ Id:-1, MinDist:3, Lv:1 }
        for kkk,vvv := range vv.(map[string]interface{}) {
            fld := reflect.ValueOf(&stru).Elem().FieldByName(kkk); convertField(vvv, fld)
        }
        conf["SkillTarget"] = append(conf["SkillTarget"], stru)
    }
}

// ### example ###
// if exist := getConf(conf["Route"],"Map","prontera"); exist != nil {
//     fmt.Printf("exist -- %v -- \n", exist.(CRoute).X)
// }
//
// if exist := getConf(conf["General"],"Key","useTP"); exist != nil {
//     TP := exist.(struct{Key string;Val int}).Val
//     fmt.Printf("TP -- %v -- \n", TP)
// }

func getConf(iiconf []interface{}, key2 string, iii interface{}) interface{} {
    for _,vv := range iiconf {

    for i := 0; i < reflect.TypeOf(vv).NumField(); i++ {
        kkk := reflect.TypeOf(vv).Field(i).Name
        vvv := reflect.ValueOf(vv).Field(i).Interface()
        if key2 == kkk{
        if (reflect.TypeOf(vvv).Kind()) == reflect.Int{
        if (reflect.TypeOf(iii).Kind()) == reflect.Int{
        if vvv.(int) == iii.(int) {
            return vv
        }}}}
        if key2 == kkk{
        if (reflect.TypeOf(vvv).Kind()) == reflect.Float64{
        if (reflect.TypeOf(iii).Kind()) == reflect.Float64{
        if vvv.(float64) == iii.(float64) {
            return vv
        }}}}
        if key2 == kkk{
        if (reflect.TypeOf(vvv).Kind()) == reflect.String{
        if (reflect.TypeOf(iii).Kind()) == reflect.String{
        if vvv.(string) == iii.(string) {
            return vv
        }}}}
    }}
    return nil
}

func convertField(ii interface{}, ff reflect.Value) {
	if ff.IsValid() {
		if (reflect.TypeOf(ii).Kind()) == reflect.Float64{ ff.SetInt(int64(ii.(float64))) }
		if (reflect.TypeOf(ii).Kind()) == reflect.Bool{ ff.SetBool(ii.(bool)) }
		if (reflect.TypeOf(ii).Kind()) == reflect.String{ ff.SetString(ii.(string)) }
	}
}
