package main

import (
	"github.com/joaofnfernandes/analytics/db"
	"github.com/joaofnfernandes/analytics/parser/pageviews"
	"github.com/joaofnfernandes/analytics/parser/votes"
	"os"
)

const (
	DbFilenameEnv    = "DB_FILENAME"
	PageViewsFileEnv = "PAGE_VIEWS_FILENAME"
	VotesFileEnv     = "VOTES_FILENAME"
)

func main() {
	opt := getSettings()

	db.CreateIfNotExists(opt.DbFilename)
	pageviews.CsvToDB(opt.DbFilename, opt.PageViewsFilename)
	votes.CsvToDb(opt.DbFilename, opt.VotesFilename)
}

type options struct {
	DbFilename        string
	PageViewsFilename string
	VotesFilename     string
}

func getSettings() options {
	const (
		defaultDbFilename        = "bin/analytics.db"
		defaultPageViewsFilename = "data/page-views.csv"
		defaultVotesFilename     = "data/votes.csv"
	)

	opt := options{}
	if opt.DbFilename = os.Getenv(DbFilenameEnv); opt.DbFilename == "" {
		opt.DbFilename = defaultDbFilename
	}
	if opt.PageViewsFilename = os.Getenv(PageViewsFileEnv); opt.PageViewsFilename == "" {
		opt.PageViewsFilename = defaultPageViewsFilename
	}
	if opt.VotesFilename = os.Getenv(VotesFileEnv); opt.VotesFilename == "" {
		opt.VotesFilename = defaultVotesFilename
	}
	return opt
}
