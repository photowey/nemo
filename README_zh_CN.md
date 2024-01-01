# `nemo`

[English](./README.md)  | 中文

`Nemo` 是一款 `Golang` 语言配置管理工具, 设计思想来源于 `Spring` `Environment `。内置的 `Binder`
支持结构体绑定。与此同时,可以通过 **多路径**、**多文件**、**多环境**、**多类型** 等高自由度的配置来实现 `Golang` `App`
的配置加载。

## 1.支持多(绝对)路径

- `/opt/data`
- `/opt/configs`
- `...`

> 高优先级

## 2.支持多搜索路径

- `./resources`
- `./config`
- `./configs`

## 3.支持多环境

- `dev`
- `test`
- `prod`
- `...`

## 4.多类型

- `yaml`

    - `yml`

- `toml`

- `properties`

## 5.支持环境变量

## 6.支持自定义 `Map` 上下文

### 6.1.加载 `Map`

```go
// env.LoadMap(...)
```

## 7.支持变量扩展

### 7.1.变量扩展

- `${a.b.c....z}`

### 7.2.变量默认值

- `${a.b.c....z:hello}`

## 8.变量操作

### 8.1.设置变量

#### 8.1.1.设置单 `Key` 变量

```go
// env.Set("key", "value")
```

#### 8.1.2.设置嵌套变量

```go
// env.Get("a.b.c....z", "value")
// env.NestedGet("a.b.c....z", "value") // 推荐 -> 这样开发者能更明确自己设置的值是否是嵌套 KEY
```

### 8.2.获取变量

#### 8.2.1.获取单 `Key` 变量

```go
// evn.Get("key", "value")
```

#### 8.2.2.获取嵌套

```go
// evn.Get("a.b.c....z", "value")
// evn.NestedGet("a.b.c....z", "value") // 推荐
```

### 8.3.上下文刷新

- `env.Refresh(...)`

### 9.支持配置中心

- `Nacos`
    - `TODO`
    - 思考中
        - 通过环境上下文加载的 **准备事件** `PrepareEvent` 向 `Nacos` 获取配置
        - 构造 `ProertySource`
            - `Map` 类型的 `PropertySource`

## 10.加载

### 10.1.支持事件

- `EnvironmentEvent`
- `StandardAnyEvent`
    - `data any`