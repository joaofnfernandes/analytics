package main

import (
	"github.com/joaofnfernandes/analytics/db"
	"github.com/joaofnfernandes/analytics/parser/votes"
)

func main() {
	db.CreateIfNotExists()
	votes.CsvToDb("analytics.db", "data/votes.csv")
}
