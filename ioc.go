// Package ioc provides a reflection based dependency injection container with
// Transient, Singleton and Scoped dependency lifetimes.
package ioc

var DefaultContainer Container

func init() {
	DefaultContainer = NewContainer()
}
