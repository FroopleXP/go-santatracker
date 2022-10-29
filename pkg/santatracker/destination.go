package santatracker

type Destination struct {
	Id                string   `json:"id"`
	Arrival           int64    `json:"arrival"`
	Departure         int64    `json:"departure"`
	Population        int64    `json:"population"`
	PresentsDelivered int64    `json:"presentsDelivered"`
	City              string   `json:"city"`
	Region            string   `json:"region"`
	Location          Location `json:"location"`
}

func (d *Destination) AdjustTime(offset int64) {
	d.Arrival = d.Arrival + offset
	d.Departure = d.Departure + offset
}
