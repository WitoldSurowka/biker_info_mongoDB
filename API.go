package main

import (
	"fmt"
	"strings"
	"time"
)

func ReceiveSMS(repoUsers *MongoDBRepository, repoOutboundMessages *MongoDBRepository,
	phoneNumber string, content string) {

	content = strings.ToLower(content)

	if content == "help" {
		InfoHelp(repoOutboundMessages, phoneNumber)
		return
	}
	if content == "miasta" {
		InfoCities(repoOutboundMessages, phoneNumber)
		return
	}
	if content == "rodo" {
		InfoRODO(repoOutboundMessages, phoneNumber)
		return
	}
	if CheckIfCityInBase(content) != true {
		InfoError(repoOutboundMessages, phoneNumber)
	}
	if CheckIfCityInBase(content) == true {
		content = synonims[content]
		repoUsers.UsersAddRecord(repoOutboundMessages, phoneNumber, content)
	}

}

func OutboundSMS(repoOutboundSMS *MongoDBRepository, phoneNumber string, message string) {
	fmt.Println("na numer:" + phoneNumber + "wiadomość:\n" + message) //docelowo przekazanie do funkcji wysyłającej
	repoOutboundSMS.OutboundSMSAddSMS(phoneNumber, message)
}
func SendFeed(repoUsers, repoOutboundMessages *MongoDBRepository) {
	start := time.Now() // zanotowanie czasu rozpoczęcia

	data, _ := repoUsers.UsersReadData()

	for _, record := range data {
		wc := NewWeatherConditions(WeatherFetcher(record.City))
		OutboundSMS(repoOutboundMessages, record.PhoneNumber, wc.WeatherConditionMessage())
	}

	duration := time.Since(start) // obliczenie czasu trwania
	fmt.Printf("someFunction took %v to execute\n", duration)
}

func CheckIfCityInBase(city string) bool {

	_, ok := synonims[city]

	return ok
}

func InfoError(repoOutboundMessages *MongoDBRepository, phoneNumber string) {
	message := "Wpisałeś złą komendę. Wyślij 'help' aby dowiedzieć się jak korzystać z usługi biker_info"
	OutboundSMS(repoOutboundMessages, phoneNumber, message)
}

func InfoHelp(repoOutboundMessages *MongoDBRepository, phoneNumber string) {
	message := "Witaj! Usługa biker_info codziennie wieczorem dostarcza  informacje o pogodzie spodziewanej na kolejny dzień. " +
		"Zaplanuj swoją drogę do pracy bez pogodowych niespodzianek!\n Aby uruchomić usługę dla danej miejscowości wyślij jej nazwę.\n" +
		"Aby wyłączyć usługę wyślij jej nazwę ponownie.\n" +
		"Aby uzyskać listę miejscowości dostępnych w ramach usługi wyślij 'miasta'." +
		"Informacja RODO wyślij 'rodo'."
	OutboundSMS(repoOutboundMessages, phoneNumber, message)
}
func InfoCities(repoOutboundMessages *MongoDBRepository, phoneNumber string) {
	var message string
	for key, _ := range YRcities {
		message += key + "\n"
	}
	message = "Usługa biker_info obecnie dostępna jest dla miejscowości:\n" + message
	OutboundSMS(repoOutboundMessages, phoneNumber, message)
}

func InfoStatusAdded(repoOutboundMessages *MongoDBRepository, phoneNumber string, city string) {
	message := "Usługa pomyślnie uruchomiona dla miejscowości " + city + "! \n Codziennie o 20 otrzymasz raport pogodowy dla jednośladów, na następny dzień. Aby wyłączyć usługę wpisz nazwę miejscowości dla której chcesz ją wyłączyć. Informacja RODO - wyślij 'rodo'.\""
	OutboundSMS(repoOutboundMessages, phoneNumber, message)
}

func InfoStatusDeleted(repoOutboundMessages *MongoDBRepository, phoneNumber string, city string) {
	message := "Usługa pomyślnie wyłączona dla miejscowości " + city + "! \n Nie będziesz otrzymasz więcej raportów dla tej miejscowiości. Aby wyłączyć kolejną usługę, wyślij nazwę miejscowości. W celu uzyskania pomocy wyślij 'help'."
	OutboundSMS(repoOutboundMessages, phoneNumber, message)
}

func InfoRODO(repoOutboundMessages *MongoDBRepository, phoneNumber string) {
	message := "Biker_info: informacja RODO \n Uruchamiając usługę zgadzasz się na automatyczne przetwarzanie Twojego numeru telefonu jedynie w celu dostarczania włączonej usługi. Jeżeli nie zgadzasz się na przetwarzanie, wyłącz usługę, a Twoje dane osobowe zostaną usunięte z bazy."
	OutboundSMS(repoOutboundMessages, phoneNumber, message)
}
