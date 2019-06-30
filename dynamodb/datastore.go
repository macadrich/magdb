package dynamodb

// DBSearchByEmail search type for email
type DBSearchByEmail string

// DBSearchByID search type for id
type DBSearchByID string

// ID convert normal string to DBSearchByID
func ID(id string) DBSearchByID {
	return DBSearchByID(id)
}

// Email convert normal string to DBSearchByEmail
func Email(email string) DBSearchByEmail {
	return DBSearchByEmail(email)
}

// FieldAttributes set series of attributes
func FieldAttributes(fields ...string) []string {
	return fields
}

// Datastore method interface
type Datastore interface {
	Get(key interface{}, attribute []string, castTo interface{}) error
	List(castTo interface{}) error
	Add(item interface{}) error
	Update(id string, item interface{}) error
}
