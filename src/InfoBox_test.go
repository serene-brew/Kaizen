package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoBoxInitialization(t *testing.T) {
	infoBox := NewInfoBox()

	// Test initial state
	assert.Equal(t, 50, infoBox.width, "Default width should be 50")
	assert.Equal(t, 25, infoBox.height, "Default height should be 25")
	assert.False(t, infoBox.hasAnimeLoaded, "Should not have anime loaded initially")
	assert.False(t, infoBox.focused, "Should not be focused initially")
	assert.Empty(t, infoBox.thumbnailURL, "Thumbnail URL should be empty initially")
}

func TestInfoBoxSetSize(t *testing.T) {
	infoBox := NewInfoBox()

	testCases := []struct {
		name           string
		width          int
		height         int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "Normal dimensions",
			width:          100,
			height:         50,
			expectedWidth:  100,
			expectedHeight: 50,
		},
		{
			name:           "Minimum dimensions",
			width:          10,
			height:         5,
			expectedWidth:  10,
			expectedHeight: 5,
		},
		{
			name:           "Large dimensions",
			width:          200,
			height:         100,
			expectedWidth:  200,
			expectedHeight: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			infoBox.SetSize(tc.width, tc.height)
			assert.Equal(t, tc.expectedWidth, infoBox.width)
			assert.Equal(t, tc.expectedHeight, infoBox.height)
		})
	}
}

func TestInfoBoxAnimeInfo(t *testing.T) {
	infoBox := NewInfoBox()

	testData := struct {
		title        string
		englishName  string
		description  string
		genres       []string
		status       string
		animeType    string
		rating       string
		score        float64
		thumbnailURL string
	}{
		title:        "Test Anime",
		englishName:  "Test Anime EN",
		description:  "Test description",
		genres:       []string{"Action", "Comedy"},
		status:       "Ongoing",
		animeType:    "TV",
		rating:       "PG-13",
		score:        8.5,
		thumbnailURL: "http://example.com/thumb.jpg",
	}

	// Test setting info without thumbnail
	infoBox.SetAnimeInfo(
		testData.title,
		testData.englishName,
		testData.description,
		testData.genres,
		testData.status,
		testData.animeType,
		testData.rating,
		testData.score,
	)

	assert.True(t, infoBox.hasAnimeLoaded)
	assert.Equal(t, testData.title, infoBox.title)
	assert.Equal(t, testData.englishName, infoBox.englishName)
	assert.Equal(t, testData.description, infoBox.description)
	assert.Equal(t, testData.genres, infoBox.genres)
	assert.Equal(t, testData.status, infoBox.status)
	assert.Equal(t, testData.animeType, infoBox.animeType)
	assert.Equal(t, testData.rating, infoBox.rating)
	assert.Equal(t, testData.score, infoBox.score)
	assert.Empty(t, infoBox.thumbnailURL)

	// Test setting info with thumbnail
	infoBox.SetAnimeInfoWithThumbnail(
		testData.title,
		testData.englishName,
		testData.description,
		testData.genres,
		testData.status,
		testData.animeType,
		testData.rating,
		testData.score,
		testData.thumbnailURL,
	)

	assert.Equal(t, testData.thumbnailURL, infoBox.thumbnailURL)
}

func TestInfoBoxFocus(t *testing.T) {
	infoBox := NewInfoBox()

	// Test focus
	infoBox.Focus()
	assert.True(t, infoBox.focused, "InfoBox should be focused after Focus()")

	// Test blur
	infoBox.Blur()
	assert.False(t, infoBox.focused, "InfoBox should not be focused after Blur()")
}
