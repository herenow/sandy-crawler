package main

import "log"
import "github.com/herenow/go-crate"
import "github.com/PuerkitoBio/gocrawl"
import "github.com/PuerkitoBio/goquery"
import "net/http"
import "io/ioutil"

// Database connection holder
var db crate.CrateConn

// Extend the crawler
type Ext struct {
	*gocrawl.DefaultExtender
}

// Function executed once the crawler finishes crawling a page
func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	url := ctx.URL()

	log.Println("Crawled", url, "indexing...")

	buf, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		return nil, true
	}

	body := string(buf)

	log.Println(body)

	json_str, err := db.Query("SELECT content last_scan, title FROM web_index WHERE uri = ?", url)

	if err != nil {
		log.Println(err)
		return nil, true
	}

	log.Println(json_str)


	return nil, true
}

func main() {
	log.Println("Starting...")

	// Connect to database
	db, err := crate.Open("http://127.0.0.1:4200/")

	if err != nil {
		log.Println(err)
		log.Fatal("Failed to connect to DB")
	}

	// Prepare database
	PrepareDatabase(db)

	// Receive pages to crawl
	// Channel
	pageCrawl := make(chan string)

	// Wait for crawl requests, sent over our text based protocol
	// Good for debugging with a simple telnet
	go TextProtocolHandler(pageCrawl)

	// Wait for crawl requests
	for {
		select {
			// Crawl request
		case url := <-pageCrawl:
			go func() {
				// Prepare crawler
				ext := &Ext{&gocrawl.DefaultExtender{}}
				// Set custom options
				opts := gocrawl.NewOptions(ext)
				opts.CrawlDelay = 1
				opts.LogFlags = gocrawl.LogError
				opts.SameHostOnly = true
				// Crate crawler
				c := gocrawl.NewCrawlerWithOptions(opts)
				c.Run(url)
			}()
		}
	}
}

