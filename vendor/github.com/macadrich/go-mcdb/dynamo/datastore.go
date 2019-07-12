package dynamo

// DBSearchByEmail search type for email
type DBSearchByEmail string

// DBSearchByID search type for id
type DBSearchByID string

// FieldsAttribute -
type FieldsAttribute interface {
	ID(id string) DBSearchByID
	Email(email string) DBSearchByEmail
	Fields(fields ...string) []string
}

// Datastore method interface
type Datastore interface {
	ID(id string) DBSearchByID
	Email(email string) DBSearchByEmail
	Fields(fields ...string) []string

	Get(key interface{}, attribute []string, castTo interface{}) error
	List(castTo interface{}) error
	Add(item interface{}) error
	Update(id string, item interface{}) error
}
