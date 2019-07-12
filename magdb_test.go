package magdb

import (
	"testing"
)

func TestMongoDB(t *testing.T) {
	mdb := NewMagDBMongo("localhost", "username", "password", "database name", "source")
	db, err := mdb.InitMongoDBConnection()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	defer db.Close()
}
