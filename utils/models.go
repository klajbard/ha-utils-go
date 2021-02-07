package utils

type Attributes struct {
	Friendly_name string `json:"friendly_name"`
	Icon          string `json:"icon"`
}

type Sensor struct {
	State      string `json:"state"`
	Attributes Attributes
}
