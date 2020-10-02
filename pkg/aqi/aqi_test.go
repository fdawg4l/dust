package aqi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAqiIndex(t *testing.T) {
	in := []struct {
		concentration, expected int
	}{
		{29, 86},
		{301, 341},
	}

	for _, test := range in {
		t.Run(fmt.Sprintf("%d", test.concentration), func(t *testing.T) {
			f, _ := I(float32(test.concentration))
			if !assert.Equal(t, test.expected, int(f)) {
				return
			}

			return
		})
	}
}
