package miles

import (
	"fmt"
	"github.com/alitto/pond"
	"github.com/dgraph-io/badger/v4"
	"github.com/hoyle1974/miles/internal/store"
	"log/slog"
	"time"
)

func Bootstrap(logger *slog.Logger) {
	batchSize := 64
	pool := pond.New(batchSize, 1000)

	logger.Info("Bootstrapping . . .")

	docStore := store.NewDocStore()
	frontier := GetFrontier()
	robots := GetRobots(docStore)

	for true {
		// Display some stats about the frontier
		size, domains := frontier.Sizes()
		if size == 0 {
			logger.Warn(("Frontier is now empty"))
			time.Sleep(time.Second)
			continue
		}
		logger.Info(fmt.Sprintf("--------- Frontier Size: %d   Domains: %d", size, domains))

		// Get a batch to process
		batch, err := frontier.GetNextURLBatch(batchSize)
		if err != nil {
			panic(err)
		}
		if len(batch) != batchSize {
			logger.Warn(fmt.Sprintf("Asked for %d but got %d", batchSize, len(batch)))
		}

		for _, url := range batch {
			pool.Submit(func() {
				log := fmt.Sprintf("Working with(%s) ", url)

				if robots.IsValid(url) {
					doc, err := docStore.GetDoc(url)
					if err != nil && err != badger.ErrKeyNotFound {
						logger.Error("Error Fetching URL from Docstore", "err", err)
						return
					}
					var data []byte
					if err != badger.ErrKeyNotFound && doc.GetError() == nil {
						log = "** CACHED ** " + log
						data = doc.GetData()
					} else {
						data, err := FetchURL(url)
						if err != nil {
							_ = docStore.Store(url, data, err)
							logger.Error("Error Fetching URL", "err", err)
							return
						}
						log = log + fmt.Sprintf("Data Size = %d", len(data))
						err = docStore.Store(url, data, err)
						if err != nil {
							logger.Error("Error Storing URL Data", "err", err)
							return
						}
					}

					urls, err := ExtractURLs(url, data)
					if err != nil {
						logger.Error("Error Extracing URLS", "err", err)
						return
					}
					log = log + fmt.Sprintf("URLS=%d ", len(urls))

					frontier.AddURLS(Filter(DeduplicateURLs(urls)))
				} else {
					log = log + "Robot Invalid "
				}
				logger.Info(log)
			})
		}

		time.Sleep(time.Second)
	}

	logger.Info("Stopping and waiting . . .")
	pool.StopAndWait()

}
