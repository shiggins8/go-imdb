// Package imdb is a wrapper around the OMDb API that will allow people to quickly
// query and get relevant information about Movies within their golang scripts.
// While this version only supports basic Movie functionality, continuing support
// will be added very soon.
package imdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// Movie representation returned by the OMDB web API. This mapping allows us to
// populate the struct based on information returned by the web query
type Movie struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	Rated      string `json:"Rated"`
	Released   string `json:"Released"`
	Runtime    string `json:"Runtime"`
	Genre      string `json:"Genre"`
	Director   string `json:"Director"`
	Writer     string `json:"Writer"`
	Actors     string `json:"Actors"`
	Plot       string `json:"Plot"`
	Language   string `json:"Language"`
	Country    string `json:"Country"`
	Awards     string `json:"Awards"`
	Poster     string `json:"Poster"`
	Ratings    []interface{}
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	DVD        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
}

// APIKey contains the API key for OMDB. Free tier, 1000 requests per day
var APIKey = os.Getenv("OMDB_API_KEY")

// OmdbURL contains the base URL for the API, formatted with two potential strings
// that will be included for requests: the API key and the title of the movie
const OmdbURL = "http://www.omdbapi.com/?apikey=%s&t=%s"

// SetOmdbAPIKey can be used to set the API key environment variable from within the go script, so that the
// user does not have to worry about manually setting the variable themselves before running the script.
func SetOmdbAPIKey(key string) {
	os.Setenv("OMDB_API_KEY", key)
	APIKey = os.Getenv("OMDB_API_KEY")
}

// FetchMovie takes 1 parameter (string, the name of the movie). Then it will
// query the OMDB API, Unmarshal the response data into a Movie struct, and
// then return a pointer of that struct
func FetchMovie(title string) (*Movie, error) {
	err := checkValidAPIKey()
	if err != nil {
		return nil, fmt.Errorf("no valid API key set")
	}
	response, err := http.Get(fmt.Sprintf(OmdbURL, APIKey, url.QueryEscape(title)))
	if err != nil {
		return nil, err
	}

	result := Movie{}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBasicInfo will return the most basic info about a Movie struct
func (m Movie) GetBasicInfo() string {
	if m.Response == "True" {
		return (fmt.Sprintf("Movie: %s\nGenre: %s\nRelease Year: %s\n", m.Title, m.Genre, m.Year))
	}

	return "No movie stored in this struct"
}

// GetRatings will return the rating of the movie based on a provided source string.
// Due to various formats of ratings, value is returned as a string. Formats for
// ratings are included with sources below:
//
// Valid sources include:
//
// - Metacritic (97/100)
//
// - Internet Movie Database (9.7/10)
//
// - Rotten Tomatoes (97%)
//
// An invalid source will return an error
func (m Movie) GetRatings(source string) (string, error) {
	for _, d := range m.Ratings {
		dict := d.(map[string]interface{})
		if dict["Source"].(string) == source {
			return dict["Value"].(string), nil
		}
	}
	return "error", fmt.Errorf("invalid source provided: %s", source)
}

func checkValidAPIKey() error {
	if APIKey == "" {
		APIKey = os.Getenv("OMDB_API_KEY")
		if APIKey == "" {
			return fmt.Errorf("please set the environment variable OMDB_API_KEY to a valid key before utilizing the imdb go package")
		}
		return nil
	}
	return nil
}
