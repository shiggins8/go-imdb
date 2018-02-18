package imdb

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSetOmdbAPIKey(t *testing.T) {
	SetOmdbAPIKey("TESTER12")
	if os.Getenv("OMDB_API_KEY") != "TESTER12" {
		t.Error("Environment variable not set correctly")
	}
}

func TestGetBasicInfoGoodResult(t *testing.T) {
	testMovie := Movie{
		Response: "False",
	}

	if testMovie.GetBasicInfo() != "No movie stored in this struct" {
		t.FailNow()
	}
}

func TestGetBasicInfoBadResult(t *testing.T) {
	testMovie := Movie{
		Title:    "TestMovie",
		Genre:    "Horror",
		Year:     "1994",
		Response: "True",
	}

	if testMovie.GetBasicInfo() != "Movie: TestMovie\nGenre: Horror\nRelease Year: 1994\n" {
		t.FailNow()
	}
}

func TestGetRatings(t *testing.T) {
	testJSON := []byte(`{"Title":"Toy Story","Year":"1995","Rated":"G","Released":"22 Nov 1995","Runtime":"81 min","Genre":"Animation, Adventure, Comedy","Director":"John Lasseter","Writer":"John Lasseter (original story by), Pete Docter (original story by), Andrew Stanton (original story by), Joe Ranft (original story by), Joss Whedon (screenplay by), Andrew Stanton (screenplay by), Joel Cohen (screenplay by), Alec Sokolow (screenplay by)","Actors":"Tom Hanks, Tim Allen, Don Rickles, Jim Varney","Plot":"A cowboy doll is profoundly threatened and jealous when a new spaceman figure supplants him as top toy in a boy's room.","Language":"English","Country":"USA","Awards":"Nominated for 3 Oscars. Another 23 wins & 17 nominations.","Poster":"https://images-na.ssl-images-amazon.com/images/M/MV5BMDU2ZWJlMjktMTRhMy00ZTA5LWEzNDgtYmNmZTEwZTViZWJkXkEyXkFqcGdeQXVyNDQ2OTk4MzI@._V1_SX300.jpg","Ratings":[{"Source":"Internet Movie Database","Value":"8.3/10"},{"Source":"Rotten Tomatoes","Value":"100%"},{"Source":"Metacritic","Value":"95/100"}],"Metascore":"95","imdbRating":"8.3","imdbVotes":"714,261","imdbID":"tt0114709","Type":"movie","DVD":"20 Mar 2001","BoxOffice":"N/A","Production":"Buena Vista","Website":"http://www.disney.com/ToyStory","Response":"True"}`)
	testMovie := Movie{}
	json.Unmarshal(testJSON, &testMovie)

	temp, _ := testMovie.GetRatings("Metacritic")
	if temp != "95/100" {
		t.FailNow()
	}

	temp, _ = testMovie.GetRatings("Internet Movie Database")
	if temp != "8.3/10" {
		t.FailNow()
	}

	temp, _ = testMovie.GetRatings("Rotten Tomatoes")
	if temp != "100%" {
		t.FailNow()
	}

	temp, _ = testMovie.GetRatings("Bogus Site")
	if temp != "error" {
		t.FailNow()
	}
}

func TestCheckValidAPIKeyAlreadySet(t *testing.T) {
	APIKey = "12345678"
	err := checkValidAPIKey()
	if err != nil {
		t.FailNow()
	}
}

func TestCheckValidAPIKeyValidSet(t *testing.T) {
	APIKey = ""
	os.Setenv("OMDB_API_KEY", "12345678")
	err := checkValidAPIKey()
	if err != nil {
		t.FailNow()
	}
}

func TestCheckValidAPIKeyInvalidSet(t *testing.T) {
	APIKey = ""
	os.Setenv("OMDB_API_KEY", "")
	checkValidAPIKey()
	err := checkValidAPIKey()
	if err == nil {
		t.FailNow()
	}
}

func TestFetchMovieNoAPIKeySet(t *testing.T) {
	os.Setenv("OMDB_API_KEY", "")
	_, err := FetchMovie("Toy Story")
	if err == nil {
		t.FailNow()
	}
}

func TestFetchMovieValidAPIKeySet(t *testing.T) {
	os.Setenv("OMDB_API_KEY", "12345678")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Title":"Toy Story","Year":"1995","Rated":"G","Released":"22 Nov 1995","Runtime":"81 min","Genre":"Animation, Adventure, Comedy","Director":"John Lasseter","Writer":"John Lasseter (original story by), Pete Docter (original story by), Andrew Stanton (original story by), Joe Ranft (original story by), Joss Whedon (screenplay by), Andrew Stanton (screenplay by), Joel Cohen (screenplay by), Alec Sokolow (screenplay by)","Actors":"Tom Hanks, Tim Allen, Don Rickles, Jim Varney","Plot":"A cowboy doll is profoundly threatened and jealous when a new spaceman figure supplants him as top toy in a boy's room.","Language":"English","Country":"USA","Awards":"Nominated for 3 Oscars. Another 23 wins & 17 nominations.","Poster":"https://images-na.ssl-images-amazon.com/images/M/MV5BMDU2ZWJlMjktMTRhMy00ZTA5LWEzNDgtYmNmZTEwZTViZWJkXkEyXkFqcGdeQXVyNDQ2OTk4MzI@._V1_SX300.jpg","Ratings":[{"Source":"Internet Movie Database","Value":"8.3/10"},{"Source":"Rotten Tomatoes","Value":"100%"},{"Source":"Metacritic","Value":"95/100"}],"Metascore":"95","imdbRating":"8.3","imdbVotes":"714,261","imdbID":"tt0114709","Type":"movie","DVD":"20 Mar 2001","BoxOffice":"N/A","Production":"Buena Vista","Website":"http://www.disney.com/ToyStory","Response":"True"}`))
	}))
	defer ts.Close()

	OmdbURL = ts.URL + "/?apikey=%s&t=%s"
	movie, err := FetchMovie("Toy Story")
	if err != nil {
		t.Error(err)
	}
	if movie.Title != "Toy Story" {
		t.Errorf("Expected title: Toy Story, actual title: %s", movie.Title)
	}
	if movie.Year != "1995" {
		t.Errorf("Expected year: 1995, actual year: %s", movie.Year)
	}
	if movie.Response != "True" {
		t.Error("Invalid response type")
	}
}

func TestFetchMovieResponseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Fake Error Code", 303)
	}))
	defer ts.Close()

	OmdbURL = ts.URL + "/?apikey=%s&t=%s"
	movie, err := FetchMovie("Toy Story")
	if err == nil {
		t.Errorf("FetchMovie() didn’t return an error with bad server response")
	}
	if movie != nil {
		t.Error("Should not have returned a movie struct with an invalid response")
	}
}

func TestFetchMovieDecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	OmdbURL = ts.URL + "/?apikey=%s&t=%s"
	movie, err := FetchMovie("Toy Story")
	if err == nil {
		t.Errorf("FetchMovie() didn’t return an error with bad server response")
	}
	if movie != nil {
		t.Error("Should not have returned a movie struct with an invalid response")
	}
}
