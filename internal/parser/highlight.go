package parser

import (
	"bytes"
	"fmt"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"

	"github.com/engineervix/kwelea/internal/config"
)

// ChromaCSS generates a combined Chroma syntax-highlighting CSS string for
// both the light and dark themes specified in themeCfg.
//
// The light theme uses standard .chroma selectors. Dark-theme rules are
// prefixed with [data-theme="dark"] so they activate only when that attribute
// is set on <html>, matching the toggle logic in _dev/docs-ui.html.
func ChromaCSS(themeCfg config.ThemeConfig) (string, error) {
	formatter := chromahtml.New(chromahtml.WithClasses(true))

	lightStyle := styles.Get(themeCfg.LightCodeTheme)
	if lightStyle == nil {
		lightStyle = styles.Fallback
	}
	darkStyle := styles.Get(themeCfg.DarkCodeTheme)
	if darkStyle == nil {
		darkStyle = styles.Fallback
	}

	var buf bytes.Buffer

	// Light theme — standard selectors.
	if err := formatter.WriteCSS(&buf, lightStyle); err != nil {
		return "", fmt.Errorf("generating light Chroma CSS (%s): %w", themeCfg.LightCodeTheme, err)
	}

	buf.WriteString("\n")

	// Dark theme — prefix every selector with [data-theme="dark"].
	var darkBuf bytes.Buffer
	if err := formatter.WriteCSS(&darkBuf, darkStyle); err != nil {
		return "", fmt.Errorf("generating dark Chroma CSS (%s): %w", themeCfg.DarkCodeTheme, err)
	}
	buf.WriteString(prefixCSSSelectors(darkBuf.String(), `[data-theme="dark"]`))

	return buf.String(), nil
}

// prefixCSSSelectors prepends prefix to every CSS selector in chroma formatter
// output. Chroma emits one rule per line in the form:
//
//	/* Token Name */ .chroma { ... }
//	/* Token Name */ .chroma .tok { ... }
//
// so finding the first '.' on each non-blank line reliably locates the selector.
func prefixCSSSelectors(css, prefix string) string {
	lines := strings.Split(css, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if i := strings.Index(line, "."); i >= 0 {
			result = append(result, line[:i]+prefix+" "+line[i:])
		} else {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}
