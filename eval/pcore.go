package eval

import (
	"context"
	"github.com/lyraproj/semver/semver"
)

type (
	URI string

	// Pcore is the interface to the evaluator runtime. The runtime is
	// a singleton available using the global variable Puppet.
	Pcore interface {
		// Reset clears all settings and loaders, except the static loaders
		Reset()

		// SystemLoader returns the loader that finds all built-ins. It's parented
		// by a static loader.
		SystemLoader() Loader

		// EnvironmentLoader returns the loader that finds things declared
		// in the environment and its modules. This loader is parented
		// by the SystemLoader
		EnvironmentLoader() Loader

		// Loader returns a loader for module.
		Loader(moduleName string) Loader

		// Logger returns the logger that this instance was created with
		Logger() Logger

		// RootContext returns a new Context that is parented by the context.Background()
		// and is initialized with a loader that is parented by the EnvironmentLoader.
		//
		RootContext() Context

		// Get returns a setting or calls the given defaultProducer
		// function if the setting does not exist
		Get(key string, defaultProducer Producer) Value

		// Set changes a setting
		Set(key string, value Value)

		// SetLogger changes the logger
		SetLogger(Logger)

		// Do executes a given function with an initialized Context instance.
		//
		// The Context will be parented by the Go context returned by context.Background()
		Do(func(Context))

		// DoWithParent executes a given function with an initialized Context instance.
		//
		// The context will be parented by the given Go context
		DoWithParent(context.Context, func(Context))

		// Try executes a given function with an initialized Context instance. If an error occurs,
		// it is caught and returned. The error returned from the given function is returned when
		// no other error is caught.
		//
		// The Context will be parented by the Go context returned by context.Background()
		Try(func(Context) error) error

		// TryWithParent executes a given function with an initialized Context instance.  If an error occurs,
		// it is caught and returned. The error returned from the given function is returned when no other
		// error is caught
		//
		// The context will be parented by the given Go context
		TryWithParent(context.Context, func(Context) error) error

		// DefineSetting defines a new setting with a given valueType and default
		// value.
		DefineSetting(key string, valueType Type, dflt Value)
	}
)

const (
	KeyPcoreUri     = `pcore_uri`
	KeyPcoreVersion = `pcore_version`

	RuntimeNameAuthority = URI(`http://puppet.com/2016.1/runtime`)
	PcoreUri             = URI(`http://puppet.com/2016.1/pcore`)
)

var PcoreVersion, _ = semver.NewVersion3(1, 0, 0, ``, ``)
var ParsablePcoreVersions, _ = semver.ParseVersionRange(`1.x`)

var Puppet Pcore = nil

func GetSetting(name string, dflt Value) Value {
	if Puppet == nil {
		return dflt
	}
	return Puppet.Get(name, func() Value { return dflt })
}