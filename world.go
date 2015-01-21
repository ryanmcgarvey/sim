package main

import (
	// "fmt"
	"math/rand"
)

type World struct {
	worldMap [][]Location
	bots     []Bot
}

func (world *World) printWorld() {
	for x := range world.worldMap {
		for y := range world.worldMap[x] {
			if world.worldMap[x][y].food > 0 {
				world.worldMap[x][y].print()
			}
		}
	}
	world.worldMap[0][0].print()
	// for b := range world.bots {
	// world.bots[b].print()
	// }
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

func (world *World) setup_locations(size_x, size_y int) {
	worldMap := make([][]Location, size_x)
	world.worldMap = worldMap
	for x := range worldMap {
		worldMap[x] = make([]Location, size_y)
	}

	for x := range worldMap {
		for y := range worldMap[x] {
			var location = &worldMap[x][y]
			location.setup(worldMap, size_x, size_y, x, y)
		}
	}
}

func (world *World) setup_bots(botCount int) {
	bots := make([]Bot, botCount)
	world.bots = bots
	for b := range bots {
		bots[b].currentLocation = &world.worldMap[0][0]
		bots[b].direction = Direction{rand.Intn(8)}
		bots[b].signal = make(chan int)
		bots[b].wait = make(chan int)
		bots[b].quit = make(chan int)
	}
}
func (world *World) setup_environment() {
	for i := 0; i < 3; i++ {
		rand_x := rand.Intn(len(world.worldMap))
		rand_y := rand.Intn(len(world.worldMap))
		world.worldMap[rand_x][rand_y].food += 50

	}

	world.worldMap[0][0].nest = true
}

func NewWorld(size_x, size_y, botCount int) *World {
	world := new(World)
	world.setup_locations(size_x, size_y)
	world.setup_bots(botCount)

	world.setup_environment()

	return world
}
