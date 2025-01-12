package jyoro

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected Config
		haveErr  bool
	}{
		{
			name: "valid",
			input: `
{
	"location": "Asia/Tokyo",
	"entries": [
		{
		  "start_at": "12:34:56",
			"duration": "1h30m"
		},
		{
			"start_at": "23:59:59",
			"duration": "1h0m0s"
		}
	]
}`,
			expected: Config{
				Location: mustLoadLocation("Asia/Tokyo"),
				Entries: []Entry{
					{
						StartAt:  HMS{Hour: 12, Minute: 34, Second: 56},
						Duration: Duration{time.Hour + 30*time.Minute},
					},
					{
						StartAt:  HMS{Hour: 23, Minute: 59, Second: 59},
						Duration: Duration{time.Hour},
					},
				},
			},
			haveErr: false,
		},
		{
			name: "empty location",
			input: `
{
	"entries": [
		{
		  "start_at": "12:34:56",
			"duration": "1h30m"
		},
		{
			"start_at": "23:59:59",
			"duration": "1h0m0s"
		}
	]
}`,
			expected: Config{
				Location: time.UTC,
				Entries: []Entry{
					{
						StartAt:  HMS{Hour: 12, Minute: 34, Second: 56},
						Duration: Duration{time.Hour + 30*time.Minute},
					},
					{
						StartAt:  HMS{Hour: 23, Minute: 59, Second: 59},
						Duration: Duration{time.Hour},
					},
				},
			},
			haveErr: false,
		},
		{
			name: "invalid time format",
			input: `
{
	"entries": [
		{
		  "start_at": "12:34:56:78",
			"duration": "1h30m"
		}
	]
}`,
			expected: Config{},
			haveErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var c Config
			err := json.Unmarshal([]byte(tc.input), &c)
			if tc.haveErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected.Location.String(), c.Location.String())
			require.Equal(t, tc.expected.Entries, c.Entries)
		})
	}
}

func TestConfigHMS(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected HMS
		haveErr  bool
	}{
		{
			name:     "valid",
			input:    `"12:34:56"`,
			expected: HMS{Hour: 12, Minute: 34, Second: 56},
			haveErr:  false,
		},
		{
			name:     "invalid time format",
			input:    `"12:34:56:78"`,
			expected: HMS{},
			haveErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var hms HMS
			err := json.Unmarshal([]byte(tc.input), &hms)
			if tc.haveErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, hms)
		})
	}
}

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
