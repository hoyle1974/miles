package url

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type Nurl struct {
	URL url.URL
}

func (m Nurl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	temp := m.URL.String()
	fmt.Fprintln(&b, temp)
	return b.Bytes(), nil
}

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (m Nurl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	temp := ""
	_, err := fmt.Fscanln(b, &temp)

	a, _ := NewURL(temp, "", "")
	m.URL = a.URL

	return err
}

func (m Nurl) String() string {
	return m.URL.String()
}

func (m Nurl) Hostname() string {
	return m.URL.Hostname()
}

func (m Nurl) Scheme() string {
	return m.URL.Scheme
}

func (m Nurl) Path() string {
	return m.URL.Path
}

func (m Nurl) CopyHostname(from Nurl) Nurl {
	m.URL.Scheme = from.URL.Scheme
	m.URL.Host = from.URL.Host
	return m
}

func NewURL(urlString string, defaultScheme string, defaultHost string) (Nurl, error) {

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return Nurl{}, err
	}

	// Lowercase scheme and host
	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
	parsedURL.Host = strings.ToLower(parsedURL.Host)

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = defaultScheme
	}
	if parsedURL.Host == "" {
		parsedURL.Host = defaultHost
	}

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

	return Nurl{URL: *parsedURL}, nil
}
