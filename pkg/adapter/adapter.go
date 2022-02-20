package adapter

import (
	"sync"

	"github.com/elliotchance/orderedmap"
)

type Adapter struct {
	mu sync.RWMutex

	storage *orderedmap.OrderedMap
}

func New() *Adapter {
	return &Adapter{
		storage: orderedmap.NewOrderedMap(),
	}
}

func (a *Adapter) Add(key, value string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.storage.Set(key, value)
}

func (a *Adapter) Del(key string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.storage.Delete(key)
}

func (a *Adapter) Get(key string) (string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	element, ok := a.storage.Get(key)
	if !ok {
		return "", false
	}

	return element.(string), true
}

func (a *Adapter) Range(callback func(key, value string) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for element := a.storage.Front(); element != nil; element = element.Next() {
		if !callback(element.Key.(string), element.Value.(string)) {
			return
		}
	}
}
