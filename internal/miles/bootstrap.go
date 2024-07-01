package miles

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/hoyle1974/miles/internal/url"
	"log/slog"
	"strings"
	"time"

	"github.com/alitto/pond"
	"github.com/hoyle1974/miles/internal/store"
)

type TimerLog struct {
	name  string
	start time.Time
	end   time.Time
}

func NewtimerLog(name string) *TimerLog {
	return &TimerLog{name: name, start: time.Now()}
}

func (t *TimerLog) Stop() {
	t.end = time.Now()
}

func (t *TimerLog) String() string {
	d := t.end.Sub(t.start).Seconds()
	return fmt.Sprintf("%s(%v)", t.name, d)
}

type BootstrapContext struct {
	logger    *slog.Logger
	batchSize int
	ioPool    *pond.WorkerPool
	cpuPool   *pond.WorkerPool
	docStore  store.DocStore
	frontier  Frontier
	robots    Robots
}

func processFetch(ctx *BootstrapContext, url url.Nurl) {
	ctx.logger.Info("processFetch", "url", url)
	data, contentType, responseCode, err := FetchURL(url)
	if err != nil {
		_ = ctx.docStore.Store(url, data, contentType, responseCode, err)
		ctx.logger.Error("Error Fetching URL", "url", url, "err", err)
		return
	}
	if !strings.Contains(contentType, "text") && contentType != "" {
		ctx.logger.Debug("Skipping URL", "url", url, "contentType", contentType)
		_ = ctx.docStore.Store(url, nil, contentType, responseCode, nil)
		return // Skip this url
	}

	err = ctx.docStore.Store(url, data, contentType, responseCode, err)
	if err != nil {
		ctx.logger.Error("Error Storing URL Data", "err", err, "url", url)
		return
	}

	ctx.cpuPool.Submit(func() { processCPU(ctx, data, url) })
}

func processCPU(ctx *BootstrapContext, data []byte, url url.Nurl) {
	ctx.logger.Info("processCPU", "url", url, "dataLength", len(data))
	urls, err := ExtractURLs(url, data)
	if err != nil {
		ctx.logger.Error("Error Extracing URLS", "err", err)
		return
	}
	ctx.frontier.AddURLS(Filter(DeduplicateURLs(urls)))
}

func processCache(ctx *BootstrapContext, url url.Nurl) {
	ctx.logger.Info("processCache", "url", url)

	doc, err := ctx.docStore.GetDoc(url)
	if err != nil && err != badger.ErrKeyNotFound {
		ctx.logger.Error("Error Fetching URL from Docstore", "err", err)
		return
	}

	if err == badger.ErrKeyNotFound {
		// Go fetch the data
		ctx.ioPool.Submit(func() { processFetch(ctx, url) })
		return
	}

	ctx.cpuPool.Submit(func() { processCPU(ctx, doc.GetData(), url) })

}

func Bootstrap(logger *slog.Logger) {

	ctx := &BootstrapContext{
		logger:    logger,
		batchSize: 256,
		ioPool:    pond.New(1, 1),
		cpuPool:   pond.New(1, 1),
	}

	ctx.logger.Info("Bootstrapping . . .")
	ctx.docStore = store.NewDocStore()
	ctx.frontier = NewFrontier()
	ctx.robots = GetRobots(ctx.docStore)

	defer ctx.docStore.Close()
	defer ctx.frontier.Close()

	// Get as much data as we can get
	for true {
		size, domains := ctx.frontier.Sizes()
		if size == 0 {
			logger.Warn("Frontier is now empty")
			time.Sleep(time.Second)
			continue
		}
		logger.Info(fmt.Sprintf("--------- Frontier Size: %d   Domains: %d", size, domains))

		// Get a batch to process
		batch, err := ctx.frontier.GetNextURLBatch(ctx.batchSize)
		if err != nil {
			panic(err)
		}
		if len(batch) != ctx.batchSize {
			logger.Warn(fmt.Sprintf("Asked for %d but got %d", ctx.batchSize, len(batch)))
		}

		for _, url := range batch {
			ctx.ioPool.Submit(func() { processCache(ctx, url) })
		}
	}

	/*
		for true {
			size, domains := frontier.Sizes()
			if size == 0 {
				logger.Warn("Frontier is now empty")
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
				ioPool.Submit(func() {
					all := NewtimerLog("all")

					//countdown.Add(-1)
					log := fmt.Sprintf("Working with(%s) ", url)

					if robots.IsValid(url) {
						diskPool.Submit(func() {
							doc, err := docStore.GetDoc(url)
							if err != nil && err != badger.ErrKeyNotFound {
								logger.Error("Error Fetching URL from Docstore", "err", err)
								return
							}
							var data []byte
							if err != badger.ErrKeyNotFound {
								if doc.GetError() == nil {
									log = "** CACHED ** " + log
									data = doc.GetData()
								} else {
									logger.Warn("Cached Error for URL", "url", url, "err", doc.GetError())
								}
							} else {
								ioPool.Submit(func() {
									var responseCode int
									var contentType = "text"

									fetchTime := NewtimerLog("fetch")
									data, contentType, responseCode, err = FetchURL(url)
									fetchTime.Stop()
									log = fetchTime.String() + " " + log
									if err != nil {
										_ = docStore.Store(url, data, contentType, responseCode, err)
										logger.Error("Error Fetching URL", "url", url, "err", err)
										return
									}
									if !strings.Contains(contentType, "text") && contentType != "" {
										logger.Debug("Skipping URL", "url", url, "contentType", contentType)
										_ = docStore.Store(url, nil, contentType, responseCode, nil)
										return // Skip this url
									}

									log = log + fmt.Sprintf("Data Size = %d ", len(data))
									err = docStore.Store(url, data, contentType, responseCode, err)
									if err != nil {
										logger.Error("Error Storing URL Data", "err", err)
										return
									}
								})

							}

							cpuPool.Submit(func() {
								processTime := NewtimerLog("process")
								urls, err := ExtractURLs(url, data)
								if err != nil {
									logger.Error("Error Extracing URLS", "err", err)
									return
								}
								log = log + fmt.Sprintf("URLS=%d ", len(urls))
								frontier.AddURLS(Filter(DeduplicateURLs(urls)))
								processTime.Stop()

								log = processTime.String() + " " + log

								all.Stop()
								log = all.String() + " " + log

								logger.Info(log)
							})
						})
					}
				})
			}
		}
	*/

	logger.Info("Stopping and waiting . . .")
	ctx.ioPool.StopAndWait()
	ctx.cpuPool.StopAndWait()

}
