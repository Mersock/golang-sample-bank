package main

import (
	"database/sql"
	"log"

	"github.com/Mersock/golang-sample-bank/api"
	db "github.com/Mersock/golang-sample-bank/db/sqlc"
	"github.com/Mersock/golang-sample-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("connat load config ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db ", err)
	}
	store := db.NewStore(conn)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("server down ", err)
	}
}
