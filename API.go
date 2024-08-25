package main

import (
	"fmt"
	"strings"
)

func ReceiveSMS(repo *MongoDBRepository, phoneNumber string, content string) {

	content = strings.ToLower(content)

	if content == "help" {
		InfoHelp(phoneNumber)
		return
	}
	if content == "miasta" {
		InfoCities(phoneNumber)
		return
	}
	if content == "rodo" {
		InfoRODO(phoneNumber)
		return
	}
	if CheckIfCityInBase(content) != true {
		InfoError(phoneNumber)
	}
	if CheckIfCityInBase(content) == true {
		content = synonims[content]
		repo.AddRecord(phoneNumber, content)
	}

}

func SendSMS(phoneNumber string, message string) {
	fmt.Println("na numer:" + phoneNumber + "wiadomość:\n" + message) //docelowo przekazanie do funkcji wysyłającej
}
func SendFeed(repo *MongoDBRepository) {
	data, _ := repo.ReadData()

	for _, record := range data {
		wc := NewWeatherConditions(WeatherFetcher(record.City))
		SendSMS(record.PhoneNumber, wc.WeatherConditionMessage())
	}
}

func CheckIfCityInBase(city string) bool {

	_, ok := synonims[city]

	return ok
}

func InfoError(phoneNumber string) {
	message := "Wpisałeś złą komendę. Wyślij 'help' aby dowiedzieć się jak korzystać z usługi biker_info"
	SendSMS(phoneNumber, message)
}

func InfoHelp(phoneNumber string) {
	message := "Witaj! Usługa biker_info codziennie wieczorem dostarcza  informacje o pogodzie spodziewanej na kolejny dzień. " +
		"Zaplanuj swoją drogę do pracy bez pogodowych niespodzianek!\n Aby uruchomić usługę dla danej miejscowości wyślij jej nazwę.\n" +
		"Aby wyłączyć usługę wyślij jej nazwę ponownie.\n" +
		"Aby uzyskać listę miejscowości dostępnych w ramach usługi wyślij 'miasta'." +
		"Informacja RODO wyślij 'rodo'."
	SendSMS(phoneNumber, message)
}
func InfoCities(phoneNumber string) {
	var message string
	for key, _ := range YRcities {
		message += key + "\n"
	}
	message = "Usługa biker_info obecnie dostępna jest dla miejscowości:\n" + message
	SendSMS(phoneNumber, message)
}

func InfoStatusAdded(phoneNumber string, city string) {
	message := "Usługa pomyślnie uruchomiona dla miejscowości " + city + "! \n Codziennie o 20 otrzymasz raport pogodowy dla jednośladów, na następny dzień. Aby wyłączyć usługę wpisz nazwę miejscowości dla której chcesz ją wyłączyć. Informacja RODO - wyślij 'rodo'.\""
	SendSMS(phoneNumber, message)
}

func InfoStatusDeleted(phoneNumber string, city string) {
	message := "Usługa pomyślnie wyłączona dla miejscowości " + city + "! \n Nie będziesz otrzymasz więcej raportów dla tej miejscowiości. Aby wyłączyć kolejną usługę, wyślij nazwę miejscowości. W celu uzyskania pomocy wyślij 'help'."
	SendSMS(phoneNumber, message)
}

func InfoRODO(phoneNumber string) {
	message := "Biker_info: informacja RODO \n Uruchamiając usługę zgadzasz się na automatyczne przetwarzanie Twojego numeru telefonu jedynie w celu dostarczania włączonej usługi. Jeżeli nie zgadzasz się na przetwarzanie, wyłącz usługę, a Twoje dane osobowe zostaną usunięte z bazy."
	SendSMS(phoneNumber, message)
}
