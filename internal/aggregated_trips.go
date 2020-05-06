package internal

import "time"

// AggregatedTripsRepository is the repository to manage the aggregated trips
type AggregatedTripsRepository interface {
	GetAll() ([]AggregatedTrips, error)
	Add(trips *AggregatedTrips)
}

// AggregatedTrips represents a document with the aggregated trips details
type AggregatedTrips struct {
	ID                     int64     `json:"id"`
	Date                   time.Time `json:"date" transl:"fecha"`
	SupplierID             string    `json:"supplier_id" transl:"id_proveedor"`
	TotalRides             int       `json:"total_rides" transl:"tot_viajes"`
	TotalVehicleRides      int       `json:"total_vehicle_rides" transl:"tot_veh_viaje"`
	TotalAvailableVehicles int       `json:"total_available_vehicles" transl:"tot_veh_disp"`
	TotalDistTraveled      float64   `json:"total_dist_traveled" transl:"dist_pasajero"`
	PassingTime            float64   `json:"passing_time" transl:"tiempo_pasajero"`
	RequestTime            float64   `json:"request_time" transl:"tiempo_solicitud"`
	EmptyTime              float64   `json:"empty_time" transl:"tiempo_vacio"`
	EodMultiplier          float64   `json:"eod_multiplier" transl:"multiplicador_eod"`
	Accessibility          float64   `json:"accessibility" transl:"accesibilidad"`
	FemaleOperator         float64   `json:"female_operator" transl:"operador_mujer"`
	EodStart               int       `json:"eod_start" transl:"inicio_eod"`
	EodEnd                 int       `json:"eod_end" transl:"fin_eod"`
	EodPassDist            float64   `json:"eod_pass_dist" transl:"dist_pasajero_eod"`
	EodPassTime            int       `json:"eod_pass_time" transl:"tiempo_pasajero_eod"`
	RequestDist            float64   `json:"request_dist" transl:"dist_solicitud"`
	EmptyDist              float64   `json:"empty_dist" transl:"dist_vacío"`
	EodRequestDist         float64   `json:"eod_request_dist" transl:"dist_solicitud_eod"`
	EodRequestTime         float64   `json:"eod_request_time" transl:"tiempo_solicitud_eod"`
	EodEmptyDist           float64   `json:"eod_empty_dist" transl:"dist_vacio_eod"`
	EodEmptyTime           float64   `json:"eod_empty_time" transl:"tiempo_vacio_eod"`
}
