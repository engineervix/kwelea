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
extra_head = """
<meta name="google-site-verification" content="abc123">
"""
extra_footer = """
<script defer src="https://example.com/analytics.js"></script>
"""

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
| `base_url` | `""` | Full URL of the deployed site (e.g. `https://yourorg.github.io/mylib`). The path component is used as a prefix for all asset and nav links — must match the host and subpath where the site is served |
| `repo` | `""` | GitHub URL — renders the GitHub icon in the header |
| `extra_head` | `""` | Verbatim HTML injected into `<head>` — custom fonts, verification tags, etc. |
| `extra_footer` | `""` | Verbatim HTML injected into the page footer — analytics scripts, banners, etc. |

Use triple-quoted TOML strings for multi-line values:

```toml
[site]
extra_head = """
<meta name="google-site-verification" content="abc123">
<link rel="preconnect" href="https://fonts.example.com">
"""

extra_footer = """
<script defer src="https://example.com/analytics.js"></script>
"""
```

::: warning
`extra_head` and `extra_footer` are rendered unescaped. Only include HTML you trust.
:::

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
