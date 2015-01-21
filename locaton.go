package main

import (
	"fmt"
	"math"
	"sync"
)

type Location struct {
	Signatures map[string]int
	Food       int
	X          int
	Y          int
	Nest       bool
	lock       sync.RWMutex `json:-`
	neighbors  [8]*Location `json:-`
}

func (loc *Location) print() {
	fmt.Println(loc.X, loc.Y, loc.Food, loc.Signatures, loc.Nest)
}

func (location *Location) neighbor_for(direction *Direction) *Location {
	return location.neighbors[direction.degree]
}

type Direction struct {
	degree int
}

func (location *Location) setup(WorldMap [][]Location, size_x, size_y, x, y int) {
	nx := 0
	ny := 0
	location.X = x
	location.Y = y
	location.lock = sync.RWMutex{}
	location.Signatures = map[string]int{"food": math.MaxInt64, "search": math.MaxInt64}
	var neighbors = &location.neighbors
	for n := range neighbors {

		switch n {
		case 0:
			nx = x
			ny = y + 1
		case 1:
			nx = x + 1
			ny = y + 1
		case 2:
			nx = x + 1
			ny = y
		case 3:
			nx = x + 1
			ny = y - 1
		case 4:
			nx = x
			ny = y - 1
		case 5:
			nx = x - 1
			ny = y - 1
		case 6:
			nx = x - 1
			ny = y
		case 7:
			nx = x - 1
			ny = y + 1
		}

		nx = nx % size_x
		ny = ny % size_y

		if nx < 0 {
			nx = size_x + nx
		}

		if ny < 0 {
			ny = size_y + ny
		}

		neighbors[n] = &WorldMap[nx][ny]
	}
}
