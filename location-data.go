package santatracker

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
