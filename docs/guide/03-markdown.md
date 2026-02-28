---
title: Markdown Extensions
description: Admonitions, D2 diagrams, and syntax highlighting
---

Kwelea extends standard Markdown with three features: admonitions, D2 diagrams, and Chroma syntax highlighting. All are rendered at build time — no JavaScript required for the output.

## Admonitions

Use `:::` container blocks to draw attention to important content.

**Syntax:**

```markdown
::: info
This is an informational note.
:::

::: tip
A helpful tip for the reader.
:::

::: warning
Something to be careful about.
:::

::: danger
A critical warning — data loss, security risk, etc.
:::

::: details Click to expand
Hidden content shown when the reader expands it.
:::
```

**Available types:** `info`, `tip`, `warning`, `danger`, `details`

The `details` type renders as a native HTML `<details>`/`<summary>` element — collapsible with no JavaScript.

::: info
The title in the admonition label is derived from the type name. You cannot customise it.
:::

::: warning
The `:::` fence must be at the start of the line with no leading whitespace.
:::

::: details How are admonitions implemented?
Admonitions are a custom goldmark block extension (`internal/parser/admonitions.go`). The `:::` fence is parsed into an AST node and rendered to `<div class="admonition admonition-{type}">`.
:::

## D2 Diagrams

Use a fenced code block with the `d2` language tag to embed diagrams.

````markdown
```d2
direction: right
user -> server -> database
```
````

**Output:** Two SVGs are rendered at build time — one for light mode, one for dark mode. The correct one is shown based on the `[data-theme]` attribute on `<html>`, with no JavaScript required.

**Example — a request flow:**

```d2
direction: right
Browser -> "Load Balancer" -> "App Server" -> "Postgres"
```

D2 diagrams support most [D2 syntax](https://d2lang.com/), including:

- Connections with labels: `a -> b: "label"`
- Shapes: `x: { shape: cylinder }`
- Groups: `group { a; b; c }`

::: tip
Keep diagrams simple. Complex diagrams with many nodes can be hard to read at documentation site widths.
:::

## Syntax highlighting

Fenced code blocks are highlighted by [Chroma](https://github.com/alecthomas/chroma) at build time. Specify the language after the opening fence:

````markdown
```go
func main() {
    fmt.Println("hello")
}
```
````

The highlighting CSS is generated once at build time from the themes configured in `kwelea.toml`:

```toml
[theme]
light_code_theme = "github"
dark_code_theme  = "github-dark"
```

Any [Chroma style](https://xyproto.github.io/splash/docs/) is valid. Light and dark themes are emitted as separate CSS classes, toggled by `[data-theme]` on `<html>`.
