package lib

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const (
	iters = 5
	D     = time.Microsecond
)

func TestSensor(t *testing.T) {
	h := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("./data.json")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer f.Close()

		io.Copy(w, f)
	}))

	defer h.Close()

	s := NewSensor(context.Background(), h.URL, D)

	var count int
	for dat := range s.Datum {
		count++
		if len(dat.A) != iters {
			t.Fatalf("dat isn't the right length")
		}

		if count == iters {
			s.Stop()
		}
	}
}
