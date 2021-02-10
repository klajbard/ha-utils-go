package dht

import (
	"fmt"
	"log"

	"../utils"
	"github.com/d2r2/go-dht"
)

func ReadDHT(pin int) {
	temp, hum, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, pin, false, 10)
	if err != nil {
		log.Fatal(err)
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
