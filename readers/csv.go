package readers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/icza/gox/timex"
)

var client = http.Client{Timeout: 10 * time.Second}

func NewCsvReader() (DataSource, error) {
	resp, err := client.Get("https://opendata.ecdc.europa.eu/covid19/testing/csv")
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(resp.Body)
	reader.Read()
	return &csvReader{reader: reader}, nil
}

// csvReader CSV data source
type csvReader struct {
	reader *csv.Reader
}

// Fetch Fetches data from CSV source
func (c *csvReader) Read() (*Record, error) {

	record, err := c.reader.Read()
	if err != nil {
		return nil, err
	}

	return c.toRecord(record)

}

const (
	dateWeekIndex    = 2
	casesIndex       = 6
	casesCountryName = 0
	casesCountryCode = 1
	positiveRate     = 10
)

func (c *csvReader) toRecord(record []string) (*Record, error) {
	cases, err := strconv.Atoi(record[casesIndex])
	if err != nil {
		return nil, fmt.Errorf("Error parsing %v: %s", record, err)
	}
	positiveRate, err := strconv.ParseFloat(record[positiveRate], 32)
	if err != nil {
		positiveRate = 0
	}
	date, err := c.toDate(record)
	if err != nil {
		return nil, fmt.Errorf("Error parsing %v: %s", record, err)
	}
	return &Record{
		Cases:        cases,
		CountryCode:  record[casesCountryCode],
		CountryName:  record[casesCountryName],
		PositiveRate: positiveRate,
		Date:         date,
	}, nil
}

func (c *csvReader) toDate(record []string) (time.Time, error) {
	yearWeek := strings.Split(record[dateWeekIndex], "-W")
	y, err := strconv.Atoi(yearWeek[0])
	if err != nil {
		return time.Now(), err
	}
	w, err := strconv.Atoi(yearWeek[1])
	if err != nil {
		return time.Now(), err
	}
	return timex.WeekStart(y, w), nil
}
