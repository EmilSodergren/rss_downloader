package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

const dateFmt = "2006-01-02 15:04:05"

var createSvtplayTableStmt = `CREATE TABLE IF NOT EXISTS svtplay (
	guid text not null primary key,
	title text,
	link text,
	description text,
	pubDate datetime,
	downloadedDate datetime
)`

var checkIfExistsStmt = "SELECT * FROM svtplay WHERE guid=?"
var insertInSvtplayTableStmt = `INSERT INTO svtplay (guid,title,link,description,pubDate,downloadedDate)
	VALUES (?,?,?,?,?,?)`

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

func downloadURL(link string) error {
	cmd := exec.Command("yt-dlp", link)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func insertIntoTable(db *sql.DB, item *gofeed.Item) error {
	_, err := db.Exec(insertInSvtplayTableStmt, item.GUID, item.Title, item.Link, item.Description, item.PublishedParsed.Format(dateFmt), time.Now().Format(dateFmt))
	return err
}

func getFromSvtPlay(db *sql.DB, rssReader io.Reader) error {
	urls := getRssUrls(bufio.NewReader(rssReader))
	_, err := db.Exec(createSvtplayTableStmt)
	if err != nil {
		return fmt.Errorf("Create table failed ERROR: %s", err)
	}
	for _, url := range urls {
		feed, err := gofeed.NewParser().ParseURL(url)
		if err != nil {
			return err
		}
		for _, item := range feed.Items {
			row := db.QueryRow(checkIfExistsStmt, item.GUID)
			if row.Scan() == sql.ErrNoRows {
				fmt.Println("Downloading", item.Title)
				err := downloadURL(item.Link)
				if err != nil {
					return err
				}
				err = insertIntoTable(db, item)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
