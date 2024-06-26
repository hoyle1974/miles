package miles

import (
	"fmt"
	"github.com/alitto/pond"
	"log/slog"
	"time"
)

func Bootstrap(logger *slog.Logger) {
	pool := pond.New(16, 1000)

	logger.Info("Bootstrapping . . .")

	frontier := GetFrontier()

	for true {
		size, domains := frontier.Sizes()
		if size == 0 {
			logger.Warn(("Frontier is now empty"))
			break
		}
		logger.Info(fmt.Sprintf("--------- Frontier Size: %d   Domains: %d", size, domains))
		batch, err := frontier.GetNextURLBatch(10)
		if err != nil {
			panic(err)
		}

		cache := GetCache()
		for _, url := range batch {
			pool.Submit(func() {
				logger.Info("Working with", "url", url.String())
				cache.UpdateURLInfo(url)

				if GetRobots().IsValid(url) {
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
