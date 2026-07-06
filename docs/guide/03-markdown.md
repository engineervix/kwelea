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

### Filename titles

Add `title="…"` to the opening fence to render a filename label above the block:

```go title="cmd/root.go"
package main

func main() {}
```

Renders as a `<figure class="code-block">` with a `<div class="code-title">cmd/root.go</div>` header bar.

**Syntax:**

````markdown
```go title="cmd/root.go"
package main
```
````

The title can be wrapped in either single or double quotes. Use the opposite quote style if the title itself contains quotes (e.g. `title='has "quotes" in it'`).

### Line highlighting

Add `{n,m-p,…}` to the opening fence to highlight specific lines (1-indexed, counted from the first line of the block). Use a comma to list individual lines, and a hyphen to specify a range:

```go {2,4-6}
line 1
line 2
line 3
line 4
line 5
line 6
```

Each highlighted line gets a `<span class="highlight-line">` class with a soft background wash and an accent bar to its left.

**Syntax:**

````markdown
```go {2,4-6}
line 1
line 2
line 3
line 4
line 5
line 6
```
````

### Combining both

The two attributes can be combined in any order:

```python title="example.py" {1,3-4}
print("a")
print("b")
print("c")
print("d")
```

**Syntax:**

````markdown
```python title="example.py" {1,3-4}
print("a")
print("b")
print("c")
print("d")
```
````

### Behaviour notes

- A plain code block (no `title=` or `{}`) is rendered exactly as before — the new extension is a no-op for the common case.
- Empty `title=""` is treated as no title.
- Invalid range tokens (`{notanumber}`) are silently ignored; the rest of the expression still applies.
- Out-of-range line numbers are silently ignored — no error, no broken page.
- Reverse ranges (`{5-2}`) are silently ignored.
- D2 diagram blocks (```` ```d2 ````) are caught by the D2 extension before the code-attribute one sees them; adding `title=` or `{}` to a D2 block does not change its behaviour.
