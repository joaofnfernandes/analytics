package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// stringToInt transforms a string into a number "1,234" => 1234
func StringToInt(s string) (int, error) {
	s = strings.Replace(s, ",", "", -1)
	v, err := strconv.Atoi(s)
	return v, err
}

// stringPercentToFloat converts "12.48%" into float32(12.48)
func StringPercentToFloat(s string) (float32, error) {
	s = strings.Replace(s, "%", "", -1)
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

// normalizeUrl makes sure a url is represented as path/to/page/
// and returns an error if string is not a url
func NormalizeUrl(url string) (result string, err error) {
	re := regexp.MustCompile(`\/?([\w\d-_]+(\/[\w\d-_]+)*)(\/|\.md)`)
	match := re.FindStringSubmatch(url)

	if len(match) >= 2 {
		result = match[1]
	} else {
		err = errors.New(fmt.Sprintf("Invalid url: %s", url))
	}
	return result, err
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
