package main

import "log"
import "github.com/herenow/go-crate"
import "github.com/PuerkitoBio/gocrawl"
import "github.com/PuerkitoBio/goquery"
import "net/http"
import "io/ioutil"
import "time"
import "encoding/json"

// Database connection holder
var db, err = crate.Open("http://127.0.0.1:4200/")

// Extend the crawler
type Ext struct {
	*gocrawl.DefaultExtender
}

// Function executed once the crawler finishes crawling a page
func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	url := ctx.URL()
	url_str := url.String()

	log.Println("Crawled", url_str, "indexing...")

	buf, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		return nil, true
	}

	body := string(buf)

	json_str, err := db.Query("SELECT content, last_scan, title FROM web_index WHERE uri = ? LIMIT 1", url_str)

	if err != nil {
		log.Println(err)
		return nil, true
	}

	// TODO, USE THIS HACKY CODE FOR NOW
	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(json_str), &dat); err != nil {
        log.Println(err)
		return nil, true
    }

	// Check for sql errors
	if dat["error"] != nil {
		log.Println("Crate sql backend error")
		log.Println(dat["error"])
		return nil, true
	}

	// Page title
	title := doc.Find("title").Text()

	// Timestmap
	timestamp := int32(time.Now().Unix())

	// Version
	version := 0

	if dat["rowcount"].(float64) == 0.0 {
		_, err := db.Query("INSERT INTO web_index (uri, domain, title, content, first_scan, last_scan, version) VALUES ($1, $2, $3, $4, $5, $5, $6)",
					url_str, url.Host, title, body, timestamp, version)

		if err != nil {
			log.Println("Failed to insert new crawl to db")
			log.Println(err)
			return nil, true
		}

		log.Println("Sucessfull crawl insert")
	} else {
		version += 1
		_, err := db.Query("UPDATE web_index SET title = $1, content = $2, last_scan = $3, version = version + 1 WHERE uri = $4",
					title, body, timestamp, url_str)
		if err != nil {
			log.Println("Failed to update crawl version in db")
			log.Println(err)
			return nil, true
		}

		log.Println("Sucessefull crawl update")
	}

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
				opts.MaxVisits = 10000
				// Crate crawler
				c := gocrawl.NewCrawlerWithOptions(opts)
				c.Run(url)
			}()
		}
	}
}
