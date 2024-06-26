package miles

import (
	"fmt"
	"github.com/hoyle1974/miles/internal/store"
	"github.com/hoyle1974/miles/internal/url"
	"io/ioutil"
	"net/http"
)

// HTML Fetcher: HTML fetcher component is responsible for downloading web pages corresponding to a given URL provided by the URL Frontier. It does this by using a network protocol like HTTP or HTTPS. In simple words, HTML fetcher retrieves the actual web page content that needs to be analyzed and stored.
var docStore = store.NewDocStore()

func FetchURL(URL url.Nurl) ([]byte, error) {

	doc, err := docStore.GetDoc(URL)
	if doc != nil {
		return doc.GetData(), doc.GetError()
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		err = fmt.Errorf("NewRequest %w", err)
		docStore.Store(URL, nil, err)
		return nil, err
	}

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("DoRequest %w", err)
		docStore.Store(URL, nil, err)
		return nil, err
	}
	defer resp.Body.Close() // Close the response body after use

	// Check for successful response status code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		docStore.Store(URL, nil, err)
		return nil, err
	}

	// Read the response body into a byte buffer
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("ReadAll %w", err)
		docStore.Store(URL, nil, err)
		return nil, err
	}

	docStore.Store(URL, body, nil)
	return body, nil
}
