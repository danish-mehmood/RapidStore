package operations

import "github.com/danish-mehmood/RapidStore/storage"

type Operations struct {
	store storage.Engine
}

func NewOperations(s storage.Engine) *Operations {
	return &Operations{store: s}
}

func (o *Operations) Get(key string) (string, error) {
	return o.store.Get(key)
}

func (o *Operations) set(key, value string) error {
	return o.store.Set(key, value)
}

func (o *Operations) Delete(key string) error {
	return o.store.Delete(key)
}

func (o *Operations) List() []string {
	return o.store.List()
}
