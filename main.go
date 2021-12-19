package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mmcdole/gofeed"
)

type Config struct {
	DbUrl   string `json:"dburl"`
	DbName  string `json:"dbname"`
	RssFile string `json:"rssfile"`
}

// GetRssUrls removes empty lines and lines beginning with #
func getRssUrls(r *bufio.Reader) (urls []string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}
	return urls
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

	urls := getRssUrls(bufio.NewReader(rssReader))

	fmt.Println(urls)

	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", config.DbUrl)
	if err != nil {
		return fmt.Errorf("Can not open connection to %s ERROR: %s", config.DbUrl, err)
	}
	defer db.Close()

	var one, two string
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		return err
	}
	for rows.Next() {
		rows.Scan(&one, &two)
	}
	if err != nil {
		return fmt.Errorf("Query failed ERROR: %s", err)
	}
	fmt.Println("Connected to:", one, two)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}
