package parser

import (
	"errors"
	"strings"
	"testing"
)

// errReader is an io.Reader that always returns an error, used to
// simulate a broken/interrupted response body.
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
	// The current traversal logic overwrites `title` every time a <title>
	// element is encountered, so the LAST one in document order wins.
	// This test pins down that (perhaps surprising) documented behavior.
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
	// <title></title> has no FirstChild, and <title><b>x</b></title> has a
	// FirstChild that is an ElementNode, not a TextNode - both should be
	// skipped without panicking.
	htmlDoc := `<html><head><title></title></head><body></body></html>`
	title, _, _, err := ExtractData(strings.NewReader(htmlDoc), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "" {
		t.Errorf("title = %q, want empty", title)
	}

	// <title> content is parsed as RCDATA (raw text), so markup written
	// inside it is NOT turned into child elements - it stays literal text
	// and IS picked up as the title.
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