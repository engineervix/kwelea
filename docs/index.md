---
title: Kwelea
description: A fast, weaving documentation generator for Go
---

Kwelea turns a `docs/` folder of Markdown files and a `kwelea.toml` config into a complete static documentation site.

Install it once globally. Consuming projects gain zero new dependencies — just a `docs/` folder.

```bash
go install github.com/engineervix/kwelea@latest
```

## What you get

- Three-column layout: sidebar nav, main content, in-page table of contents
- Full-text search, built in — no external service
- Light and dark mode, configurable Chroma code themes
- Live reload during writing
- Syntax highlighting, admonitions, and D2 diagrams in Markdown
- A single self-contained binary; no Node.js, no Python, no CDN calls

## Where to go next

- [Getting Started](getting-started/) — up and running in five minutes
- [Configuration Reference](guide/configuration/) — every `kwelea.toml` option
- [Markdown Extensions](guide/markdown/) — admonitions, diagrams, code highlighting
- [Deployment](guide/deployment/) — publish to GitHub Pages
