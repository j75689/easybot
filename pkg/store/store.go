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
	SaveWithFilter(collection string, data interface{}, filter map[string]interface{}) error
	Save(collection string, data interface{}) error
	LoadAllWithFilter(collection string, filter map[string]interface{}, callback func(id string, value []byte)) error
	LoadWithFilter(collection string, filter map[string]interface{}) ([]byte, error)
	Load(collection string, id string) ([]byte, error)
	LoadAll(collection string, callback func(id string, value []byte)) error
	DeleteWithFilter(collection string, filter map[string]interface{}) error
	Delete(collection string, id string) error
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
