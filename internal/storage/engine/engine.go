package engine

import "sync"

type Engine interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Del(key string)
}

type engine struct {
	mx *sync.RWMutex
	kv map[string]string
}

func NewEngine() Engine {
	return &engine{
		mx: &sync.RWMutex{},
		kv: make(map[string]string),
	}
}

func (e *engine) Set(key, value string) {
	defer e.mx.Unlock()
	e.mx.Lock()
	e.kv[key] = value
}
func (e *engine) Get(key string) (string, bool) {
	defer e.mx.RUnlock()
	e.mx.RLock()
	val, ok := e.kv[key]
	return val, ok
}

func (e *engine) Del(key string) {
	defer e.mx.Unlock()
	e.mx.Lock()
	delete(e.kv, key)
}
