package integration_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
)

const (
	// Attempts connection
	host       = "0.0.0.0:9500"
	healthPath = "http://" + host + "/health"
	attempts   = 20

	// HTTP REST
	basePath = "http://" + host
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

// HTTP GET: /movies
func TestGetMovies(t *testing.T) {

	page := 1
	pageSize := 10

	Test(t,
		Description("GetMovies Success"),
		Get("http://%s/%s?page=%d&pageSize=%d ", host, "movies", page, pageSize),
		Expect().Headers("Content-Type").NotEmpty(),
		Expect().Headers("Content-Type").Equal("application/json"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Contains("data", "message"),
	)
}

// HTTP GET: /movies/{movie_id}
func TestGetMovie(t *testing.T) {

	page := 1
	pageSize := 10
	movieID := 1

	Test(t,
		Description("GetMovie Success"),
		Get("http://%s/%s/%d?page=%d&pageSize=%d ", host, "movies", movieID, page, pageSize),
		Expect().Headers("Content-Type").NotEmpty(),
		Expect().Headers("Content-Type").Equal("application/json"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Contains("data", "message"),
	)
}

// HTTP GET: /characters/{movie_id}
func TestGetMovieCharacters(t *testing.T) {

	page := 1
	pageSize := 10
	movieID := 1

	Test(t,
		Description("GetMovieCharacters Success"),
		Get("http://%s/%s/%d?page=%d&pageSize=%d ", host, "characters", movieID, page, pageSize),
		Expect().Headers("Content-Type").NotEmpty(),
		Expect().Headers("Content-Type").Equal("application/json"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Contains("data", "message"),
	)
}

// HTTP POST: /comments/{movie_id}
func TestAddCommentToMovie(t *testing.T) {

	movieID := 1
	body := `{
		"message": "A great Movie!"
	}`

	Test(t,
		Description("AddCommentToMovie Success"),
		Post("http://%s/%s/%d", host, "comments", movieID),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Headers("Content-Type").NotEmpty(),
		Expect().Headers("Content-Type").Equal("application/json"),
		Expect().Status().Equal(http.StatusCreated),
		Expect().Body().JSON().Contains("code", "message", "data"),
	)
}

// HTTP GET: /comments/{movie_id}
func TestGetCommentsFromMovie(t *testing.T) {

	movieID := 1
	page := 1
	pageSize := 10

	Test(t,
		Description("GetCommentsFromMovie Success"),
		Get("http://%s/%s/%d?page=%d&pageSize=%d ", host, "comments", movieID, page, pageSize),
		Expect().Headers("Content-Type").NotEmpty(),
		Expect().Headers("Content-Type").Equal("application/json"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Contains("data", "message", "code"),
	)
}
