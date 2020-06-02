package view

// AggregatedTrips is the transport layer view of the Aggregated Trips.
type AggregatedTrips struct {
	ID                     int64    `json:"id"`
	Date                   string   `transl:"date" json:"fecha"`
	SupplierID             *string  `transl:"supplier_id" json:"id_proveedor"`
	TotalRides             *int     `transl:"total_rides" json:"tot_viajes"`
	TotalVehicleRides      *int     `transl:"total_vehicle_rides" json:"tot_veh_viaje"`
	TotalAvailableVehicles *int     `transl:"total_available_vehicles" json:"tot_veh_disp"`
	TotalDistTraveled      *float64 `transl:"total_dist_traveled" json:"dist_pasajero"`
	PassingTime            *float64 `transl:"passing_time" json:"tiempo_pasajero"`
	RequestTime            *float64 `transl:"request_time" json:"tiempo_solicitud"`
	EmptyTime              *float64 `transl:"empty_time" json:"tiempo_vacio"`
	EodMultiplier          *float64 `transl:"eod_multiplier" json:"multiplicador_eod"`
	Accessibility          *float64 `transl:"accessibility" json:"accesibilidad"`
	FemaleOperator         *float64 `transl:"female_operator" json:"operador_mujer"`
	EodStart               *int     `transl:"eod_start" json:"inicio_eod"`
	EodEnd                 *int     `transl:"eod_end" json:"fin_eod"`
	EodPassDist            *float64 `transl:"eod_pass_dist" json:"dist_pasajero_eod"`
	EodPassTime            *int     `transl:"eod_pass_time" json:"tiempo_pasajero_eod"`
	RequestDist            *float64 `transl:"request_dist" json:"dist_solicitud"`
	EmptyDist              *float64 `transl:"empty_dist" json:"dist_vacío"`
	EodRequestDist         *float64 `transl:"eod_request_dist" json:"dist_solicitud_eod"`
	EodRequestTime         *float64 `transl:"eod_request_time" json:"tiempo_solicitud_eod"`
	EodEmptyDist           *float64 `transl:"eod_empty_dist" json:"dist_vacio_eod"`
	EodEmptyTime           *float64 `transl:"eod_empty_time" json:"tiempo_vacio_eod"`
}

// OperatorStats is the view model of the Operator Stats.
type OperatorStats struct {
	ID             int64
	Date           string
	OperatorID     *string
	Gender         *int
	CompletedTrips *int
	DaysSince      *int
	AgeRange       *string
	HoursConnected *string
	TripHours      *string
	TotRevenue     *string
}

// TrafficIncident is the view model of the Traffic Incident.
type TrafficIncident struct {
	ID             int64
	Date           string
	Type           *int
	Plates         *string
	Licence        *string
	TravelDistance *string
	TravelTime     *string
	Coordinates    *string
}
