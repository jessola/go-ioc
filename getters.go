package ioc

func MustGet[T any](c Container) T {
	svc, err := c.Get(UnnamedKey[T]())
	if err != nil {
		panic(err)
	}
	return svc.(T)
}

func MustGetNamed[T any](c Container, name string) T {
	svc, err := c.Get(NamedKey[T](name))
	if err != nil {
		panic(err)
	}
	return svc.(T)
}
