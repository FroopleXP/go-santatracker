package types

type Location struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type LocationData struct {
	Status       string        `json:"status"`
	Language     string        `json:"language"`
	Destinations []Destination `json:"destinations"`
}

func (l *LocationData) AdjustDepArrTimes(offset int64) {
	for idx, _ := range l.Destinations {
		l.Destinations[idx].AdjustTime(offset)
	}
}

type LocationReport struct {
	Position          Location
	PresentsDelivered int64
	Previous          Destination
	Next              Destination
	LastSeen          Destination
	CurrentTimeMs     int64
	Status            Status
}
