package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@localhost:5433/simplebank?sslmode=disable"
)

var testQueries *Queries
var connection *pgx.Conn

func TestMain(m *testing.M) {
	var err error
	connection, err = pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(connection)

	os.Exit(m.Run())
}
