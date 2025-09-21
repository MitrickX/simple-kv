package os

import "os"

type File interface {
	Write(b []byte) (n int, err error)
	Sync() error
	Close() error
	Name() string
}

type OS interface {
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
}

type osImpl struct{}

func NewOS() OS {
	return &osImpl{}
}

func (o *osImpl) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(name, flag, perm)
}
