package main

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
