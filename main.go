package main

import (
	"encoding/json"
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
				log.Fatal(err)
			}
			return c, err
		},
	}
}

func subscribe(dmmItems chan<- model.DmmItem) {
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
					log.Print(err)
				}

				dmmItems <- item

				count++
				log.Println(count)
			case redis.Subscription:
				log.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				log.Println(v)
			}
		}
		conn.Close()
	}
}

func process(dmmItems <-chan model.DmmItem, items chan<- model.Item) {
	ips := service.NewItemProcessService()

	for dmmItem := range dmmItems {
		log.Print("process")
		item := ips.Process(dmmItem)
		items <- item
	}
	close(items)
}

func feed(items <-chan model.Item) {
	solr := infrastructure.NewSolrIndexUpdater("otknoy.dip.jp", 80, "items")
	count := 0
	for item := range items {
		log.Println(item.ID)
		err := solr.AddItem(item)
		if err != nil {
			log.Fatalln(err)
		}

		count++
		if count%1000 == 0 {
			solr.Commit()
		}
	}
}
