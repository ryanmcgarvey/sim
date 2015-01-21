package main

import (
	"encoding/json"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"sync"
	"time"
)

func api() {
	handler := rest.ResourceHandler{
		PreRoutingMiddlewares: []rest.Middleware{
			&rest.CorsMiddleware{
				RejectNonCorsRequests: false,
				OriginValidator: func(origin string, request *rest.Request) bool {
					return true
				},
				AllowedMethods: []string{"GET", "POST", "PUT"},
				AllowedHeaders: []string{
					"Accept", "Content-Type", "X-Custom-Header", "Origin"},
				AccessControlAllowCredentials: true,
				AccessControlMaxAge:           3600,
			},
		},
	}
	err := handler.SetRoutes(
		&rest.Route{"GET", "/world", func(w rest.ResponseWriter, req *rest.Request) {
			lock.Lock()
			start := time.Now()
			resp, err := json.Marshal(*world)
			if err != nil {
				fmt.Println("error:", err)
			}
			elapsed := time.Since(start)
			lock.Unlock()
			w.WriteJson(string(resp))
			fmt.Print(resp)
			fmt.Printf("JSON took %s", elapsed)
		}},
	)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/api/", http.StripPrefix("/api", &handler))

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("."))))

	log.Fatal(http.ListenAndServe(":8081", nil))
}

var (
	world *World
	lock  sync.Mutex
)

func main() {
	lock = sync.Mutex{}
	go api()
	for {
		wins := 0
		total := 20
		for i := 0; i < total; i++ {
			start := time.Now()
			world = NewWorld(500, 500, 100)
			success := world.execute(100000, lock)
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

// func main() {
// type ColorGroup struct {
// ID     int
// Name   string
// Colors []string
// }
// group := ColorGroup{
// ID:     1,
// Name:   "Reds",
// Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
// }
// b, err := json.Marshal(group)
// if err != nil {
// fmt.Println("error:", err)
// }
// os.Stdout.Write(b)
// }
