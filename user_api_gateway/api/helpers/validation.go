package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

func ValidatePhone(phone string) error {
	re := regexp.MustCompile(`^[+][9][9][8]\d{9}$`)
	if !re.MatchString(phone) {
		return errors.New("invalid phone number: " + phone)
	}

	return nil
}

func ValidateDates(startDate, endDate string) error {
	var layout = "2006-01-02 15:04:05"

	from, err := time.Parse(layout, startDate)
	if err != nil {
		return errors.New("start_date is invalid")
	}

	to, err := time.Parse(layout, endDate)
	if err != nil {
		return errors.New("end_time is invalid")
	}

	if !from.Before(to) {
		return errors.New("start_time can not be greater than end_time")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be blank")
	}
	if len(password) < 8 || len(password) > 30 {
		return errors.New("password length should be 8 to 30 characters")
	}
	if validation.Validate(password, validation.Match(regexp.MustCompile("^[A-Za-z0-9$_@.#]+$"))) != nil {
		return errors.New("password should contain only alphabetic characters, numbers and special characters(@, $, _, ., #)")
	}
	if validation.Validate(password, validation.Match(regexp.MustCompile("[0-9]"))) != nil {
		return errors.New("password should contain at least one number")
	}
	if validation.Validate(password, validation.Match(regexp.MustCompile("[A-Za-z]"))) != nil {
		return errors.New("password should contain at least one alphabetic character")
	}
	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username cannot be blank")
	}
	if len(username) < 5 || len(username) > 30 {
		return errors.New("username length should be 6 to 30 characters")
	}
	if validation.Validate(username, validation.Match(regexp.MustCompile("^[A-Za-z0-9$@_.#]+$"))) != nil {
		return errors.New("username should contain only alphabetic characters, numbers and special characters(@, $, _, ., #)")
	}
	return nil
}

// func ValidateEmailAddress(email string) error {
// 	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// 	if !emailRegex.MatchString(email) {
// 		return fmt.Errorf("email address %s is not valid", email)
// 	}

// 	return nil
// }

// CheckEventRegistration checks if there are at least 3 hours remaining until the event
func CheckEventRegistration(eventTime time.Time) error {

	now := time.Now()

	// Combine the date part of 'now' with the time part of 'eventTime'
	eventDateTime := time.Date(now.Year(), now.Month(), now.Day(), eventTime.Hour(), eventTime.Minute(), eventTime.Second(), 0, now.Location())

	// If the event time for today has already passed, assume it's for the next day
	if eventDateTime.Before(now) {
		eventDateTime = eventDateTime.Add(24 * time.Hour)
	}

	remainingTime := eventDateTime.Sub(now)

	if remainingTime < 3*time.Hour {
		return fmt.Errorf("error: less than 3 hours remaining until the event")
	}

	return nil
}

func IsSunday(date string) error {
	layout := "2006-01-02"
	t, err := time.Parse(layout, date)

	if err != nil {
		return err
	}
	if t.Weekday() != time.Sunday {
		return fmt.Errorf("date is not Sunday")
	}
	return err
}
