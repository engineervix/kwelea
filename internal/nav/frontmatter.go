package nav

import "strings"

// Frontmatter holds the YAML front matter fields recognised by kwelea.
// All fields are optional; their zero values are sensible defaults.
type Frontmatter struct {
	Title       string
	Description string
	Draft       bool
}

// parseFrontmatter splits a Markdown file's content into its front matter
// and body. It handles the standard "---"-delimited YAML block.
//
// If the file does not start with "---\n", the whole content is returned as
// the body with an empty Frontmatter. Malformed delimiters (opening "---"
// with no closing one) are treated the same way.
func parseFrontmatter(content string) (Frontmatter, string, error) {
	const delim = "---"
	if !strings.HasPrefix(content, delim+"\n") {
		return Frontmatter{}, content, nil
	}

	rest := content[len(delim)+1:] // skip the opening "---\n"
	end := strings.Index(rest, "\n"+delim)
	if end < 0 {
		// No closing delimiter — treat the whole file as body.
		return Frontmatter{}, content, nil
	}

	block := rest[:end]
	body := rest[end+1+len(delim):] // skip "\n---"
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:] // drop the blank line that follows the closing "---"
	}

	var fm Frontmatter
	for _, line := range strings.Split(block, "\n") {
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		// Strip surrounding double-quotes, e.g. title: "My: Guide"
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		switch key {
		case "title":
			fm.Title = val
		case "description":
			fm.Description = val
		case "draft":
			fm.Draft = val == "true"
		}
	}

	return fm, body, nil
}
