package readers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var client = http.Client{Timeout: 10 * time.Second}

func NewCsvReader() (DataSource, error) {
	resp, err := client.Get("https://opendata.ecdc.europa.eu/covid19/casedistribution/csv")
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
	dateIndex        = 0
	casesIndex       = 2
	deathsIndex      = 3
	casesCountryName = 4
	casesCountryCode = 5
)

func (c *csvReader) toRecord(record []string) (*Record, error) {
	cases, err := strconv.Atoi(record[casesIndex])
	if err != nil {
		return nil, err
	}
	deaths, err := strconv.Atoi(record[deathsIndex])
	if err != nil {
		return nil, err
	}
	date, err := c.toDate(record)
	if err != nil {
		return nil, err
	}
	return &Record{
		Cases:       cases,
		CountryCode: record[casesCountryCode],
		CountryName: record[casesCountryName],
		Deaths:      deaths,
		Date:        date,
	}, nil
}

func (c *csvReader) toDate(record []string) (time.Time, error) {
	dmy := strings.Split(record[dateIndex], "/")
	d, err := strconv.Atoi(dmy[0])
	if err != nil {
		return time.Now(), err
	}
	m, err := strconv.Atoi(dmy[1])
	if err != nil {
		return time.Now(), err
	}
	y, err := strconv.Atoi(dmy[2])
	if err != nil {
		return time.Now(), err
	}
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), nil
}
