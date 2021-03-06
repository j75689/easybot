package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/j75689/easybot/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB driver using go.mongodb.org/mongo-driver/mongo
type MongoDB struct {
	info     *Connection
	instance *mongo.Client
}

func (db *MongoDB) SaveWithFilter(collection string, data interface{}, filter map[string]interface{}) (err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var updated = 0

	update := bson.D{
		{"$set", data},
	}
	v := util.ReflectFieldValue(data, "ID")
	db.LoadAllWithFilter(collection, filter, func(id string, value []byte) {
		v.SetString(id)
		var result *mongo.UpdateResult
		result, err = coll.UpdateMany(ctx, filter, update)

		updated += int(result.ModifiedCount)
	})

	// New One
	if updated == 0 {
		err = db.Save(collection, data)
	}
	return
}

func (db *MongoDB) Save(collection string, data interface{}) (err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	v := util.ReflectFieldValue(data, "ID")
	id := v.String()

	if id == "" {
		v.SetString(primitive.NewObjectID().Hex())
		_, err = coll.InsertOne(ctx, data)
	} else {
		filter := bson.M{"_id": id}
		update := bson.D{
			{"$set", data},
		}
		_, err = coll.UpdateMany(ctx, filter, update)
	}

	return
}

func (db *MongoDB) LoadAllWithFilter(collection string, filter map[string]interface{}, callback func(id string, value []byte)) (err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := coll.Find(ctx, bson.M(filter))
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
		objectID := result["_id"]
		result["id"] = objectID
		if objectID != nil && result != nil {
			data, err := json.Marshal(result)
			if err != nil {
				continue
			}
			callback(objectID.(string), data)
		}
	}
	err = cur.Err()
	return
}

func (db *MongoDB) LoadWithFilter(collection string, filter map[string]interface{}) (data []byte, err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var d map[string]interface{}
	err = coll.FindOne(ctx, filter).Decode(&d)
	if err != nil {
		return
	}
	d["id"] = d["_id"]
	data, err = json.Marshal(d)
	return
}

func (db *MongoDB) LoadAll(collection string, callback func(id string, value []byte)) (err error) {
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
		objectID := result["_id"]
		result["id"] = objectID
		if objectID != nil && result != nil {
			data, err := json.Marshal(result)
			if err != nil {
				continue
			}
			callback(objectID.(string), data)
		}
	}
	err = cur.Err()
	return
}

func (db *MongoDB) Load(collection string, id string) (data []byte, err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)

	filter := bson.M{"_id": id}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var d map[string]interface{}
	err = coll.FindOne(ctx, filter).Decode(&d)
	if err != nil {
		return
	}
	d["id"] = d["_id"]
	data, err = json.Marshal(d)
	return
}

func (db *MongoDB) DeleteWithFilter(collection string, filter map[string]interface{}) (err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err = coll.DeleteMany(ctx, filter)
	return
}

func (db *MongoDB) Delete(collection string, id string) (err error) {
	coll := db.instance.Database(db.info.DBName).Collection(collection)
	filter := bson.M{"_id": id}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err = coll.DeleteMany(ctx, filter)
	return
}

func (db *MongoDB) Connect(conn *Connection, args ...interface{}) (err error) {
	db.info = conn

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	url := db.info.Host
	if db.info.Port != "" {
		url = fmt.Sprintf("%s:%s", db.info.Host, db.info.Port)
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
