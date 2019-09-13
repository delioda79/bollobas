package driver

import "time"

type Driver struct {
	ID string `json:"driver_id"`
	Email string
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	PhonePrefix int `json:"phone_prefix"`
	AreaPrefix string `json:"area_prefix"`
	PhoneNo string `json:"phone"`
	ReferralCode string `json:"registration_id_reference"`
	RegistrationDate time.Time `json:"registration_date"`
	Action string `json:"event_action"`
}
