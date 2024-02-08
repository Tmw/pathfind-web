package main

import (
	"math"
)

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c Coordinate) Dist(d Coordinate) int {
	return int(math.Floor(float64(c.X - d.X + c.Y - d.Y)))
}
