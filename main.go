package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DbUrl string `json:"dburl"`
}

func run() error {
	confFile, err := os.Open("config.json")
	if err != nil {
		return err
	}

	var config Config
	err = json.NewDecoder(confFile).Decode(&config)
	if err != nil {
		return err
	}
	// Create the database handle, confirm driver is present
	db, _ := sql.Open("mysql", config.DbUrl)
	defer db.Close()
	return nil

}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
	}
}
