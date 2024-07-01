package miles

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/hoyle1974/miles/internal/store"
	"github.com/hoyle1974/miles/internal/url"
	"log/slog"
	"os"
	"testing"

	"github.com/lmittmann/tint"
)

func BenchmarkBootstrap(b *testing.B) {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	for i := 0; i < b.N; i++ {
		fmt.Printf("%d of %d\n", i, b.N)
		Bootstrap(logger)
	}
}

func TestBootstrap(t *testing.T) {
	docStore := store.NewDocStore()

	temp, err := url.NewURL("http://www.stackoverflow.com", "http", "www.stackoverflow.com")
	if err != nil {
		t.Error(err)
		return
	}

	err = docStore.Del(temp)
	if err != nil {
		t.Error(err)
		return
	}

	data, contentType, responseCode, err := FetchURL(temp)
	if err != nil {
		t.Error(err)
		return
	}

	doc, err := docStore.GetDoc(temp)
	if err != badger.ErrKeyNotFound {
		t.Errorf("expected key not to found")
		return
	}

	err = docStore.Store(temp, data, contentType, responseCode, nil)
	if err != nil {
		t.Error(err)
		return
	}

	doc, err = docStore.GetDoc(temp)
	if err != nil {
		t.Error(err)
		return
	}
	if string(doc.GetData()) != string(data) {
		t.Errorf("Data blob did not match")
	}
	if doc.GetError() != nil {
		t.Errorf("Data error did not match")
	}
	if doc.GetResponse() != 200 {
		t.Errorf("Data reponse did not match")
	}
}
