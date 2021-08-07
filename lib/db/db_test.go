package db

import (
	"context"
	"fmt"
	"testing"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

func TestDB(t *testing.T) {
	c, err := NewClient(context.Background(), "http://127.0.0.1:8086", "my-bucket", "my-org", "UQmCM9yFYT-B4UlQ0S2xLHSUQeSkEu-sm_9JObGKk-NGpvy56aU1bLAjgRduPBFZ0jXIuBBVgYE8AwgHYBSeNg==")
	if err != nil {
		t.Skip("no db")
		t.Fatalf(err.Error())
	}

	if err = c.Write(context.Background(), "foo.host", map[string]interface{}{"foo_val": 25.0}); err != nil {
		t.Fatalf(err.Error())
	}

	clnt := influxdb.NewClient(c.Host, c.Token)
	defer clnt.Close()

	queryAPI := clnt.QueryAPI(c.Org)
	result, err := queryAPI.Query(context.Background(), `from(bucket:"my-bucket")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "air_quality")`)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// Iterate over query response
	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		// Access data
		fmt.Printf("value: %v\n", result.Record().Value())
	}
	// check for an error
	if result.Err() != nil {
		fmt.Printf("query parsing error: %s\n", result.Err().Error())
	}
}
