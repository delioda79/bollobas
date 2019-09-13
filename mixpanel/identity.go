package mixpanel

import "time"

type Identity struct {
	FirstName string `json:"$first_name,omitempty"`
	LastName string `json:"$last_name,omitempty"`
	RegistrationDate time.Time `json:"$created,omitempty"`
	Type string `json:"type,omitempty"`
	Email string `json:"$email,omitempty"`
	ReferralCode string `json:"referral_code,omitempty"`
	Phone string `json:"$phone,omitempty"`
}
