package main

import (
	"fmt"
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
	err = InboundSendUnreadMongoSMSarray(UnreadMongoSMSarray, repoUsers, repoInboundSMS, repoOutboundSMS)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		err = repoInboundSMS.InboundOutboundUpdateReadMongoSMSProcessedField(UnreadMongoSMSarray)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
	fmt.Println(UnreadMongoSMSarray)
}
func InboundSendUnreadMongoSMSarray(UnreadMongoSMSarray []MongoSMS, repoUsers, repoInboundSMS, repoOutboundSMS *MongoDBRepository) error {
	for _, UnreadMongoSMS := range UnreadMongoSMSarray {
		err := ReceiveSMS(repoUsers, repoOutboundSMS, UnreadMongoSMS.PhoneNumber, UnreadMongoSMS.Message)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
