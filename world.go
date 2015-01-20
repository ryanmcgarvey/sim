package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

type World struct {
	worldMap [][]Location
	bots     []Bot
}

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

type Direction struct {
	degree int
}

func (world *World) printWorld() {
	fmt.Println("State of the World")
	// for x := range world.worldMap {
	// for y := range world.worldMap[x] {
	// world.worldMap[x][y].print()
	// }
	// }
	world.worldMap[0][0].print()
	world.worldMap[5][5].print()
	for b := range world.bots {
		world.bots[b].print()
	}
}

func (world *World) execute(rounds int) bool {
	for b := range world.bots {
		go world.bots[b].run()
	}
	for r := 0; r < rounds; r++ {
		for b := range world.bots {
			world.bots[b].signal <- 1
		}
		for b := range world.bots {
			<-world.bots[b].wait
		}
		if world.worldMap[0][0].food >= 10 {
			for b := range world.bots {
				world.bots[b].quit <- 1
			}
			return true
		}
	}
	for b := range world.bots {
		world.bots[b].quit <- 1
	}
	return false
}

func NewWorld(size_x, size_y, botCount int) *World {
	world := new(World)

	worldMap := make([][]Location, size_x)

	for x := range worldMap {
		worldMap[x] = make([]Location, size_y)
	}

	var nx = 0
	var ny = 0
	for x := range worldMap {
		for y := range worldMap[x] {
			var location = &worldMap[x][y]
			location.x = x
			location.y = y
			location.nest = x == 0 && y == 0
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
	}

	bots := make([]Bot, botCount)
	for b := range bots {
		bots[b].currentLocation = &worldMap[0][0]
		bots[b].direction = Direction{rand.Intn(8)}
		bots[b].signal = make(chan int)
		bots[b].wait = make(chan int)
		bots[b].quit = make(chan int)
	}

	worldMap[5][5].food = 10

	world.worldMap = worldMap
	world.bots = bots
	return world
}
