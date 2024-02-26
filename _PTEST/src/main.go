package main

import(
    "fmt"
    "net"
    "encoding/json"
    "github.com/cimgui-go"
    "io/ioutil"
    "strings"
    "time"
    "reflect"
)


type Coord struct {
	X,Y int
}

var (
    err error
    proxyCo net.Conn
    proxyCoClient net.Conn
    exit = make(chan bool)
    // gatMaps = map[string]ROGatMap{}
    lgatMaps = map[string]ROLGatMap{}
    packetsMap = map[string]int{}
    profil map[string][]interface{}
    conf = map[string][]interface{}{}
    mobDB []map[string]interface{}
)


func main() {

    fmt.Println("#--- ROBOTGO START ---#")
    fmt.Printf("current dir -- %v -- \n", CurDir())

    maps, _ := ioutil.ReadDir(CurDir()+"data/flds/")
    for _, m := range maps {
        if !m.IsDir() {
            name := strings.Split(m.Name(), ".fld")[0]
            // fmt.Printf("read -- %v -- \n", name)
            loadFLD(name)
        }
    }


    // ########################
    backend := imgui.CreateBackend(imgui.NewGLFWBackend())
    backend.SetAfterCreateContextHook(func () {  })

    backend.SetBeforeDestroyContextHook(func () {  })
    backend.SetBgColor(imgui.NewVec4(0.45, 0.55, 0.6, 1.0))
    backend.CreateWindow("ROBOTGO", 500, 800)

    targetFPS := 15
	frameTime := time.Second / time.Duration(targetFPS)

	lastFrameTime := time.Now()

    MAP := "gl_step"
    start := Coord{X:20,Y:100}
    end := Coord{X:20,Y:105}
    movePath := pathfind(start, end, lgatMaps[MAP])

    fmt.Printf("movePath -- %v -- \n", movePath)

    // return

    backend.Run(func () {
        currentTime := time.Now()
		elapsedTime := currentTime.Sub(lastFrameTime)
		if elapsedTime < frameTime {
            time.Sleep(frameTime - elapsedTime);
        }

        basePos := imgui.MainViewport().Pos()
        baseSize := imgui.MainViewport().Size()

        imgui.SetNextWindowPosV(imgui.NewVec2(basePos.X+1, basePos.Y + 500+1), 0, imgui.NewVec2(0, 0))
        imgui.SetNextWindowSize(imgui.Vec2{X: baseSize.X-2, Y: baseSize.Y - 500-2})
        imgui.Begin("Info")


        imgui.End()

        drawList := imgui.BackgroundDrawListNil()

        if _, exist := lgatMaps[MAP]; exist {
            lgatMap := lgatMaps[MAP]
            for x := 0; x < lgatMap.width; x++{
            for y := 0; y < lgatMap.height; y++{
                bbcolor := []byte{111,111,111,255}
                if lgatMap.cells[x][y] == 0 || lgatMap.cells[x][y] == 3 {
                    bbcolor[0] = 255; bbcolor[1] = 255; bbcolor[2] = 255;
                }
                if movePath != nil {
                    for _,vv := range movePath {
                        if vv.X == x && vv.Y == y{
                            bbcolor[0] = 50; bbcolor[1] = 100; bbcolor[2] = 150;
                        }
                    }
                }

                xpos := float32(x)
                ypos := float32(lgatMap.height - 1 - y)

                xpos *= 2
                ypos *= 2

                xpos += basePos.X
                ypos += basePos.Y

                drawList.AddRectFilled(imgui.Vec2{X: xpos, Y:  ypos}, imgui.Vec2{X: xpos+2,Y: ypos+2}, byteArrayToUInt32(bbcolor))

            }}
        }

        imgui.Render()
        lastFrameTime = currentTime
    })

    <-exit
}

type CRoute struct{ Map string; X int; Y int; WarpPortal string; }
type CMob struct{ Priority int; Id int; Name string;
                  AtkName string; AtkId int; AtkLv int; MinDist int; MinHP int;}
type CItemLoot struct{ Priority int; Id int; Name string; }
type CItemUse struct{ Id int; Name string; MinHP int; MinSP int; BuffId int; }
type CSKillSelf struct{ Id int; Lv int; Name string; MinHP int; MinSP int; BuffId int; }

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
        stru := CMob{ Priority:-1, AtkLv:1, MinDist:4 }
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
