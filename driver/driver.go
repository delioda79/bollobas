package driver

type Driver struct {
	ID string `json:"id"`
	Email string
	Fname string
	Lname string
	PhonePrefix int `json:"phone_prefix"`
	AreaPrefix int `json:"area_prefix"`
	PhoneNo string `json:"phone_no"`
	ReferralCode string `json:"registration_id_reference"`
	RegistrationDate int64 `json:"created"`
}
