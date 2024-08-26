package main

import (
	"log"
	"time"
)

func main() {

	//repoUsers, repoOutboundSMS, err := NewMongoDBRepository("X", "Y", "Z")
	repoUsers, err := NewMongoDBRepository("X", "Y", "Z")
	repoOutboundSMS, err := NewMongoDBRepository("X", "Y", "Z")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	currentTime := time.Now()

	ticker := time.NewTicker(1 * time.Second)

	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "ikasfhbgdiasfgiasdf")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "Rodo")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "help")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "miasta")
	ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "Ropica")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "Ropica")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "Ropica")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "Kraków")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123124", "Ropica")

	for range ticker.C {
		currentTime = time.Now()
		//fmt.Println(currentTime)
		//if currentTime.Hour() == 20 && currentTime.Minute() == 22 && currentTime.Second() == 50 {
		//	wc := NewWeatherConditions(WeatherFetcher("Ropica Górna"))
		//	fmt.Println(wc.WeatherConditionMessage())
		//}
		if currentTime.Second() == 50 {

			go SendFeed(repoUsers, repoOutboundSMS)

		}

	}

}
