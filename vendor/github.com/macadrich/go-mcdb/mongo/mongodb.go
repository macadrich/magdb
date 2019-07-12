package mongo

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DB -
type DB struct {
	Conn         *mgo.Session
	DatabaseName string
	TableName    string
	Collection   map[string]*mgo.Collection
}

// NewMongoDB initialize mongodb with credential and source
func NewMongoDB(host, username, password, database, source string) (*DB, error) {
	Host := []string{host}
	conn, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    Host,
		Username: username,
		Password: password,
		Database: database,
		Source:   source,
	})
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &DB{
		Conn:         conn,
		DatabaseName: database,
		Collection:   make(map[string]*mgo.Collection),
	}, nil
}

// Close mongo database.
func (db *DB) Close() {
	db.Conn.Close()
}

// AddTable add table in mongo database
func (db *DB) AddTable(tableName string) {
	db.Collection[tableName] = db.Conn.DB(db.DatabaseName).C(tableName)
}

// Table query specific table
func (db *DB) Table(name string) *DB {
	db.TableName = name
	return db
}

// Delete delete user in database
func (db *DB) Delete(key string) error {
	return db.Collection[db.TableName].Remove(bson.D{{Name: key, Value: key}})
}

// Get get user using by id
// @param id search key
// @param user object to cast
func (db *DB) Get(key string, item interface{}) error {
	if err := db.Collection[db.TableName].Find(bson.D{{Name: key, Value: key}}).One(item); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// Add saves a given user, assigning it a new ID.
func (db *DB) Add(item interface{}) (id string, err error) {
	if err := db.Collection[db.TableName].Insert(item); err != nil {
		return "", errors.New(err.Error())
	}
	return id, nil
}

// Update update fields in database
func (db *DB) Update(key string, item interface{}) error {
	return db.Collection[db.TableName].Update(bson.D{{Name: key, Value: key}}, item)
}

// List view list of data
func (db *DB) List(items interface{}) error {
	if err := db.Collection[db.TableName].Find(nil).All(items); err != nil {
		return err
	}
	return nil
}
