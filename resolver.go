package ioc

type Resolver func(Container) any

func ScopedFunc(container Container, fn any, keys ...Key) error {
	s := container.NewScope()
	defer s.Destroy()

	var handler Func
	if len(keys) == 0 {
		fn, err := MakeFunc(fn)
		if err != nil {
			return err
		}
		handler = fn
	} else {
		fn, err := MakeFuncFromKeys(fn, keys)
		if err != nil {
			return err
		}
		handler = fn
	}

	handler(s)

	return nil
}
