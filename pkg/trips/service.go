package trips

type Service interface {
	CreateTripPlan(name string)
	ReadTripPlan(id string) TripPlan

	// Transportation

	CreateFlight() error
	DeleteFlight() error
	UpdateFlight() error

	CreateTransit() error
	DeleteTransit() error
	UpdateTransit() error

	// Lodging

	CreateLodging() error
	DeleteLodging() error
	UpdateLodging() error
}
