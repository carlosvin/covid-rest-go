package repo

import (
	"fmt"
	"io"
	"time"

	"github.com/carlosvin/covid-rest-go/readers"
	constants "github.com/carlosvin/covid-rest-go/utils"
)

// Repo retrieve all statistics
type Repo interface {
	Countries() map[string]Country
	// country(code string) string
	CountryDates(code string) map[time.Time]RecordInfo
	CountryDate(code string, date time.Time) RecordInfo
}

type repoImpl struct {
	countries map[string]Country
	info      RecordInfo
}

func (r *repoImpl) Countries() map[string]Country {
	return r.countries
}
func (r *repoImpl) CountryDates(code string) map[time.Time]RecordInfo {
	country, found := r.countries[code]
	if found {
		return country.Dates()
	}
	return nil
}

func (r *repoImpl) CountryDate(code string, date time.Time) RecordInfo {
	country, found := r.countries[code]
	if found {
		r, _ := country.Dates()[date]
		return r
	}
	return nil
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
		country.Add(record)
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
			path:      fmt.Sprintf("%s/dates/%s", c.info.Path(), record.Date.Format(constants.DateFormat)),
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
