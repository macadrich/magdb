package dynamo

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DB interface with common DynamoDB operations
type DB struct {
	table string
	conn  *dynamodb.DynamoDB
}

// GenerateExpression extract type name of a structure and concantenate string
// to updateExpression in dynamodb
func GenerateExpression(f interface{}, attributes string, omitkey string) string {
	var expression string
	fm := f.(map[string]interface{})
	i := 0
	for k, v := range fm {
		if k != omitkey {
			if len(v.(string)) > 0 {
				if i == 0 {
					expression = expression + attributes + " " + k + " = " + ":" + k + ","
				} else {
					expression = expression + " " + k + " = " + ":" + k + ","
				}
				i++
			}
		}
	}
	expression = expression[:len(expression)-1]
	return expression
}

// AddAttributeValue map request value to dynamodb.AttributeValue
func AddAttributeValue(f interface{}, omitkey string) map[string]*dynamodb.AttributeValue {
	const cl = ":"
	eav := map[string]*dynamodb.AttributeValue{}
	fm := f.(map[string]interface{})
	for k, v := range fm {
		if k != omitkey {
			if len(v.(string)) > 0 {
				eav[cl+k] = &dynamodb.AttributeValue{
					S: aws.String(v.(string)),
				}
			}
		}
	}
	return eav
}

// CreateConnection to dynamodb
func CreateConnection(region string) (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}

// NewDynamoDB create new dynamodb instance
func NewDynamoDB(conn *dynamodb.DynamoDB, table string) *DB {
	return &DB{
		conn: conn, table: table,
	}
}

// query functions for ID
// keyattrib default at index 0
func (db *DB) queryByID(search string, attribute []string, obj interface{}) error {
	var keyattrib string
	if len(attribute) == 1 {
		keyattrib = attribute[0]
	} else {
		return errors.New("invalid attributes")
	}

	result, err := db.conn.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(db.table),
		Key: map[string]*dynamodb.AttributeValue{
			keyattrib: {
				S: aws.String(search),
			},
		},
	})

	if err != nil {
		return err
	}

	if result.Item != nil {
		if err := dynamodbattribute.UnmarshalMap(result.Item, &obj); err != nil {
			return err
		}
	} else {
		return errors.New("user id not found")
	}

	return nil
}

// query functions for email
// indexName should be at index 0
// keyattrib should be at index 1
func (db *DB) queryByEmail(search string, attribute []string, obj interface{}) error {
	var keyattrib string
	var indexName string
	if len(attribute) == 2 {
		indexName = attribute[0]
		keyattrib = attribute[1]
	} else {
		return errors.New("invalid attributes")
	}

	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(db.table),
		IndexName: aws.String(indexName),
		KeyConditions: map[string]*dynamodb.Condition{
			keyattrib: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(search),
					},
				},
			},
		},
	}

	result, err := db.conn.Query(queryInput)
	if err != nil {
		return err
	}

	if len(result.Items) >= 1 {
		if err := dynamodbattribute.UnmarshalMap(result.Items[0], &obj); err != nil {
			return err
		}
	} else {
		return errors.New("user not exist")
	}

	return nil
}

// ID convert normal string to DBSearchByID
func (db *DB) ID(id string) DBSearchByID {
	return DBSearchByID(id)
}

// Email convert normal string to DBSearchByEmail
func (db *DB) Email(email string) DBSearchByEmail {
	return DBSearchByEmail(email)
}

// Fields set series of attributes
func (db *DB) Fields(fields ...string) []string {
	return fields
}

// Get user search by (ID | email)
// @param key search key type (DBSearchByEmail | DBSearchByID)
// @param indexName - dynamodb GSI
// @param attribute - key element
// @param castTo - object to cast
func (db *DB) Get(key interface{}, attribute []string, castTo interface{}) error {
	switch v := key.(type) {
	case DBSearchByEmail:
		err := db.queryByEmail(string(v), attribute, castTo)
		return err
	case DBSearchByID:
		err := db.queryByID(string(v), attribute, castTo)
		return err
	default:
		return errors.New("invalid search key type should be (DBSearchByEmail | DBSearchByID) ")
	}
}

// List collection of users
func (db *DB) List(castTo interface{}) error {
	results, err := db.conn.Scan(&dynamodb.ScanInput{
		TableName: aws.String(db.table),
	})
	if err != nil {
		return err
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &castTo); err != nil {
		return err
	}
	return nil
}

// Add a new user
func (db *DB) Add(item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.table),
	}
	_, err = db.conn.PutItem(input)
	if err != nil {
		return err
	}
	return err
}

// Update update user details
func (db *DB) Update(id string, item interface{}) error {
	attribValue := AddAttributeValue(item, "id")
	expression := GenerateExpression(item, "set", "id")
	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		UpdateExpression:          aws.String(expression),
		ExpressionAttributeValues: attribValue,
		ReturnValues:              aws.String("UPDATED_NEW"),
		TableName:                 aws.String(db.table),
	}
	_, err := db.conn.UpdateItem(input)
	if err != nil {
		return err
	}
	return nil
}
