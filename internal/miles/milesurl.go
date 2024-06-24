package miles

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type MilesURL struct {
	URL url.URL
}

func (m MilesURL) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	temp := m.URL.String()
	fmt.Fprintln(&b, temp)
	return b.Bytes(), nil
}

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (m *MilesURL) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	temp := ""
	_, err := fmt.Fscanln(b, &temp)

	a, _ := NewURL(temp)
	m.URL = a.URL

	return err
}

func (m MilesURL) String() string {
	return m.URL.String()
}

func (m MilesURL) Hostname() string {
	return m.URL.Hostname()
}

func (m MilesURL) Path() string {
	return m.URL.Path
}

func (m MilesURL) CopyHostname(from MilesURL) MilesURL {
	m.URL.Scheme = from.URL.Scheme
	m.URL.Host = from.URL.Host
	return m
}

func NewURL(urlString string) (MilesURL, error) {

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return MilesURL{}, err
	}

	// Lowercase scheme and host
	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
	parsedURL.Host = strings.ToLower(parsedURL.Host)

	// Remove default port (if present)
	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if parsedURL.Scheme == "http" && port == "80" {
		port = ""
	}
	if parsedURL.Scheme == "https" && port == "443" {
		port = ""
	}
	if port != "" {
		host = fmt.Sprintf("%s:%s", host, port)
	}
	parsedURL.Host = host

	// Remove trailing slash (except for root path)
	if parsedURL.Path != "/" && strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path = parsedURL.Path[:len(parsedURL.Path)-1]
	}

	query := parsedURL.Query()

	keys := make([]string, len(query))
	i := 0
	for k := range query {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	for _, key := range keys {
		values := query[key]
		// Sort values within the key (optional)
		sort.Strings(values)
		query[key] = values
	}
	parsedURL.RawQuery = query.Encode()

	return MilesURL{URL: *parsedURL}, nil
}
