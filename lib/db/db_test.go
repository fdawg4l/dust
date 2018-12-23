package db

import (
	"context"
	"testing"
)

func TestDB(t *testing.T) {
	c, err := NewClient(context.Background(), "http://127.0.0.1:8086", "air_quality", "", "")
	if err != nil {
		t.Skip("no db")
		t.Fatalf(err.Error())
	}

	if err = c.Write("foo.host", map[string]interface{}{"foo_val": 25.0}); err != nil {
		t.Fatalf(err.Error())
	}
}
