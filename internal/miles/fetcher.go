package miles

import (
	"fmt"
	"github.com/hoyle1974/miles/internal/store"
	"github.com/hoyle1974/miles/internal/url"
	"io/ioutil"
	"net/http"
)

// HTML Fetcher: HTML fetcher component is responsible for downloading web pages corresponding to a given URL provided by the URL Frontier. It does this by using a network protocol like HTTP or HTTPS. In simple words, HTML fetcher retrieves the actual web page content that needs to be analyzed and stored.

func FetchURL(URL url.Nurl) ([]byte, error) {
	data, err := store.GetKVStore().Get(URL.String())
	if data != nil {
		return data, nil
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("NewRequest %w", err)
	}

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("DoRequest %w", err)
	}
	defer resp.Body.Close() // Close the response body after use

	// Check for successful response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body into a byte buffer
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll %w", err)
	}

	store.GetKVStore().Put(URL.String(), body)

	return body, nil
}
