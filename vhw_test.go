package nmea

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var vhw = []struct {
	name string
	raw  string
	err  string
	msg  VHW
}{
	{
		name: "good sentence",
		raw:  "$VWVHW,45.0,T,43.0,M,3.5,N,6.4,K*56",
		msg: VHW{
			TrueHeading:            45.0,
			MagneticHeading:        43.0,
			SpeedThroughWaterKnots: 3.5,
			SpeedThroughWaterKPH:   6.4,
		},
	},
	{
		name: "bad sentence",
		raw:  "$VWVHW,T,45.0,43.0,M,3.5,N,6.4,K*56",
		err:  "nmea: VWVHW invalid true heading: T",
	},
}

func TestVHW(t *testing.T) {
	for _, tt := range vhw {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				vhw := m.(VHW)
				vhw.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, vhw)
			}
		})
	}
}

func TestVHW_NaNForEmptyFloat(t *testing.T) {
	// Test that SentenceParser with NaNForEmptyFloat returns NaN for empty float fields
	p := SentenceParser{
		NaNForEmptyFloat: true,
	}

	t.Run("partial sentence with NaN", func(t *testing.T) {
		m, err := p.Parse("$INVHW,187.9,T,,,19.6,N,36.3,K*3E")
		assert.NoError(t, err)
		vhw := m.(VHW)

		assert.Equal(t, 187.9, vhw.TrueHeading)
		assert.True(t, math.IsNaN(vhw.MagneticHeading), "empty MagneticHeading should be NaN")
		assert.Equal(t, 19.6, vhw.SpeedThroughWaterKnots)
		assert.Equal(t, 36.3, vhw.SpeedThroughWaterKPH)
	})

	t.Run("full sentence unaffected", func(t *testing.T) {
		m, err := p.Parse("$VWVHW,45.0,T,43.0,M,3.5,N,6.4,K*56")
		assert.NoError(t, err)
		vhw := m.(VHW)

		assert.Equal(t, 45.0, vhw.TrueHeading)
		assert.Equal(t, 43.0, vhw.MagneticHeading)
		assert.Equal(t, 3.5, vhw.SpeedThroughWaterKnots)
		assert.Equal(t, 6.4, vhw.SpeedThroughWaterKPH)
	})
}
