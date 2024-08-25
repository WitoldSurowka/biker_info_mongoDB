package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
	"time"
)

// initializing a data structure to keep the scraped data
type WeatherStatus struct {
	date, url, precip, tempMin, wind string
}

func WeatherFetcher(city string) (float64, int, float64, string) {
	var status WeatherStatus
	c := colly.NewCollector()
	shouldStop := false
	dateTomorrow := time.Now().Add(time.Hour * 24)
	dateTomorrowString := fmt.Sprintf(dateTomorrow.Format("2006-01-02"))

	c.OnHTML(".daily-weather-list-item", func(e *colly.HTMLElement) {
		if shouldStop {
			return
		}
		status.date = e.ChildAttr("time", "datetime")
		status.url = e.ChildAttr("a", "href")
		status.precip = e.ChildText(".Precipitation-module__main-sU6qN[data-tone=primary]")
		if len(e.ChildText("span.temperature.min-max-temperature__min.temperature--warm-primary")) != 0 {
			status.tempMin = e.ChildText("span.temperature.min-max-temperature__min.temperature--warm-primary")
		}
		if len(e.ChildText("span.temperature.min-max-temperature__min.temperature--cold-primary")) != 0 {
			status.tempMin = e.ChildText("span.temperature.min-max-temperature__min.temperature--cold-primary")
		}
		status.wind = e.ChildText("div.daily-weather-list-item__wind")
		//c.OnHTML scrape in a loop, so after the desired data is fetched, we do not process data no more
		//fmt.Println(status)
		if strings.EqualFold(dateTomorrowString, status.date) {
			shouldStop = true
		}
	})

	c.Visit(YRcities[city])
	c.Wait() // Wait until scraping is complete

	//process precipitation string->float
	precipStringLong := status.precip[:len(status.precip)-2]
	precipStringShort := precipStringLong[8:]
	precipStringShort = strings.Replace(precipStringShort, ",", ".", 1)
	precip, err := strconv.ParseFloat(precipStringShort, 32)
	if err != nil {
		fmt.Println("Precip conversion error:", err)
		precip = 0
	}
	precip = roundToPlaces(precip, 2)
	//process tempMin string->int
	tempMinString := status.tempMin[:len(status.tempMin)-2]
	tempMin, err := strconv.Atoi(tempMinString)
	if err != nil {
		fmt.Println("TempMin conversion error:", err)
		tempMin = 15
	}

	//process wind string->int
	windStringLong := status.wind[:len(status.wind)-3]
	windStringShort := windStringLong[5:]
	wind, err := strconv.ParseFloat(windStringShort, 64)
	if err != nil {
		fmt.Println("Wind conversion error:", err)
		wind = 0
	}
	return precip, tempMin, wind, city
}
