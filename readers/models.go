package readers

import (
	"time"
)

type Record struct {
	Date        time.Time
	Cases       int
	Deaths      int
	CountryCode string
	CountryName string
}
