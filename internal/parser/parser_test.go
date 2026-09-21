package parser

import (
	"errors"
	"strings"
	"testing"
)

type errReader struct{}

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("simulated read failure")
}

func TestExtract_Basic(t *testing.T) {
	htmlDoc := `
	<html>
		<head><title>Example Domain</title></head>
		<body>
			<a href="https://example.com/about">About</a>
			<a href="/contact">Contact</a>
			<img src="/logo.png">
			<img src="https://cdn.example.com/banner.jpg">
		</body>
	</html>`

	title, links, images, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "Example Domain" {
		t.Errorf("title = %q, want %q", title, "Example Domain")
	}

	wantLinks := []string{"https://example.com/about", "https://example.com/contact"}
	if !equalSlices(links, wantLinks) {
		t.Errorf("links = %v, want %v", links, wantLinks)
	}

	wantImages := []string{"https://example.com/logo.png", "https://cdn.example.com/banner.jpg"}
	if !equalSlices(images, wantImages) {
		t.Errorf("images = %v, want %v", images, wantImages)
	}
}

func TestExtract_NoTitle(t *testing.T) {
	htmlDoc := `<html><body><p>No title here</p></body></html>`

	title, links, images, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "" {
		t.Errorf("title = %q, want empty string", title)
	}
	if links != nil {
		t.Errorf("links = %v, want nil", links)
	}
	if images != nil {
		t.Errorf("images = %v, want nil", images)
	}
}

func TestExtract_Empty(t *testing.T) {
	title, links, images, err := ExtractData(strings.NewReader(""), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "" || links != nil || images != nil {
		t.Errorf("got title=%q links=%v images=%v, want all empty", title, links, images)
	}
}

func TestExtract_MultipleTitles(t *testing.T) {
	htmlDoc := `<html><head><title>First</title></head><body><title>Second</title></body></html>`

	title, _, _, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "Second" {
		t.Errorf("title = %q, want %q (last <title> wins)", title, "Second")
	}
}

func TestExtract_TitleNodes(t *testing.T) {
	htmlDoc := `<html><head><title></title></head><body></body></html>`
	title, _, _, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "" {
		t.Errorf("title = %q, want empty", title)
	}

	htmlDoc2 := `<html><head><title><b>Bold Title</b></title></head><body></body></html>`
	title2, _, _, err2 := ExtractData(strings.NewReader(htmlDoc2), "https://example.com")
	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}
	want2 := "<b>Bold Title</b>"
	if title2 != want2 {
		t.Errorf("title = %q, want %q (title content is RCDATA, kept as literal text)", title2, want2)
	}
}
func TestExtract_NoHref(t *testing.T) {
	htmlDoc := `<html><body><a name="anchor">No href here</a></body></html>`
	_, links, _, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if links != nil {
		t.Errorf("links = %v, want nil", links)
	}
}

func TestExtract_NoSrc(t *testing.T) {
	htmlDoc := `<html><body><img alt="no src"></body></html>`
	_, _, images, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if images != nil {
		t.Errorf("images = %v, want nil", images)
	}
}

func TestExtract_EmptyAttrs(t *testing.T) {
	htmlDoc := `<html><body><a href="">Self</a><img src=""></body></html>`
	_, links, images, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com/page")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 1 || links[0] != "https://example.com/page" {
		t.Errorf("links = %v, want [https://example.com/page]", links)
	}
	if len(images) != 1 || images[0] != "https://example.com/page" {
		t.Errorf("images = %v, want [https://example.com/page]", images)
	}
}

func TestExtract_URLResolution(t *testing.T) {
	tests := []struct {
		name    string
		base    string
		href    string
		want    string
	}{
		{"relative path", "https://example.com/blog/post1", "../about", "https://example.com/about"},
		{"root-relative path", "https://example.com/blog/post1", "/contact", "https://example.com/contact"},
		{"same-dir relative", "https://example.com/blog/", "post2.html", "https://example.com/blog/post2.html"},
		{"query string preserved", "https://example.com", "/search?q=go", "https://example.com/search?q=go"},
		{"fragment preserved", "https://example.com/page", "#section2", "https://example.com/page#section2"},
		{"protocol-relative", "https://example.com", "//cdn.example.com/a.js", "https://cdn.example.com/a.js"},
		{"already absolute", "https://example.com", "http://other.com/x", "http://other.com/x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htmlDoc := `<html><body><a href="` + tt.href + `">link</a></body></html>`
			_, links, _, err := ExtractData(strings.NewReader(htmlDoc), tt.base)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(links) != 1 || links[0] != tt.want {
				t.Errorf("links = %v, want [%s]", links, tt.want)
			}
		})
	}
}

func TestExtract_BadBaseURL(t *testing.T) {
	htmlDoc := `<html><body><a href="/relative">link</a></body></html>`
	_, links, _, err := ExtractData(strings.NewReader(htmlDoc), "://bad-scheme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("links = %v, want exactly one entry", links)
	}
}

func TestExtract_EmptyBaseURL(t *testing.T) {
	htmlDoc := `<html><body><a href="/relative">link</a><a href="https://absolute.com/x">abs</a></body></html>`
	_, links, _, err := ExtractData(strings.NewReader(htmlDoc), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantLinks := []string{"/relative", "https://absolute.com/x"}
	if !equalSlices(links, wantLinks) {
		t.Errorf("links = %v, want %v", links, wantLinks)
	}
}

func TestExtract_BadHref(t *testing.T) {

	htmlDoc := "<html><body><a href=\"http://a b.com/\">bad</a></body></html>"
	_, links, _, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("links = %v, want exactly one (possibly-unresolved) entry", links)
	}
}

func TestExtract_DuplicateLinks(t *testing.T) {
	htmlDoc := `<html><body>
		<a href="/a">A</a>
		<a href="/a">A again</a>
		<a href="/b">B</a>
	</body></html>`
	_, links, _, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Parser does not deduplicate - duplicates should be preserved in order.
	want := []string{"https://example.com/a", "https://example.com/a", "https://example.com/b"}
	if !equalSlices(links, want) {
		t.Errorf("links = %v, want %v", links, want)
	}
}

func TestExtract_Nested(t *testing.T) {
	htmlDoc := `<html><body>
		<div>
			<ul>
				<li><a href="/one">One</a></li>
				<li><div><a href="/two"><img src="/two.png"></a></div></li>
			</ul>
		</div>
	</body></html>`
	_, links, images, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantLinks := []string{"https://example.com/one", "https://example.com/two"}
	if !equalSlices(links, wantLinks) {
		t.Errorf("links = %v, want %v", links, wantLinks)
	}
	wantImages := []string{"https://example.com/two.png"}
	if !equalSlices(images, wantImages) {
		t.Errorf("images = %v, want %v", images, wantImages)
	}
}
