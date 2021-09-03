package driver

//Driver store driver methods
type Driver interface {
	//Set create/update new key/value
	Set(key, value string) error
	//Get get value by given key
	Get(key string) (string, error)
}
