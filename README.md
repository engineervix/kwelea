# Kwelea

A documentation generator for Go projects. Write Markdown in `docs/`, run one command, get a beautiful static site.

Install it once globally — consuming projects gain zero new dependencies.

```bash
go install github.com/engineervix/kwelea@latest
```

## Quickstart

```bash
# In your Go project root:
kwelea init    # create docs/ and kwelea.toml
kwelea serve   # dev server at http://localhost:4000
kwelea build   # build the static site to site/
```

**[Full documentation →](https://engineervix.github.io/kwelea)**

## Features

- **Markdown** — syntax highlighting (Chroma), admonitions (`::: warning`), D2 diagrams
- **Search** — full-text via FlexSearch, built into the binary, no external service
- **Live reload** — WebSocket-based dev server, zero config
- **Self-contained** — all CSS, JS, and fonts embedded in the binary; no CDN calls in the output site
- **Navigation** — auto-ordered from filesystem (`01-intro.md` → "Intro") or manual `[[nav]]` in config
- **Themes** — light/dark mode, configurable Chroma code themes

## Configuration

`kwelea.toml` at your project root (created by `kwelea init`):

```toml
[site]
title    = "mylib"
version  = "v1.0.0"
base_url = "https://yourorg.github.io/mylib"
repo     = "https://github.com/yourorg/mylib"

[build]
docs_dir   = "docs"
output_dir = "site"
```

See the [configuration reference](https://engineervix.github.io/kwelea/guide/configuration/) for all options.

## Publishing to GitHub Pages

Create `.github/workflows/docs.yml`:

```yaml
name: Deploy Docs
on:
  push:
    branches: [main]
permissions:
  contents: read
  pages: write
  id-token: write
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: go install github.com/engineervix/kwelea@latest
      - run: kwelea build
      - uses: actions/upload-pages-artifact@v3
        with:
          path: site/
  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/deploy-pages@v4
```

## Named after

The [Quelea](https://en.wikipedia.org/wiki/Quelea) — a genus of small weaver birds. Kwelea *weaves* Markdown, templates, and assets into documentation sites.

## Licence

MIT — see [LICENSE](LICENSE).
