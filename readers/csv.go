package readers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"
)

func NewCsvReader() (DataSource, error) {
	resp, err := http.Get("https://opendata.ecdc.europa.eu/covid19/casedistribution/csv")
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
	casesIndex       = 4
	casesCountryCode = 7
	casesCountryName = 6
	deathsIndex      = 5
	dayIndex         = 1
	monthIndex       = 2
	yearIndex        = 3
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
	y, err := strconv.Atoi(record[yearIndex])
	if err != nil {
		return time.Now(), err
	}
	m, err := strconv.Atoi(record[monthIndex])
	if err != nil {
		return time.Now(), err
	}
	d, err := strconv.Atoi(record[dayIndex])
	if err != nil {
		return time.Now(), err
	}
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), nil
}
