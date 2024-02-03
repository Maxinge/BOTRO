package main

import(
    "unsafe"
    // "fmt"
    "encoding/binary"
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

func loadGatMap(mapName string){
	gatMap := parseGat([]byte(readFileString(CurDir()+"data/gats/"+mapName+".gat")))
	gatMaps[mapName] = gatMap
}

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
    for x := 0; x < int(width); x++{
    for y := 0; y < int(height); y++{
        cell_list[x][y] = byte(data[index])
        index += 1
    }}
    return ROLGatMap{cells : cell_list, width : int(width), height : int(height)}
}

func parseGat(data []byte) ROGatMap {
    // var format string = string(data[0:4])
    // var vermajor byte = *(*byte)(unsafe.Pointer(&data[4]))
    // var verminor byte = *(*byte)(unsafe.Pointer(&data[5]))
    var width int32 =  *(*int32)(unsafe.Pointer(&data[6]))
    var height int32 =  *(*int32)(unsafe.Pointer(&data[10]))
    cell_list := make([][]ROGatCell, width)
    for i := range cell_list {
        cell_list[i] = make([]ROGatCell, height)
    }
    var index int32 = 14;
    for x := int32(0); x < width; x++{
    for y := int32(0); y < height; y++{
        cell_list[x][y] = *(*ROGatCell)(unsafe.Pointer(&data[index]))
        index += int32(unsafe.Sizeof(ROGatCell{}))
    }}
    return ROGatMap{cells : cell_list, width : int(width), height : int(height)}
}
