# `nemo`

English | [中文](./README_zh_CN.md)

`nemo` is a lightweight configuration runtime for Go, inspired by Spring Environment but designed to feel like a native Go capability.

Its goals are:

- strong capability
- simple usage
- reliable behavior
- low cognitive overhead

## What It Does

`nemo` supports:

- multiple search paths
- multiple config names
- multiple config formats
- active profiles
- system environment ingestion
- in-memory property sources
- nested property lookup and mutation
- struct binding
- refresh lifecycle

Supported formats:

- `ini`
- `properties`
- `toml`
- `yaml`
- `yml`

## Public API

The recommended entrypoint is the root package:

```go
import nemo "github.com/photowey/nemo"
```

Key public API:

- `nemo.New(...)`
- `nemo.Environment`
- `nemo.PropertySource`
- `nemo.WithSearchPaths(...)`
- `nemo.WithProfiles(...)`
- `nemo.WithProperties(...)`
- `nemo.WithThreshold(...)`
- `nemo.Bind[T](env, prefix)`
- `nemo.MustBind[T](env, prefix)`
- `nemo.AsBindError(err)`
- `nemo.IsBindErrorKind(err, kind)`
- `nemo.FormatBindDiagnostic(err)`

## Quick Start

```go
package main

import (
    "fmt"

    nemo "github.com/photowey/nemo"
)

type FeatureConfig struct {
    Enabled bool `binder:"enabled" required:"true"`
    Port    int  `binder:"port" default:"8080"`
}

func main() {
    env := nemo.New()
    err := env.Start(
        nemo.WithSearchPaths("config", "configs"),
        nemo.WithProfiles("dev"),
        nemo.WithProperties(nemo.MixedMap{
            "app": nemo.MixedMap{
                "feature": nemo.MixedMap{
                    "enabled": "true",
                },
            },
        }),
    )
    if err != nil {
        panic(err)
    }

    cfg := nemo.MustBind[FeatureConfig](env, "app.feature")
    fmt.Println(cfg.Enabled, cfg.Port)
}
```

## Binding Tags

`nemo` binder supports:

- ``binder:"field"`` maps a property key
- ``required:"true"`` makes the property mandatory
- ``default:"value"`` provides a fallback

Supported binding targets include:

- scalar values
- `time.Duration`
- slices like `[]string` and `[]int`
- pointer fields
- nested pointer structs

## Property Sources

`nemo` merges property sources in priority order.

Typical source kinds:

- explicit absolute file paths
- search paths plus config names
- in-memory maps
- system environment variables

## Profiles

Profiles participate in config candidate generation. For example, with profile `dev`, `nemo` may consider names such as:

- `application.yml`
- `application-dev.yml`

## Error Model

Binding errors are structured and expose categories such as:

- invalid target
- invalid tag
- missing required field
- unsupported type
- conversion failure

This makes it practical to integrate `nemo` into larger frameworks like `iocgo`.

## Status

Implemented and tested:

- environment lifecycle
- profile-aware loading
- multi-format loaders
- nested binding
- required/default semantics
- typed root package API

Planned, but not fully implemented:

- placeholder expansion like `${a.b.c}`
- remote config center support such as Nacos
