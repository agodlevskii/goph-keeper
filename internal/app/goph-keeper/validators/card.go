package validators

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrCardNumberFormat = errors.New("card number must match the format xxxx xxxx xxxx xxxx")
	ErrCVVFormat        = errors.New("card's CardCVV must be a 3-digits value")
	ErrExpDateFormat    = errors.New("card's expire date must match the format mm/yy")
	ErrExpDateInPast    = errors.New("card's expire date must not be in the past")
)

func CardNumber(num string) error {
	if ok, err := regexp.MatchString("(\\d{4}( )?){4}", num); err != nil || !ok {
		return ErrCardNumberFormat
	}
	return nil
}

func CardCVV(cvv string) error {
	if ok, err := regexp.MatchString("\\d{3}", cvv); err != nil || !ok {
		return ErrCVVFormat
	}
	return nil
}

func CardExpDate(date string) error {
	if ok, err := regexp.MatchString("\\d{2}/\\d{2}", date); err != nil || !ok {
		return ErrExpDateFormat
	}

	dateTime, err := time.Parse("01/06", date)
	if err != nil {
		return ErrExpDateFormat
	}

	if time.Now().After(dateTime) {
		return ErrExpDateInPast
	}
	return nil
}
