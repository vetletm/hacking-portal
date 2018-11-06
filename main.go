package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/vetletm/hacking-portal/db"
	"github.com/vetletm/hacking-portal/routes"
)

func main() {
	// initialize the database connection
	db.Init(os.Getenv("DB_URL"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"))

	// set up routing
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", routes.GetLogin)
		r.Mount("/student", routes.StudentRouter())
		r.Mount("/admin", routes.AdminRouter())
	})

	// attempt to get the port from the environment
	port := os.Getenv("PORT")
	if port == "" {
		// set a default port since it wasn't provided
		port = "8080"
	}

	// start webserver
	log.Printf("Serving on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
