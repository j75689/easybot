package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/j75689/easybot/pkg/util"

	"go.etcd.io/bbolt"
)

// BoltDB driver using go.etcd.io/bbolt
type BoltDB struct {
	info     *Connection
	instance *bbolt.DB
}

func (db *BoltDB) SaveWithFilter(collection string, data interface{}, filter map[string]interface{}) (err error) {
	var updated = 0
	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		c := b.Cursor()

		for k, v := c.Seek(nil); v != nil; k, v = c.Next() {
			for filterKey, filterValue := range filter {
				filterString := fmt.Sprintf("\"%v\":\"%v\"", filterKey, filterValue)
				if bytes.Index(v, []byte(filterString)) > -1 {
					if byteData, err := json.Marshal(data); err == nil {
						err = b.Put(k, byteData)
						updated++
					}
				}
			}
		}

		return err
	})

	// New One
	if updated == 0 {
		err = db.Save(collection, data)
	}

	return
}

func (db *BoltDB) Save(collection string, data interface{}) (err error) {

	err = db.instance.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}

		v := util.ReflectFieldValue(data, "ID")
		if v.String() == "" {
			id, _ := b.NextSequence()
			v.SetString(strconv.Itoa(int(int64(id))))
		}

		if byteData, err := json.Marshal(data); err == nil {
			err = b.Put([]byte(v.String()), byteData)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return
}

func (db *BoltDB) LoadAllWithFilter(collection string, filter map[string]interface{}, callback func(id string, value []byte)) (err error) {
	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.Seek(nil); v != nil; k, v = c.Next() {
			for filterKey, filterValue := range filter {
				filterString := fmt.Sprintf("\"%v\":\"%v\"", filterKey, filterValue)
				if bytes.Index(v, []byte(filterString)) > -1 {
					callback(string(k), v)
				}
			}
		}

		return err
	})
	return
}

func (db *BoltDB) LoadWithFilter(collection string, filter map[string]interface{}) (data []byte, err error) {
	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for _, v := c.Seek(nil); v != nil; _, v = c.Next() {
			for filterKey, filterValue := range filter {
				filterString := fmt.Sprintf("\"%v\":\"%v\"", filterKey, filterValue)
				if bytes.Index(v, []byte(filterString)) > -1 {
					data = v
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

func (db *BoltDB) LoadAll(collection string, callback func(id string, value []byte)) (err error) {
	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		err = b.ForEach(func(k, v []byte) (err error) {
			callback(string(k), v)
			return
		})

		return err
	})
	return
}

func (db *BoltDB) Load(collection string, id string) (data []byte, err error) {

	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		data = b.Get([]byte(id))
		if data == nil {
			err = fmt.Errorf("data [%v] not found.", id)
		}
		return err
	})

	return
}

func (db *BoltDB) DeleteWithFilter(collection string, filter map[string]interface{}) (err error) {

	err = db.instance.Batch(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.Seek(nil); v != nil; k, v = c.Next() {
			for filterKey, filterValue := range filter {
				filterString := fmt.Sprintf("\"%v\":\"%v\"", filterKey, filterValue)
				if bytes.Index(v, []byte(filterString)) > -1 {
					err = b.Delete(k)
				}
			}
		}

		return err
	})
	return
}

func (db *BoltDB) Delete(collection string, id string) (err error) {
	err = db.instance.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		err = b.Delete([]byte(id))
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
