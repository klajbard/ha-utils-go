package utils

type Attributes struct {
	Friendly_name       string `json:"friendly_name"`
	Icon                string `json:"icon"`
	Unit_of_measurement string `json:"unit_of_measurement"`
}

type Sensor struct {
	State      string `json:"state"`
	Attributes Attributes
}
