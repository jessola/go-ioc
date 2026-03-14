package ioc

import (
	"fmt"
	"reflect"
)

type resolverFunc func(*collection) Service

type collection struct {
	parent    *collection
	cache     map[Key]Service
	resolvers map[Key]resolverFunc
}

type ErrDuplicateKey struct {
	key Key
}

func createResolverFunc(fn any) resolverFunc {
	// var params []reflect.Type
	// t := reflect.TypeOf(fn)

	// // fmt.Println("VAL------", reflect.ValueOf(fn))

	// for i := 0; i < t.NumIn(); i++ {
	// 	params = append(params, t.In(i))
	// }

	return func(c *collection) Service {
		return MakeFunc(c, fn)()[0].Interface()
		// var deps []reflect.Value

		// // Resolve each dependency in turn
		// for _, p := range params {
		// 	svc := c.MustGet(Key(p))
		// 	deps = append(deps, reflect.ValueOf(svc))
		// }

		// output := reflect.ValueOf(fn).Call(deps)

		// return output[0].Interface()
	}
}

func (e ErrDuplicateKey) Error() string {
	return fmt.Sprintf("duplicate key:%s", e.key)
}

func (c *collection) AddSingleton(k Key, fn any) error {
	// TODO: Check fn is a function
	if _, exists := c.resolvers[k]; exists {
		return ErrDuplicateKey{k}
	}

	resolve := createResolverFunc(fn)

	c.resolvers[k] = func(child *collection) Service {
		if c.cache[k] == nil {
			c.cache[k] = resolve(child)
		}
		return c.cache[k]
	}

	return nil
}

func (c *collection) AddScoped(k Key, fn any) error {
	// TODO: Check fn is a function
	if _, exists := c.resolvers[k]; exists {
		return ErrDuplicateKey{k}
	}

	resolve := createResolverFunc(fn)

	c.resolvers[k] = func(child *collection) Service {
		if child.cache[k] == nil {
			child.cache[k] = resolve(child)
		}
		return child.cache[k]
	}

	return nil
}

func (c *collection) AddTransient(k Key, fn any) error {
	// TODO: Check fn is a function
	if _, exists := c.resolvers[k]; exists {
		return ErrDuplicateKey{k}
	}

	resolve := createResolverFunc(fn)

	c.resolvers[k] = func(child *collection) Service {
		return resolve(child)
	}

	return nil
}

func (c *collection) Get(k Key) (Service, error) {
	// Check cache
	if svc, ok := c.cache[k]; ok {
		return svc, nil
	}

	// Check local resolvers
	if rsv, ok := c.resolvers[k]; ok {
		return rsv(c), nil
	}

	// Check parent resolvers
	if c.parent == nil {
		return nil, fmt.Errorf("no service with key:%s", k)

	}

	if rsv, ok := c.parent.resolvers[k]; ok {
		c.resolvers[k] = rsv
		return rsv(c), nil
	}

	return nil, fmt.Errorf("no service with key:%s", k)
}

func (c *collection) MustGet(k Key) Service {
	svc, err := c.Get(k)
	if err != nil {
		panic(err)
	}
	return svc
}

func (c *collection) NewScope() Scope {
	coll := &collection{
		parent:    c,
		cache:     map[Key]Service{},
		resolvers: map[Key]resolverFunc{},
	}

	coll.AddScoped(reflect.TypeFor[Collection](), func() Collection { return coll })

	return coll
}

func (c *collection) Destroy() error {
	clear(c.cache)
	clear(c.resolvers)
	return nil
}
