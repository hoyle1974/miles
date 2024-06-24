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
				//logger.Info("Working with: ", url)
				//info, hostInfo, hostCount :=
				cache.UpdateURLInfo(url)
				//info := urlCache.GetURLInfo(url)

				if GetRobots().IsValid(url) {
					//logger.Info(fmt.Sprintf("	Seeds: %s  Hits: %d  Host Hits: %d  Host Count: %d", url, info.Hits, hostInfo.Hits, hostCount))
					data, err := FetchURL(url)
					if err != nil {
						return
					}

					urls, err := ExtractURLs(url, data)
					if err != nil {
						return
					}

					frontier.AddURLS(Filter(DeduplicateURLs(urls)))
				}
			})
		}

		time.Sleep(time.Second)

	}

	pool.StopAndWait()

}
