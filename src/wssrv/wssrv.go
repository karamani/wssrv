package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net/http"
)

type Point struct {
	Long float64 `json:"longitude"`
	Lat  float64 `json:"latitude"`
}

type PlayerState struct {
	Point Point
}

type Object struct {
	Name  string
	Type  string
	Point Point
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var objects = [...]Object{
	{
		Name:  "Мухомор",
		Type:  "fung",
		Point: Point{Long: 38.888758, Lat: 47.221259},
	},
	{
		Name:  "Контейнер",
		Type:  "container",
		Point: Point{Long: 38.888338, Lat: 47.222534},
	},
	{
		Name:  "Поганка",
		Type:  "fung",
		Point: Point{Long: 38.889159, Lat: 47.221054},
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Received: %s\n", string(p))

		point := &Point{}
		if err := json.Unmarshal(p, point); err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("%v\n", point)

		var visibleObjects []Object
		for _, eachObject := range objects {
			if math.Abs(eachObject.Point.Long-point.Long) <= 0.0003 ||
				math.Abs(eachObject.Point.Lat-point.Lat) <= 0.0003 {

				visibleObjects = append(visibleObjects, eachObject)
			}
		}

		resp, err := json.Marshal(visibleObjects)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Resp: %s\n", string(resp))

		err = conn.WriteMessage(messageType, resp)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	log.Println("wssrv starting...")

	http.HandleFunc("/echo", wsHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
