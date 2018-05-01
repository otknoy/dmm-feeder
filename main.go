package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/otknoy/dmm-crawler/model"
)

func main() {
	subscribe()
}

func subscribe() {
	log.Println("Subscriber")

	conn, _ := redis.Dial("tcp", "localhost:6379")
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe("dmm-items")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			item := &model.DmmItem{}
			if err := json.Unmarshal(v.Data, item); err != nil {
				fmt.Println(err)
			}

			// fmt.Println(string(v.Data))
			fmt.Println(item.Title)
			// fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Println("error")
		}
	}
}
