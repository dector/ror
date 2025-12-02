package io

import "os"

type FileSystem interface {
	ReadFile(name string) ([]byte, error)
	Stat(name string) (os.FileInfo, error)
	IsNotExist(err error) bool
}

type RealFiles struct{}

func (self *RealFiles) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (self *RealFiles) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (self *RealFiles) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
