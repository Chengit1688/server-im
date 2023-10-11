package sms

type ISms interface {
	Send(phoneNumber string, countryCode string, code string) error
}
