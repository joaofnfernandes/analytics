package parser

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

// stringToInt transforms a string into a number "1,234" => 1234
func StringToInt(s string) (int, error) {
	s = strings.Replace(s, ",", "", -1)
	v, err := strconv.Atoi(s)
	return v, err
}

// TODO: needs to be more resilient to handle google analytics
// csv files without manual cleanup
func ImportCSV(filename string) (record [][]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment = '#'
	records, err := reader.ReadAll()

	return records, err
}
