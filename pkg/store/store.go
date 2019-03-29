package store

import (
	"fmt"
)

// Connection struct
type Connection struct {
	DBName string
	Host   string
	Port   int
	User   string
	Pass   string
}

// Storage interface
type Storage interface {
	Save(key string, data interface{}) error
	Load(key string) (interface{}, error)
	LoadAll(func(key string, value interface{})) error
	Delete(key string) error
	Connect(conn *Connection, args ...interface{}) error
	Close() error
}

var (
	supports = map[string]Storage{
		"bolt": new(BoltDB),
	}
)

func NewStoreage(driver string, conn *Connection, args ...interface{}) (Storage, error) {

	if constructor := supports[driver]; constructor != nil {
		return constructor, constructor.Connect(conn, args...)
	}
	return nil, fmt.Errorf("not supported driver [%s]", driver)
}
