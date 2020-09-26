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
	//r.GET("/countries/:country", countryHandler)
	//r.GET("/countries/:country/dates", countryDatesHandler)
	//r.GET("/countries/:country/dates/:date", countryDateHandler)
	r.Run()
}
