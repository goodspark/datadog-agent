# Using other components

Once your component has been created you might have to "link" it to other component such as flare, status page, health,
...

This page document what needs to be done for each of them so you component is fully integrated in the Agent life cycle

## Flare

You need to migrate the code related to your component's domain from `pkg/flare` to your component. Once the migration
of the Agent to components we should be able to delete `pkg/flare`.

### Creating you callback

In order to add data to a flare you need to register a provider.

First create a `flare.go` file in your component (the file name is a convention) and create a `func (c *yourComp) fillFlare(fb flarehelpers.FlareBuilder) error` function.

This function will be called everytime the agent generate a flare, either from the CLI or from the running agent. This
callback will receive a `comp/core/flare/helpers:FlareBuilder`. The `FlareBuilder` is an interface that provides all the
required helpers needed to add data to a flare (add files, copy directory, scrub data, ...).

Example of `flare.go`:

```golang
import (
	yaml "gopkg.in/yaml.v2"

	flarehelpers "github.com/DataDog/datadog-agent/comp/core/flare/helpers"
)

func (c *cfg) fillFlare(fb flarehelpers.FlareBuilder) error {
	fb.AddFileFromFunc(
		"runtime_config_dump.yaml",
		func () ([]byte, error) {
			return yaml.Marshal(c.AllSettings())
		},
	)

	fb.CopyFile("/etc/datadog-agent/datadog.yaml")
	return nil
}
```

Read the package documentation for `FlareBuilder` for more information. But note that all errors will automatically be
added to a log file shipped within the flare. You should try to ship as much data as possible in a flare instead of
stopping at the first error. Returning an error will not stop the flare from being created nor sent.

### Migrating your code

The code in `pkg/flare` already use the `FlareBuilder`, so migration should be simple. You need to locate the code
related to your component domain from `pkg/flare` and move it to your `fillFlare` function.

### Register you callback

Finally, to Register you callback you need to provide a new `comp/core/flare/helpers:Provider`. The function `comp/core/flare/helpers:NewProvider` is there for that.

Example from the `config` component:

In `component.go`:
```golang
// Module defines the fx options for this component.
var Module = fxutil.Component(
	fx.Provide(newConfig),
)
```

In `config.go`:
```golang
import (
	flarehelpers "github.com/DataDog/datadog-agent/comp/core/flare/helpers"
)

type provides struct {
	fx.Out

	// [...]
	FlareProvider flarehelpers.Provider
	// [...]
}

func newConfig(deps dependencies) (provides, error) {
	// [...]
	return provides{
		// [...]
		FlareProvider: flarehelpers.NewProvider(myComponent.fillFlare),
		// [...]
	}, nil
}
```
