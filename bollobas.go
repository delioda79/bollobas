package bollobas

import "time"

// Identity represents a basic Analytics identity, passengers and drivers
type Identity struct {
	ID               string
	FirstName        string
	LastName         string
	RegistrationDate time.Time
	Type             string
	Email            string
	ReferralCode     string
	Phone            string
}
