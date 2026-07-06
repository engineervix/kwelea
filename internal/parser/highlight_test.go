package parser

import (
	"reflect"
	"strings"
	"testing"
)

// ----- parseCodeBlockInfo -----

func TestParseCodeBlockInfoTitleOnly(t *testing.T) {
	got := parseCodeBlockInfo(`go title="cmd/root.go"`)
	if got.title != "cmd/root.go" {
		t.Errorf("title: got %q want %q", got.title, "cmd/root.go")
	}
	if got.hasHighlights() {
		t.Errorf("expected no highlights, got %v", got.highlight)
	}
}

func TestParseCodeBlockInfoHighlightsOnly(t *testing.T) {
	got := parseCodeBlockInfo(`go {2,4-6}`)
	if got.title != "" {
		t.Errorf("title: got %q want empty", got.title)
	}
	want := []int{2, 4, 5, 6}
	if !reflect.DeepEqual(got.highlight, want) {
		t.Errorf("highlight: got %v want %v", got.highlight, want)
	}
}

func TestParseCodeBlockInfoTitleAndHighlights(t *testing.T) {
	got := parseCodeBlockInfo(`python title="x.py" {1,3-4}`)
	if got.title != "x.py" {
		t.Errorf("title: got %q want %q", got.title, "x.py")
	}
	want := []int{1, 3, 4}
	if !reflect.DeepEqual(got.highlight, want) {
		t.Errorf("highlight: got %v want %v", got.highlight, want)
	}
}

func TestParseCodeBlockInfoOrderIndependent(t *testing.T) {
	a := parseCodeBlockInfo(`go {1} title="a"`)
	b := parseCodeBlockInfo(`go title="a" {1}`)
	if a.title != b.title || !reflect.DeepEqual(a.highlight, b.highlight) {
		t.Errorf("order changed result: a=%+v b=%+v", a, b)
	}
}

func TestParseCodeBlockInfoSingleQuotedTitle(t *testing.T) {
	got := parseCodeBlockInfo(`sh title='has spaces'`)
	if got.title != "has spaces" {
		t.Errorf("single-quoted title: got %q want %q", got.title, "has spaces")
	}
}

func TestParseCodeBlockInfoUnquotedTitleDropped(t *testing.T) {
	// title=foo (no quotes) is ambiguous; per the doc comment it is
	// silently dropped to avoid false positives.
	got := parseCodeBlockInfo(`go title=foo`)
	if got.title != "" {
		t.Errorf("unquoted title should be dropped, got %q", got.title)
	}
}

func TestParseCodeBlockInfoInvalidRangeTokens(t *testing.T) {
	got := parseCodeBlockInfo(`go {abc,5,xyz,2-4}`)
	want := []int{5, 2, 3, 4}
	if !reflect.DeepEqual(got.highlight, want) {
		t.Errorf("invalid tokens should be skipped, got %v want %v", got.highlight, want)
	}
}

func TestParseCodeBlockInfoEmpty(t *testing.T) {
	got := parseCodeBlockInfo("")
	if got.hasTitle() || got.hasHighlights() {
		t.Errorf("empty input should yield no attrs, got %+v", got)
	}
}

func TestParseCodeBlockInfoLanguageOnly(t *testing.T) {
	// A bare language with no other attributes must come out empty.
	got := parseCodeBlockInfo("go")
	if got.hasTitle() || got.hasHighlights() {
		t.Errorf("language-only should yield no attrs, got %+v", got)
	}
}

// ----- parseHighlightRanges -----

func TestParseHighlightRangesSingleNumber(t *testing.T) {
	got := parseHighlightRanges("{5}")
	if !reflect.DeepEqual(got, []int{5}) {
		t.Errorf("got %v want [5]", got)
	}
}

func TestParseHighlightRangesRangeExpands(t *testing.T) {
	got := parseHighlightRanges("{3-5}")
	if !reflect.DeepEqual(got, []int{3, 4, 5}) {
		t.Errorf("got %v want [3 4 5]", got)
	}
}

func TestParseHighlightRangesMixed(t *testing.T) {
	got := parseHighlightRanges("{1,3-4,7}")
	if !reflect.DeepEqual(got, []int{1, 3, 4, 7}) {
		t.Errorf("got %v want [1 3 4 7]", got)
	}
}

func TestParseHighlightRangesReverseDropped(t *testing.T) {
	got := parseHighlightRanges("{5-2}")
	if len(got) != 0 {
		t.Errorf("reverse range should be empty, got %v", got)
	}
}

