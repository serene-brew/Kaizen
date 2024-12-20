package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

/*
AnimeResponse is a struct that represents the response format for anime search queries.
It contains a single field, Result, which is a slice of slices containing interface{} elements.
This is used to store the search results from the API response.
*/
type AnimeResponse struct {
	Result [][]interface{} `json:"result"`
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

resp -> string(animeID), string(animeName), int(subEpisodes), int(dubEpisodes) -> [][]interface{}
*/
func extractInfo(query string) ([][]interface{}, error) {
	apiUrl := "https://heavenscape.vercel.app/api/anime/search/" + strings.ReplaceAll(query, " ", "+")

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response AnimeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Result, nil
}

/*
getStreamLink is a function that retrieves the direct streaming link for a specific anime episode.
It takes the anime ID, episode type (e.g., "sub" or "dub"), and episode number as arguments.
The function constructs the API URL, sends an HTTP GET request, and parses the JSON response to extract the link.
If an error occurs at any point, it is returned along with an empty string.

resp -> string [Stream link]
*/
func getStreamLink(id string, espisodeType string, episodeNumber string) (string, error) {
	apiUrl := "https://heavenscape.vercel.app/api/anime/search/" + id + "/" + espisodeType + "/" + episodeNumber

	resp, err := http.Get(apiUrl)
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

