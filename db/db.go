package db

import (
	"log"

	"github.com/globalsign/mgo"
)

var db *mgo.Database

// Init initializes a new database session to a MongoDB instance at a given path.
// The database session is used throughout the package, so the database must be initialized for the rest of the package to work correctly.
func Init(path, name, user, pass string) {
	// bail out if any of the environment variables were missing
	if path == "" {
		log.Fatal("Missing env DB_URL")
	}
	if name == "" {
		log.Fatal("Missing env DB_NAME")
	}
	if user == "" {
		log.Fatal("Missing env DB_USER")
	}
	if pass == "" {
		log.Fatal("Missing env DB_PASS")
	}

	// connect to the mongo database with the given credentials
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{path},
		Database: name,
		Username: user,
		Password: pass,
	})

	if err != nil {
		// we failed to connect, bail out!
		log.Fatal("Failed to connect to database")
	} else {
		// log a successful connection
		log.Println("Database connected")
	}

	// store the database session so we're able to use it later
	db = session.DB(name)
}
