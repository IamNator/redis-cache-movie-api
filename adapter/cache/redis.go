package cache

import (
	"context"
	redis "github.com/go-redis/redis/v8"
	"strconv"

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

	//opts.TLSConfig = &tls.Config{
	//	InsecureSkipVerify: true,
	//
	//}
	client := redis.NewClient(opts)

	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

var _ ports.ICache = (*RedisCache)(nil)

type (
	cacheTag string
)

func (tag cacheTag) String() string {
	return string(tag)
}

func (tag cacheTag) computeKey(id int) string {
	return tag.String() + ":" + strconv.Itoa(id)
}

func (tag cacheTag) computeParentKey(id int) cacheTag {
	return cacheTag(tag.computeKey(id))
}

const (
	movieTag     cacheTag = "movie"
	characterTag cacheTag = "character"

	// DefaultTTLSec is the default TTL for cache entries
	DefaultTTLSec = 0
)

func (r RedisCache) SetMovies(movies []model.MovieDetails) error {

	pipe := r.client.TxPipeline()

	for _, movie := range movies {
		if err := pipe.Set(context.Background(), movieTag.computeKey(movie.ID), movie, DefaultTTLSec).Err(); err != nil {
			return err
		}
	}

	_, err := pipe.Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r RedisCache) SetMovieByID(id int, movie model.MovieDetails) error {

	if err := r.client.Set(context.Background(), movieTag.computeKey(id), movie, DefaultTTLSec).Err(); err != nil {
		return err
	}

	return nil
}

func (r RedisCache) SetCharactersByMovieID(movieID int, characters []model.Character) error {

	pipe := r.client.TxPipeline()

	for _, character := range characters {
		if err := pipe.Set(context.Background(),
			characterTag.computeParentKey(movieID).computeKey(character.ID),
			character, DefaultTTLSec).Err(); err != nil {
			return err
		}
	}

	_, err := pipe.Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r RedisCache) GetMovies(page, pageSize int) ([]model.MovieDetails, int64, error) {

	keys, count, err := fetchData(context.Background(), r.client, movieTag, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var movies []model.MovieDetails

	var movie model.MovieDetails

	for _, key := range keys {

		if err := r.client.Get(context.Background(), key).Scan(&movie); err != nil {
			return nil, 0, err
		}

		movies = append(movies, movie)
	}

	return movies, count, nil
}

func (r RedisCache) GetMovieByID(id int) (*model.MovieDetails, error) {
	var movie model.MovieDetails

	if err := r.client.Get(context.Background(), movieTag.computeKey(id)).Scan(&movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r RedisCache) GetCharactersByMovieID(movieID int, page, pageSize int) ([]model.Character, int64, error) {

	keys, count, err := fetchData(context.Background(), r.client, characterTag, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var characters []model.Character

	var character model.Character

	for _, key := range keys {

		if err := r.client.Get(context.Background(), key).Scan(&character); err != nil {
			return nil, 0, err
		}

		characters = append(characters, character)
	}

	return characters, count, nil
}

func fetchData(ctx context.Context, client *redis.Client, tag cacheTag, page int, pageSize int) ([]string, int64, error) {
	var keys []string
	var err error

	if page < 1 {
		page = 1
	}

	// Calculate start and end positions based on page and pageSize
	start := (page - 1) * pageSize
	end := start + pageSize - 1

	//Get the total number of keys in the list
	total, err := client.LLen(ctx, tag.String()).Result()
	if err != nil {
		return nil, 0, err
	}

	// Use Redis command to retrieve keys in the specified range
	keys, err = client.LRange(ctx, tag.String()+":*", int64(start), int64(end)).Result()
	if err != nil {
		return nil, 0, err
	}

	return keys, total, nil
}
