package cache

import (
	"context"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
	"github.com/rueian/rueidis"

	goredis "github.com/gomodule/redigo/redis"
	"github.com/iamnator/movie-api/model"
	"github.com/iamnator/movie-api/service/ports"
)

type RedisCache struct {
	characterIndex *redisearch.Client
	movieIndex     *redisearch.Client
	rueidisClient  rueidis.Client
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

	if err := createMovieSchema(context.Background(), redisClient); err != nil {
		return nil, err
	}

	if err := createCharacterSchema(context.Background(), redisClient); err != nil {
		return nil, err
	}

	clientNN, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{opts.Addr},
	})
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		characterIndex: getRedisSearchClient(pool, "idx:characters"),
		movieIndex:     getRedisSearchClient(pool, "idx:movies"),
		rueidisClient:  clientNN,
	}, nil
}

func getRedisSearchClient(pool *goredis.Pool, index string) *redisearch.Client {
	return redisearch.NewClientFromPool(pool, index)
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
	if err := r.movieIndex.IndexOptions(redisearch.IndexingOptions{
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
	if err := r.movieIndex.IndexOptions(redisearch.DefaultIndexingOptions, doc); err != nil {
		return err
	}

	return nil
}

func (r RedisCache) SetCharactersByMovieID(movieID int, characters []model.Character) error {

	// Create a document from movies
	var docs redisearch.DocumentList
	var doc redisearch.Document

	for _, character := range characters {
		doc = redisearch.NewDocument(computeCharacterKey(character.ID), 1.0).
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

	//Add the document to the index
	if err := r.characterIndex.IndexOptions(opts, docs...); err != nil {
		return err
	}

	return nil
}

///
//
//   					GETTERS
//
//

func parse_movie_data(text string) []map[string]interface{} {

	// Regular expression to match each movie record
	pattern := regexp.MustCompile(`\[(\d+)\s+movie:(\d+)\s+\[(.*?)\]\s*\]`)
	movies := []map[string]interface{}{}

	// Loop through all matches of the pattern in the text
	for _, match := range pattern.FindAllStringSubmatch(text, -1) {

		//movie_id := parseInt(match[2])
		// Split the movie record into key-value pairs
		pairs_str := strings.Split(match[3], ",")
		pairs := make([][]string, len(pairs_str))
		for i, s := range pairs_str {
			pairs[i] = strings.SplitN(strings.TrimSpace(s), " ", 2)
		}
		// Convert release date and created/updated timestamps to time.Time objects
		for i := range pairs {
			key, value := pairs[i][0], pairs[i][1]
			if key == "release_date" || key == "created_at" || key == "updated_at" {
				t, err := time.Parse("2006-01-02T15:04:05Z", value[:len(value)-1])
				if err == nil {
					pairs[i][1] = t.UTC().String()
				}
			}
		}
		// Create a map for the movie record
		movie := map[string]interface{}{}
		for _, pair := range pairs {
			key, value := pair[0], pair[1]
			//switch key {
			//case "title", "overview", "poster_path", "homepage":
			//	movie[key] = value
			//case "release_date", "created_at", "updated_at":
			//	movie[key] = value
			//default:
			// ignore unknown key
			//}
			movie[key] = value
		}

		// Add the movie map to the list of movies
		movies = append(movies, movie)
	}

	return movies
}

//
//func parseInt(s string) int {
//	var n int
//	_, err := fmt.Sscan(s, &n)
//	if err != nil {
//		return 0
//	}
//	return n
//}

func (r RedisCache) GetMovies(page, pageSize int) ([]model.MovieDetails, int64, error) {

	if page < 1 {
		page = 1
	}

	qq := redisearch.NewQuery("*").
		Limit((page-1)*pageSize, pageSize).
		SetSortBy("release_date", false)

	docs, count, err := r.movieIndex.Search(qq)
	if err != nil {
		return nil, 0, err
	}

	var movies []model.MovieDetails

	for _, doc := range docs {
		var movie model.MovieDetails

		//movie.ID = doc.Properties["id"].(int)
		movie.Name = doc.Properties["name"].(string)
		movie.ReleaseDate, _ = time.Parse(time.RFC3339, doc.Properties["release_date"].(string))
		movie.Director = doc.Properties["director"].(string)
		movie.Producer = doc.Properties["producer"].(string)
		movie.OpeningCrawl = doc.Properties["opening_crawl"].(string)
		movie.CreatedAt, _ = time.Parse(time.RFC3339, doc.Properties["created_at"].(string))
		movie.UpdatedAt, _ = time.Parse(time.RFC3339, doc.Properties["updated_at"].(string))

		movies = append(movies, movie)
	}

	log.Info().Msgf("movies: %v", len(movies))

	return movies, int64(count), nil
}

func (r RedisCache) GetMovieByID(id int) (*model.MovieDetails, error) {
	var movie model.MovieDetails

	docs, err := r.movieIndex.Get(computeMovieKey(id))
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

	query := redisearch.NewQuery("*").
		Limit((page-1)*pageSize, pageSize).
		SetSortBy("name", false)

	docs, count, err := r.characterIndex.Search(query)
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
