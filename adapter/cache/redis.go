package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/iamnator/movie-api/model"
	"net/http"
	"time"
)

type RedisCache struct {
}

//go:generate mockgen -destination=../mocks/cache.go -package=mocks github.com/iamnator/movie-api/adapter/cache ICache
type ICache interface {
	SetMovies([]model.MovieDetails) error
	SetMovieByID(id int, movie model.MovieDetails) error
	SetCharactersByMovieID(id int, characters []model.Character) error

	GetMovies() ([]model.MovieDetails, error)
	GetMovieByID(id int) (*model.MovieDetails, error)
	GetCharactersByMovieID(id int) ([]model.Character, error)
}

func (r RedisCache) SetMovies(movies []model.MovieDetails) error {
	return nil
}

func (r RedisCache) SetMovieByID(id int, movie model.MovieDetails) error {
	return nil
}

func (r RedisCache) SetCharactersByMovieID(id int, characters []model.Character) error {
	return nil
}

func (r RedisCache) GetMovies() ([]model.MovieDetails, error) {
	return []model.MovieDetails{}, nil
}

func (r RedisCache) GetMovieByID(id int) (*model.MovieDetails, error) {
	return nil, nil
}

func (r RedisCache) GetCharactersByMovieID(id int) ([]model.Character, error) {
	return []model.Character{}, nil
}

func getMovies(redisClient *redis.Client) ([]model.MovieDetails, error) {

	movieData, err := redisClient.Get("movies").Result()
	if err == redis.Nil {
		// Cache miss, fetch data from swapi.dev
		resp, err := http.Get("https://swapi.dev/api/films/")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var data struct {
			Results []model.MovieDetails `json:"results"`
		}

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}

		movieData, err := json.Marshal(data.Results)
		if err != nil {
			return nil, err
		}

		redisClient.Set("movies", movieData, 1*time.Hour)

	} else if err != nil {
		return nil, err
	}

	var movies []model.MovieDetails
	err = json.Unmarshal([]byte(movieData), &movies)
	if err != nil {
		return nil, err
	}

	return movies, nil

}

func getCommentCount(redisClient *redis.Client, movieID int) (int, error) {
	commentCount, err := redisClient.Get(fmt.Sprintf("comment_count:%d", movieID)).Int()
	if err == redis.Nil {
		// Cache miss, fetch data from
		return commentCount, nil
	}

	return 1, nil
}
