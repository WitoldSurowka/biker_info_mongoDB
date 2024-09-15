package main

import "time"

func eventSchedule(repoUsers, repoInboundSMS, repoOutboundSMS *MongoDBRepository) {

	currentTime := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			currentTime = time.Now()

			//if currentTime.Second() == 40 {
			//	repoInboundSMS.InboundOutboundSMSDeleteProcessedSMS()
			//	repoOutboundSMS.InboundOutboundSMSDeleteProcessedSMS()
			//}

			if currentTime.Minute() == 10 && currentTime.Second() == 00 {
				go SendFeed(repoUsers, repoOutboundSMS)
			}
		}
	}()
}
