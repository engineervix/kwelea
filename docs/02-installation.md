---
title: Installation
description: Installing kwelea and verifying the setup
---

## go install (recommended)

```bash
go install github.com/engineervix/kwelea@latest
```

This is the standard way to install any Go CLI tool globally. The binary lands in `$GOPATH/bin` (usually `~/go/bin`). Make sure that directory is on your `PATH`.

## Install a specific version

```bash
go install github.com/engineervix/kwelea@v0.2.0
```

## Build from source

```bash
git clone https://github.com/engineervix/kwelea
cd kwelea
go install .
```

Or to build a local binary without installing:

```bash
go build -o kwelea .
```

## Verify

```bash
kwelea --help
```

Expected output:

```
A fast, weaving documentation generator for Go

Usage:
  kwelea [command]

Available Commands:
  build       Build the documentation site
  help        Help about any command
  init        Scaffold a kwelea.toml and docs/ folder in the current project
  serve       Start the development server

Flags:
  --config string   config file (default: kwelea.toml)
  -h, --help        help for kwelea

Use "kwelea [command] --help" for more information about a command.
```

::: info
Kwelea requires **Go 1.21 or later**. Run `go version` to check.
:::
