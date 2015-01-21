package main

import (
	"fmt"
	"time"
)

func main() {
	go run_server()
	for {
		wins := 0
		total := 20
		for i := 0; i < total; i++ {
			start := time.Now()
			var world = NewWorld(1000, 1000, 100)
			success := world.execute(10000)
			elapsed := time.Since(start)

			if success {
				wins++
				fmt.Println(elapsed, success, wins)
				world.printWorld()
				fmt.Printf("\n\n")
			}
		}
		fmt.Printf("%d wins of %d tries\n", wins, total)
	}
}
