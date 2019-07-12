package magdb

import (
	"github.com/macadrich/go-mcdb/dynamo"
	"github.com/macadrich/go-mcdb/mongo"
)

// DataStoreDynamo inheret to dynamodb
type DataStoreDynamo interface {
	dynamo.Datastore
}

// DataStoreMongo inheret to mongodb
type DataStoreMongo interface {
	mongo.Datastore
}

// MagDB contains credential need to initialize for database
type MagDB struct {
	Region    string `json:"region"`
	TableName string `json:"tablename"`
	Host      string `json:"host"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	Source    string `json:"source"`
}

// NewMagDB initialize new MagDB
func NewMagDB(region string, tableName string) *MagDB {
	return &MagDB{
		Region:    region,
		TableName: tableName,
	}
}

// NewMagDBMongo initialize mongo database with credentials
func NewMagDBMongo(host, username, password, database, source string) *MagDB {
	return &MagDB{
		Host:     host,
		Username: username,
		Password: password,
		Database: database,
		Source:   source,
	}
}

// InitDynamoDBConnection initialize dynamodb connection
// with region and table name
func (mag *MagDB) InitDynamoDBConnection() (*dynamo.DB, error) {
	conn, err := dynamo.CreateConnection(mag.Region)
	if err != nil {
		return nil, err
	}
	return dynamo.NewDynamoDB(conn, mag.TableName), nil
}

// InitMongoDBConnection -
func (mag *MagDB) InitMongoDBConnection() (*mongo.DB, error) {
	conn, err := mongo.NewMongoDB(mag.Host, mag.Username, mag.Password, mag.Database, mag.Source)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
