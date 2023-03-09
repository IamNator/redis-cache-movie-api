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
		return
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

	page := "1 "
	pageSize := "10"

	Test(t,
		Description("GetMovies Success"),
		Get(basePath+"/movies"),
		Send().Custom(func(hit Hit) error {
			hit.Request().URL.Query().Add("page", page)
			hit.Request().URL.Query().Add("pageSize", pageSize)
			return nil
		}),

		Expect().Headers("Content-Type").Equal("application/json"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Contains("data", "message"),
	)
}

// HTTP GET: /movies/{movie_id}
func TestGetMovie(t *testing.T) {

	page := "1"
	pageSize := "10"
	movieID := 1

	Test(t,
		Description("GetMovie Success"),
		Get(basePath+"/%s/%d", "movies", movieID),
		Send().Custom(func(hit Hit) error {
			hit.Request().URL.Query().Add("page", page)
			hit.Request().URL.Query().Add("pageSize", pageSize)
			return nil
		}),
		Expect().Headers("Content-Type").NotEmpty(),
		Expect().Headers("Content-Type").Equal("application/json"),
		Expect().Status().OneOf(http.StatusOK, http.StatusNotFound),
		Expect().Body().JSON().Contains("code", "message"),
	)
}

// HTTP GET: /characters/{movie_id}
func TestGetMovieCharacters(t *testing.T) {

	page := "1"
	pageSize := "10"
	movieID := 1

	Test(t,
		Description("GetMovieCharacters Success"),
		Get(basePath+"/%s/%d", "characters", movieID),
		Send().Custom(func(hit Hit) error {
			hit.Request().URL.Query().Add("page", page)
			hit.Request().URL.Query().Add("pageSize", pageSize)
			return nil
		}),
		Expect().Headers("Content-Type").Equal("application/json"),
		Expect().Status().OneOf(http.StatusOK, http.StatusNotFound),
		Expect().Body().JSON().Contains("code", "message"),
	)
}

// HTTP POST: /comments/{movie_id}
func TestAddCommentToMovie(t *testing.T) {

	movieID := 1
	body := map[string]interface{}{
		"message": "A great movie!",
	}

	Test(t,
		Description("AddCommentToMovie Success"),
		Post(basePath+"/%s/%d", "comments", movieID),
		Send().Headers("Accept").Add("application/json"),
		Send().Body().JSON(body),
		Expect().Headers("Content-Type").Equal("application/json"),
		Expect().Status().OneOf(http.StatusCreated, http.StatusNotFound, http.StatusOK),
		Expect().Body().JSON().Contains("code", "message"),
	)
}

// HTTP GET: /comments/{movie_id}
func TestGetCommentsFromMovie(t *testing.T) {

	movieID := 1
	page := "1"
	pageSize := "10"

	Test(t,
		Description("GetCommentsFromMovie Success"),
		Get(basePath+"/%s/%d", "comments", movieID),
		Send().Custom(func(hit Hit) error {
			hit.Request().URL.Query().Add("page", page)
			hit.Request().URL.Query().Add("pageSize", pageSize)
			return nil
		}),
		Expect().Headers("Content-Type").NotEmpty(),
		Expect().Headers("Content-Type").Equal("application/json"),
		Expect().Status().OneOf(http.StatusOK, http.StatusNotFound),
		Expect().Body().JSON().Contains("message", "code"),
	)
}
