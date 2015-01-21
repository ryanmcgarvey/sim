package main

import (
	"fmt"
	"math"
	"sync"
)

type Location struct {
	signatures map[string]int
	neighbors  [8]*Location
	food       int
	x          int
	y          int
	lock       sync.RWMutex
	nest       bool
}

func (loc *Location) print() {
	fmt.Println(loc.x, loc.y, loc.food, loc.signatures, loc.nest)
}

func (location *Location) neighbor_for(direction *Direction) *Location {
	return location.neighbors[direction.degree]
}

type Direction struct {
	degree int
}

func (location *Location) setup(worldMap [][]Location, size_x, size_y, x, y int) {
	nx := 0
	ny := 0
	location.x = x
	location.y = y
	location.lock = sync.RWMutex{}
	location.signatures = map[string]int{"food": math.MaxInt64, "search": math.MaxInt64}
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

		neighbors[n] = &worldMap[nx][ny]
	}
}
