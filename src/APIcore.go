package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

/*
AnimeResponse is a struct that represents the response format for anime search queries.
It contains a single field, Result, which is a slice of slices containing interface{} elements.
This is used to store the search results from the API response.
*/
type Anime struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Thumbnail   string   `json:"thumbnail"`
	SubCount    float64  `json:"subCount"`
	DubCount    float64  `json:"dubCount"`
	Episodes    Episodes `json:"episodes"`
	EnglishName string   `json:"englishName"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	Status      string   `json:"status"`
	Type        string   `json:"type"`
	Rating      string   `json:"rating"`
	Score       float64  `json:"score"`
}

// Episodes represents the subtitled and dubbed episode lists
type Episodes struct {
	Sub []string `json:"sub"`
	Dub []string `json:"dub"`
}

// ApiResponse represents the overall API response structure
type AnimeResponse struct {
	Result []Anime `json:"result"`
}

/*
StreamUtils is a struct that represents the response format for retrieving a direct streaming link.
It contains a single field, Link, which holds the direct URL to the stream.
*/
type StreamUtils struct {
	Link string `json:"direct"`
}

/*
extractInfo is a function that fetches information about an anime based on a given query string.
The query string is used to build the API URL, and an HTTP GET request is sent to fetch the data.
The function parses the JSON response and returns the Result field as a slice of slices of interface{}.
If an error occurs at any stage, it is returned.

resp -> string(animeID), string(animeName), float64(subEpisodes), float64(dubEpisodes), []string, []string, string(englishName), string(description), []string(genres), string(status), string(type), string(rating) -> [][]interface{}
*/
func extractInfo(query string) ([][]any, error) {
	apiURL := "https://heavenscape.vercel.app/api/anime/search/" + strings.ReplaceAll(query, " ", "+")
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the JSON response into ApiResponse struct
	var apiResponse AnimeResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Process the data into [][]interface{}
	var result [][]any
	for _, anime := range apiResponse.Result {
		// Preserve the original column ordering expected elsewhere in the app
		// (id, title, subCount, dubCount, subEpisodes, dubEpisodes, englishName,
		// description, genres, status, type, rating, score) and append thumbnail
		// as the last column so we don't break existing index-based access.
		row := []any{
			anime.ID,
			anime.Title,
			anime.SubCount,
			anime.DubCount,
			anime.Episodes.Sub,
			anime.Episodes.Dub,
			anime.EnglishName,
			anime.Description,
			anime.Genres,
			anime.Status,
			anime.Type,
			anime.Rating,
			anime.Score,
			anime.Thumbnail,
		}
		result = append(result, row)
	}

	return result, nil
}

/*
getStreamLink is a function that retrieves the direct streaming link for a specific anime episode.
It takes the anime ID, episode type (e.g., "sub" or "dub"), and episode number as arguments.
The function constructs the API URL, sends an HTTP GET request, and parses the JSON response to extract the link.
If an error occurs at any point, it is returned along with an empty string.

resp -> string [Stream link]
*/
func getStreamLink(id string, espisodeType string, episodeNumber string) (string, error) {
	apiURL := "https://heavenscape.vercel.app/api/anime/search/" + id + "/" + espisodeType + "/" + episodeNumber

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response StreamUtils
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Link, nil
}
