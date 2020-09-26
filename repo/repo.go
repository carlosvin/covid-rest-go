package repo

import (
	"fmt"
	"io"
	"time"

	"github.com/carlosvin/covid-rest-go/readers"
)

// Repo retrieve all statistics
type Repo interface {
	Countries() map[string]Country
	// country(code string) string
	// countryDates(code string) string
	// countryDate(code string, date time.Time) string
}

type repoImpl struct {
	countries map[string]Country
	info      RecordInfo
}

func (r *repoImpl) Countries() map[string]Country {
	return r.countries
}

var c Country = &recordCountry{}

func NewRepo(source readers.DataSource) Repo {
	repo := repoImpl{
		countries: make(map[string]Country),
		info: &recordInfo{
			confirmed: 0,
			deaths:    0,
			path:      "",
		},
	}
	for {
		record, err := source.Read()
		if err == io.EOF {
			break
		}
		repo.info.Add(record)
		country, found := repo.countries[record.CountryCode]
		if !found {
			country = &recordCountry{
				code:  record.CountryCode,
				name:  record.CountryName,
				dates: make(map[time.Time]RecordInfo),
				info: &recordInfo{
					confirmed: 0,
					deaths:    0,
					path:      fmt.Sprintf("%s/countries/%s", repo.info.Path(), record.CountryCode),
				},
			}
			repo.countries[record.CountryCode] = country
		}
		country.Info().Add(record)
	}
	return &repo
}

type recordInfo struct {
	confirmed int
	path      string
	deaths    int
}

type recordCountry struct {
	info  RecordInfo
	code  string
	name  string
	dates map[time.Time]RecordInfo
}

func (c *recordCountry) Info() RecordInfo {
	return c.info
}
func (c *recordCountry) CountryCode() string {
	return c.code
}
func (c *recordCountry) CountryName() string {
	return c.name
}
func (c *recordCountry) Dates() map[time.Time]RecordInfo {
	return c.dates
}
func (c *recordCountry) Add(record *readers.Record) {
	c.info.Add(record)
	date, found := c.dates[record.Date]
	if !found {
		date = &recordInfo{
			confirmed: 0,
			deaths:    0,
			path:      fmt.Sprintf("%s/dates/%s", c.info.Path(), record.Date.Format("2006-01-02")),
		}
		c.dates[record.Date] = date
	}
	date.Add(record)
}

func (r *recordInfo) ConfirmedCases() int {
	return r.confirmed
}

func (r *recordInfo) DeathsNumber() int {
	return r.deaths
}
func (r *recordInfo) Path() string {
	return r.path
}
func (r *recordInfo) Add(r2 *readers.Record) {
	r.confirmed += r2.Cases
	r.deaths += r2.Deaths
}
