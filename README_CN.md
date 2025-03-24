<p align="left">
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# gonectr

`gonectr` 是一个命令行工具，用于创建和管理 [Gone](https://github.com/gone-io/gone) 项目。它提供代码生成、编译和运行 Gone 项目的功能。

## 安装

运行以下命令安装 `gonectr`：

```bash
go install github.com/gone-io/gonectr@latest
```

安装完成后，`gonectr` 将位于 `$GOPATH/bin` 目录下。确保 `$GOPATH/bin` 已加入 `$PATH`，以便全局使用 `gonectr`。

## 使用方法

### 1. 创建新的 Gone 项目

使用以下命令创建一个新的 Gone 项目：

```bash
gonectr create my-gone-project
```

这将创建一个名为 `my-gone-project` 的新 Gone 项目目录。

### 2. 生成 Gone 项目的 Helper Code

`gonectr` 生成的 Helper Code 文件通常以 `.gone.go` 结尾，建议在 `.gitignore` 文件中添加 `*.gone.go`，以避免版本控制这些自动生成的文件。

运行以下命令生成 Helper Code，并了解更多使用说明：

```bash
gonectr generate -h
```

### 3. 生成 Priest 函数代码

此命令生成 Priest 函数代码，功能源自 [gone tools](https://github.com/gone-io/gone/tree/feature/1.x/tools/gone)。

运行以下命令了解更多：

```bash
gonectr priest -h
```

### 4. 为接口生成 Mock 代码

此命令用于生成接口的 Mock 代码，便于测试。它依赖于 `mockgen` 工具。安装 `mockgen` 的命令如下：

```bash
go install go.uber.org/mock/mockgen@latest
```

运行以下命令了解更多：

```bash
gonectr mock -h
```

### 5. 构建 Gone 项目

`build` 命令会首先调用 `gonectr generate ...` 生成所需的 Helper Code，然后使用 `go build` 构建项目。

运行以下命令了解更多：

```bash
gonectr build -h
```

### 6. 运行 Gone 项目

`run` 命令会在运行项目之前，生成必要的 Helper Code（如有需要），然后使用 `go run` 执行项目。

运行以下命令了解更多：

```bash
gonectr run -h
```