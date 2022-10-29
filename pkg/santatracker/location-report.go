package santatracker

type LocationReport struct {
	Position          Location
	PresentsDelivered int64
	Previous          Destination
	Next              Destination
	LastSeen          Destination
	CurrentTimeMs     int64
	Status            Status
}
