package view

// AggregatedTrips is the transport layer view of the Aggregated Trips.
type AggregatedTrips struct {
	ID                     int64    `json:"id"`
	Date                   string   `json:"fecha"`
	SupplierID             *string  `json:"id_proveedor"`
	TotalRides             *int     `json:"tot_viajes"`
	TotalVehicleRides      *int     `json:"tot_veh_viaje"`
	TotalAvailableVehicles *int     `json:"tot_veh_disp"`
	TotalDistTraveled      *float64 `json:"dist_pasajero"`
	PassingTime            *float64 `json:"tiempo_pasajero"`
	RequestTime            *float64 `json:"tiempo_solicitud"`
	EmptyTime              *float64 `json:"tiempo_vacio"`
	EodMultiplier          *float64 `json:"multiplicador_eod"`
	Accessibility          *float64 `json:"accesibilidad"`
	FemaleOperator         *float64 `json:"operador_mujer"`
	EodStart               *int     `json:"inicio_eod"`
	EodEnd                 *int     `json:"fin_eod"`
	EodPassDist            *float64 `json:"dist_pasajero_eod"`
	EodPassTime            *int     `json:"tiempo_pasajero_eod"`
	RequestDist            *float64 `json:"dist_solicitud"`
	EmptyDist              *float64 `json:"dist_vacio"`
	EodRequestDist         *float64 `json:"dist_solicitud_eod"`
	EodRequestTime         *float64 `json:"tiempo_solicitud_eod"`
	EodEmptyDist           *float64 `json:"dist_vacio_eod"`
	EodEmptyTime           *float64 `json:"tiempo_vacio_eod"`
}

// OperatorStats is the view model of the Operator Stats.
type OperatorStats struct {
	ID             int64   `json:"id"`
	OperatorID     *string `json:"id_operador"`
	Gender         *int    `json:"genero"`
	CompletedTrips *int    `json:"cant_viajes"`
	DaysSince      *int    `json:"tiempo_registro"`
	AgeRange       *string `json:"edad"`
	HoursConnected *string `json:"horas_conectado"`
	TripHours      *string `json:"horas_viaje"`
	TotRevenue     *string `json:"ingreso_totales"`
}

// TrafficIncident is the view model of the Traffic Incident.
type TrafficIncident struct {
	ID             int64   `json:"id"`
	Type           *int    `json:"hecho_trans"`
	Plates         *string `json:"placa"`
	Licence        *string `json:"licencia"`
	TravelDistance *string `json:"distancia_viaje"`
	TravelTime     *string `json:"tiempo_viaje"`
	Coordinates    *string `json:"ubicacion"`
	Date           string  `json:"tiempo_hecho"`
}

// ErrorSwagger view model
type ErrorSwagger struct {
	Error string `json:"error"`
}
