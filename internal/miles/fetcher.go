package miles

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTML Fetcher: HTML fetcher component is responsible for downloading web pages corresponding to a given URL provided by the URL Frontier. It does this by using a network protocol like HTTP or HTTPS. In simple words, HTML fetcher retrieves the actual web page content that needs to be analyzed and stored.

func FetchURL(URL MilesURL) ([]byte, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		return nil, err
	}

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Close the response body after use

	// Check for successful response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body into a byte buffer
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
