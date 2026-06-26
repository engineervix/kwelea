package builder

import (
	"testing"

	"github.com/engineervix/kwelea/internal/nav"
)

func TestOGURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		pagePath string
		want     string
	}{
		{
			name:     "https base_url",
			baseURL:  "https://example.com",
			pagePath: "/getting-started/",
			want:     "https://example.com/getting-started/",
		},
		{
			name:     "base_url with trailing slash",
			baseURL:  "https://example.com/kwelea/",
			pagePath: "/guide/",
			want:     "https://example.com/kwelea/guide/",
		},
		{
			name:     "root page",
			baseURL:  "https://example.com",
			pagePath: "/",
			want:     "https://example.com/",
		},
		{
			name:     "no base_url — empty result",
			baseURL:  "",
			pagePath: "/getting-started/",
			want:     "",
		},
		{
			name:     "relative base_url — empty result (matches sitemap behaviour)",
			baseURL:  "/kwelea",
			pagePath: "/getting-started/",
			want:     "",
		},
		{
			name:     "scheme-only base_url — empty result",
			baseURL:  "https://",
			pagePath: "/getting-started/",
			want:     "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			site := &nav.Site{BaseURL: tc.baseURL}
			page := &nav.Page{Path: tc.pagePath}
			if got := ogURL(site, page); got != tc.want {
				t.Errorf("ogURL() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestOGImage(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		siteOGImg string
		pageImg   string
		want      string
	}{
		{
			name:      "page image overrides site default",
			baseURL:   "https://example.com",
			siteOGImg: "https://example.com/default.png",
			pageImg:   "https://example.com/custom.png",
			want:      "https://example.com/custom.png",
		},
		{
			name:      "falls back to site default",
			baseURL:   "https://example.com",
			siteOGImg: "https://example.com/default.png",
			pageImg:   "",
			want:      "https://example.com/default.png",
		},
		{
			name:    "neither set — empty",
			want:    "",
		},
		{
			name:    "page image only, no site default",
			pageImg: "https://example.com/page-specific.png",
			want:    "https://example.com/page-specific.png",
		},
		{
			name:    "absolute site image is returned verbatim",
			pageImg: "https://cdn.example.com/og.png",
			want:    "https://cdn.example.com/og.png",
		},
		{
			name:      "relative site image resolved against base_url",
			baseURL:   "https://engineervix.github.io/kwelea",
			siteOGImg: "/assets/og.png",
			want:      "https://engineervix.github.io/kwelea/assets/og.png",
		},
		{
			name:      "relative site image with base_url trailing slash",
			baseURL:   "https://engineervix.github.io/kwelea/",
			siteOGImg: "/assets/og.png",
			want:      "https://engineervix.github.io/kwelea/assets/og.png",
		},
		{
			name:      "relative page image resolved against base_url",
			baseURL:   "https://example.com",
			pageImg:   "/custom/og.png",
			want:      "https://example.com/custom/og.png",
		},
		{
			name:      "relative image with relative base_url — left as-is",
			baseURL:   "/kwelea",
			siteOGImg: "/assets/og.png",
			want:      "/assets/og.png",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			site := &nav.Site{BaseURL: tc.baseURL, OGImage: tc.siteOGImg}
			page := &nav.Page{Image: tc.pageImg}
			if got := ogImage(site, page); got != tc.want {
				t.Errorf("ogImage() = %q, want %q", got, tc.want)
			}
		})
	}
}