package main

import (
	_ "github.com/mmcdole/gofeed"
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

	urls := getRssUrls(bufio.NewReader(rssReader))
