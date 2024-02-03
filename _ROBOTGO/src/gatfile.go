package main

import(
    // "unsafe"
    // "fmt"
    "encoding/binary"
    "github.com/cimgui-go"
    "image"
    "image/color"
)

type ROGatCell struct {
    heights [4]float32    // 4 corners
    cell_type uint32
    //  0 - walkable block / 1 - non-walkable block
    //  2 - non-walkable water (not snipable) / 3 - walkable water
    //  4 - non-walkable water (snipable) / 5 - cliff (snipable)
    //  6 - cliff (not snipable)
}

type ROGatMap struct {
    cells [][]ROGatCell
    cells2 []byte
    width int
    height int
    // coord ysystem is difeerent from image coords so be carefull
    // x / y start from bottom left, so y reversed
}


type ROLGatMap struct {
    cells [][]uint8
    width int
    height int
}

// func loadGatMap(mapName string){
// 	gatMap := parseGat([]byte(readFileString(CurDir()+"data/gats/"+mapName+".gat")))
// 	gatMaps[mapName] = gatMap
// }

var(
    mapTextures = map[string]*imgui.Texture{}
    maskTexture *imgui.Texture
)

func loadLGatMap(mapName string){
	lgatMap := parseLGat([]byte(readFileString(CurDir()+"data/lgats/"+mapName+".lgat")))
	lgatMaps[mapName] = lgatMap
}

func parseLGat(data []byte) ROLGatMap{
    width := int(binary.BigEndian.Uint16(data[0:2]))
    height := int(binary.BigEndian.Uint16(data[2:4]))
    cell_list := make([][]uint8, width)
    for i := range cell_list {
        cell_list[i] = make([]uint8, height)
    }
    index := 4;
    for y := 0; y < int(height); y++{
    for x := 0; x < int(width); x++{
        cell_list[x][y] = uint8(byte(data[index]))
        index += 1
    }}
    return ROLGatMap{cells : cell_list, width : int(width), height : int(height)}
}

func loadGatTexture(mapName string)  {
    lgatMap := lgatMaps[mapName]
    img := image.NewRGBA(image.Rect(0, 0, lgatMap.width, lgatMap.height))
    for x := 0; x < lgatMap.width; x++{
    for y := 0; y < lgatMap.height; y++{
        R := 111; G := 111; B := 111;
        if lgatMap.cells[x][y] == 0 || lgatMap.cells[x][y] == 3 {
            R = 255; G = 255; B = 255;
        }
        c := color.RGBA{ R:uint8(R), G:uint8(G), B:uint8(B), A:255 }
        // flip y cause coord system
        img.SetRGBA(x , (lgatMap.height-1) - y , c)
    }}
    mapTextures[mapName] = imgui.NewTextureFromRgba(img)
}

// func parseGat(data []byte) ROGatMap {
//     // var format string = string(data[0:4])
//     // var vermajor byte = *(*byte)(unsafe.Pointer(&data[4]))
//     // var verminor byte = *(*byte)(unsafe.Pointer(&data[5]))
//     var width int32 =  *(*int32)(unsafe.Pointer(&data[6]))
//     var height int32 =  *(*int32)(unsafe.Pointer(&data[10]))
//     cell_list := make([][]ROGatCell, width)
//     for i := range cell_list {
//         cell_list[i] = make([]ROGatCell, height)
//     }
//     var index int32 = 14;
//     cells2 := []byte{}
//     for y := int32(0); y < height; y++{
//     for x := int32(0); x < width; x++{
//         cell_list[x][y] = *(*ROGatCell)(unsafe.Pointer(&data[index]))
//         cells2 = append(cells2, byte(data[index+16]))
//         index += int32(unsafe.Sizeof(ROGatCell{}))
//     }}
//     return ROGatMap{cells2: cells2, cells : cell_list, width : int(width), height : int(height)}
// }
