package bollobas

import (
	"nanomsg.org/go/mangos/v2"
	"time"
)

const (
	PASSENGER_REQUESTS = "request"
	PASSENGER_CANCEL = "cancel"
	PSSENGER_ACCEPTANCE_NOTIFIED = "notified"

	RIDE_CONFIRMED = "passenger_ack_notified"
)

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


type RideRequest struct {
	UserID   string
	RquestID int
}

type RideRequestCancellation RideRequest

type RideRequestConfirmed struct {
	UserID   string
	RquestID int
	Date time.Time
}


type AnalyticsHandler interface{
	mangos.Socket
	Run()
}