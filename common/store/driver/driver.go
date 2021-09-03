package driver

//Driver store driver methods
type Driver interface {
	//Set create/update new key/value
	Set(key string, value interface{}) error
	//Get get value by given key
	Get(key string) (interface{}, error)
}
