package cache

import (
	"context"
	"crypto/tls"
	redis "github.com/go-redis/redis/v8"

	"github.com/iamnator/movie-api/model"
	"github.com/iamnator/movie-api/service/ports"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(url string) (*RedisCache, error) {

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	opts.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	client := redis.NewClient(opts)

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

var _ ports.ICache = (*RedisCache)(nil)

func (r RedisCache) SetMovies(movies []model.MovieDetails) error {
	return nil
}

func (r RedisCache) SetMovieByID(id int, movie model.MovieDetails) error {
	return nil
}

func (r RedisCache) SetCharactersByMovieID(id int, characters []model.Character) error {
	return nil
}

func (r RedisCache) GetMovies(page, pageSize int) ([]model.MovieDetails, int64, error) {
	return []model.MovieDetails{}, 0, nil
}

func (r RedisCache) GetMovieByID(id int) (*model.MovieDetails, error) {
	return nil, nil
}

func (r RedisCache) GetCharactersByMovieID(movieID int, page, pageSize int) ([]model.Character, int64, error) {
	return []model.Character{}, 0, nil
}
