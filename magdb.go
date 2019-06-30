package magdb

import (
	"magdb/dynamodb"
)

// DataStore inheret to dynamodb
type DataStore interface {
	dynamodb.Datastore
}

// MagDB instance
type MagDB struct {
	Region    string `json:"region"`
	TableName string `json:"tablename"`
}

// NewMagDB initialize new MagDB
func NewMagDB(region string, tableName string) *MagDB {
	return &MagDB{
		Region:    region,
		TableName: tableName,
	}
}

// InitDynamoDBConnection initialize dynamodb connection
// with region and table name
func (mag *MagDB) InitDynamoDBConnection() (*dynamodb.DB, error) {
	conn, err := dynamodb.CreateConnection(mag.Region)
	if err != nil {
		return nil, err
	}
	return dynamodb.NewDynamoDB(conn, mag.TableName), nil
}
