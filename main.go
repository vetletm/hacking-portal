package main

import (
	"log"
	"net/http"
	"os"

	"hacking-portal/db"
	"hacking-portal/openstack"
	"hacking-portal/routes"

	"github.com/go-chi/chi"
)

func main() {
	// initialize the database connection
	db.Init(
		os.Getenv("DB_URL"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"))

	// initialize session routing
	routes.InitAuthentication(
		os.Getenv("LDAP_ADDR"),
		os.Getenv("LDAP_DC"),
		os.Getenv("COURSE_CODE"),
		os.Getenv("ADMINS"))

	// initialize openstack connection
	openstack.Init()

	// set up routing
	r := chi.NewRouter()
	r.Use(routes.SessionHandler)
	r.Get("/login", routes.GetLogin)
	r.Post("/login", routes.PostLogin)
	r.Get("/logout", routes.GetLogout)

	// let the remaining sub-routes handle themselves
	r.Mount("/groups", routes.GroupsRouter())
	r.Mount("/group", routes.GroupRouter())
	r.Mount("/admin", routes.AdminRouter())

	// serve static files directly
	fs := http.FileServer(http.Dir("/var/www/static/"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// set default route
	r.Get("/*", routes.RedirectLogin)

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
