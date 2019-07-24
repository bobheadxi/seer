package riot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeltas_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		d    Deltas
		want string
	}{
		{"keys in order", map[string]float64{
			"0-10":   6.6,
			"10-20":  8.2,
			"20-30":  6.4,
			"30-end": 1.4,
		}, "[6.6,8.2,6.4,1.4]"},
		{"keys sorted", map[string]float64{
			"30-end": 1.4,
			"10-20":  8.2,
			"0-10":   6.6,
			"20-30":  6.4,
		}, "[6.6,8.2,6.4,1.4]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}
