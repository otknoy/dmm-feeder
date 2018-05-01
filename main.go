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

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 1000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func subscribe(dmmItems chan<- model.DmmItem) {
	log.Println("Subscriber")

	pool := newPool()

	count := 0

	for {
		conn := pool.Get()
		psc := redis.PubSubConn{Conn: conn}

		psc.Subscribe("dmm-items")

		for conn.Err() == nil {
			switch v := psc.Receive().(type) {
			case redis.Message:
				var item model.DmmItem
				if err := json.Unmarshal(v.Data, &item); err != nil {
					fmt.Println(err)
				}

				dmmItems <- item

				count++
				log.Println(count)

				// fmt.Println(string(v.Data))
				// fmt.Println(item.Title)
				// fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				fmt.Println(v)
			}
		}
		conn.Close()
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
