package votes

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/joaofnfernandes/analytics/parser"
	_ "github.com/mattn/go-sqlite3"
)

// csvToDb imports a polldaddy csv file and inserts into a sqlite db
func CsvToDb(dataSourceName string, csvFilePath string) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	votes := importFromCsv(csvFilePath)
	for _, vote := range votes {
		if !vote.isDefault() {
			vote.update(db)
		}
	}
}

type vote struct {
	Url           string
	TotalVotes    int
	PositiveVotes int
	NegativeVotes int
}

// isDefault checks if a vote is set with its default values
func (v *vote) isDefault() bool {
	empty := vote{}
	if v.Url != empty.Url {
		return false
	}
	if v.TotalVotes != empty.TotalVotes {
		return false
	}
	if v.PositiveVotes != empty.PositiveVotes {
		return false
	}
	if v.NegativeVotes != empty.NegativeVotes {
		return false
	}
	return true
}

// getRating computes the rating sentiment for a page
// -1 if sentiment is negative, +1 if positive. 0 if neutral or no data
func (v *vote) getRating() float32 {
	rating := float32(0.0)
	if v.TotalVotes > 0 && v.PositiveVotes > 0 {
		rating = (float32(v.PositiveVotes) / float32(v.TotalVotes)) - (float32(v.NegativeVotes) / float32(v.TotalVotes))
	}
	return rating
}

// NewVote is a constructor for a vote but returns error if the vote is set
// with default values
func newVote(csvRecord []string) (vote, error) {
	v := vote{}
	var err error

	if len(csvRecord) < 7 {
		return v, errors.New(fmt.Sprintf("Trying to create vote with from invalid csv: %s", csvRecord))
	}

	v.Url, err = normalizeUrl(csvRecord[1])
	if err != nil {
		return v, err
	}
	v.TotalVotes, err = parser.StringToInt(csvRecord[4])
	if err != nil {
		return v, err
	}
	v.PositiveVotes, err = parser.StringToInt(csvRecord[5])
	if err != nil {
		return v, err
	}
	v.NegativeVotes, err = parser.StringToInt(csvRecord[6])
	if err != nil {
		return v, err
	}

	if v.isDefault() {
		err = errors.New("Created vote with default values")
	}
	return v, err
}

// normalizeUrl takes a url and transforms it into /path/to/resource/
// returns error if url is invalid or empty
func normalizeUrl(url string) (string, error) {
	var err error
	url = strings.Replace(url, ".md", "", 1)
	url = fmt.Sprintf("/%s/", url)

	if url == "" {
		err = errors.New(fmt.Sprintf("Trying to normalize invalid url: %s", url))
	}
	return url, err
}

// TODO: consider returning error
func (v *vote) update(db *sql.DB) {
	const sqlStmt = `update page
	set votes = ?, rating = ?
	where url = ?`

	// TODO: do we really need to start a transaction?
	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
	}
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(v.TotalVotes, v.getRating(), v.Url)
	if err != nil {
		log.Print(err)
	}
	tx.Commit()
}

// importFromCsv parses the polldaddy votes.csv file into a typed structure
func importFromCsv(filename string) []vote {
	records, _ := parser.ImportCSV(filename)

	votes := make([]vote, 1)
	// Skip the CSV headers
	for _, record := range records[1:] {
		vote, err := newVote(record)
		if err != nil {
			log.Print(err)
			continue
		}
		votes = append(votes, vote)
	}
	return votes
}
