package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"time"
)

type Data struct {
	Records   []Record `json:"records"`
	CurrentID int      `json:"currentId"`
}

type Record struct {
	ID               int       `bson:"id"`
	PhoneNumber      string    `bson:"phoneNumber"`
	City             string    `bson:"city"`
	SubscriptionDate time.Time `bson:"subscriptionDate"` // New field for subscription date
}

type MongoDBRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type WeatherConditions struct {
	Precip  float64 // Precipitation in mm
	TempMin int     // Minimum temperature in °C
	Wind    float64 // Wind speed in km/h
	City    string
}

var weekdays = map[string]string{
	"Monday":       "Poniedziałek",
	"Tuesday":      "Wtorek",
	"Wednesday":    "Środa",
	"Thursday":     "Czwartek",
	"Friday":       "Piątek",
	"Saturday":     "Sobota",
	"%!Weekday(7)": "Niedziela",
}

var YRcities = map[string]string{
	"Kraków":           "https://www.yr.no/nb/v%C3%A6rvarsel/daglig-tabell/2-3094802/Polen/Ma%C5%82opolskie/Krak%C3%B3w/Krak%C3%B3w",
	"Górna Wieś":       "https://www.yr.no/nb/v%C3%A6rvarsel/daglig-tabell/2-3098807/Polen/Ma%C5%82opolskie/Powiat%20krakowski/G%C3%B3rna%20Wie%C5%9B",
	"Kędzierzyn-Koźle": "https://www.yr.no/nb/v%C3%A6rvarsel/daglig-tabell/2-3096372/Polen/Opolskie/Powiat%20k%C4%99dzierzy%C5%84sko-kozielski/K%C4%99dzierzyn-Ko%C5%BAle",
	"Ropica Górna":     "https://www.yr.no/nb/v%C3%A6rvarsel/daglig-tabell/2-760338/Polen/Ma%C5%82opolskie/Powiat%20gorlicki/Ropica%20G%C3%B3rna",
	"Lipowa k.Żywca":   "https://www.yr.no/nb/v%C3%A6rvarsel/daglig-tabell/2-3093263/Poland/Silesia/%C5%BBywiec%20County/Lipowa",
	"Gdańsk":           "https://www.yr.no/nb/v%C3%A6rvarsel/daglig-tabell/2-3099434/Polen/Pomorskie/Gda%C5%84sk/Gda%C5%84sk",
	"Pcim":             "https://www.yr.no/nb/v%C3%A6rvarsel/daglig-tabell/2-3089310/Poland/Lesser%20Poland/My%C5%9Blenice%20County/Pcim",
	"Katowice":         "https://www.yr.no/nb/v%C3%A6rvarsel/daglig-tabell/2-3096472/Poland/Silesia/Katowice/Katowice",
	"Jędrzejów":        "https://www.yr.no/nb/v%C3%A6rvarsel/daglig-tabell/2-770157/Polen/%C5%9Awi%C4%99tokrzyskie/Powiat%20j%C4%99drzejowski/J%C4%99drzej%C3%B3w",
}

var synonims = map[string]string{
	"kraków":            "Kraków",
	"krakow":            "Kraków",
	"górna wieś":        "Górna Wieś",
	"gorna wies":        "Górna Wieś",
	"górna wies":        "Górna Wieś",
	"gorna wieś":        "Górna Wieś",
	"kędierzyn-koźle":   "Kędierzyn-Koźle",
	"kedierzyn-kozle":   "Kędierzyn-Koźle",
	"kędierzyn-kozle":   "Kędierzyn-Koźle",
	"kedierzyn-koźle":   "Kędierzyn-Koźle",
	"kedierzyn kozle":   "Kędierzyn-Koźle",
	"kędierzyn kozle":   "Kędierzyn-Koźle",
	"kedierzyn koźle":   "Kędierzyn-Koźle",
	"ropica górna":      "Ropica Górna",
	"ropica":            "Ropica Górna",
	"ropica gorna":      "Ropica Górna",
	"lipowa k.Żywca":    "Lipowa k.Żywca",
	"lipowa":            "Lipowa k.Żywca",
	"lipowa k.zywca":    "Lipowa k.Żywca",
	"lipowa k. zywca":   "Lipowa k.Żywca",
	"lipowa kolo zywca": "Lipowa k.Żywca",
	"lipowa koło zywca": "Lipowa k.Żywca",
	"lipowa kolo żywca": "Lipowa k.Żywca",
	"gdańsk":            "Gdańsk",
	"gdansk":            "Gdańsk",
	"pcim":              "Pcim",
	"Katowice":          "Katowice",
	"jędrzejów":         "Jędrzejów",
	"jedrzejow":         "Jędrzejów",
	"jędrzejow":         "Jędrzejów",
	"jedrzejów":         "Jędrzejów",
}

func roundToPlaces(value float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(value*shift+0.5) / shift
}