func TestParseHighlightRangesInvalidNumbersSkipped(t *testing.T) {
	got := parseHighlightRanges("{abc,0,5,-1,3}")
	// 0 and -1 are below 1 → dropped; "abc" → dropped; 5, 3 kept.
	if !reflect.DeepEqual(got, []int{5, 3}) {
		t.Errorf("got %v want [5 3]", got)
	}
}

func TestParseHighlightRangesBadBraces(t *testing.T) {
	cases := []string{"5", "{5", "5}", "{}", "{ }"}
	for _, c := range cases {
		if got := parseHighlightRanges(c); len(got) != 0 {
			t.Errorf("parseHighlightRanges(%q) = %v, want empty", c, got)
		}
	}
}

// ----- injectHighlightClass -----

func TestInjectHighlightClassSkipsEmpty(t *testing.T) {
	in := `<span class="line"><span class="cl">x</span></span>`
	out := injectHighlightClass(in, nil)
	if out != in {
		t.Errorf("nil set should leave input untouched")
	}
}

func TestInjectHighlightClassAddsClass(t *testing.T) {
	in := `<span class="line"><span class="cl">a</span></span>` +
		`<span class="line"><span class="cl">b</span></span>`
	out := injectHighlightClass(in, map[int]bool{2: true})
	if !strings.Contains(out, `<span class="line highlight-line"><span class="cl">b</span></span>`) {
		t.Errorf("second line should gain highlight-line; got:\n%s", out)
	}
	// First line should be untouched.
	if !strings.Contains(out, `<span class="line"><span class="cl">a</span></span>`) {
		t.Errorf("first line should be untouched; got:\n%s", out)
	}
}

func TestInjectHighlightClassPreservesChromaHl(t *testing.T) {
	// When chroma's HighlightLines is on, the span already has "line hl".
	// Our rewriter should add "highlight-line" without dropping "hl".
	in := `<span class="line hl"><span class="cl">a</span></span>` +
		`<span class="line"><span class="cl">b</span></span>`
	out := injectHighlightClass(in, map[int]bool{1: true})
	if !strings.Contains(out, `<span class="line highlight-line hl"><span class="cl">a</span></span>`) {
		t.Errorf("hl should be preserved; got:\n%s", out)
	}
}

func TestInjectHighlightClassIdempotent(t *testing.T) {
	// Running inject twice should be the same as running it once.
	in := `<span class="line"><span class="cl">x</span></span>`
	once := injectHighlightClass(in, map[int]bool{1: true})
	twice := injectHighlightClass(once, map[int]bool{1: true})
	if once != twice {
		t.Errorf("not idempotent:\nonce:  %s\ntwice: %s", once, twice)
	}
}

// ----- rewriteLineOpenTag -----

func TestRewriteLineOpenTagBare(t *testing.T) {
	got := rewriteLineOpenTag(`<span class="line">`)
	want := `<span class="line highlight-line">`
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestRewriteLineOpenTagWithHl(t *testing.T) {
	got := rewriteLineOpenTag(`<span class="line hl">`)
	want := `<span class="line highlight-line hl">`
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestRewriteLineOpenTagAlreadyHighlighted(t *testing.T) {
	// If "highlight-line" is already there, leave the tag alone.
	in := `<span class="line highlight-line">`
	if got := rewriteLineOpenTag(in); got != in {
		t.Errorf("got %q want %q (no change)", got, in)
	}
}

func TestRewriteLineOpenTagUnrelatedSpanUntouched(t *testing.T) {
	in := `<span class="other">`
	if got := rewriteLineOpenTag(in); got != in {
		t.Errorf("non-line span should be left alone; got %q", got)
	}
}

// ----- tokenizeInfo -----

// tokenizeInfo is invoked on the post-language portion of a fenced
// code-block info string. Tests below therefore pass strings without
// the leading language token.

func TestTokenizeInfoKeepsQuotedTitle(t *testing.T) {
	got := tokenizeInfo(`title="my file.go" {2}`)
	want := []string{`title="my file.go"`, "{2}"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestTokenizeInfoSingleQuotedTitle(t *testing.T) {
	got := tokenizeInfo(`title='has "double" quotes'`)
	want := []string{`title='has "double" quotes'`}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestTokenizeInfoHandlesExtraSpaces(t *testing.T) {
	got := tokenizeInfo("  {2,4}   ")
	want := []string{"{2,4}"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
