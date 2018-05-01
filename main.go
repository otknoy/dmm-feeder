package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"

	"github.com/otknoy/dmm-feeder/infrastructure"
	"github.com/otknoy/dmm-feeder/model"
	"github.com/otknoy/dmm-feeder/service"
)

func main() {
	dmmItems := make(chan model.DmmItem)
	items := make(chan model.Item)

	go subscribe(dmmItems)

	go process(dmmItems, items)
	feed(items)
}

func subscribe(dmmItems chan<- model.DmmItem) {
	log.Println("Subscriber")

	conn, _ := redis.Dial("tcp", "localhost:6379")
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe("dmm-items")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			var item model.DmmItem
			if err := json.Unmarshal(v.Data, &item); err != nil {
				fmt.Println(err)
			}

			dmmItems <- item

			// fmt.Println(string(v.Data))
			// fmt.Println(item.Title)
			// fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Println(v)
		}
	}
}

func process(dmmItems <-chan model.DmmItem, items chan<- model.Item) {
	ips := service.NewItemProcessService()

	for dmmItem := range dmmItems {
		item := ips.Process(dmmItem)
		items <- item
	}
	close(items)
}

func feed(items <-chan model.Item) {
	solr := infrastructure.NewSolrRepository()
	for item := range items {
		log.Println(item.ID)
		err := solr.Add(item)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
