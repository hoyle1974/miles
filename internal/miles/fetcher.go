package miles

import (
	"fmt"
	"github.com/hoyle1974/miles/internal/url"
	"io/ioutil"
	"net/http"
)

// HTML Fetcher: HTML fetcher component is responsible for downloading web pages corresponding to a given URL provided by the URL Frontier. It does this by using a network protocol like HTTP or HTTPS. In simple words, HTML fetcher retrieves the actual web page content that needs to be analyzed and stored.
func GetHeaders(URL url.Nurl) (int, map[string][]string, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request (using HEAD method to fetch only headers)
	req, err := http.NewRequest(http.MethodHead, URL.String(), nil)
	if err != nil {
		return -1, nil, err
	}

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close() // Close the body after use

	// Return the headers as a map
	return resp.StatusCode, resp.Header, nil
}

func FetchURL(URL url.Nurl) ([]byte, string, int, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		err = fmt.Errorf("NewRequest %w", err)
		return nil, "", -1, err
	}

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("DoRequest %w", err)
		return nil, "", -1, err
	}
	defer resp.Body.Close() // Close the response body after use

	// Check for successful response status code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return nil, "", resp.StatusCode, err
	}

	// Read the response body into a byte buffer
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("ReadAll %w", err)
		return nil, "", resp.StatusCode, err
	}

	contentType := resp.Header.Get("Content-Type")

	return body, contentType, resp.StatusCode, nil
}
