package main

import (
	// "fmt"
	"math/rand"
)

type World struct {
	WorldMap [][]Location `json:"WorldMap"`
	Bots     []Bot        `json:"Bots"`
}

func (world *World) printWorld() {
	for x := range world.WorldMap {
		for y := range world.WorldMap[x] {
			if world.WorldMap[x][y].Food > 0 {
				world.WorldMap[x][y].print()
			}
		}
	}
	world.WorldMap[0][0].print()
	// for b := range world.Bots {
	// world.Bots[b].print()
	// }
}

func (world *World) execute(rounds int) bool {
	for b := range world.Bots {
		go world.Bots[b].run()
	}
	for r := 0; r < rounds; r++ {
		for b := range world.Bots {
			world.Bots[b].signal <- 1
		}
		for b := range world.Bots {
			<-world.Bots[b].wait
		}
		if world.WorldMap[0][0].Food >= 100 {
			for b := range world.Bots {
				world.Bots[b].quit <- 1
			}
			return true
		}
	}
	for b := range world.Bots {
		world.Bots[b].quit <- 1
	}
	return false
}

func (world *World) setup_locations(size_x, size_y int) {
	WorldMap := make([][]Location, size_x)
	world.WorldMap = WorldMap
	for x := range WorldMap {
		WorldMap[x] = make([]Location, size_y)
	}

	for x := range WorldMap {
		for y := range WorldMap[x] {
			var location = &WorldMap[x][y]
			location.setup(WorldMap, size_x, size_y, x, y)
		}
	}
}

func (world *World) setup_bots(botCount int) {
	bots := make([]Bot, botCount)
	world.Bots = bots
	for b := range bots {
		bots[b].currentLocation = &world.WorldMap[0][0]
		bots[b].direction = Direction{rand.Intn(8)}
		bots[b].signal = make(chan int)
		bots[b].wait = make(chan int)
		bots[b].quit = make(chan int)
	}
}
func (world *World) setup_environment() {
	for i := 0; i < 5; i++ {
		rand_x := rand.Intn(len(world.WorldMap))
		rand_y := rand.Intn(len(world.WorldMap))
		world.WorldMap[rand_x][rand_y].Food += 50

	}

	world.WorldMap[0][0].Nest = true
}

func NewWorld(size_x, size_y, botCount int) *World {
	world := new(World)
	world.setup_locations(size_x, size_y)
	world.setup_bots(botCount)

	world.setup_environment()

	return world
}
