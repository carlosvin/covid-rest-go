package main

import (
	"fmt"
	"io"

	"github.com/carlosvin/covid-rest-go/repo"

	"github.com/carlosvin/covid-rest-go/handlers"
	"github.com/carlosvin/covid-rest-go/readers"
	"github.com/gin-gonic/gin"
)

func main() {
	csvReader, err := readers.NewCsvReader()
	if err != nil {
		panic(err)
	}
	router := handlers.NewRouter(repo.NewRepo(csvReader))
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(record)
	}
	r := gin.Default()
	r.GET("/countries", router.Countries)
	r.GET("/countries/:code", router.Country)
	r.GET("/countries/:code/dates", router.CountryDates)
	r.GET("/countries/:code/dates/:date", router.CountryDate)
	r.Run()
}
