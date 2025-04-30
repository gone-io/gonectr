<p align="left">
   English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# gonectl
- [gonectl](#gonectl)
  - [Introduction](#introduction)
  - [Installation](#installation)
    - [Method 1: Using go install (Recommended)](#method-1-using-go-install-recommended)
    - [Method 2: Direct Binary Download](#method-2-direct-binary-download)
  - [Feature Overview](#feature-overview)
  - [Detailed Usage Guide](#detailed-usage-guide)
    - [1. create Command: Create Gone Projects from Templates](#1-create-command-create-gone-projects-from-templates)
      - [View Help:](#view-help)
      - [Basic Usage: Create Project](#basic-usage-create-project)
      - [Create Project with Specific Template](#create-project-with-specific-template)
      - [List All Available Templates](#list-all-available-templates)
      - [Create Project with Module Name](#create-project-with-module-name)
      - [Create Project from Remote Git Repository Template](#create-project-from-remote-git-repository-template)
    - [2. install Command: Install Gone Modules and Generate `module.load.go`](#2-install-command-install-gone-modules-and-generate-moduleloadgo)
      - [View Help:](#view-help-1)
      - [Basic Usage: Install Module](#basic-usage-install-module)
      - [Specify LoadFunc](#specify-loadfunc)
      - [Real Example](#real-example)
      - [Uninstall/Modify Module](#uninstallmodify-module)
    - [gone-io Official Modules, Supporting Short Names](#gone-io-official-modules-supporting-short-names)
    - [3. generate Command: Generate `*.gone.go` Files for Gone Projects](#3-generate-command-generate-gonego-files-for-gone-projects)
      - [Functionality](#functionality)
      - [Specify Scan Directory](#specify-scan-directory)
      - [Specify Main Function Directory](#specify-main-function-directory)
      - [Advanced Usage: Generate `import.gone.go` for Non-main Package](#advanced-usage-generate-importgonego-for-non-main-package)
      - [Advanced Usage: Support Multiple Gone Instances](#advanced-usage-support-multiple-gone-instances)
      - [Use with go generate](#use-with-go-generate)
    - [4. mock Command: Generate Mock Code](#4-mock-command-generate-mock-code)
      - [View Help:](#view-help-2)
      - [Basic Usage](#basic-usage)
      - [More Options](#more-options)
    - [5. build Command: Build Gone Projects](#5-build-command-build-gone-projects)
      - [Features](#features)
      - [View Help:](#view-help-3)
      - [Basic Usage](#basic-usage-1)
    - [6. run Command: Run Gone Projects](#6-run-command-run-gone-projects)
      - [Features](#features-1)
      - [View Help:](#view-help-4)
      - [Basic Usage](#basic-usage-2)
  - [FAQ](#faq)
    - [Q: What is the relationship between gonectl and standard Go tools?](#q-what-is-the-relationship-between-gonectl-and-standard-go-tools)
    - [Q: How to upgrade gonectl to the latest version?](#q-how-to-upgrade-gonectl-to-the-latest-version)
    - [Q: Should \*.gone.go files be included in version control?](#q-should-gonego-files-be-included-in-version-control)
  - [More Resources](#more-resources)


> Command-line tool for Gone framework, simplifying project creation, module management, and code generation

## Introduction

`gonectl` is the official command-line tool for the Gone framework, designed to streamline the development process of Gone projects. It provides a series of convenient commands to help developers quickly create projects, manage modules, generate code, and build applications. Whether you're new to Gone or an experienced developer, `gonectl` can significantly improve your development efficiency.

## Installation

### Method 1: Using go install (Recommended)

Run the following command to install `gonectl`:

```bash
go install github.com/gone-io/gonectl@latest
```

After installation, `gonectl` will be located in the `$GOPATH/bin` directory. Make sure this directory is added to your system's `$PATH` environment variable for global access to the `gonectl` command.

> **Tip**: If you're unsure about the location of `$GOPATH`, you can check it by running `go env GOPATH`.

### Method 2: Direct Binary Download

You can also visit the [gonectl/releases](https://github.com/gone-io/gonectl/releases) page to download the latest version binary for your operating system, then:

1. Extract the downloaded file
2. Copy the extracted `gonectl` executable to a directory in your system PATH
3. Ensure the file has execution permissions (on Linux/macOS, you may need to run `chmod +x gonectl`)

## Feature Overview

`gonectl` provides the following core features:

- **Project Creation**: Quickly scaffold Gone project architecture from templates
- **Module Installation**: Integrate Gone modules and automatically generate loading code
- **Code Generation**: Automatically generate necessary Gone framework integration code
- **Mock Generation**: Create Mock implementations for interfaces, facilitating unit testing
- **Build and Run**: Simplify project building and running processes

## Detailed Usage Guide

### 1. create Command: Create Gone Projects from Templates

The `create` command helps you quickly create Gone projects based on preset or custom templates.

#### View Help:
```bash
gonectl create -h
```

#### Basic Usage: Create Project
```bash
gonectl create demo-project
```
This will create a basic Gone project named `demo-project` in the current directory.

#### Create Project with Specific Template
```bash
gonectl create demo-project -t template-name
```

#### List All Available Templates
```bash
gonectl create -ls
```
This command lists all built-in project templates with their brief descriptions.

#### Create Project with Module Name
```bash
gonectl create demo-project -t template-name -m github.com/gone-io/my-module
```
This is particularly useful when creating projects that will be published as public packages.

#### Create Project from Remote Git Repository Template
```bash
gonectl create demo-project -t https://github.com/gone-io/template-v2-web-mysql
```
You can directly use any Git repository that follows the Gone template specification as a project template.

### 2. install Command: Install Gone Modules and Generate `module.load.go`

The `install` command integrates Gone modules into your project and automatically generates the necessary loading code.

> **Gone Module Best Practice**: We recommend each Gone module to provide one or more `gone.LoadFunc` functions, such as:
> ```go
> func Load(gone.Loader) error {
>     // Load related Goners
>     return nil
> }
> ```

#### View Help:
```bash
gonectl install -h
```

#### Basic Usage: Install Module
```bash
gonectl install demo-module
```
This adds `demo-module` to your project and generates the corresponding loading code.

#### Specify LoadFunc
```bash
# Specify LoadA and LoadB functions for generating loading code
gonectl install module LoadA,LoadB
```

#### Real Example
```bash
gonectl install github.com/gone-io/goner/nacos RegistryLoad
```
This installs the nacos module and uses its `RegistryLoad` function for initialization.

#### Uninstall/Modify Module
When executing `gonectl install module` command:
- If the module is not installed, it will be installed
- If already installed, an interactive selection list will be displayed where you can uncheck unwanted LoadFunc to remove them from `module.load.go`

### gone-io Official Modules, Supporting Short Names

```bash
gonectl install goner/nacos
```
> **Note**: For unofficial modules, you need to use the complete Golang module name.

### 3. generate Command: Generate `*.gone.go` Files for Gone Projects

The `generate` command scans project directories and automatically generates integration code files needed by the Gone framework.

#### Functionality

This command will:

1. Scan all packages in specified directories
2. Create `init.gone.go` file for packages containing **Goner** or **LoadFunc**, generating automatic loading code:
   ```go
   func init() {
       gone.
           Loads(Load).  // Load LoadFunc
           Load(&MyGoner{})  // Load Goner
           // ... Load more Goners
   }
   ```
   > **Note**: If a package defines `LoadFunc`, it will only load `LoadFunc` and not directly load Goners, indicating that the user has chosen to manually manage Goners.

3. Create `import.gone.go` file in the main package directory to import all discovered Goner packages:
   ```go
   package main

   import (
       _ "test"
       _ "test/modules/a"
       _ "test/modules/b"
   )
   ```

> **Important**: Do not manually modify `*.gone.go` files, as they will be automatically overwritten by `gonectl`.

#### Specify Scan Directory
```bash
# Can specify multiple directories simultaneously
gonectl generate -s ./test -s ./test2
```

#### Specify Main Function Directory
```bash
gonectl generate -m cmd/server
```

#### Advanced Usage: Generate `import.gone.go` for Non-main Package
```bash
gonectl generate -m for_import --main-package-name for_import
```

#### Advanced Usage: Support Multiple Gone Instances
When using multiple Gone instances in the same program, you can use `--preparer-code` and `--preparer-package` parameters:

```bash
# Goners in gone1 directory use instance-1 instance
gonectl generate -s gone1 --preparer-code 'g.App("instance-1")' --preparer-package 'github.com/gone-io/goner/g'

# Goners in gone2 directory use instance-2 instance
gonectl generate -s gone2 --preparer-code 'g.App("instance-2")' --preparer-package 'github.com/gone-io/goner/g'
```

#### Use with go generate
Create a `generate.go` file in the project root directory and add the following code:
```go
//go:generate gonectl generate -m main-package-dir
```
Then execute `go generate ./...` to automatically run the gonectl command.

### 4. mock Command: Generate Mock Code

The `mock` command generates Mock implementations for interfaces and registers them as Goners, facilitating integration into the Gone framework for testing.

> **Prerequisites**: This feature depends on the `uber mockgen` tool, please install it first:
> ```bash
> go install go.uber.org/mock/mockgen@latest
> ```

#### View Help:
```bash
gonectl mock -h
```

#### Basic Usage
```bash
# Generate Mock implementation for UserService interface in service package
gonectl mock -package service -interfaces UserService
```

#### More Options
```bash
# Generate Mock implementations for multiple interfaces and specify output directory
gonectl mock -package service -interfaces "UserService,OrderService" -output ./mocks
```

### 5. build Command: Build Gone Projects

The `build` command is an enhanced wrapper around the standard `go build`, specifically designed for Gone projects.

#### Features
- Automatically executes `go generate ./...` before compilation to ensure all auxiliary code is updated
- Supports all standard `go build` parameters and options

#### View Help:
```bash
gonectl build -h
```

#### Basic Usage
```bash
# Build Gone project in current directory
gonectl build

# Specify output filename
gonectl build -o myapp

# Use other go build parameters
gonectl build -v -ldflags="-s -w"
```

### 6. run Command: Run Gone Projects

The `run` command is similar to `build`, serving as an enhanced wrapper around `go run`.

#### Features
- Automatically runs `go generate ./...` before execution to update all auxiliary code
- Supports all standard `go run` parameters and options

#### View Help:
```bash
gonectl run -h
```

#### Basic Usage
```bash
# Run Gone project in current directory
gonectl run

# Run specific file
gonectl run main.go

# Run with parameters
gonectl run . -config=dev.yaml
```

## FAQ

### Q: What is the relationship between gonectl and standard Go tools?
A: gonectl is a complement to standard Go tools, specifically designed for the Gone framework. It simplifies Gone-specific code generation and project management processes but still internally calls standard Go commands.

### Q: How to upgrade gonectl to the latest version?
A: Execute `go install github.com/gone-io/gonectl@latest` to update to the latest version.

### Q: Should *.gone.go files be included in version control?
A: It's recommended to include these files in version control as they are part of the project structure. However, they can also be dynamically generated in CI/CD pipelines.

## More Resources

- [Gone Framework Official Documentation](https://github.com/gone-io/gone)
- [Gone Project Templates List](https://github.com/gone-io/goner/tree/main/examples)
- [Issue Feedback](https://github.com/gone-io/gonectl/issues)
