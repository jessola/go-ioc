package ioc

import "fmt"

type AddOptions struct {
	Key      Key
	Scope    Scope
	Resolver Resolver
}

type Container interface {
	Add(opts AddOptions) error
	Remove(k Key) error
	Get(k Key) (any, error)
	Destroy() error
	NewScope() Container
}

type container struct {
	parent    *container
	services  map[Key]any
	resolvers map[Key]func(*container) any
}

func NewContainer() Container {
	return newContainer()
}

func newContainer() *container {
	c := &container{
		parent:    nil,
		services:  map[Key]any{},
		resolvers: make(map[Key]func(*container) any),
	}

	c.Add(AddOptions{
		Key:   UnnamedKey[Container](),
		Scope: Singleton,
		Resolver: func(c Container) any {
			return c
		},
	})

	return c
}

func newContainerWithParent(p *container) *container {
	c := newContainer()
	c.parent = p
	return c
}

func (c *container) Add(opts AddOptions) error {
	if _, exists := c.resolvers[opts.Key]; exists {
		return fmt.Errorf("duplicate key %s", opts.Key)
	}

	k := opts.Key

	switch opts.Scope {
	case Singleton:
		c.resolvers[k] = func(child *container) any {
			if _, ok := c.services[k]; !ok {
				c.services[k] = opts.Resolver(child)
			}
			return c.services[k]
		}

	case Scoped:
		c.resolvers[k] = func(child *container) any {
			if _, ok := child.services[k]; !ok {
				child.services[k] = opts.Resolver(child)
			}
			return child.services[k]
		}

	case Transient:
		c.resolvers[k] = func(child *container) any {
			return opts.Resolver(child)
		}
	}

	return nil
}

func (c *container) Get(k Key) (any, error) {
	// Try local cache
	if svc, ok := c.services[k]; ok {
		return svc, nil
	}
	// Try local resolver
	if r, ok := c.resolvers[k]; ok {
		return r(c), nil
	}
	// Try parent resolver
	if c.parent != nil {
		if r, ok := c.parent.resolvers[k]; ok {
			c.resolvers[k] = c.parent.resolvers[k]
			return r(c), nil
		}
	}
	return nil, fmt.Errorf("no service with key %s", k)
}

func (c *container) Remove(k Key) error {
	delete(c.services, k)
	delete(c.resolvers, k)
	return nil
}

func (c *container) Destroy() error {
	clear(c.services)
	clear(c.resolvers)
	return nil
}

func (c *container) NewScope() Container {
	return newContainerWithParent(c)
}

func _(c Container) {
	c.Add(AddOptions{
		Key:   UnnamedKey[string](),
		Scope: Singleton,
		Resolver: func(c Container) any {
			c.Get(UnnamedKey[string]())
			c.Get(NamedKey[int]("foo"))
			return "hello"
		},
	})
}
