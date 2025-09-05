package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/bobafetch/dashboard/data"
	"github.com/bobafetch/dashboard/routes"
	"github.com/bobafetch/dashboard/utils"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := data.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	mux := http.NewServeMux()

	routes.RegisterWorkcenterRoutes(mux)

	mux.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(config.ADDRESS, mux))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
