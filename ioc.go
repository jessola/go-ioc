// Package ioc provides a reflection based dependency injection container with
// Transient, Singleton and Scoped dependency lifetimes.
package ioc

import (
	"fmt"
	"reflect"
)

var DefaultContainer Container

type Func func(Container) []reflect.Value

func init() {
	DefaultContainer = NewContainer()
}

func Service[T any]() func(Container) T {
	return func(c Container) T {
		return MustGet[T](c)
	}
}

func ServiceNamed[T any](name string) func(Container) T {
	return func(c Container) T {
		return MustGetNamed[T](c, name)
	}
}

func MakeFunc(fn any) (Func, error) {
	if err := ensureFunc(fn); err != nil {
		return nil, err
	}
	return func(c Container) []reflect.Value {
		var deps []reflect.Value
		t := reflect.TypeOf(fn)

		for i := 0; i < t.NumIn(); i++ {
			svc, err := c.Get(Key{Type: t.In(i)})
			if err != nil {
				panic(err)
			}
			deps = append(deps, reflect.ValueOf(svc))
		}

		return reflect.ValueOf(fn).Call(deps)
	}, nil
}

func MakeFuncFromKeys(fn any, keys []Key) (Func, error) {
	t := reflect.TypeOf(fn)

	if err := ensureFunc(fn); err != nil {
		return nil, err
	}
	if err := checkArgLength(t, keys); err != nil {
		return nil, err
	}

	return func(c Container) []reflect.Value {
		var deps []reflect.Value

		for _, key := range keys {
			svc, err := c.Get(key)
			if err != nil {
				panic(err)
			}
			deps = append(deps, reflect.ValueOf(svc))
		}

		return reflect.ValueOf(fn).Call(deps)
	}, nil
}

func checkArgLength(fnType reflect.Type, keys []Key) error {
	if len(keys) != fnType.NumIn() {
		return fmt.Errorf("number of keys does not match number of function args")
	}
	return nil
}
