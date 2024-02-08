package main

import (
	"slices"

	"github.com/tmw/pathfind"
)

var (
	InitialPositionStart  = Coordinate{X: 6, Y: 6}
	InitialPositionFinish = Coordinate{X: 18, Y: 12}
	InitialPositionWall   = Coordinate{X: 4, Y: 2}
)

type State struct {
	Grid       Grid       `json:"grid"`
	TileSize   int        `json:"tile_size"`
	CursorPos  Coordinate `json:"cursor_position"`
	CursorTile TileType   `json:"cursor_tile"`

	MouseDown bool
}

func NewState(tileSize int) State {
	return State{
		TileSize:   tileSize,
		CursorTile: TileTypeWall,
		CursorPos:  Coordinate{X: 12, Y: 5},
	}
}

// Takes pixel values and transforms to Coordinate on the grid.
func (s *State) OnMouseMove(x, y int) {
	s.CursorPos.X = min(x/s.TileSize, s.Grid.Width-1)
	s.CursorPos.Y = min(y/s.TileSize, s.Grid.Height-1)

	// placing wall
	if s.MouseDown &&
		s.CursorTile == TileTypeWall &&
		s.Grid.GetTileType(s.CursorPos) == TileTypeOpen {
		s.Grid.SetTileType(s.CursorPos, TileTypeWall)
	}

	// clearing wall
	if s.MouseDown &&
		s.CursorTile == TileTypeOpen &&
		s.Grid.GetTileType(s.CursorPos) == TileTypeWall {
		s.Grid.SetTileType(s.CursorPos, TileTypeOpen)
	}
}

func (s *State) OnMouseDown() {
	if s.MouseDown {
		return
	}

	s.MouseDown = true

	switch s.Grid.GetTileType(s.CursorPos) {
	case TileTypeStart:
		s.Grid.SetTileType(s.CursorPos, TileTypeOpen)
		s.CursorTile = TileTypeStart

	case TileTypeFinish:
		s.Grid.SetTileType(s.CursorPos, TileTypeOpen)
		s.CursorTile = TileTypeFinish

	case TileTypeOpen:
		s.CursorTile = TileTypeWall

	case TileTypeWall:
		s.CursorTile = TileTypeOpen
	}
}

func (s *State) OnMouseUp() {
	s.MouseDown = false

	switch s.CursorTile {
	case TileTypeStart:
		s.Grid.SetTileType(s.CursorPos, TileTypeStart)
		s.CursorTile = TileTypeSelector

	case TileTypeFinish:
		s.Grid.SetTileType(s.CursorPos, TileTypeFinish)
		s.CursorTile = TileTypeSelector

	case TileTypeOpen:
		s.Grid.SetTileType(s.CursorPos, TileTypeOpen)
		s.CursorTile = TileTypeSelector

	case TileTypeWall:
		s.Grid.SetTileType(s.CursorPos, TileTypeWall)
		s.CursorTile = TileTypeSelector
	}
}

type Grid struct {
	Width  int          `json:"width"`
	Height int          `json:"height"`
	Tiles  [][]TileType `json:"tiles"`

	start  Coordinate
	finish Coordinate
}

func (g *Grid) Solve() {
	pf := pathfind.NewSolver[Coordinate](pathfind.AlgorithmAStar, g.start, &pathfind.FuncAdapter[Coordinate]{
		NeighboursFn: g.Neighbours,
		CostToFinishFn: func(c Coordinate) int {
			return c.Dist(g.finish)
		},
		IsFinishFn: func(c Coordinate) bool {
			return c == g.finish
		},
	})

	steps := pf.Walk()

	for _, c := range steps {
		g.SetTileType(c, TileTypePath)
	}
}

func (g *Grid) Coordinates() []Coordinate {
	c := make([]Coordinate, g.Width*g.Height)
	for idx := range c {
		c[idx] = Coordinate{
			X: idx % g.Width,
			Y: idx / g.Width,
		}
	}

	return c
}

func (g *Grid) Clear() {
	for _, coord := range g.Coordinates() {
		if g.GetTileType(coord) == TileTypeStart ||
			g.GetTileType(coord) == TileTypeFinish {
			continue
		}

		g.SetTileType(coord, TileTypeOpen)
	}
}

func (g *Grid) SetTileType(coord Coordinate, typ TileType) {
	g.Tiles[coord.Y][coord.X] = typ

	if typ == TileTypeFinish {
		g.SetTileType(g.finish, TileTypeOpen)
		g.finish = coord
	}

	if typ == TileTypeStart {
		g.SetTileType(g.start, TileTypeOpen)
		g.start = coord
	}
}

func (g *Grid) GetTileType(coord Coordinate) TileType {
	return g.Tiles[coord.Y][coord.X]
}

func (g *Grid) Resize(w, h int) {
	tiles := make([][]TileType, h)
	for idx := range tiles {
		tiles[idx] = make([]TileType, w)
	}

	g.Width = w
	g.Height = h
	g.Tiles = tiles

	g.Clear()
	g.SetTileType(InitialPositionStart, TileTypeStart)
	g.SetTileType(InitialPositionFinish, TileTypeFinish)
	g.SetTileType(InitialPositionWall, TileTypeWall)
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
		if i.X < 0 || i.X >= g.Width || i.Y < 0 || i.Y >= g.Height {
			return true
		}

		// delete neighbours that are not walkable
		return g.Tiles[i.Y][i.X] == TileTypeWall

	}

	return slices.DeleteFunc(n, invalid)
}
