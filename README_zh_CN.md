# `nemo`

[English](./README.md) | 中文

`nemo` 是一个面向 Go 的轻量配置运行时。它参考了 Spring Environment 的一些思想，但目标不是做一个“很重的框架”，而是做一个“强但无感”的原生能力。

它强调：

- 能力足够强
- 使用足够简单
- 行为足够稳定
- 尽量降低心智负担

## 核心能力

`nemo` 当前支持：

- 多搜索路径
- 多配置名
- 多配置格式
- Active Profile
- 系统环境变量加载
- 内存属性源
- 嵌套属性读写
- 结构体绑定
- 刷新生命周期

支持的配置格式：

- `ini`
- `properties`
- `toml`
- `yaml`
- `yml`

## 推荐入口

推荐直接使用根包：

```go
import nemo "github.com/photowey/nemo"
```

常用公开 API：

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

## 快速开始

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

## 绑定标签

Binder 当前支持：

- ``binder:"field"``：字段映射
- ``required:"true"``：必须项
- ``default:"value"``：默认值

目前支持的目标类型包括：

- 基础标量类型
- `time.Duration`
- 切片，例如 `[]string`、`[]int`
- 指针字段
- 嵌套指针结构体

## 配置来源

`nemo` 会按优先级合并属性源。

常见来源：

- 显式绝对文件路径
- 搜索路径 + 配置名
- 内存 `Map`
- 系统环境变量

## Profile

Profile 会参与候选配置文件的推导。例如 profile 为 `dev` 时，可能会考虑：

- `application.yml`
- `application-dev.yml`

## 错误模型

绑定错误支持结构化分类，例如：

- invalid target
- invalid tag
- missing required field
- unsupported type
- conversion failure

这让 `nemo` 很适合作为更大框架中的配置基础能力。

## 当前状态

已经实现并测试：

- 环境生命周期
- profile 感知加载
- 多格式 loader
- 嵌套绑定
- required/default 语义
- 根包公共 API

规划中但尚未完整实现：

- `${a.b.c}` 形式的占位符展开
- Nacos 等远程配置中心支持
