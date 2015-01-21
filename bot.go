package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type Bot struct {
	currentLocation *Location
	food            int
	direction       Direction
	steps           int
	signal          chan int
	wait            chan int
	quit            chan int
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(10)
}

func (bot *Bot) run() {
	for {
		select {
		case <-bot.signal:
			bot.step()
			bot.wait <- 1
		case <-bot.quit:
			return
		}
	}
}

func (bot *Bot) step() {
	bot.leaveScent()
	bot.findFood()
	bot.leaveScent()
	bot.move()
}

func (bot *Bot) findFood() {
	if bot.food > 0 {
		if bot.currentLocation.Nest {
			bot.currentLocation.lock.Lock()
			bot.currentLocation.Food += bot.food
			bot.currentLocation.lock.Unlock()
			// fmt.Printf("Food NEST in %d steps\n", bot.steps)
			bot.food = 0
			bot.steps = 0
			bot.direction.degree = (4 + bot.direction.degree) % 8
		}
		return
	}

	if !bot.currentLocation.Nest {
		bot.currentLocation.lock.RLock()
		loc_food := bot.currentLocation.Food
		bot.currentLocation.lock.RUnlock()

		if loc_food > 0 {
			bot.currentLocation.lock.Lock()
			bot.currentLocation.Food--
			bot.currentLocation.lock.Unlock()
			// fmt.Printf("Found FOOD in %d steps\n", bot.steps)
			bot.food++
			bot.steps = 0
			bot.direction.degree = (4 + bot.direction.degree) % 8
		}
	}
}

func (bot *Bot) possible_directions() []int {
	dir := bot.direction.degree

	dirs := []int{dir - 2, dir - 1, dir, dir + 1, dir + 2}
	for d := range dirs {
		if dirs[d] > 7 {
			dirs[d] = dirs[d] % 8
		}
		if dirs[d] < 0 {
			dirs[d] = 8 + dirs[d]
		}
	}

	return dirs
}

func (bot *Bot) possible_locations_to_move_to() []*Location {
	directions := bot.possible_directions()
	locs := make([]*Location, len(directions))
	neighbors := &bot.currentLocation.neighbors
	for i := range directions {
		locs[i] = neighbors[directions[i]]
	}
	return locs
}

func (bot *Bot) min_signal(signal string) int {
	moves := bot.possible_locations_to_move_to()
	directions := bot.possible_directions()

	lowest_directions := []int{0}
	min := moves[0]

	for l := 1; l < len(moves); l++ {
		min_val := min.Signatures[signal]
		pot_min_val := moves[l].Signatures[signal]
		if min_val > pot_min_val {
			min = moves[l]
			lowest_directions = []int{l}
		} else {
			if min_val == pot_min_val {
				lowest_directions = append(lowest_directions, l)
			}
		}
	}
	return directions[lowest_directions[rand.Intn(len(lowest_directions))]]
}

func (bot *Bot) move() {
	var direction int

	if bot.food > 0 {
		direction = bot.min_signal("search")
	} else {
		direction = bot.min_signal("food")
	}
	bot.currentLocation = bot.currentLocation.neighbors[direction]
	bot.direction.degree = direction
	bot.steps += 1
}

func (bot *Bot) leaveScent() {
	var loc = bot.currentLocation
	if bot.food > 0 {
		if loc.Signatures["food"] > bot.steps {
			loc.Signatures["food"] = bot.steps
		}
		return
	}
	if loc.Signatures["search"] > bot.steps {
		loc.Signatures["search"] = bot.steps
	}
}

func (bot *Bot) print() {
	fmt.Println(bot.food, bot.direction)
	bot.currentLocation.print()
}
