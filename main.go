package main

import (
	"context"
	"log"

	"github.com/ahmadfarhanstwn/backend-masterclass/api"
	db "github.com/ahmadfarhanstwn/backend-masterclass/db/sqlc"
	"github.com/ahmadfarhanstwn/backend-masterclass/util"
	"github.com/jackc/pgx/v5"
)

func main() {
	var err error

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	connection, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(connection)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create new server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
