package readers

type Factory interface {
	NewReader() (DataSource, error)
}

func NewFactory() Factory {
	return &factoryImpl{}
}

type factoryImpl struct{}

func (f *factoryImpl) NewReader() (DataSource, error) {
	return NewCsvReader()
}
