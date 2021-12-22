package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mmcdole/gofeed"
)

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

func getFromSvtPlay(rssReader io.Reader) error {
	urls := getRssUrls(bufio.NewReader(rssReader))
	for _, url := range urls {
		feed, err := gofeed.NewParser().ParseURL(url)
		if err != nil {
			return err
		}
		for _, item := range feed.Items {
			fmt.Println(item.GUID)
		}
	}
	return nil
}
