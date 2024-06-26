package miles

import (
	"fmt"
	"github.com/alitto/pond"
	"github.com/hoyle1974/miles/internal/store"
	"github.com/hoyle1974/miles/internal/url"
	"log/slog"
	"time"
)

func Bootstrap(logger *slog.Logger) {
	pool := pond.New(1, 1000)

	logger.Info("Bootstrapping . . .")

	docStore := store.NewDocStore()
	frontier := GetFrontier()
	robots := GetRobots(docStore)

	url, err := url.NewURL("www.google.com")
	if err != nil {
		panic(err)
	}
	doc, err := docStore.GetDoc(url)
	if err != nil {
		panic(err)
	}
	if doc.GetError() != nil {
		panic(doc.GetError())
	}
	fmt.Printf("DocSize: %d\n", len(doc.GetData()))

	for true {
		// Display some stats about the frontier
		size, domains := frontier.Sizes()
		if size == 0 {
			logger.Warn(("Frontier is now empty"))
			break
		}
		logger.Info(fmt.Sprintf("--------- Frontier Size: %d   Domains: %d", size, domains))

		// Get a batch to process
		batch, err := frontier.GetNextURLBatch(10)
		if err != nil {
			panic(err)
		}

		//cache := GetCache()
		for _, url := range batch {
			pool.Submit(func() {
				logger.Info("Working with", "url", url.String())
				cache.UpdateURLInfo(url)

				if robots.IsValid(url) {
					data, err := FetchURL(url)
					if err != nil {
						logger.Error("Error Fetching URL", "err", err)
						return
					}
					logger.Debug("Data", "size", len(data))

					urls, err := ExtractURLs(url, data)
					if err != nil {
						return
					}
					logger.Debug("	URLS", "count", len(urls))

					frontier.AddURLS(Filter(DeduplicateURLs(urls)))
				} else {
					logger.Warn("Robot Invalid", "url", url)
				}
			})
		}

		time.Sleep(time.Second)
	}

	logger.Info("Stopping and waiting . . .")
	pool.StopAndWait()

}
