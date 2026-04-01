package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/surojuvenkatesh20/bank-mgmt/api"
	db "github.com/surojuvenkatesh20/bank-mgmt/db/sqlc"
	"github.com/surojuvenkatesh20/bank-mgmt/utils"
)

func main() {
	config, err := utils.LoadConfigFile(".")
	if err != nil {
		log.Fatalln("cannot load config file: ", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatalln("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}

	log.Println("server listing to: ", config.ServerAddress)
	err = server.StartServer(config.ServerAddress)

	if err != nil {
		log.Fatalln("cannot start server: ", err)
	}

}
