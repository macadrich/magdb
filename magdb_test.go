package magdb

import (
	"testing"
)

type User struct {
	ID         string `json:"id"`
	UserName   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	MiddleName string `json:"middlename"`
}

func TestMongoDB(t *testing.T) {
	mdb := NewMagDBMongo("localhost", "username", "password", "database name", "source")
	db, err := mdb.InitMongoDBConnection()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	defer db.Close()
}
