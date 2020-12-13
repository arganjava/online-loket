package main

import (
	"database/sql"
	"fmt"
	"github.com/arganjava/online-loket/src/routers"
	_ "github.com/lib/pq"
	"os"
	"strconv"
)

func getEnv(value string) string {
	return os.Getenv(value)
}

func main() {
	port, err := strconv.Atoi(getEnv("TICKET_DB_PORT"))
	if err != nil {
		panic(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		getEnv("TICKET_DB_HOST"),
		port,
		getEnv("TICKET_DB_USERNAME"),
		getEnv("TICKET_DB_PASSWORD"),
		getEnv("TICKET_DB_NAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	routers.SetupServer(db).Run(getEnv("TICKET_APP_PORT"))

}
