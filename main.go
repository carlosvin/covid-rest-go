package main

import (
	"time"

	"github.com/carlosvin/covid-rest-go/repo"

	"github.com/carlosvin/covid-rest-go/handlers"
	"github.com/carlosvin/covid-rest-go/readers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/swaggo/gin-swagger/example/basic/docs"
)

// @title COVID cases tracking API
// @version 1.0
// @description Exposes COVID-19 statistics extracted from daily published data by European Centre for Disease Prevention and Control at: https://www.ecdc.europa.eu/en/publications-data/download-todays-data-geographic-distribution-covid-19-cases-worldwide.
// @contact.name API support
// @contact.url https://github.com/carlosvin/covid-rest-go/issues
// @license.name MIT
// @license.url ./LICENSE

// @host covid-rest.appspot.com
// @BasePath /
func main() {
	url := ginSwagger.URL("http://localhost:8080/docs/doc.json") // The url pointing to API definition

	repository := repo.NewRepo(readers.NewReaderFactory())
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
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.Run()
}
