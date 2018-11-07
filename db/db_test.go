package db

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/globalsign/mgo/dbtest"
)

func TestMain(m *testing.M) {
	// temp directory to store test database
	tempDir, _ := ioutil.TempDir("", "testing")

	// start the test database server and get a session
	var server dbtest.DBServer
	server.SetPath(tempDir)
	session := server.Session()

	// set the database variable to the database in a session
	db = session.DB("testing")

	// run the tests
	ret := m.Run()

	// cleanup
	db.DropDatabase()
	session.Close()
	server.Stop()
	os.Exit(ret)
}
