// Gets DHT22 sensor data of Raspberry GPIO
package dht

import (
	"fmt"

	"../utils"
	"github.com/d2r2/go-dht"
)

// Read from DHT22 sensor which should
// be connected to GPIO 4 (Data pin) and updates
// homeassistant sensor.rpi_temperature value
func ReadDHT(pin int) {
	temp, hum, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, pin, false, 10)
	if err != nil {
		utils.PrintError(err)
		return
	}
	setTemperature(temp)
	setHumidity(hum)
}

func setTemperature(temp float32) {
	payloadTemp := utils.Sensor{
		State: fmt.Sprintf("%.2f", temp),
		Attributes: utils.Attributes{
			Friendly_name:       "RPI DHT temperature",
			Unit_of_measurement: "°C",
			Device_class:        "temperature",
		},
	}
	utils.SetState("sensor.rpi_temperature", payloadTemp)
}

func setHumidity(hum float32) {
	payloadHum := utils.Sensor{
		State: fmt.Sprintf("%.2f", hum),
		Attributes: utils.Attributes{
			Friendly_name:       "RPI DHT humidity",
			Unit_of_measurement: "%",
			Device_class:        "humidity",
		},
	}
	utils.SetState("sensor.rpi_humidity", payloadHum)
}
