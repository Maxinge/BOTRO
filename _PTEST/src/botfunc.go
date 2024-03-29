package main

import(
	"math"
	"math/rand"
	"time"
	// "fmt"
)


func getDist(from Coord, to Coord) float64 {
	return math.Sqrt(math.Pow(float64(to.X-from.X), 2) + math.Pow(float64(to.Y-from.Y), 2))
}

func isValidGatCell(cell uint8)  bool{
	if cell == 0 || cell == 3 { return true } ; return false
}

func isIn(coord Coord, list []Coord) bool{
	for _, v := range list { if coord == v{ return true } } ; return false
}

func isInPos(coord Coord, list []Coord) int{
	for k, v := range list { if coord == v { return k} } ; return -1
}

func generatePoints(dist int) []Coord {
	var points []Coord
	for x := -dist; x <= dist; x++ {
		for y := -dist; y <= dist; y++ {
			if (int( getDist(Coord{X: x, Y: y},Coord{X: 0, Y: 0}) )) == dist  {
				points = append(points, Coord{X: x, Y: y})
			}
		}
	}
	return points
}


func firstCircleVectors() []Coord {
	points := []Coord{  {-1, 1},  {0, 1},  {1, 1},
					    {-1, 0},   /*X*/   {1, 0},
				   	    {-1, -1}, {0, -1}, {1, -1} }
	return points
}

func firstCircle(point Coord) []Coord {
	points := firstCircleVectors()
    list := []Coord{}
	for _,v := range points {
		list = append(list,Coord{point.X + v.X, point.Y + v.Y})
	}
	return list
}

func secondCircleVectors() []Coord {
	points := []Coord{ {-2, 2}, {-1, 2}, {0, 2}, {1, 2}, {2, 2},
					   {-2, 1},							 {2, 1},
				   	   {-2, 0}, 		 /*X*/			 {2, 0},
				       {-2, -1},                         {2, -1},
				       {-2, -2},{-1, -2},{0, -2},{1, -2},{2, -2}  }
	return points
}

func secondCircle(point Coord) []Coord {
	points := secondCircleVectors()
    list := []Coord{}
	for _,v := range points {
		list = append(list,Coord{point.X + v.X, point.Y + v.Y})
	}
	return list
}

func directionTo(from Coord, to Coord) Coord {
	dx := float64(to.X) - float64(from.X)
	dy := float64(to.Y) - float64(from.Y)
	xtot := math.Sqrt(dx*dx)
	ytot := math.Sqrt(dy*dy)
	vx := 0
	vy := 0
	if xtot != 0 { vx = int( xtot/dx) }
	if ytot != 0 { vy = int( ytot/dy) }
	return Coord{X: vx, Y: vy}
}

func linearInterpolation(from Coord, to Coord) []Coord {
	var path []Coord
	deltaX := to.X - from.X
	deltaY := to.Y - from.Y
	dist := int(getDist(from,to))
	for i := 1; i <= dist; i++ {
		t := float64(i) / float64(dist)
		x := float64(from.X) + t*float64(deltaX)
		y := float64(from.Y) + t*float64(deltaY)
		path = append(path, Coord{X: int(x), Y: int(y)})
	}
	return path
}


func randomPoint(lgatMap ROLGatMap, from Coord, dist int) Coord{
	for {
		rand.Seed(time.Now().UnixNano())
		rX := rand.Intn(lgatMap.width)
		// rand.Seed(time.Now().UnixNano())
		rY := rand.Intn(lgatMap.height)
		gatCell := lgatMap.cells[rX][rY]
		if isValidGatCell(gatCell) {
		if getDist(from,Coord{X:rX, Y:rY}) < float64(dist) {
		if (Coord{X:rX, Y:rY}) != from {
			return Coord{X:rX, Y:rY}
		}}}
	}
}

