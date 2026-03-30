package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	url        = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	driverName = "postgres"
)

var testQueries *Queries
var testDB *sql.DB
var err error

func TestMain(m *testing.M) {
	testDB, err = sql.Open(driverName, url)
	if err != nil {
		log.Fatalln("Error in connecting to db: ", err)
	}

	testQueries = New(testDB)
	code := m.Run()
	fmt.Println(code)

	os.Exit(code)

}
