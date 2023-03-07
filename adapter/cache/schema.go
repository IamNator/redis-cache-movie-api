package cache

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/rs/zerolog/log"
)

func createMovieSchema(client *redisearch.Client) error {

	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("name")).
		AddField(redisearch.NewNumericField("id")).
		AddField(redisearch.NewNumericField("episode_id")).
		AddField(redisearch.NewTextField("opening_crawl")).
		AddField(redisearch.NewTextField("director")).
		AddField(redisearch.NewTextField("producer")).
		AddField(redisearch.NewTextFieldOptions("release_date", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextField("created_at")).
		AddField(redisearch.NewTextField("updated_at"))

	// Create a RedisSearch index definition
	indexDefinition := redisearch.NewIndexDefinition().AddPrefix("movie:")

	if er := client.Drop(); er != nil {
		log.Info().Msgf("Error dropping index: %v", er)
	}

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

	client.Drop()

	// Create the RedisSearch index
	err := client.CreateIndexWithIndexDefinition(schema, indexDefinition)
	if err != nil {
		return err
	}

	return nil
}
