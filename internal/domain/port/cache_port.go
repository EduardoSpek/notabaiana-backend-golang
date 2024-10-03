package port

type CachePort interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
}
