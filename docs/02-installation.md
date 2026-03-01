---
title: Installation
description: Installing kwelea — binary download, go install, or build from source
---

Kwelea is distributed as a single self-contained binary. Pick the method that suits your setup.

## Download binary

No Go toolchain required. Grab the binary for your platform from [GitHub Releases](https://github.com/engineervix/kwelea/releases/latest).

**Linux (x86_64)**

```bash
curl -fsSL https://github.com/engineervix/kwelea/releases/latest/download/kwelea-linux-amd64 \
  -o /usr/local/bin/kwelea && chmod +x /usr/local/bin/kwelea
```

**Linux (ARM64)**

```bash
curl -fsSL https://github.com/engineervix/kwelea/releases/latest/download/kwelea-linux-arm64 \
  -o /usr/local/bin/kwelea && chmod +x /usr/local/bin/kwelea
```

**macOS (Apple Silicon)**

```bash
curl -fsSL https://github.com/engineervix/kwelea/releases/latest/download/kwelea-darwin-arm64 \
  -o /usr/local/bin/kwelea && chmod +x /usr/local/bin/kwelea
```

**macOS (Intel)**

```bash
curl -fsSL https://github.com/engineervix/kwelea/releases/latest/download/kwelea-darwin-amd64 \
  -o /usr/local/bin/kwelea && chmod +x /usr/local/bin/kwelea
```

::: tip macOS: Gatekeeper warning
After downloading, macOS may show:

> _"Apple could not verify 'kwelea' is free of malware that may harm your Mac or compromise your privacy."_

This happens because the binary is not yet notarized with Apple. To clear it, run:

```bash
xattr -d com.apple.quarantine /usr/local/bin/kwelea
```

Or right-click the binary in Finder → **Open** → **Open** to allow it once.

If you have Go installed, `go install` compiles from source on your machine and bypasses Gatekeeper entirely.
:::

**Windows**

Download [`kwelea-windows-amd64.exe`](https://github.com/engineervix/kwelea/releases/latest/download/kwelea-windows-amd64.exe), rename it to `kwelea.exe`, and place it in a directory on your `PATH`.

::: info
You can install kwelea anywhere on your `PATH`, not just `/usr/local/bin`. On Linux, `~/.local/bin` is a common user-local alternative that doesn't require `sudo`.
:::

To pin a specific release, replace `latest` in the URL with a tag — e.g. `.../releases/download/v0.2.0/kwelea-linux-amd64`.

## go install

If your project already uses Go:

```bash
go install github.com/engineervix/kwelea@latest
```

The binary lands in `$GOPATH/bin` (usually `~/go/bin`). Make sure that directory is on your `PATH`.

To pin a specific version:

```bash
go install github.com/engineervix/kwelea@v0.2.0
```

::: info
`go install` requires **Go 1.25 or later**. Run `go version` to check.
:::

## Build from source

```bash
git clone https://github.com/engineervix/kwelea
cd kwelea
go install .
```

Or to build a local binary without installing globally:

```bash
go build -o kwelea .
```

## Verify

```bash
kwelea --help
```

Expected output:

```
Kwelea weaves Markdown, templates, and assets into beautiful,
fast documentation sites for Go projects.

Install once globally, use across all your Go projects.
Zero runtime dependencies in consuming projects.

Usage:
  kwelea [command]

Available Commands:
  build       Build the documentation site
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Scaffold a kwelea.toml and docs/ folder in the current project
  serve       Start the live-reloading development server

Flags:
      --config string   path to kwelea.toml config file (default "kwelea.toml")
  -h, --help            help for kwelea
  -v, --version         version for kwelea

Use "kwelea [command] --help" for more information about a command.
```
