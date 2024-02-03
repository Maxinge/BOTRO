package main

import(
    "github.com/cimgui-go"
    "image"
	"image/color"
    // "github.com/go-gl/gl/v2.1/gl"
    // "github.com/go-gl/glfw/v3.3/glfw"

)


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
        img.SetRGBA(x , lgatMap.height - 1 - y , c) // flip y cause coord system
    }}
    mapTextures[mapName] = imgui.NewTextureFromRgba(img)
}
