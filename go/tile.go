package main

type TileType string

const (
	TileTypeSelector  = "selector"
	TileTypeOpen      = "open"
	TileTypeWall      = "wall"
	TileTypeStart     = "start"
	TileTypeFinish    = "finish"
	TileTypeCandidate = "candidate"
	TileTypeVisited   = "visited"
	TileTypePath      = "path"
)
