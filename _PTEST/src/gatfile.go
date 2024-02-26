package main

import(
    // "unsafe"
    // "fmt"
    "encoding/binary"
    // "github.com/cimgui-go"
    // "image"
    // "image/color"
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

// fld is a lighweight format of gat files
func loadFLD(mapName string){
	lgatMap := parseFLD([]byte(readFileString(CurDir()+"data/flds/"+mapName+".fld")))
	lgatMaps[mapName] = lgatMap
}


func parseFLD(data []byte) ROLGatMap{
    width := int(binary.LittleEndian.Uint16(data[0:2]))
    height := int(binary.LittleEndian.Uint16(data[2:4]))
    cell_list := make([][]uint8, width )
    for i := range cell_list {
        cell_list[i] = make([]uint8, height )
    }
    index := 4;
    for y := 0; y < int(height ); y++{
    for x := 0; x < int(width ); x++{
        cell_list[x][y] = uint8(byte(data[index]))
        index += 1
    }}
    return ROLGatMap{cells : cell_list, width : int(width), height : int(height)}
}
