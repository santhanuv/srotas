package contract

type Store interface {
	Set(key string, value any)
	Get(key string) (any, bool)
}
