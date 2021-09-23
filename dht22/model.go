package dht22

import (
	"fmt"
	"time"

	"github.com/MichaelS11/go-dht"
)

type MyDHT22 struct {
	Temperature float64
	Humidity    float64
	Timestamp   time.Time
}

func (d MyDHT22) Read() MyDHT22 {
	err := dht.HostInit()
	if err != nil {
		fmt.Println("HostInit error:", err)
		return d
	}

	dht, err := dht.NewDHT("GPIO2", dht.Celsius, "")
	if err != nil {
		fmt.Println("NewDHT error:", err)
		return d
	}

	humidity, temperature, err := dht.ReadRetry(11)
	if err != nil {
		fmt.Println("Read error:", err)
		return d
	}
	d.Humidity = humidity
	d.Temperature = temperature
	d.Timestamp = time.Now()

	fmt.Printf("humidity: %v\n", d.Humidity)
	fmt.Printf("temperature: %v\n", d.Temperature)
	fmt.Printf("temperature: %v\n", d.Timestamp)
	return d
}
