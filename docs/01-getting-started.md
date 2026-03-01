---
title: Getting Started
description: Up and running in five minutes
---

## Prerequisites

- Any project with a `docs/` folder — no Go toolchain required

## Install

```bash
# With Go:
go install github.com/engineervix/kwelea@latest

# Without Go — download the binary for your platform from GitHub Releases:
# https://github.com/engineervix/kwelea/releases/latest
```

See the [Installation page](../installation/) for platform-specific download commands.

Verify:

```bash
kwelea --version
```

## Scaffold your docs

In your project root, run:

```bash
kwelea init
```

This creates:

```
your-project/
├── docs/
│   └── index.md     ← starter homepage
└── kwelea.toml      ← config file
```

Open `kwelea.toml` and fill in `base_url` and `repo`:

```toml
[site]
title    = "your-project"
version  = "v0.1.0"
base_url = "https://yourorg.github.io/your-project"
repo     = "https://github.com/yourorg/your-project"
```

## Write docs

Add Markdown files to `docs/`. Prefix filenames with numbers to control order:

```
docs/
├── index.md              ← homepage
├── 01-getting-started.md
├── 02-installation.md
└── guide/
    ├── index.md          ← section landing page
    └── 01-configuration.md
```

Each file can have optional frontmatter:

```markdown
---
title: Getting Started
description: Up and running in five minutes
draft: false
---

Your content here.
```

::: tip
Omit the `# Title` heading from your Markdown body. Kwelea renders the title from frontmatter — if you include both, it will appear twice.
:::

## Start the dev server

```bash
kwelea serve
```

The browser opens automatically at `http://localhost:4000`. Any file save triggers a full rebuild and live reload.

## Build for production

```bash
kwelea build
```

This writes a complete static site to `site/`. The output is ready to deploy to any static host — GitHub Pages, Netlify, S3, a plain web server.
