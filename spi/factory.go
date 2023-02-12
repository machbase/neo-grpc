package spi

import (
	"errors"
	"fmt"
)

type FactoryFunc func() (Database, error)

var factories = map[string]FactoryFunc{}

func RegisterFactory(name string, f FactoryFunc) {
	factories[name] = f
}

func New() (Database, error) {
	count := len(factories)
	if count > 0 {
		for _, f := range factories {
			return f()
		}
	}
	return nil, errors.New("no database factory found")
}

func NewDatabase(name string) (Database, error) {
	if f, ok := factories[name]; ok {
		return f()
	} else {
		return nil, fmt.Errorf("database factory '%s' not found", name)
	}
}
