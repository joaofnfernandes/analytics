package pageviews

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/joaofnfernandes/analytics/parser"
	_ "github.com/mattn/go-sqlite3"
)

func CsvToDB(dataSourceName string, csvFilePath string) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pageViews := importFromCsv(csvFilePath)
	for _, pageView := range pageViews {
		if !pageView.isDefault() {
			pageView.insert(db)
		}
	}
}

type pageView struct {
	Url             string
	PageViews       int
	UniquePageViews int
	AvgTime         int
	BounceRate      float32
}

func (p *pageView) isDefault() bool {
	empty := pageView{}
	if p.Url != empty.Url {
		return false
	}
	if p.PageViews != empty.PageViews {
		return false
	}
	if p.UniquePageViews != empty.UniquePageViews {
		return false
	}
	if p.AvgTime != empty.AvgTime {
		return false
	}
	if p.BounceRate != empty.BounceRate {
		return false
	}
	return true
}

func (p *pageView) insert(db *sql.DB) {
	const sqlStmt = `
	insert into page(url, pageviews, unique_pageviews, avg_time, bounce_rate)
	values (?, ?, ?, ?, ?)
	`

	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		return
	}
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Print(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Url, p.PageViews, p.UniquePageViews, p.AvgTime, p.BounceRate)
	if err != nil {
		log.Printf("Failed to insert page view. Error: %v", err)
	}
	tx.Commit()
}

// newPageView is a constructor for a page view that returns error
// if the page view has just default values
func newPageView(csvRecord []string) (pageView, error) {
	p := pageView{}
	var err, currErr error

	if len(csvRecord) < 8 {
		return p, errors.New(fmt.Sprintf("Trying to create page view from invalid csv: %v", csvRecord))
	}

	p.Url, currErr = normalizeUrl(csvRecord[0])
	if err == nil {
		err = currErr
	}

	p.PageViews, err = parser.StringToInt(csvRecord[1])
	if err == nil {
		err = currErr
	}

	p.UniquePageViews, err = parser.StringToInt(csvRecord[2])
	if err == nil {
		err = currErr
	}

	p.AvgTime, err = normalizeDuration(csvRecord[3])
	if err == nil {
		err = currErr
	}
	p.BounceRate, err = stringPercentToFloat(csvRecord[5])
	if err == nil {
		err = currErr
	}

	if p.isDefault() {
		err = errors.New("Created page view with default values")
	}
	return p, err
}

// normalizeUrl takes a url and transforms it into /path/to/resource/
// returns error if url is invalid or empty
func normalizeUrl(url string) (string, error) {
	var err error
	if url == "" {
		err = errors.New(fmt.Sprintf("Trying to normalize invalid url: %s", url))
	}
	return url, err
}

// normalizeDuration removes characters that are not part of a time
// <00:01:01 => 61
func normalizeDuration(time string) (int, error) {
	time = strings.Replace(time, "<", "", -1)
	re := regexp.MustCompile(`(\d{1,2}):(\d{2}):(\d{2})`)
	matches := re.FindStringSubmatch(time)

	if len(matches) < 4 {
		return 0, errors.New(fmt.Sprintf("Invalid duration: %s", time))
	}
	var hour, min, sec int
	var err error
	if hour, err = strconv.Atoi(matches[1]); err != nil {
		return 0, errors.New(fmt.Sprintf("Invalid duration: %s", time))
	}
	if min, err = strconv.Atoi(matches[2]); err != nil {
		return 0, errors.New(fmt.Sprintf("Invalid duration: %s", time))
	}
	if sec, err = strconv.Atoi(matches[3]); err != nil {
		return 0, errors.New(fmt.Sprintf("Invalid duration: %s", time))
	}
	return (sec + 60*min + 3600*hour), nil
}

// stringPercentToFloat converts "12.48%" into float32(12.48)
func stringPercentToFloat(s string) (float32, error) {
	s = strings.Replace(s, "%", "", -1)
	v, err := strconv.ParseFloat(s, 32)
	return (float32(v) / 100), err
}

func importFromCsv(filename string) []pageView {
	records, err := parser.ImportCSV(filename)
	if err != nil {
		log.Fatalf("Cannot import csv file: %v", err)
	}

	pageViews := make([]pageView, 1)
	//skip the CSV headers
	for _, record := range records[1:] {
		pageView, err := newPageView(record)
		if err != nil {
			log.Printf("Failed to create page view. page: %v, err: %v", pageView, err)
			continue
		}
		pageViews = append(pageViews, pageView)
	}
	return pageViews
}
