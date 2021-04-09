package types

type SensorCons struct {
	EntityId   string `yaml:"entity_id"`
	State      string `yaml:"state"`
	Attributes struct {
		UnitOfMeasurement string `yaml:"unit_of_measurement"`
		FriendlyName      string `yaml:"friendly_name"`
		Icon              string `yaml:"icon"`
	} `json:"attributes"`
	LastChanged string `yaml:"last_changed"`
	LastUpdated string `yaml:"last_updated"`
	Context     struct {
		Id       string `yaml:"id"`
		ParentId string
		UserId   string
	} `json:"context"`
}
