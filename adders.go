package ioc

import (
	"fmt"
	"reflect"
)

func Value[T any](v T) func() T {
	return func() T {
		return v
	}
}

type ResolverOpts struct {
	Deps []Key
	Fn   any
}

func MakeResolver(opts ResolverOpts) Resolver {
	return func(c Container) any {
		var deps []reflect.Value
		fn := reflect.ValueOf(opts.Fn)

		for _, key := range opts.Deps {
			svc, err := c.Get(key)
			if err != nil {
				panic(err)
			}
			deps = append(deps, reflect.ValueOf(svc))
		}

		out := fn.Call(deps)
		return out[0].Interface()
	}
}

func AddSingleton[T any](c Container, resolver any) error {
	return add[T](c, Singleton, resolver)
}
func AddScoped[T any](c Container, resolver any) error {
	return add[T](c, Scoped, resolver)
}
func AddTransient[T any](c Container, resolver any) error {
	return add[T](c, Transient, resolver)
}

func AddSingletonNamed[T any](c Container, name string, resolver any) error {
	return addNamed[T](c, name, Singleton, resolver)
}
func AddScopedNamed[T any](c Container, name string, resolver any) error {
	return addNamed[T](c, name, Scoped, resolver)
}
func AddTransientNamed[T any](c Container, name string, resolver any) error {
	return addNamed[T](c, name, Transient, resolver)
}

func ensureFunc(fn any) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("resolver must be a function that returns exactly 1 value")
	}
	return nil
}

func add[T any](c Container, scope Scope, resolver any) error {
	if err := ensureFunc(resolver); err != nil {
		return err
	}

	return c.Add(AddOptions{
		Key:   UnnamedKey[T](),
		Scope: scope,
		Resolver: func(c Container) any {
			var deps []reflect.Value
			t := reflect.TypeOf(resolver)
			fn := reflect.ValueOf(resolver)

			for i := 0; i < t.NumIn(); i++ {
				svc, err := c.Get(Key{Type: t.In(i)})
				if err != nil {
					panic(err)
				}
				deps = append(deps, reflect.ValueOf(svc))
			}

			out := fn.Call(deps)
			return out[0].Interface()
		},
	})
}

func addNamed[T any](c Container, name string, scope Scope, resolver any) error {
	if err := ensureFunc(resolver); err != nil {
		return err
	}

	return c.Add(AddOptions{
		Key:   NamedKey[T](name),
		Scope: scope,
		Resolver: func(c Container) any {
			var deps []reflect.Value
			t := reflect.TypeOf(resolver)
			fn := reflect.ValueOf(resolver)

			for i := 0; i < t.NumIn(); i++ {
				svc, err := c.Get(Key{Type: t.In(i)})
				if err != nil {
					panic(err)
				}
				deps = append(deps, reflect.ValueOf(svc))
			}

			out := fn.Call(deps)
			return out[0].Interface()
		},
	})
}
