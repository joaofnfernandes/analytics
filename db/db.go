package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const dbFilename = "analytics.db"

// CreateIfNotExists creates a DB with a single table named page
// TODO: return error instead of exiting
func CreateIfNotExists() {

	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
  create table page (
    url text not null primary key,
    pageviews integer,
    unique_pageviews integer,
    avg_time text,
    bouce_rate real,
    rating integer
  );
  `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
