# gonectr

## Install

```bash
go install github.com/gone-io/gonectr@latest
```


## Usage

## 1. Generate Gone Project helper code, which is used for loading goners to project.

```bash
gonectr generate -h
```

### 2. Generate Gone Project Priest function code, which is planted from `https://github.com/gone-io/gone/tree/feature/1.x/tools/gone`

```bash
gonectr priest -h 
```

### 3. Generate Goner Mock Code for interface in Gone Project.

```bash
gonectr mock -h
```

### 4. Build Gone Project, which can call `gonectr generate ...` to generate gone project helper code first and then call `go build` to build gone project.
```bash
gonectr build -h
```

### 5. Run Gone Project, which can call `gonectr generate ...` to generate gone project helper code first and then call `go run` to run gone project.
```yaml
gonectr run -h
```