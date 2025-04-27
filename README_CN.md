<p align="left">
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# gonectr

- [gonectr](#gonectr)
    - [简介](#简介)
    - [安装](#安装)
        - [方法一：使用 go install（推荐）](#方法一使用-go-install推荐)
        - [方法二：直接下载二进制文件](#方法二直接下载二进制文件)
    - [功能概览](#功能概览)
    - [详细使用指南](#详细使用指南)
        - [1. create 子命令：从模板创建 Gone 项目](#1-create-子命令从模板创建-gone-项目)
            - [查看帮助：](#查看帮助)
            - [基本用法：创建项目](#基本用法创建项目)
            - [使用指定模板创建项目](#使用指定模板创建项目)
            - [查看所有可用模板](#查看所有可用模板)
            - [创建项目时指定模块名](#创建项目时指定模块名)
            - [从远程 Git 仓库模板创建项目](#从远程-git-仓库模板创建项目)
        - [2. install 子命令：安装 Gone 模块，生成 `module.load.go`](#2-install-子命令安装-gone-模块生成-moduleloadgo)
            - [查看帮助：](#查看帮助-1)
            - [基本用法：安装模块](#基本用法安装模块)
            - [指定 LoadFunc](#指定-loadfunc)
            - [实际示例](#实际示例)
            - [卸载/修改模块](#卸载修改模块)
        - [3. generate 子命令：生成 Gone 项目的 `*.gone.go` 文件](#3-generate-子命令生成-gone-项目的-gonego-文件)
            - [功能说明](#功能说明)
            - [指定扫描目录](#指定扫描目录)
            - [指定 main 函数所在目录](#指定-main-函数所在目录)
            - [高级用法：为非 main 包生成 `import.gone.go`](#高级用法为非-main-包生成-importgonego)
            - [高级用法：支持多个 Gone 实例](#高级用法支持多个-gone-实例)
            - [配合 go generate 使用](#配合-go-generate-使用)
        - [4. mock 子命令：生成 Mock 代码](#4-mock-子命令生成-mock-代码)
            - [查看帮助：](#查看帮助-2)
            - [基本用法](#基本用法)
            - [更多选项](#更多选项)
        - [5. build 子命令：构建 Gone 项目](#5-build-子命令构建-gone-项目)
            - [特点](#特点)
            - [查看帮助：](#查看帮助-3)
            - [基本用法](#基本用法-1)
        - [6. run 子命令：运行 Gone 项目](#6-run-子命令运行-gone-项目)
            - [特点](#特点-1)
            - [查看帮助：](#查看帮助-4)
            - [基本用法](#基本用法-2)
    - [常见问题解答](#常见问题解答)
        - [Q: gonectr 与标准 Go 工具的关系是什么？](#q-gonectr-与标准-go-工具的关系是什么)
        - [Q: 如何升级 gonectr 到最新版本？](#q-如何升级-gonectr-到最新版本)
        - [Q: 生成的 \*.gone.go 文件应该纳入版本控制吗？](#q-生成的-gonego-文件应该纳入版本控制吗)
    - [更多资源](#更多资源)



> Gone框架的命令行工具，简化项目创建、模块管理与代码生成

## 简介

`gonectr` 是 Gone 框架的官方命令行工具，旨在简化 Gone 项目的开发流程。它提供了一系列便捷命令，帮助开发者快速创建项目、管理模块、生成代码和构建应用程序。无论您是 Gone 新手还是有经验的开发者，`gonectr` 都能大幅提高您的开发效率。

## 安装

### 方法一：使用 go install（推荐）

运行以下命令安装 `gonectr`：

```bash
go install github.com/gone-io/gonectr@latest
```

安装完成后，`gonectr` 将位于 `$GOPATH/bin` 目录下。请确保该目录已添加到系统环境变量 `$PATH` 中，以便全局使用 `gonectr` 命令。

> **提示**：如果不确定 `$GOPATH` 的位置，可以通过运行 `go env GOPATH` 命令查看。

### 方法二：直接下载二进制文件

您也可以访问 [gonectr/releases](https://github.com/gone-io/gonectr/releases) 页面，下载适合您操作系统的最新版本二进制文件，然后：

1. 解压下载的文件
2. 将解压后的 `gonectr` 可执行文件复制到系统 PATH 路径下的某个目录
3. 确保文件具有执行权限（Linux/macOS 下可能需要运行 `chmod +x gonectr`）

## 功能概览

`gonectr` 提供以下核心功能：

- **创建项目**：从模板快速搭建 Gone 项目架构
- **安装模块**：集成 Gone 模块并自动生成加载代码
- **代码生成**：自动生成必要的 Gone 框架集成代码
- **生成 Mock**：为接口创建 Mock 实现，方便单元测试
- **构建与运行**：简化项目构建和运行过程

## 详细使用指南

### 1. create 子命令：从模板创建 Gone 项目

`create` 命令帮助您快速创建基于预设模板或自定义模板的 Gone 项目。

#### 查看帮助：
```bash
gonectr create -h
```

#### 基本用法：创建项目
```bash
gonectr create demo-project
```
这会在当前目录下创建名为 `demo-project` 的基础 Gone 项目。

#### 使用指定模板创建项目
```bash
gonectr create demo-project -t template-name
```

#### 查看所有可用模板
```bash
gonectr create -ls
```
该命令会列出所有内置的项目模板及其简要描述。

#### 创建项目时指定模块名
```bash
gonectr create demo-project -t template-name -m github.com/gone-io/my-module
```
这对于创建需要发布为公共包的项目特别有用。

#### 从远程 Git 仓库模板创建项目
```bash
gonectr create demo-project -t https://github.com/gone-io/template-v2-web-mysql
```
您可以直接使用任何符合 Gone 模板规范的 Git 仓库作为项目模板。

### 2. install 子命令：安装 Gone 模块，生成 `module.load.go`

`install` 命令用于集成 Gone 模块到您的项目中，自动生成必要的加载代码。

> **Gone 模块最佳实践**：我们建议每个 Gone 模块提供一个或多个 `gone.LoadFunc` 函数，如：
> ```go
> func Load(gone.Loader) error {
>     // 加载相关 Goner
>     return nil
> }
> ```

#### 查看帮助：
```bash
gonectr install -h
```

#### 基本用法：安装模块
```bash
gonectr install demo-module
```
这会添加 `demo-module` 到项目中，并生成相应的加载代码。

#### 指定 LoadFunc
```bash
# 指定使用 LoadA 和 LoadB 函数生成加载代码
gonectr install module LoadA,LoadB
```

#### 实际示例
```bash
gonectr install github.com/gone-io/goner/nacos RegistryLoad
```
这会安装 nacos 模块，并使用其 `RegistryLoad` 函数进行初始化。

#### 卸载/修改模块
执行 `gonectr install module` 命令时：
- 如果模块未安装，会进行安装
- 如果已安装，会显示交互式选择列表，您可以取消勾选不需要的 LoadFunc，将其从 `module.load.go` 中移除

### gone-io官方模块，支持短名称

```bash
gonectr install goner/nacos
```
> **注意**：非官方模块，需要使用完整golang 模块名

### 3. generate 子命令：生成 Gone 项目的 `*.gone.go` 文件

`generate` 命令扫描项目目录，自动生成 Gone 框架需要的集成代码文件。

#### 功能说明

该命令会：

1. 扫描指定目录中的所有包
2. 为包含 **Goner** 或 **LoadFunc** 的包创建 `init.gone.go` 文件，生成自动加载代码：
   ```go
   func init() {
       gone.
           Loads(Load).  // 加载 LoadFunc
           Load(&MyGoner{})  // 加载 Goner
           // ... 加载更多 Goner
   }
   ```
   > **注意**：如果包中定义了 `LoadFunc`，则只会加载 `LoadFunc` 而不会直接加载 Goner，这表示用户选择了手动管理 Goner。

3. 在 main 包目录创建 `import.gone.go` 文件，导入所有发现的 Goner 包：
   ```go
   package main

   import (
       _ "test"
       _ "test/modules/a"
       _ "test/modules/b"
   )
   ```

> **重要提示**：请不要手动修改 `*.gone.go` 文件，这些文件会被 `gonectr` 自动覆盖。

#### 指定扫描目录
```bash
# 可同时指定多个目录
gonectr generate -s ./test -s ./test2
```

#### 指定 main 函数所在目录
```bash
gonectr generate -m cmd/server
```

#### 高级用法：为非 main 包生成 `import.gone.go`
```bash
gonectr generate -m for_import --main-package-name for_import
```

#### 高级用法：支持多个 Gone 实例
在同一程序中使用多个 Gone 实例时，可以使用 `--preparer-code` 和 `--preparer-package` 参数：

```bash
# gone1 目录下的 Goner 使用 instance-1 实例
gonectr generate -s gone1 --preparer-code 'g.App("instance-1")' --preparer-package 'github.com/gone-io/goner/g'

# gone2 目录下的 Goner 使用 instance-2 实例
gonectr generate -s gone2 --preparer-code 'g.App("instance-2")' --preparer-package 'github.com/gone-io/goner/g'
```

#### 配合 go generate 使用
在项目根目录创建 `generate.go` 文件，添加以下代码：
```go
//go:generate gonectr generate -m main-package-dir
```
然后执行 `go generate ./...` 即可自动运行 gonectr 命令。

### 4. mock 子命令：生成 Mock 代码

`mock` 命令用于为接口生成 Mock 实现，并将这些实现注册为 Goner，便于集成到 Gone 框架中进行测试。

> **前提条件**：此功能依赖 `uber mockgen` 工具，请先安装：
> ```bash
> go install go.uber.org/mock/mockgen@latest
> ```

#### 查看帮助：
```bash
gonectr mock -h
```

#### 基本用法
```bash
# 为 service 包中的 UserService 接口生成 Mock 实现
gonectr mock -package service -interfaces UserService
```

#### 更多选项
```bash
# 为多个接口生成 Mock 实现，并指定输出目录
gonectr mock -package service -interfaces "UserService,OrderService" -output ./mocks
```

### 5. build 子命令：构建 Gone 项目

`build` 命令是对标准 `go build` 的增强封装，专为 Gone 项目设计。

#### 特点
- 在编译前自动执行 `go generate ./...`，确保所有辅助代码都已更新
- 支持所有标准 `go build` 的参数和选项

#### 查看帮助：
```bash
gonectr build -h
```

#### 基本用法
```bash
# 构建当前目录下的 Gone 项目
gonectr build

# 指定输出文件名
gonectr build -o myapp

# 使用其他 go build 参数
gonectr build -v -ldflags="-s -w"
```

### 6. run 子命令：运行 Gone 项目

`run` 命令类似于 `build`，是对 `go run` 的增强封装。

#### 特点
- 在执行前自动运行 `go generate ./...`，更新所有辅助代码
- 支持所有标准 `go run` 的参数和选项

#### 查看帮助：
```bash
gonectr run -h
```

#### 基本用法
```bash
# 运行当前目录的 Gone 项目
gonectr run

# 运行指定文件
gonectr run main.go

# 带参数运行
gonectr run . -config=dev.yaml
```

## 常见问题解答

### Q: gonectr 与标准 Go 工具的关系是什么？
A: gonectr 是对标准 Go 工具的补充，专为 Gone 框架设计。它简化了 Gone 特有的代码生成和项目管理流程，但内部仍然调用标准的 Go 命令。

### Q: 如何升级 gonectr 到最新版本？
A: 执行 `go install github.com/gone-io/gonectr@latest` 即可更新到最新版本。

### Q: 生成的 *.gone.go 文件应该纳入版本控制吗？
A: 建议将这些文件纳入版本控制，它们是项目结构的一部分。但也可以在 CI/CD 流程中动态生成。

## 更多资源

- [Gone 框架官方文档](https://github.com/gone-io/gone)
- [Gone 项目模板列表](https://github.com/gone-io/goner/tree/main/examples)
- [问题反馈](https://github.com/gone-io/gonectr/issues)