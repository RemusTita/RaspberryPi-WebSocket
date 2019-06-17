package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"net/http"
	"time"
)

var (
	r   = raspi.NewAdaptor()
	adc = spi.NewMCP3008Driver(r)
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type sensor struct {
	Humidity int `json:"humidity"`
}

func main() {
	if err := adc.Start(); err != nil {
		log.Println("Unable to initialize ADC Driver ", err)
		return
	}
	http.HandleFunc("/ws", handler)
	log.Println("WebSocket Initialized")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Connection Establish")
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		analogRead, err := adc.AnalogRead("A0")
		if err != nil {
			log.Println(err)
		}
		percentConv := mapPercent(analogRead, 250, 1023, 100, 0)
		read := sensor{
			Humidity: percentConv,
		}

		if err := conn.WriteJSON(read); err != nil {
			log.Println(err)
			return
		}
		time.Sleep(2 * time.Second)
	}
}

func mapPercent(x, inMin, inMax, outMin, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}
