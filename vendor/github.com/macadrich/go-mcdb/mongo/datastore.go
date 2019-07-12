package mongo

// Datastore mongo database interface
type Datastore interface {
	AddTable(tableName string)
	Table(name string) Datastore
	Delete(key string) error
	Get(key string, item interface{}) error
	Add(item interface{}) (id string, err error)
	Update(item interface{}) error
	List(item interface{}) error
}
