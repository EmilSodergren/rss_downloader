package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DbUrl   string `json:"dburl"`
	DbName  string `json:"dbname"`
	RssFile string `json:"rssfile"`
}

func run() error {
	confFile, err := os.Open("config.json")
	if err != nil {
		return fmt.Errorf("Can not open config.json ERROR: %s", err)
	}

	var config Config
	err = json.NewDecoder(confFile).Decode(&config)
	if err != nil {
		return fmt.Errorf("Can not parse config.json ERROR: %s", err)
	}

	rssReader, err := os.Open(config.RssFile)
	if err != nil {
		return fmt.Errorf("Can not open %s ERROR: %s", config.RssFile, err)
	}

	err = getFromSvtPlay(rssReader)
	if err != nil {
		return err
	}

	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", config.DbUrl)
	if err != nil {
		return fmt.Errorf("Can not open connection to %s ERROR: %s", config.DbUrl, err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		return fmt.Errorf("Query failed ERROR: %s", err)
	}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println("Db has", id, name)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}
