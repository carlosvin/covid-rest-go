package readers

// DataSource in charge of fetching the data from its source
type DataSource interface {
	Read() (*Record, error)
}
