package main

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"

	"syscall/js"

	"github.com/tmw/pathfind"
)

func main() {
	// scope all functions under window.bridge
	bridge := map[string]interface{}{
		"solve": makeSolveFunc(),
	}

	js.Global().Set("bridge", js.ValueOf(bridge))

	// ensure Go doesn't exit before we have a chance to call it from JS.
	<-make(chan struct{})
}

type TileType int

const (
	TileTypeSelector = iota
	TileTypeOpen
	TileTypeWall
	TileTypeStart
	TileTypeFinish
	TileTypeCandidate
	TileTypeVisited
	TileTypePath
)

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c Coordinate) Dist(d Coordinate) int {
	return int(math.Floor(float64(c.X - d.X + c.Y - d.Y)))
}

type SolveRequest struct {
	GridWidth  int          `json:"gridWidth"`
	GridHeight int          `json:"gridHeight"`
	Grid       [][]TileType `json:"grid"`
}

type Grid struct {
	width  int
	height int
	tiles  [][]TileType
	start  Coordinate
	finish Coordinate
}

func (g *Grid) Solve() []Coordinate {
	pf := pathfind.NewSolver[Coordinate](pathfind.AlgorithmAStar, g.start, &pathfind.FuncAdapter[Coordinate]{
		NeighboursFn: g.Neighbours,
		CostToFinishFn: func(c Coordinate) int {
			return c.Dist(g.finish)
		},
		IsFinishFn: func(c Coordinate) bool {
			return c == g.finish
		},
	})

	return pf.Walk()
}

func (g *Grid) Neighbours(c Coordinate) []Coordinate {
	n := []Coordinate{
		{X: c.X - 1, Y: c.Y},
		{X: c.X + 1, Y: c.Y},
		{X: c.X, Y: c.Y - 1},
		{X: c.X, Y: c.Y + 1},
	}

	invalid := func(i Coordinate) bool {
		// delete neighbours that fall outside of the map
		if i.X < 0 || i.X >= g.width || i.Y < 0 || i.Y >= g.height {
			return true
		}

		// delete neighbours that are not walkable
		return g.tiles[i.Y][i.X] == TileTypeWall

	}

	return slices.DeleteFunc(n, invalid)
}

func parseGrid(req SolveRequest) Grid {
	g := Grid{
		tiles:  req.Grid,
		width:  req.GridWidth,
		height: req.GridHeight,
	}

	// find start and finish
	for y := range req.Grid {
		for x := range req.Grid[y] {
			if req.Grid[y][x] == TileTypeStart {
				g.start = Coordinate{X: x, Y: y}
			}
			if req.Grid[y][x] == TileTypeFinish {
				g.finish = Coordinate{X: x, Y: y}
			}
		}
	}

	return g
}

func makeSolveFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		var req SolveRequest
		json.Unmarshal([]byte(args[0].String()), &req)

		grid := parseGrid(req)
		path := grid.Solve()

		response, err := json.Marshal(path)
		if err != nil {
			fmt.Println("error marshalling json.", err)
		}

		return string(response)
	})
}
