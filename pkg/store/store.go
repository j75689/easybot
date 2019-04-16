package store

import (
	"fmt"
)

// Connection struct
type Connection struct {
	DBName string
	Host   string
	Port   string
	User   string
	Pass   string
}

// Storage interface
type Storage interface {
	Save(collection, key string, data interface{}) error
	LoadAllWithFilter(collection string, filter map[string]interface{}, callback func(key string, value interface{})) error
	LoadWithFilter(collection string, filter map[string]interface{}) (interface{}, error)
	Load(collection, key string) (interface{}, error)
	LoadAll(collection string, callback func(key string, value interface{})) error
	Delete(collection, key string) error
	Connect(conn *Connection, args ...interface{}) error
	Close() error
}

var (
	supports = map[string]Storage{
		"bolt":  new(BoltDB),
		"mongo": new(MongoDB),
	}
)

func NewStoreage(driver string, conn *Connection, args ...interface{}) (Storage, error) {

	if constructor := supports[driver]; constructor != nil {
		return constructor, constructor.Connect(conn, args...)
	}
	return nil, fmt.Errorf("not supported driver [%s]", driver)
}
