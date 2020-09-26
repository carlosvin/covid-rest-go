package handlers

import (
	"github.com/carlosvin/covid-rest-go/repo"
	"github.com/gin-gonic/gin"
)

type Router interface {
	Countries(c *gin.Context)
}

type routerImpl struct {
	repo repo.Repo
}

func NewRouter(repo repo.Repo) Router {
	return &routerImpl{repo: repo}
}

func (r *routerImpl) Countries(c *gin.Context) {
	c.JSON(200, gin.H{"countries": r.countries()})
}

func (r *routerImpl) countries() map[string]*countryResponse {
	countries := make(map[string]*countryResponse)
	for code, country := range r.repo.Countries() {
		countries[code] = r.toCountryResponse(country)
	}
	return countries
}

func (r *routerImpl) toCountryResponse(country repo.Country) *countryResponse {
	return &countryResponse{
		Deaths:      country.Info().DeathsNumber(),
		Confirmed:   country.Info().ConfirmedCases(),
		CountryCode: country.CountryCode(),
		CountryName: country.CountryName(),
		Path:        country.Info().Path(),
	}
}

type countryResponse struct {
	Deaths      int    `json:"deathsNumber"`
	Confirmed   int    `json:"confirmedCases"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	Path        string `json:"path"`
}
