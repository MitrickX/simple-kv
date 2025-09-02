package storage

import "github.com/MitrickX/simple-kv/internal/storage/engine"

type Storage interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Del(key string)
}

func NewStorage(engine engine.Engine) Storage {
	return &storage{
		engine: engine,
	}
}

type storage struct {
	engine engine.Engine
}

func (s *storage) Set(key, value string) {
	s.engine.Set(key, value)
}
func (s *storage) Get(key string) (string, bool) {
	return s.engine.Get(key)
}
func (s *storage) Del(key string) {
	s.engine.Del(key)
}
