// Package ioc provides a lightweight dependency injection container, with request scoping.
package ioc

import "reflect"

var (
	DefaultCollection Collection
)

func init() {
	DefaultCollection = NewCollection()
}

type Key any
type Service any

type Collection interface {
	Get(Key) (Service, error)
	MustGet(Key) Service
	AddSingleton(k Key, fn any) error
	AddScoped(k Key, fn any) error
	AddTransient(k Key, fn any) error
	NewScope() Scope
	Destroy() error
}

type Scope interface {
	Get(Key) (Service, error)
	MustGet(Key) Service
	AddScoped(k Key, fn any) error
	AddTransient(k Key, fn any) error
	NewScope() Scope
	Destroy() error
}

func NewCollection() Collection {
	col := &collection{
		cache:     map[Key]Service{},
		resolvers: map[Key]resolverFunc{},
	}
	col.AddScoped(reflect.TypeFor[Collection](), func() Collection { return col })
	return col
}

// C() returns the `DefaultCollection`
func C() Collection {
	return DefaultCollection
}

func AddSingleton[T Key](c Collection, fn any) error {
	return c.AddSingleton(reflect.TypeFor[T](), fn)
}

func AddScoped[T Key](s Scope, fn any) error {
	return s.AddScoped(reflect.TypeFor[T](), fn)
}

func AddTransient[T Key](s Scope, fn any) error {
	return s.AddTransient(reflect.TypeFor[T](), fn)
}

func Get[T Key](s Scope) (T, error) {
	svc, err := s.Get(reflect.TypeFor[T]())
	if err != nil {
		return (svc).(T), err
	}
	return svc.(T), nil
}

func MustGet[T Key](s Scope) T {
	return s.MustGet(reflect.TypeFor[T]()).(T)
}

func NewScope() Scope {
	return DefaultCollection.NewScope()
}

func MakeFunc(s Scope, fn any) func() []reflect.Value {
	return func() []reflect.Value {
		var deps []reflect.Value
		t := reflect.TypeOf(fn)

		for i := 0; i < t.NumIn(); i++ {
			deps = append(deps, reflect.ValueOf(s.MustGet(t.In(i))))
		}

		return reflect.ValueOf(fn).Call(deps)
	}
}
