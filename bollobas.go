package bollobas

import "time"

type Identity struct {
	ID string
	FirstName string
	LastName string
	RegistrationDate time.Time
	Type string
	Email string
	ReferralCode string
	Phone string
}