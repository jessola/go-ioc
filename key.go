package ioc

import (
	"fmt"
	"reflect"
)

type Key struct {
	Type reflect.Type
	Name string
}

func NamedKey[T any](name string) Key {
	return Key{Type: reflect.TypeFor[T](), Name: name}
}

func UnnamedKey[T any]() Key {
	return Key{Type: reflect.TypeFor[T]()}
}

func (k Key) String() string {
	return fmt.Sprintf("%s:%s", k.Type.String(), k.Name)
}
