package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

func InboundLoadUnprocessedSMSes(repoUsers, repoInboundSMS, repoOutboundSMS *MongoDBRepository) {

	//lets lock this function not to let it run concurrently
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	UnreadMongoSMSarray, err := repoInboundSMS.InboundOutboundMakeUnreadMongoSMSarrayFromDB()
	if err != nil {
		fmt.Print(err)
		return
	}
	if len(UnreadMongoSMSarray) == 0 {
		fmt.Println("no Unread SMSes in inboundSMS collection - no actions in inboundSMS collection taken")
		return
	}
	successfullySentArray, err := InboundSendUnreadMongoSMSarray(UnreadMongoSMSarray, repoUsers, repoOutboundSMS)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		err = repoInboundSMS.InboundOutboundUpdateReadMongoSMSProcessedField(successfullySentArray)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
	fmt.Println(UnreadMongoSMSarray)
}
func InboundSendUnreadMongoSMSarray(UnreadMongoSMSarray []MongoSMS, repoUsers,
	repoOutboundSMS *MongoDBRepository) ([]primitive.ObjectID, error) {
	var successfullySentArray []primitive.ObjectID
	var lastError error
	for _, UnreadMongoSMS := range UnreadMongoSMSarray {
		lastError = ReceiveSMS(repoUsers, repoOutboundSMS, UnreadMongoSMS.PhoneNumber, UnreadMongoSMS.Message)
		if lastError != nil {
			fmt.Println(lastError)
		} else {
			successfullySentArray = append(successfullySentArray, UnreadMongoSMS.ID)
		}
	}
	return successfullySentArray, lastError
}
