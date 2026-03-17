package ioc

type Scope int

const (
	Singleton Scope = iota
	Scoped
	Transient
)
