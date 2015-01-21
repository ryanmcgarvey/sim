package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
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
			w.WriteJson(world)
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
)

func main() {
	go api()
	for {
		wins := 0
		total := 20
		for i := 0; i < total; i++ {
			start := time.Now()
			world = NewWorld(500, 500, 100)
			success := world.execute(100000)
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
