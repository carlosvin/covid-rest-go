package handlers

import (
	"strings"
	"time"

	"github.com/carlosvin/covid-rest-go/repo"
	"github.com/gin-gonic/gin"
)

type Router interface {
	Countries(c *gin.Context)
	Country(c *gin.Context)
	CountryDates(c *gin.Context)
	CountryDate(c *gin.Context)
}

type routerImpl struct {
	repo repo.Repo
}

func NewRouter(repo repo.Repo) Router {
	return &routerImpl{repo: repo}
}

func (r *routerImpl) Countries(c *gin.Context) {
	c.JSON(200, r.countries())
}

func (r *routerImpl) Country(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	resp := r.country(code)
	if resp != nil {
		c.JSON(200, resp)
	} else {
		c.JSON(404, gin.H{"message": "not found"})
	}
}

func (r *routerImpl) CountryDates(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	resp := r.countryDates(code)
	if resp != nil {
		c.JSON(200, resp)
	} else {
		c.JSON(404, gin.H{"message": "not found"})
	}
}

func (r *routerImpl) CountryDate(c *gin.Context) {
	codeParam := strings.ToUpper(c.Param("code"))
	dateParam := c.Param("date")
	resp := r.countryDate(codeParam, dateParam)
	if resp != nil {
		c.JSON(200, r.countryDate(codeParam, dateParam))
	} else {
		c.JSON(404, gin.H{"message": "not found"})
	}
}

func (r *routerImpl) country(code string) *countryResponse {
	country, found := r.repo.Countries()[code]
	if !found {
		return nil
	}
	return r.toCountryResponse(country)
}

func (r *routerImpl) countries() map[string]*countryResponse {
	countries := make(map[string]*countryResponse)
	for code, country := range r.repo.Countries() {
		countries[code] = r.toCountryResponse(country)
	}
	return countries
}

func (r *routerImpl) countryDates(code string) map[string]*dateResponse {
	dates := r.repo.CountryDates(code)
	if dates == nil {
		return nil
	}
	datesResp := make(map[string]*dateResponse)
	for t, date := range dates {
		datesResp[t.Format(repo.DateFormat)] = r.toDateResponse(t, date)
	}
	return datesResp
}

func (r *routerImpl) countryDate(code string, date string) *dateResponse {
	t, err := time.Parse(repo.DateFormat, date)
	if err != nil {
		return nil
	}
	record := r.repo.CountryDate(code, t)
	if record == nil {
		return nil
	}
	return r.toDateResponse(t, record)
}

func (r *routerImpl) toCountryResponse(country repo.Country) *countryResponse {
	return &countryResponse{
		response:    r.toResponse(country.Info()),
		CountryCode: country.CountryCode(),
		CountryName: country.CountryName(),
	}
}

func (r *routerImpl) toDateResponse(date time.Time, info repo.RecordInfo) *dateResponse {
	return &dateResponse{
		response:  r.toResponse(info),
		EpochDays: date.Unix() / (int64(time.Hour.Seconds()) * 24),
		Date:      date.Format(repo.DateFormat),
	}
}

func (r *routerImpl) toResponse(rec repo.RecordInfo) *response {
	return &response{
		Deaths:    rec.DeathsNumber(),
		Confirmed: rec.ConfirmedCases(),
		Path:      rec.Path(),
	}
}

type response struct {
	Deaths    int    `json:"deathsNumber"`
	Confirmed int    `json:"confirmedCases"`
	Path      string `json:"path"`
}

type countryResponse struct {
	*response
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
}

type dateResponse struct {
	*response
	Date      string `json:"date"`
	EpochDays int64  `json:"epochDay"`
}
