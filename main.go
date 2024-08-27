package main

import (
	"log"
)

func main() {

	//repoUsers, repoOutboundSMS, err := NewMongoDBRepository("X", "Y", "Z")
	repoUsers, err := NewMongoDBRepository("mongodb://biker_witold:sylvia@localhost:27017/", "biker_info_DB", "users")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	repoOutboundSMS, err := NewMongoDBRepository("mongodb://biker_witold:sylvia@localhost:27017/", "biker_info_DB", "outboundSMS")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	repoInboundSMS, err := NewMongoDBRepository("mongodb://biker_witold:sylvia@localhost:27017/", "biker_info_DB", "inboundSMS")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "ikasfhbgdiasfgiasdf")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "Rodo")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "help")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "miasta")
	ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "Ropica")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "Ropica")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "Ropica")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123123", "Kraków")
	//ReceiveSMS(repoUsers, repoOutboundSMS, "123123124", "Ropica")

	go eventSchedule(repoUsers, repoInboundSMS, repoOutboundSMS)

	runServer(repoUsers, repoInboundSMS, repoOutboundSMS)

}
