package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB driver using go.mongodb.org/mongo-driver/mongo
type MongoDB struct {
	info     *Connection
	instance *mongo.Client
}

func (db *MongoDB) Save(collection string, key string, data interface{}) (err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	filter := bson.M{"id": key}
	if _, err := db.Load(collection, key); err != nil {
		_, err = coll.InsertOne(ctx, data)
	} else {
		update := bson.D{
			{"$set", data},
		}
		_, err = coll.UpdateMany(ctx, filter, update)
	}

	return
}

func (db *MongoDB) LoadAll(collection string, callback func(key string, value interface{})) (err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result map[string]interface{}
		err = cur.Decode(&result)
		if err != nil {
			continue
		}
		key := result["id"]
		if key != nil {
			if id, ok := key.(string); ok {
				callback(id, result)
			}
		}
	}
	err = cur.Err()
	return
}

func (db *MongoDB) Load(collection string, key string) (data interface{}, err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)

	filter := bson.M{"id": key}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var d map[string]interface{}
	err = coll.FindOne(ctx, filter).Decode(&d)

	return d, err
}
func (db *MongoDB) Delete(collection string, key string) (err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)
	filter := bson.M{"id": key}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err = coll.DeleteMany(ctx, filter)
	return
}

func (db *MongoDB) Connect(conn *Connection, args ...interface{}) (err error) {
	db.info = conn

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	url := fmt.Sprintf("%s", db.info.Host)
	if db.info.Port != "" {
		url = fmt.Sprintf("%s:%s", db.info.Port)
	}
	opts := options.Client().ApplyURI(url)

	if db.info.User != "" && db.info.Pass != "" {
		opts = opts.SetAuth(options.Credential{
			Username: db.info.User,
			Password: db.info.Pass,
		})
	}

	db.instance, err = mongo.Connect(ctx, opts)

	return
}
func (db *MongoDB) Close() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	return db.instance.Disconnect(ctx)
}