package cache

import (
	"github.com/RediSearch/redisearch-go/redisearch"
)

func createMovieSchema(client *redisearch.Client) error {

	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("name")).
		AddField(redisearch.NewNumericField("id")).
		AddField(redisearch.NewNumericField("episode_id")).
		AddField(redisearch.NewTextField("opening_crawl")).
		AddField(redisearch.NewTextField("director")).
		AddField(redisearch.NewTextField("producer")).
		AddField(redisearch.NewTextFieldOptions("release_date", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewSortableTextField("created_at", 1)).
		AddField(redisearch.NewSortableTextField("updated_at", 1))

	// Create a RedisSearch index definition
	indexDefinition := redisearch.NewIndexDefinition().AddPrefix("movie:")

	client.DropIndex(false)

	// Create the RedisSearch index
	err := client.CreateIndexWithIndexDefinition(schema, indexDefinition)
	if err != nil {
		return err
	}

	return nil
}

func createCharacterSchema(client *redisearch.Client) error {

	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewNumericField("id")).
		AddField(redisearch.NewTextFieldOptions("name", redisearch.TextFieldOptions{Weight: 5.0})).
		AddField(redisearch.NewNumericFieldOptions("movie_id", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewSortableTextField("gender", 5)).
		AddField(redisearch.NewNumericField("height_cm"))

	// Create a RedisSearch index definition
	indexDefinition := redisearch.NewIndexDefinition().AddPrefix("character:")

	client.DropIndex(false)

	// Create the RedisSearch index
	err := client.CreateIndexWithIndexDefinition(schema, indexDefinition)
	if err != nil {
		return err
	}

	return nil
}
