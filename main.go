package main

// import "time"
import "fmt"

func main() {
	for j := 0; j < 100; j++ {
		wins := 0
		total := 20
		for i := 0; i < total; i++ {
			// start := time.Now()
			var world = NewWorld(100, 100, 10)
			success := world.execute(1000)
			// elapsed := time.Since(start)
			// world.printWorld()
			// fmt.Println(elapsed, success)
			if success {
				wins++
			}
		}
		fmt.Printf("%d wins of %d tries\n", wins, total)
	}
}
