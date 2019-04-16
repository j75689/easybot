package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.etcd.io/bbolt"
)

// BoltDB driver using go.etcd.io/bbolt
type BoltDB struct {
	info     *Connection
	instance *bbolt.DB
}

func (db *BoltDB) Save(collection string, key string, data interface{}) (err error) {

	err = db.instance.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		if byteData, err := json.Marshal(data); err == nil {
			err = b.Put([]byte(key), byteData)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return
}

func (db *BoltDB) LoadAllWithFilter(collection string, filter map[string]interface{}, callback func(key string, value interface{})) (err error) {
	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.Seek(nil); v != nil; k, v = c.Next() {
			for filterKey, filterValue := range filter {
				filterString := fmt.Sprintf(`"$v":"$v"`, filterKey, filterValue)
				if bytes.Index(v, []byte(filterString)) > -1 {
					var (
						data interface{}
					)
					if err = json.Unmarshal(v, &data); err == nil {
						callback(string(k), data)
					}
				}
			}
		}

		return err
	})
	return
}

func (db *BoltDB) LoadWithFilter(collection string, filter map[string]interface{}) (data interface{}, err error) {
	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for _, v := c.Seek(nil); v != nil; _, v = c.Next() {
			for filterKey, filterValue := range filter {
				filterString := fmt.Sprintf(`"$v":"$v"`, filterKey, filterValue)
				if bytes.Index(v, []byte(filterString)) > -1 {
					var (
						data interface{}
					)
					if err = json.Unmarshal(v, &data); err == nil {
						break
					}
				}
			}
		}
		if data == nil {
			err = fmt.Errorf("data notfound [%v]", filter)
		}
		return err
	})
	return
}

func (db *BoltDB) LoadAll(collection string, callback func(key string, value interface{})) (err error) {
	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		err = b.ForEach(func(k, v []byte) (err error) {
			var (
				data interface{}
			)
			if err = json.Unmarshal(v, &data); err == nil {
				callback(string(k), data)
			}
			return
		})

		return err
	})
	return
}

func (db *BoltDB) Load(collection string, key string) (data interface{}, err error) {

	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		v := b.Get([]byte(key))
		if v != nil {
			json.Unmarshal(v, &data)
		} else {
			return fmt.Errorf("data [%s] not found.", key)
		}
		return nil
	})

	return
}
func (db *BoltDB) Delete(collection string, key string) (err error) {
	err = db.instance.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (db *BoltDB) Connect(conn *Connection, args ...interface{}) (err error) {
	db.info = conn
	// check directory
	path := conn.Host
	if strings.LastIndex(path, "/") > -1 {
		path = path[0:strings.LastIndex(path, "/")]
	}
	os.MkdirAll(path, 0755)

	db.instance, err = bbolt.Open(conn.Host, 0644, nil)
	return
}
func (db *BoltDB) Close() (err error) {
	return db.instance.Close()
}
