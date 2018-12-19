package lib

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestMarshalling(t *testing.T) {
	f, err := ioutil.ReadFile("./data.json")
	if err != nil {
		t.Fatal(err)
	}

	d := new(Datum)
	if err := json.Unmarshal(f, d); err != nil {
		t.Fatal(err)
	}

	if len(d.A) != 5 {
		t.Fatalf("should be 5")
	}

	if d.T.Humidity != 59 {
		t.Fatalf("should be 59")
	}

	if d.A[3].SP_10_0 != 12 {
		t.Fatalf("should be 12")
	}
}
