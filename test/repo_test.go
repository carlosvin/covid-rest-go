package yours

import (
	"testing"

	"github.com/carlosvin/covid-rest-go/readers"
	"github.com/carlosvin/covid-rest-go/repo"
	"github.com/stretchr/testify/assert"
)

func TestES(t *testing.T) {

	repository := repo.NewRepo(readers.NewReaderFactory())
	repository.Fetch()
	assert.NotEmpty(t, repository.Countries()["ES"], "country should be found")
}