func cleanPath(coordList []Coord, sighDist int, lgatMap ROLGatMap, banned []Coord) []Coord{
	if len(coordList) < 2  { return coordList }
	k := 0

	for {
		k += 1 ; if k > len(coordList)-1 { break }
		curCoord := coordList[k]

		best := 0
		for _,vect := range secondCircleVectors() {
			beamLine := Coord{ X:curCoord.X + (vect.X * sighDist), Y: curCoord.Y + (vect.Y * sighDist)}
			for _, beamCell := range linearInterpolation(curCoord, beamLine) {
				if !isValidCell(beamCell,lgatMap) { break }
				if isIn(beamCell, banned) { break }
				pos := isInPos(beamCell, coordList)
				if pos > -1 && pos > best && pos > k{
					best = pos
				}
			}
		}


		if best > 0 {
			temp := []Coord{}
			temp = append(temp,coordList[0:k+1]...)
			line := linearInterpolation(coordList[k], coordList[best])
			temp = append(temp,line[:]...)
			temp = append(temp,coordList[best:len(coordList)]...)
			coordList = temp
			k++
		}
	}

	return coordList
}


func isValidLine(start Coord, dest Coord, lgatMap ROLGatMap) bool{
    line := linearInterpolation(start, dest)
    for _,vv := range line {
    	gatcell := lgatMap.cells[vv.X][vv.Y]
    	if !isValidGatCell(gatcell) { return false }
    }
    return true
}

func isValidCell(cc Coord, lgatMap ROLGatMap) bool{
	if cc.X > lgatMap.width -1 || cc.Y > lgatMap.height -1 { return false }
	if cc.X < 0 || cc.Y < 0 { return false }
	gatCell := lgatMap.cells[cc.X][cc.Y]
	if isValidGatCell(gatCell) { return true }
	return false
}



func walkback(ccc Coord,paths *map[Coord]Coord,result *[]Coord){
	if cc, exist := (*paths)[ccc]; exist {
		*result = append(*result, cc)
		if cc != ccc { walkback(cc,paths,result) }
	}
}


func pathfind(start Coord, finish Coord, lgatMap ROLGatMap, banned []Coord) []Coord {

	if !isValidCell(finish, lgatMap) { return []Coord{start} }
	if start == finish{ return []Coord{start} }


	if isIn(finish, banned) {
		ccc := []Coord{}
		ccc = append(ccc, firstCircle(finish)...)
		ccc = append(ccc, secondCircle(finish)...)
		for _,vv := range ccc{
			if isValidCell(vv, lgatMap) && !isIn(vv, banned) {
				finish = vv; break;
			}
		}
	}

	// rand.Seed(time.Now().UnixNano())
	// rX := rand.Intn(lgatMap.width)
	// // rand.Seed(time.Now().UnixNano())
	// rY := rand.Intn(lgatMap.height)
	// gatCell := lgatMap.cells[rX][rY]

	paths := map[Coord]Coord{}
	heads := []Coord{}

	paths[start] = start
	heads = append(heads, start)

	PFstartTime := time.Now()
	PFelapsed := time.Now()
	found:
	for {
		if PFelapsed.Sub(PFstartTime).Milliseconds() > 500 { return []Coord{start} }
		banned := []Coord{}
		for _,vv := range heads {
			count := 0
			for _,vvvv := range firstCircle(vv) {
				if isValidCell(vvvv,lgatMap) {
				if !isIn(finish, banned) {
				if _, ok := paths[vvvv]; !ok {
					paths[vvvv] = vv
					heads = append(heads, vvvv)
					count ++
					if finish == vvvv{ break found }
				}}}
			}
			if count == 0 { banned = append(banned,vv) }
		}
		for _,vv := range banned {
			if isIn(vv,heads) {
				for kk,vvvv := range heads {
					if vvvv == vv {
						heads = append(heads[:kk], heads[kk+1:]...)
					}
				}
			}
		}
		if len(heads) == 0 { break }
		PFelapsed = time.Now()
	}

	result := []Coord{}
	if _, exist := paths[finish]; exist {
		walkback(finish,&paths,&result)
	}

	length := len(result)
	for i := 0; i < length/2; i++ {
		result[i], result[length-i-1] = result[length-i-1], result[i]
	}
	result = append(result, finish)


	cleanp := result
	cleanp = cleanPath(cleanp, 100, lgatMap, banned)
	cleanp = cleanp[1:]

	return cleanp

}
