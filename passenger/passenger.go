package passenger

import "time"

type Passenger struct {
	ID string `json:"passenger_id"`
	Email string
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	PhoneNo string `json:"phone"`
	PhonePrefix string `json:"phone_prefix"`
	RegistrationDate time.Time `json:"registration_date"`
	Action string `json:"event_action"`
}
