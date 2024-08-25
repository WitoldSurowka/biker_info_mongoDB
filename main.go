package main

import (
	"log"
	"time"
)

func main() {

	repo, err := NewMongoDBRepository("X", "Y", "Z")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	currentTime := time.Now()

	ticker := time.NewTicker(1 * time.Second)

	ReceiveSMS(repo, "123123123", "ikasfhbgdiasfgiasdf")
	ReceiveSMS(repo, "123123123", "Rodo")
	ReceiveSMS(repo, "123123123", "help")
	ReceiveSMS(repo, "123123123", "miasta")
	ReceiveSMS(repo, "123123123", "Ropica")
	ReceiveSMS(repo, "123123123", "Ropica")
	ReceiveSMS(repo, "123123123", "Ropica")
	ReceiveSMS(repo, "123123123", "Kraków")
	ReceiveSMS(repo, "123123124", "Ropica")

	for range ticker.C {
		currentTime = time.Now()
		//fmt.Println(currentTime)
		//if currentTime.Hour() == 20 && currentTime.Minute() == 22 && currentTime.Second() == 50 {
		//	wc := NewWeatherConditions(WeatherFetcher("Ropica Górna"))
		//	fmt.Println(wc.WeatherConditionMessage())
		//}
		if currentTime.Second() == 50 {
			SendFeed(repo)
		}

	}

}
