package magdb

import (
	"github.com/macadrich/go-mcdb/mcdb"
)

// DataStore inheret to dynamodb
type DataStore interface {
	mcdb.Datastore
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
func (mag *MagDB) InitDynamoDBConnection() (*mcdb.DB, error) {
	conn, err := mcdb.CreateConnection(mag.Region)
	if err != nil {
		return nil, err
	}
	return mcdb.NewDynamoDB(conn, mag.TableName), nil
}
