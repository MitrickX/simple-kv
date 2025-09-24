package engine

type Engine interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Del(key string)
}

type engine struct {
	kv map[string]string
}

func NewEngine() Engine {
	return &engine{
		kv: make(map[string]string),
	}
}

func (e *engine) Set(key, value string) {
	e.kv[key] = value
}
func (e *engine) Get(key string) (string, bool) {
	val, ok := e.kv[key]
	return val, ok
}

func (e *engine) Del(key string) {
	delete(e.kv, key)
}
