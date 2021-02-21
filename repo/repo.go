package repo

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/carlosvin/covid-rest-go/readers"
)

const DateFormat = "2006-01-02"

// Repo retrieve all statistics
type Repo interface {
	Countries() map[string]Country
	// country(code string) string
	CountryDates(code string) map[time.Time]RecordInfo
	CountryDate(code string, date time.Time) RecordInfo
	Fetch() error
	Watch(period time.Duration)
}

type repoImpl struct {
	countries map[string]Country
	info      RecordInfo
	factory   readers.Factory
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

func (r *repoImpl) Watch(period time.Duration) {
	go func() {
		for true {
			time.Sleep(period)
			err := r.Fetch()
			if err == nil {
				log.Println("Remote data successfully fetched")
			} else {
				log.Println(err)
			}
		}
	}()
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

func NewRepo(factory readers.Factory) Repo {
	repo := repoImpl{
		countries: make(map[string]Country),
		info: &recordInfo{
			confirmed:    0,
			positiveRate: 0,
			path:         "",
		},
		factory: factory,
	}
	return &repo
}

func (r *repoImpl) Fetch() error {
	source, err := r.factory.NewReader()
	if err != nil {
		return err
	}
	for {
		var record, err = source.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		r.info.Add(record)
		country, found := r.countries[record.CountryCode]
		if !found {
			country = &recordCountry{
				code:  record.CountryCode,
				name:  record.CountryName,
				dates: make(map[time.Time]RecordInfo),
				info: &recordInfo{
					confirmed:    0,
					positiveRate: 0,
					path:         fmt.Sprintf("%s/countries/%s", r.info.Path(), record.CountryCode),
				},
			}
			r.countries[record.CountryCode] = country
		}
		country.Add(record)
	}
	return nil
}

type recordInfo struct {
	confirmed    int
	path         string
	positiveRate float64
	total        int
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
			confirmed:    0,
			positiveRate: 0,
			total:        0,
			path:         fmt.Sprintf("%s/dates/%s", c.info.Path(), record.Date.Format(DateFormat)),
		}
		c.dates[record.Date] = date
	}
	date.Add(record)
}

func (r *recordInfo) ConfirmedCases() int {
	return r.confirmed
}

func (r *recordInfo) PositiveRate() float64 {
	return r.positiveRate / float64(r.total)
}
func (r *recordInfo) Path() string {
	return r.path
}
func (r *recordInfo) Add(r2 *readers.Record) {
	r.confirmed += r2.Cases
	r.positiveRate += r2.PositiveRate
	r.total++
}
