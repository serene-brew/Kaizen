package src

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnimeStructSerialization(t *testing.T) {
	// Test anime JSON serialization/deserialization
	testCases := []struct {
		name     string
		jsonData string
		want     Anime
	}{
		{
			name: "Complete Anime Data",
			jsonData: `{
				"id": "123",
				"title": "Test Anime",
				"thumbnail": "http://example.com/thumb.jpg",
				"subCount": 12,
				"dubCount": 10,
				"episodes": {
					"sub": ["1", "2", "3"],
					"dub": ["1", "2"]
				},
				"englishName": "Test Anime EN",
				"description": "Test description",
				"genres": ["Action", "Comedy"],
				"status": "Ongoing",
				"type": "TV",
				"rating": "PG-13",
				"score": 8.5
			}`,
			want: Anime{
				ID:          "123",
				Title:       "Test Anime",
				Thumbnail:   "http://example.com/thumb.jpg",
				SubCount:    12,
				DubCount:    10,
				Episodes:    Episodes{Sub: []string{"1", "2", "3"}, Dub: []string{"1", "2"}},
				EnglishName: "Test Anime EN",
				Description: "Test description",
				Genres:      []string{"Action", "Comedy"},
				Status:      "Ongoing",
				Type:        "TV",
				Rating:      "PG-13",
				Score:       8.5,
			},
		},
		{
			name: "Minimal Anime Data",
			jsonData: `{
				"id": "456",
				"title": "Minimal Anime",
				"subCount": 1,
				"dubCount": 0,
				"episodes": {
					"sub": ["1"],
					"dub": []
				}
			}`,
			want: Anime{
				ID:       "456",
				Title:    "Minimal Anime",
				SubCount: 1,
				DubCount: 0,
				Episodes: Episodes{Sub: []string{"1"}, Dub: []string{}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got Anime
			err := json.Unmarshal([]byte(tc.jsonData), &got)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(got)
			assert.NoError(t, err)

			var unmarshaled Anime
			err = json.Unmarshal(marshaled, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, unmarshaled)
		})
	}
}
