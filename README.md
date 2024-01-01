# `nemo`

English | [中文](./README_zh_CN.md)

`Nemo` is a `Golang` language configuration management tool, similar to Spring Environment. At the same time, the
built-in Binder supports structure binding. Supports loading of **`multi-path`**, **`multi-file`**
, **`multi-environment`** and **`multi-type` **configurations

## 1.`Multi-path`

- `/opt/data`
- `/opt/configs`
- `...`

> High priority

## 2.`Multi-search-path`

- `./resources`
- `./config`
- `./configs`

## 3.`Multi-environment`

- `dev`
- `test`
- `prod`
- `...`

## 4.`Multi-type`

- `yaml`

  - `yml`
- `toml`
- `properties`

## 5.`System Environment`

- `os.Env`

## 6.`Custem Map context`

```go
// env.LoadMap(...)
```

## 7.`Expand`

### 7.1.`Expand`

- `${a.b.c....z}`

### 7.2.`Expand/default`

- `${a.b.c....z:hello}`

## 8.`Operation`

### 8.1.`Set`

#### 8.1.1.`Simple Set`

```go
// env.Set("key", "value")
```

#### 8.1.2.`Nested Set`

```go
// env.Get("a.b.c....z", "value")
// env.NestedGet("a.b.c....z", "value")
```

### 8.2.`Get`

#### 8.2.1.`Simple Get`

```go
// evn.Get("key", "value")
```

#### 8.2.2.`Nested Get`

```go
// evn.Get("a.b.c....z", "value")
// evn.NestedGet("a.b.c....z", "value")
```

### 8.3.`Refresh Context`

- `env.Refresh(...)`

## 9.`ConfigCenter`

- `Nacos`
  - `TODO`

## 10.`Load`

### 10.1.`Support Event`

- `EnvironmentEvent`
- `StandardAnyEvent`
  - `data any`

