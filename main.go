package main

import (
	"time"

	"github.com/carlosvin/covid-rest-go/repo"

	"github.com/carlosvin/covid-rest-go/handlers"
	"github.com/carlosvin/covid-rest-go/readers"
	"github.com/gin-gonic/gin"
)

func main() {
	repository := repo.NewRepo(readers.NewFactory())
	router := handlers.NewRouter(repository)
	err := repository.Fetch()
	if err != nil {
		panic(err)
	}
	repository.Watch(time.Hour)
	r := gin.Default()
	r.GET("/countries", router.Countries)
	r.GET("/countries/:code", router.Country)
	r.GET("/countries/:code/dates", router.CountryDates)
	r.GET("/countries/:code/dates/:date", router.CountryDate)
	r.Run()
}
