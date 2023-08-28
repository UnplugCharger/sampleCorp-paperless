package db

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/qwetu_petro/backend/utils"
)

//  Since the tests need a db connections we will create one in the TestMain function.

var testQueries *Queries

func TestMain(m *testing.M) {
	//  Connect to the database and create a test database.
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)

	}

	conn, err := pgx.Connect(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	userID := 1 // Set this to the user ID you want to use
	_, err = conn.Exec(context.Background(), "SET audit.current_user_id = "+strconv.Itoa(userID))
	if err != nil {
		log.Fatal("Error setting configuration parameter:", err)
	}

	testQueries = New(conn)

	//  Run the tests.
	os.Exit(m.Run())

}
