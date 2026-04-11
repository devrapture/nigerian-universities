package models

type Institution struct {
	Name                string `json:"name"`
	ViceChancellor      string `json:"vice_chancellor"`
	YearOfEstablishment string `json:"year_of_establishment"`
	Type                string `json:"type"`
	Url                 string `json:"url"`
}
