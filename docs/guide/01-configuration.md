---
title: Configuration
description: kwelea.toml reference — every option explained
---

All configuration lives in `kwelea.toml` at your project root. `kwelea init` creates this file with sensible defaults.

## Full reference

```toml
[site]
title    = "mylib"
version  = "v1.4.2"
base_url = "https://yourorg.github.io/mylib"
repo     = "https://github.com/yourorg/mylib"

[build]
docs_dir   = "docs"
output_dir = "site"

[serve]
port         = 4000
open_browser = true

[theme]
light_code_theme = "github"
dark_code_theme  = "github-dark"

# Optional — omit entirely for automatic filesystem-based nav
[[nav]]
section = "Introduction"
pages   = ["index.md", "01-getting-started.md", "02-installation.md"]

[[nav]]
section = "Core Concepts"
pages   = ["configuration.md", "routing.md"]
```

## [site]

| Key | Default | Description |
|-----|---------|-------------|
| `title` | `""` | Site name, shown in the header and `<title>` |
| `version` | `""` | Optional version badge shown next to the title |
| `base_url` | `""` | Canonical URL — used in sitemap and meta tags |
| `repo` | `""` | GitHub URL — renders the GitHub icon in the header |

## [build]

| Key | Default | Description |
|-----|---------|-------------|
| `docs_dir` | `"docs"` | Directory containing your Markdown files |
| `output_dir` | `"site"` | Directory to write the built site |

## [serve]

| Key | Default | Description |
|-----|---------|-------------|
| `port` | `4000` | Dev server port. If occupied, kwelea auto-increments up to `port+10` |
| `open_browser` | `true` | Open the browser automatically on `kwelea serve` |

::: tip
Set `open_browser = false` if you run `kwelea serve` in a headless environment or prefer to open the browser yourself.
:::

## [theme]

| Key | Default | Description |
|-----|---------|-------------|
| `light_code_theme` | `"github"` | Chroma theme for light mode |
| `dark_code_theme` | `"github-dark"` | Chroma theme for dark mode |

Any [Chroma style name](https://xyproto.github.io/splash/docs/) is valid. The CSS is generated at build time and embedded in the output — no runtime dependency.

## [[nav]] (optional)

When `[[nav]]` sections are present, they override the automatic filesystem nav entirely.

```toml
[[nav]]
section = "Introduction"
pages   = ["index.md", "01-getting-started.md"]

[[nav]]
section = "Reference"
pages   = ["guide/configuration.md", "guide/navigation.md"]
```

Each `pages` entry is a file path relative to `docs_dir`. Pages not listed in any `[[nav]]` section are excluded from the build.

See [Navigation](navigation/) for the full comparison of auto vs manual mode.

## Per-page frontmatter

Each Markdown file can include optional YAML frontmatter:

```markdown
---
title: Getting Started
description: Up and running in five minutes
draft: false
---
```

| Key | Description |
|-----|-------------|
| `title` | Overrides the display name derived from the filename |
| `description` | Used in `<meta name="description">` and search results |
| `draft` | Set `true` to exclude a page from builds and the search index |
