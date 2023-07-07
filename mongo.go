package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func getMongoEnv() string {
	host := getEnv("MONGO_HOST")
	port := getEnv("MONGO_PORT")
	url := fmt.Sprintf("mongodb://%s:%s", host, port)
	return url
}

func MongoMain() {

	url := getMongoEnv()
	counter := 1
	for {
		MongoWorker(url)
		time.Sleep(time.Duration(counter) * time.Second)
		if counter <= 10 {
			counter++
		}
	}
}

func MongoWorker(url string) {
	ctx := context.TODO()
	opts := options.Client().ApplyURI(url)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Println("Ошибка подключения к MongoDB: ", err)
		return
	}

	defer client.Disconnect(ctx)
	fmt.Printf("%T\n", client)

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println(fmt.Sprintf("Подключние по url %s, не состоялось: ", url), err)
		return
	}
}
