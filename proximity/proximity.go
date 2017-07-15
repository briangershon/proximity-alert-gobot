package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	board := firmata.NewAdaptor(os.Args[1])
	blue := gpio.NewLedDriver(board, "3")
	soundSensor := aio.NewAnalogSensorDriver(board, "5")
	irSensor := aio.NewAnalogSensorDriver(board, "4")

	work := func() {
		soundSensor.On(aio.Data, func(data interface{}) {
			influxData := bytes.NewBufferString(fmt.Sprintf("sound_level value=%v", data))
			_, influxErr := http.Post("http://localhost:8086/write?db=gobot", "application/x-www-form-urlencoded", influxData)
			if influxErr != nil {
				log.Println("influxErr", influxErr.Error())
			}

			if data.(int) < 150 {
				blue.Off()
			} else {
				blue.On()
			}
		})

		irSensor.On(aio.Data, func(data interface{}) {
			influxData := bytes.NewBufferString(fmt.Sprintf("ir_light_level value=%v", data))
			_, influxErr := http.Post("http://localhost:8086/write?db=gobot", "application/x-www-form-urlencoded", influxData)
			if influxErr != nil {
				log.Println("influxErr", influxErr.Error())
			}

			if data.(int) < 150 {
				blue.Off()
			} else {
				blue.On()
			}
		})
	}

	robot := gobot.NewRobot("proximity",
		[]gobot.Connection{board},
		[]gobot.Device{soundSensor, irSensor, blue},
		work,
	)

	robot.Start()
}
