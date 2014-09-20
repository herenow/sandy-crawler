// Database management related functions
// Prepare the database for input, etc
package main

import "log"
import "github.com/herenow/go-crate"

// Database schema
const (
	SCHEMA_WEB_INDEX = `Create Table web_index (
	domain string,
	uri string,
	title string,
	first_scan timestamp,
	last_scan timestamp,
	content string INDEX using fulltext,
	versions object(dynamic) as (
		version integer	
	)
)`
)

// Analyzer
// TODO

func PrepareDatabase(con crate.CrateConn) {
	res, err := con.Query(SCHEMA_WEB_INDEX)

	if err != nil {
		log.Println(err)
		log.Fatal("Failed creating schema, while preparing database")
	}

	log.Println("Database schema creating response:", res)
}

