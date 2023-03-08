package cache

import (
	"context"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
	goredis "github.com/gomodule/redigo/redis"
	"github.com/iamnator/movie-api/model"
	"github.com/iamnator/movie-api/service/ports"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type RedisCache struct {
	redisearchClient *redisearch.Client
	//redisClient      *redis.Client
}

func NewRedisCache(url string) (*RedisCache, error) {

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(opts)

	if _, err := redisClient.Ping(context.TODO()).Result(); err != nil {
		return nil, err
	}

	pool := &goredis.Pool{Dial: func() (goredis.Conn, error) {
		return goredis.Dial(opts.Network, opts.Addr, goredis.DialPassword(opts.Password))
	}}

	// Create a RedisSearch redis-searchClient
	client := redisearch.NewClientFromPool(pool, "busha_movie_api")

	if err := createMovieSchema(client); err != nil {
		return nil, err
	}

	if err := createCharacterSchema(client); err != nil {
		return nil, err
	}

	return &RedisCache{redisearchClient: client}, nil
}

var _ ports.ICache = (*RedisCache)(nil)

func (r RedisCache) SetMovies(movies []model.MovieDetails) error {

	// Create a document from movies

	var docs redisearch.DocumentList
	var doc redisearch.Document

	for _, movie := range movies {

		doc = redisearch.NewDocument(computeMovieKey(movie.ID), 1.0).
			Set("id", movie.ID).
			Set("name", movie.Name).
			Set("release_date", movie.ReleaseDate.UTC().Format(time.RFC3339)).
			Set("director", movie.Director).
			Set("producer", movie.Producer).
			Set("opening_crawl", movie.OpeningCrawl).
			Set("created_at", movie.CreatedAt.UTC().Format(time.RFC3339)).
			Set("updated_at", movie.UpdatedAt.UTC().Format(time.RFC3339))

		docs = append(docs, doc)
	}

	// Add the document to the index
	if err := r.redisearchClient.IndexOptions(redisearch.IndexingOptions{
		Language:         redisearch.DefaultIndexingOptions.Language,
		NoSave:           redisearch.DefaultIndexingOptions.NoSave,
		Replace:          true,
		Partial:          true,
		ReplaceCondition: redisearch.DefaultIndexingOptions.ReplaceCondition, // replace only if the document exists
	}, docs...); err != nil {
		return err
	}

	return nil
}

func (r RedisCache) SetMovieByID(movieID int, movie model.MovieDetails) error {

	// Create a document from movies
	doc := redisearch.NewDocument(computeMovieKey(movieID), 1.0)

	doc.
		Set("id", movieID).
		Set("name", movie.Name).
		Set("release_date", movie.ReleaseDate.UTC().Format(time.RFC3339)).
		Set("director", movie.Director).
		Set("producer", movie.Producer).
		Set("opening_crawl", movie.OpeningCrawl).
		Set("created_at", movie.CreatedAt.UTC().Format(time.RFC3339)).
		Set("updated_at", movie.UpdatedAt.UTC().Format(time.RFC3339))

	// Add the document to the index
	if err := r.redisearchClient.IndexOptions(redisearch.DefaultIndexingOptions, doc); err != nil {
		return err
	}

	return nil
}

func (r RedisCache) SetCharactersByMovieID(movieID int, characters []model.Character) error {

	// Create a document from movies
	var docs redisearch.DocumentList
	var doc redisearch.Document

	for _, character := range characters {
		doc = redisearch.NewDocument(computeCharacterKey(character.MovieID, character.ID), 1.0).
			Set("id", character.ID).
			Set("name", character.Name).
			Set("movie_id", character.MovieID).
			Set("gender", character.Gender).
			Set("height_cm", character.HeightCm)

		docs = append(docs, doc)
	}

	opts := redisearch.DefaultIndexingOptions
	opts.Replace = true
	opts.Partial = true

	// Add the document to the index
	if err := r.redisearchClient.IndexOptions(opts, docs...); err != nil {
		return err
	}

	return nil
}

func (r RedisCache) GetMovies(page, pageSize int) ([]model.MovieDetails, int64, error) {

	var movies []model.MovieDetails

	if page < 1 {
		page = 1
	}

	query := redisearch.NewQuery("release_date")

	docs, count, err := r.redisearchClient.Search(query)
	if err != nil {
		return movies, 0, err
	}

	for _, doc := range docs {
		var movie model.MovieDetails
		movie.ID = doc.Properties["id"].(int)
		movie.Name = doc.Properties["name"].(string)
		movie.ReleaseDate, _ = time.Parse(time.RFC3339, doc.Properties["release_date"].(string))
		movie.Director = doc.Properties["director"].(string)
		movie.Producer = doc.Properties["producer"].(string)
		movie.OpeningCrawl = doc.Properties["opening_crawl"].(string)
		movie.CreatedAt, _ = time.Parse(time.RFC3339, doc.Properties["created_at"].(string))
		movie.UpdatedAt, _ = time.Parse(time.RFC3339, doc.Properties["updated_at"].(string))

		movies = append(movies, movie)
	}

	log.Info().Msgf("Found %d movies", count)

	return movies, int64(count), nil
}

func (r RedisCache) GetMovieByID(id int) (*model.MovieDetails, error) {
	var movie model.MovieDetails

	docs, err := r.redisearchClient.Get(computeMovieKey(id))
	if err != nil {
		return nil, err
	}

	movie.ID = docs.Properties["id"].(int)
	movie.Name = docs.Properties["name"].(string)
	movie.ReleaseDate, _ = time.Parse(time.RFC3339, docs.Properties["release_date"].(string))
	movie.Director = docs.Properties["director"].(string)
	movie.Producer = docs.Properties["producer"].(string)
	movie.OpeningCrawl = docs.Properties["opening_crawl"].(string)
	movie.CreatedAt, _ = time.Parse(time.RFC3339, docs.Properties["created_at"].(string))
	movie.UpdatedAt, _ = time.Parse(time.RFC3339, docs.Properties["updated_at"].(string))

	return &movie, nil
}

func (r RedisCache) GetCharactersByMovieID(movieID int, page, pageSize int) ([]model.Character, int64, error) {

	query := redisearch.NewQuery(computeCharacterKey(strconv.Itoa(movieID), "*")).
		Limit((page-1)*pageSize, pageSize).
		SetSortBy("name", false)

	docs, count, err := r.redisearchClient.Search(query)
	if err != nil {
		return nil, 0, err
	}

	var characters []model.Character
	var character model.Character

	for _, doc := range docs {

		character = model.Character{
			ID:       doc.Properties["id"].(int),
			Name:     doc.Properties["name"].(string),
			MovieID:  doc.Properties["movie_id"].(int),
			Gender:   doc.Properties["gender"].(string),
			HeightCm: doc.Properties["height_cm"].(int),
		}

		characters = append(characters, character)
	}

	return characters, int64(count), nil
}
