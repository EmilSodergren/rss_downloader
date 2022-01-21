package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DbUrl      string `json:"dburl"`
	DbName     string `json:"dbname"`
	SvtRssFile string `json:"rssfile"`
}

func findConfigFile(path string) (string, error) {
	if filepath.IsAbs(path) {
		_, err := os.Stat(path)
		if err != nil {
			return "", err
		}
		return path, nil
	}
	// Look in curdir
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	// Look in /etc
	newPath := filepath.Join("/", "etc", "rss_downloader", path)
	if _, err := os.Stat(newPath); err != os.ErrNotExist {
		return newPath, nil
	}
	return "", fmt.Errorf("Could not find %s anywhere", path)
}

func run() error {
	configPath, err := findConfigFile("config.json")
	if err != nil {
		return err
	}
	confFile, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("Can not open config.json ERROR: %s", err)
	}

	var config Config
	err = json.NewDecoder(confFile).Decode(&config)
	if err != nil {
		return fmt.Errorf("Can not parse config.json ERROR: %s", err)
	}

	svtRssPath, err := findConfigFile(config.SvtRssFile)
	if err != nil {
		return err
	}
	rssReader, err := os.Open(svtRssPath)
	if err != nil {
		return fmt.Errorf("Can not open %s ERROR: %s", config.SvtRssFile, err)
	}

	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", config.DbUrl)
	if err != nil {
		return fmt.Errorf("Can not open connection to %s ERROR: %s", config.DbUrl, err)
	}
	defer db.Close()

	err = getFromSvtPlay(db, rssReader)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}
