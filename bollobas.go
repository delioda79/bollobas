package bollobas

import (
	"time"

	"nanomsg.org/go/mangos/v2"
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
	Market           string
}

// RideRequest represents an internal message notifying a ride request
type RideRequest struct {
	UserID   string
	RquestID int
}

// RideRequestCancellation represents an internal message notifying of a ride request cancellation
type RideRequestCancellation RideRequest

// RideRequestConfirmed represents an internal message notifying of a ride request confirmed by a driver
type RideRequestConfirmed struct {
	UserID   string
	RquestID int
	Date     time.Time
}

// AnalyticsHandler is an interface for analytics handlers. A handler will register to internal producers and receive messages
type AnalyticsHandler interface {
	mangos.Socket
	Run()
}

// ConfigurationManager is a gneral interface which can be used in order to modify a configuration and check for validity
type ConfigurationManager interface {
	Configure(map[string]interface{})
	Check([]byte) bool
}
