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
		name       string
		siteOGImg  string
		pageImg    string
		want       string
	}{
		{
			name:      "page image overrides site default",
			siteOGImg: "https://example.com/default.png",
			pageImg:   "https://example.com/custom.png",
			want:      "https://example.com/custom.png",
		},
		{
			name:      "falls back to site default",
			siteOGImg: "https://example.com/default.png",
			pageImg:   "",
			want:      "https://example.com/default.png",
		},
		{
			name:      "neither set — empty",
			siteOGImg: "",
			pageImg:   "",
			want:      "",
		},
		{
			name:      "page image only, no site default",
			siteOGImg: "",
			pageImg:   "https://example.com/page-specific.png",
			want:      "https://example.com/page-specific.png",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			site := &nav.Site{OGImage: tc.siteOGImg}
			page := &nav.Page{Image: tc.pageImg}
			if got := ogImage(site, page); got != tc.want {
				t.Errorf("ogImage() = %q, want %q", got, tc.want)
			}
		})
	}
}