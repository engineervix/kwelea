---
title: Deployment
description: Publishing your docs to GitHub Pages and other static hosts
---

`kwelea build` writes a complete static site to `output_dir` (default: `site/`). Any static host works — GitHub Pages, Netlify, Cloudflare Pages, S3, or a plain web server.

## GitHub Pages

The recommended approach uses GitHub Actions to build and deploy on every push to `main`.

### 1. Enable GitHub Pages

In your repository, go to **Settings → Pages** and set the source to **GitHub Actions**.

### 2. Create the workflow

Create `.github/workflows/docs.yml`:

```yaml
name: Deploy Docs

on:
  push:
    branches: [main]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true

      - name: Install kwelea
        run: go install github.com/engineervix/kwelea@latest

      - name: Build docs
        run: kwelea build

      - name: Upload Pages artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: site/

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

### 3. Set base_url

Update `kwelea.toml` with your Pages URL:

```toml
[site]
base_url = "https://yourorg.github.io/your-repo"
```

GitHub Pages serves project sites under a subpath (`/your-repo`). kwelea extracts that path from `base_url` and uses it as the prefix for all asset and nav links. Getting this wrong produces broken CSS and JS — set it exactly.

Push to `main` — the workflow builds your docs and publishes them automatically.

## Cloudflare Pages

Install [Wrangler](https://developers.cloudflare.com/workers/wrangler/install-and-update/) and deploy directly:

```bash
kwelea build --base-url https://your-project.pages.dev
wrangler pages deploy site/
```

`--base-url` overrides the `base_url` in `kwelea.toml` for this build. This is useful when the same config targets multiple hosts — for example, `kwelea.toml` points to a GitHub Pages URL but you want a clean Cloudflare build without editing the file.

## Netlify

Drop-in config. Create `netlify.toml` at your project root:

```toml
[build]
command   = "go install github.com/engineervix/kwelea@latest && kwelea build"
publish   = "site"
```

::: info
Netlify's build environment includes Go. The `go install` step fetches kwelea on each build — pin a version if you want reproducible builds:

```
go install github.com/engineervix/kwelea@v0.1.0
```
:::

## Serving locally

`kwelea serve` is a full development server — it's not intended for production. For a quick local preview of the built site, use any static file server:

```bash
kwelea build
cd site && python3 -m http.server 8080
```

## .gitignore

Add the build output to `.gitignore`:

```
site/
```

The `site/` directory is regenerated on every build; committing it is unnecessary.
