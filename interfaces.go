package main

type Repo interface {
	AddRecord(phoneNumber string, city string) error
	DeleteRecord(id int) error
	ReadData() (Data, error)
}
