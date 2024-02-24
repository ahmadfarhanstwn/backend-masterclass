package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ahmadfarhanstwn/backend-masterclass/util"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var connection *pgx.Conn

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connection, err = pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(connection)

	os.Exit(m.Run())
}
