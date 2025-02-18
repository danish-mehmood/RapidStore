package storage

type Engine interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	List() []string
}

func NewEngine(engineType string) (Engine, error) {
	switch engineType {
	case "btree":
		return NewBTreeEngine(), nil
	default:
		return NewBTreeEngine(), nil
	}
}
