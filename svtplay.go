package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

const dateFmt = "2006-01-02 15:04:05"

var createSvtplayTableStmt = `CREATE TABLE IF NOT EXISTS svtplay (
	id INTEGER not null primary key AUTO_INCREMENT,
	guid TEXT(512) not null,
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
	cmd := exec.Command("yt-dlp", "-S", "codec:h264", link)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func insertIntoTable(db *sql.DB, item *gofeed.Item) error {
	_, err := db.Exec(insertInSvtplayTableStmt, item.GUID, item.Title, item.Link, item.Description, item.PublishedParsed.Format(dateFmt), time.Now().Format(dateFmt))
	return err
}

func getEnclosure(e *gofeed.Enclosure, name string) error {
	response, err := http.Get(e.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	name = strings.Replace(name, ":", " -", -1)
	file, err := os.Create(fmt.Sprintf("%s%s", name, ".jpeg"))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
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
			// Verify that the episode does not exist in the database
			row := db.QueryRow(checkIfExistsStmt, item.GUID)
			if row.Scan() == sql.ErrNoRows {
				fmt.Println("Downloading", item.Title)
				err := downloadURL(item.Link)
				if err != nil {
					return err
				}
				for _, e := range item.Enclosures {
					err = getEnclosure(e, item.Title)
					if err != nil {
						return err
					}
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
