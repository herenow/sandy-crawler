// Crawler manager
// Wait for crawl requests
// Dispatch them
// And process them
package main

import "log"
import "net/http"
import "net"

// Crawled page struct
type PageCraweled struct {
	Url string
	ResponseHeaders string
	Body string
	Conn net.Conn
}


func CrawlReceive() (crawl PageCraweled)
