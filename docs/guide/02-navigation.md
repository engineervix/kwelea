---
title: Navigation
description: Automatic filesystem nav vs manual [[nav]] sections
---

Kwelea supports two navigation modes. Choose the one that fits your project.

## Automatic mode

When your `kwelea.toml` has no `[[nav]]` sections, kwelea walks `docs_dir` and builds the nav tree automatically.

**Ordering** is controlled by numeric filename prefixes:

```
docs/
├── index.md              → "Home"        /
├── 01-getting-started.md → "Getting Started"  /getting-started/
├── 02-installation.md    → "Installation"     /installation/
└── guide/
    ├── index.md          → "Guide"       /guide/
    └── 01-configuration.md → "Configuration" /guide/configuration/
```

Kwelea strips the numeric prefix (`01-`) when computing the display name and URL slug. `01-getting-started.md` becomes title "Getting Started" at path `/getting-started/`.

::: info
Subdirectories become nav sections. The directory name is title-cased to form the section label (e.g. `guide/` → "Guide"). A `index.md` inside a subdirectory is the section's landing page.
:::

**Title overrides** — set `title` in frontmatter to override the filename-derived name:

```markdown
---
title: Quick Start
---
```

## Manual mode

Add one or more `[[nav]]` sections to `kwelea.toml` to take full control:

```toml
[[nav]]
section = "Introduction"
pages   = ["index.md", "01-getting-started.md", "02-installation.md"]

[[nav]]
section = "Reference"
pages   = ["guide/configuration.md", "guide/navigation.md"]
```

- The `section` field sets the sidebar section label.
- The `pages` array lists file paths relative to `docs_dir`, in display order.
- Pages not listed in any section are **excluded** from the build entirely.

::: warning
Once any `[[nav]]` section exists in `kwelea.toml`, automatic mode is disabled. All pages you want to include must be listed explicitly.
:::

## Draft pages

Set `draft: true` in frontmatter to exclude a page from builds, nav, and the search index:

```markdown
---
title: Upcoming Feature
draft: true
---
```

Draft pages are never written to `output_dir`, even if they appear in `[[nav]]`.

## Prev / Next

Kwelea derives a flat ordered list of all non-draft pages and wires up previous/next links automatically. The order follows either the filesystem walk (auto mode) or the `[[nav]]` declaration order (manual mode).
