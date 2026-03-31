package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/surojuvenkatesh20/bank-mgmt/utils"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfigFile("../../")
	if err != nil {
		log.Fatalln("cannot load config file: ", err)
	}
	testDB, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatalln("Error in connecting to db: ", err)
	}

	testQueries = New(testDB)
	code := m.Run()
	fmt.Println(code)

	os.Exit(code)
}
