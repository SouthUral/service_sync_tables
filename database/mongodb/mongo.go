package mongodb

import (
	"context"
	"fmt"
	"regexp"

	"time"

	tools "github.com/SouthUral/service_sync_tables/tools"

	log "github.com/sirupsen/logrus"
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
	host := tools.GetEnv("MONGO_HOST")
	port := tools.GetEnv("MONGO_PORT")
	mdb.collection = tools.GetEnv("MONGO_COLLECTION")
	mdb.database = tools.GetEnv("MONGO_DATABASE")
	mdb.url = fmt.Sprintf("mongodb://%s:%s", host, port)
}

func MDBInit(ch_input chan MessCommand, ch_output chan MessCommand) {
	mdb := MongoBase{}
	log.Debug("Запуск MDB")
	go mdb.MongoMain(ch_input, ch_output)
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

	collection := mdb.client.Database(mdb.database).Collection(mdb.collection)

	for msg := range ch_input {
		switch msg.Info {
		case GetAll:
			data, err := mdb.getAllData(collection)
			if err != nil {
				answer := MessCommand{
					Info:  GetAll,
					Data:  StateMess{},
					Error: err,
				}
				ch_output <- answer
				return
			}
			for _, messDb := range data {
				answer := MessCommand{
					Info:  GetAll,
					Data:  *messDb,
					Error: nil,
				}

				ch_output <- answer
			}
		case InputData:
			// метод используется для добавления нового состояния в БД
			var answer MessCommand
			err, answDB := mdb.inputData(msg.Data, collection)
			if err != nil {
				answer = MessCommand{
					Info:  InputData,
					Data:  msg.Data,
					Error: err,
				}
				ch_output <- answer
				return
			}
			answer = MessCommand{
				Info:  InputData,
				Data:  answDB,
				Error: nil,
			}
			ch_output <- answer
		case DropCollection:
			err := mdb.dropCollection(collection)
			if err != nil {
				return
			}
		case UpdateData:
			var answer MessCommand
			err := mdb.updateData(msg.Data, collection)
			if err != nil {
				answer = MessCommand{
					Info:  UpdateData,
					Data:  msg.Data,
					Error: err,
				}
				ch_output <- answer
				return
			}
			answer = MessCommand{
				Info:  UpdateData,
				Data:  msg.Data,
				Error: nil,
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

	filter := bson.M{"table": data.Table,
		"database": data.DataBase}
	updated := bson.M{
		"$set": bson.M{
			"id":       data.Oid,
			"table":    data.Table,
			"schema":   data.Schema,
			"database": data.DataBase,
			"offset":   data.Offset,
			"isactive": data.IsActive,
		},
	}
	updateRes, err := collection.UpdateOne(mdb.ctx, filter, updated)
	if err != nil {
		log.Error("updateData error: ", err)
		return err
	}
	if updateRes.MatchedCount == 0 {
		log.Error("updateData error: ", "Данные не обновлены")
		return "updateData error"
	}

	log.Println("data updated: ", data.Oid)
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
		log.Error("Insert error: ", err)
		return err, DbObject
	}
	resMess = insertResult.InsertedID

	fmt.Println("Inserted a single document: ", resMess)
	id := getId(fmt.Sprintf("%s", resMess))
	// if err != nil {
	// 	log.Error(err)
	// }
	DbObject = StateMess{
		Oid:      id,
		Table:    data.Table,
		Schema:   data.Schema,
		DataBase: data.DataBase,
		Offset:   data.Offset,
		IsActive: data.IsActive,
	}

	return nil, DbObject
}

// метод для удаления документа
func (mdb *MongoBase) dropData() {

}

// метод для удалении коллекции
func (mdb *MongoBase) dropCollection(collection *mongo.Collection) interface{} {

	// проверка подключения
	if err := mdb.checkConn(); err != nil {
		return err
	}

	err := collection.Drop(mdb.ctx)
	if err != nil {
		log.Error("Ошибка удаления коллекции: ", err)
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
		var bitem bson.M
		err := cursor.Decode(&bitem)
		if err != nil {
			log.Error(err)
			return res, err
		}
		err = cursor.Decode(&elem)
		if err != nil {
			return res, err
		}
		oid, _ := bitem["_id"]
		elem.Oid = getId(fmt.Sprintf("%s", oid))
		res = append(res, &elem)
	}
	return res, nil
}

// метод для проверки подключения
func (mdb *MongoBase) checkConn() interface{} {
	if err := mdb.client.Ping(mdb.ctx, readpref.Primary()); err != nil {
		log.Error(fmt.Sprintf("Подключние по url %s, не состоялось: ", mdb.url), err)
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
		log.Error("Ошибка подключения к MongoDB: ", err)
		return err
	}

	if err := mdb.checkConn(); err != nil {
		return err
	}

	return nil
}

func getId(oid string) string {
	re := regexp.MustCompile(`"([a-fA-F0-9]{24})"`)
	match := re.FindStringSubmatch(oid)
	id := fmt.Sprintf("%s", match[1])
	return id
}
