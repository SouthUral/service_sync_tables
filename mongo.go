package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoBase struct {
	collection string
	database   string
	url        string
	ctx        context.Context
	client     *mongo.Client
}

func (mdb *MongoBase) getMongoEnv() {
	host := getEnv("MONGO_HOST")
	port := getEnv("MONGO_PORT")
	mdb.collection = getEnv("MONGO_COLLECTION")
	mdb.database = getEnv("MONGO_DATABASE")
	mdb.url = fmt.Sprintf("mongodb://%s:%s", host, port)
}

// Запускает MongoWorker в цикле, если происходит дисконнект, то MongoWorker будет
// постоянно перезапускаться
func (mdb *MongoBase) MongoMain(ch_input chan MessCommand, ch_output chan MessCommand) {
	mdb.getMongoEnv()
	counter := 1
	for {
		err := mdb.connectMongo()
		if err != nil {
			time.Sleep(time.Duration(counter) * time.Second)
			if counter <= 10 {
				counter++
			}
			continue
		}
		mdb.mongoWorker(ch_input, ch_output)
	}
}

func (mdb *MongoBase) mongoWorker(ch_input chan MessCommand, ch_output chan MessCommand) {

	// defer mdb.client.Disconnect(mdb.ctx)
	// fmt.Printf("%T\n", mdb.client)

	collection := mdb.client.Database(mdb.database).Collection(mdb.collection)

	for msg := range ch_input {
		switch msg.Info {
		case "get_all":
			data, err := mdb.getAllData(collection)
			if err != nil {
				answer := MessCommand{
					Info: "False",
					Data: StateMess{},
				}
				ch_output <- answer
				return
			}
			for _, messDb := range data {
				answer := MessCommand{
					Info: "True",
					Data: *messDb}

				ch_output <- answer
			}
		case "input_data":
			// метод используется для добавления нового состояния в БД
			var answer MessCommand
			err, answDB := mdb.inputData(msg.Data, collection)
			if err != nil {
				answer = MessCommand{
					Info: "False",
					Data: msg.Data,
				}
				ch_output <- answer
				return
			}
			answer = MessCommand{
				Info: "True",
				Data: answDB}
			ch_output <- answer
		case "drop_colllection":
			err := mdb.dropCollection(collection)
			if err != nil {
				return
			}
		case "update":
			var answer MessCommand
			err := mdb.updateData(msg.Data, collection)
			if err != nil {
				answer = MessCommand{
					Info: "False",
					Data: msg.Data,
				}
				ch_output <- answer
				return
			}
			answer = MessCommand{
				Info: "True",
				Data: msg.Data,
			}
			ch_output <- answer
		}
	}

}

// метод для обновления документа
func (mdb *MongoBase) updateData(data StateMess, collection *mongo.Collection) interface{} {
	// проверка подключения
	if err := mdb.checkConn(); err != nil {
		return err
	}

	filter := bson.M{"_id": data.oid}
	updated := bson.M{
		"$set": bson.M{
			"table":    data.Table,
			"database": data.DataBase,
			"offset":   data.Offset,
			"isactive": data.IsActive,
		},
	}
	_, err := collection.UpdateOne(mdb.ctx, filter, updated)
	if err != nil {
		log.Println("updateData error: ", err)
		return err
	}
	log.Println("data updated: ", data.oid)
	return nil

}

// метод для добавления документа, возвращает ошибку и заполенную структуру с oid из монго
func (mdb *MongoBase) inputData(data StateMess, colection *mongo.Collection) (interface{}, StateMess) {
	var resMess interface{}
	var DbObject StateMess

	// проверка подключения
	if err := mdb.checkConn(); err != nil {
		return err, DbObject
	}

	insertResult, err := colection.InsertOne(mdb.ctx, data)
	if err != nil {
		log.Println("Insert error: ", err)
		return err, DbObject
	}
	resMess = insertResult.InsertedID

	fmt.Println("Inserted a single document: ", resMess)
	DbObject = StateMess{
		oid:      resMess,
		Table:    data.Table,
		DataBase: data.DataBase,
		Offset:   data.Offset,
		IsActive: data.IsActive,
	}

	return nil, DbObject
}

// метод для удаления документа
func (mdb *MongoBase) deleteData() {

}

// метод для удалении коллекции
func (mdb *MongoBase) dropCollection(collection *mongo.Collection) interface{} {

	// проверка подключения
	if err := mdb.checkConn(); err != nil {
		return err
	}

	err := collection.Drop(mdb.ctx)
	if err != nil {
		log.Println("Ошибка удаления коллекции: ", err)
		return err
	}
	return nil
}

// метод для получения всех документов
func (mdb *MongoBase) getAllData(collection *mongo.Collection) ([]*StateMess, interface{}) {
	res := []*StateMess{}

	// проверка подключения
	if err := mdb.checkConn(); err != nil {
		return res, err
	}

	cursor, err := collection.Find(mdb.ctx, bson.M{})
	if err != nil {
		return res, err
	}
	for cursor.Next(mdb.ctx) {
		var elem StateMess
		err := cursor.Decode(&elem)
		if err != nil {
			return res, err
		}
		res = append(res, &elem)
	}
	return res, nil
}

// метод для проверки подключения
func (mdb *MongoBase) checkConn() interface{} {
	if err := mdb.client.Ping(mdb.ctx, readpref.Primary()); err != nil {
		log.Println(fmt.Sprintf("Подключние по url %s, не состоялось: ", mdb.url), err)
		return err
	}
	return nil
}

// метод для установления коннекта с монго
func (mdb *MongoBase) connectMongo() interface{} {
	var err error

	mdb.ctx = context.TODO()
	opts := options.Client().ApplyURI(mdb.url)

	mdb.client, err = mongo.Connect(mdb.ctx, opts)
	if err != nil {
		log.Println("Ошибка подключения к MongoDB: ", err)
		return err
	}

	if err := mdb.checkConn(); err != nil {
		return err
	}

	return nil
}
