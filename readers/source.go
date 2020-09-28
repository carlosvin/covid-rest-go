package readers

import "time"

// DataSource in charge of fetching the data from its source
type DataSource interface {
	Read() (*Record, error)
}

type Record struct {
	Date        time.Time
	Cases       int
	Deaths      int
	CountryCode string
	CountryName string
}
