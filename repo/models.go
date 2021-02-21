package repo

import (
	"time"

	"github.com/carlosvin/covid-rest-go/readers"
)

type RecordInfo interface {
	ConfirmedCases() int
	PositiveRate() float64
	Path() string
	Add(r *readers.Record)
}

type Countries interface {
	Info() RecordInfo
	Entries() map[string]Country
}

type Country interface {
	Info() RecordInfo
	CountryCode() string
	CountryName() string
	Dates() map[time.Time]RecordInfo
	Add(rec *readers.Record)
}
