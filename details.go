package santatracker

type Details struct {
	Timezone int     `json:"timezone"`
	Weather  Weather `json:"weather"`
}
