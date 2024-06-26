package miles

import (
	"bytes"
	"fmt"
	"github.com/hoyle1974/miles/internal/url"
	"golang.org/x/net/html"
	"strings"
)

/*
// ExtractText parses HTML from a byte buffer and returns the text content as a string.
func extractText(buffer []byte) (string, error) {
	reader := bytes.NewReader(buffer)
	doc, err := html.Parse(reader)
	if err != nil {
		return "", err
	}

	var text string
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			text += node.Data + " " // Add space between text nodes
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(doc)
	return strings.TrimSpace(text), nil
}
*/

func extractText(htmlBytes []byte) (string, error) {
	var textBuffer bytes.Buffer
	tokenizer := html.NewTokenizer(bytes.NewReader(htmlBytes))

	for {
		tt := tokenizer.Next()

		switch tt {
		case html.ErrorToken:
			if tokenizer.Err().Error() == "EOF" {
				return strings.TrimSpace(textBuffer.String()), nil
			}
			return "", fmt.Errorf("error parsing HTML: %w", tokenizer.Err())
		case html.TextToken:
			text := strings.TrimSpace(string(tokenizer.Text()))
			if text != "" {
				textBuffer.WriteString(text)
				textBuffer.WriteRune(' ') // Add space between text nodes
			}
		case html.StartTagToken, html.EndTagToken:
			name, _ := tokenizer.TagName()
			if string(name) == "style" || string(name) == "head" || string(name) == "form" || string(name) == "meta" || string(name) == "script" || string(name) == "img" || string(name) == "svg" || string(name) == "style" {
				// Skip script content
				for tt := tokenizer.Next(); tt != html.ErrorToken && tt != html.EndTagToken; tt = tokenizer.Next() {
				}
			}
		}
	}

	return strings.TrimSpace(textBuffer.String()), nil
}

// ExtractURLs takes a byte array containing HTML and returns a slice of extracted URLs
func extractURLs(htmlData []byte) ([]string, error) {
	var urls []string
	doc, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return nil, err
	}

	var extract func(*html.Node)
	extract = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					urls = append(urls, attr.Val)
					break
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			extract(child)
		}
	}
	extract(doc)
	return urls, nil
}

// ExtractURLs finds all URLs within an HTML byte array.
func ExtractURLs(currentURL url.Nurl, data []byte) ([]url.Nurl, error) {
	/*
		// Regex pattern for finding URLs (can be improved for specific needs)
		urlRegex := regexp.MustCompile(`(?i)(href|src)=["'](?P<url>[^"\s]+)["']`)

		// Find all matches of the regex pattern
		matches := urlRegex.FindAllSubmatch(data, -1)

		// Extract URLs from the matches
		urls := make([]url.Nurl, len(matches))
		for i, match := range matches {
			m, _ := url.NewURL(string(match[2]), currentURL.Scheme(), currentURL.Hostname())
			urls[i] = m // Access captured group (index 2)
		}
	*/

	surls, err := extractURLs(data)
	if err != nil {
		return []url.Nurl{}, nil
	}

	text, err := extractText(data)
	if err == nil {
		fmt.Println("----------------------: " + currentURL.String() + "\n" + text + "------------------\n")
	}

	var urls []url.Nurl

	for _, surl := range surls {
		m, err := url.NewURL(surl, currentURL.Scheme(), currentURL.Hostname())
		if err == nil {
			urls = append(urls, m)
		}
	}

	return urls, nil
}
