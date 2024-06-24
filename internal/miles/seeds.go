package miles

// Seed URLs: To begin the crawl process, we need to provide a set of seed URLs to the web crawler. One way to do
// this is to use a website's domain name to crawl all of its web pages. To make system more efficient, we should
// be strategic in choosing the seed URL because it can impact the number of web pages that are crawled. The
// selection of the seed URL can depend on factors like geographical location, categories (entertainment, education,
// sports, food), content type, etc.

func GetSeeds() []MilesURL {
	url, _ := NewURL("http://www.stackoverflow.com")

	return []MilesURL{
		url,
	}
}
