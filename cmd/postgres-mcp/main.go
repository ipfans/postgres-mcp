package main

import (
	"flag"
	"fmt"
	"os"

	postgresmcp "github.com/ipfans/postgres-mcp"
	"github.com/joho/godotenv"
)

func main() {
	var dbString string
	flag.StringVar(&dbString, "db", "", "Database URL. (e.g. postgres://postgres:postgres@localhost:5432/postgres )")
	flag.Parse()

	if dbString == "" {
		godotenv.Load()
		dbString = os.Getenv("DATABASE_URL")
		if dbString == "" {
			fmt.Print("Please provide a database URL, DATABASE_URL environment variable or a dotenv file.\n\n")
			flag.PrintDefaults()
			return
		}
	}

	postgresmcp.Server(dbString)
}
