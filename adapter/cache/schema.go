package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

const (
	MovieIndexName     = "idx:movies"
	CharacterIndexName = "idx:characters"
)

func createMovieSchema(ctx context.Context, client *redis.Client) error {

	//drop the index if it exists
	_ = client.Do(ctx, "FT.DROPINDEX", MovieIndexName, "DD").Err()

	err := client.Do(ctx, "FT.CREATE", MovieIndexName, "ON", "HASH", "PREFIX", "1", "movie:", "SCHEMA",
		"name", "TEXT",
		"id", "NUMERIC",
		"episode_id", "NUMERIC",
		"opening_crawl", "TEXT",
		"director", "TEXT",
		"producer", "TEXT",
		"release_date", "TEXT", "WEIGHT", "5.0", "SORTABLE",
		"created_at", "TEXT",
		"updated_at", "TEXT",
	).
		Err()

	if err != nil {
		return err
	}

	return nil
}

func createCharacterSchema(ctx context.Context, client *redis.Client) error {

	//drop the index if it exists
	_ = client.Do(ctx, "FT.DROPINDEX", CharacterIndexName, "DD").Err()

	err := client.Do(ctx, "FT.CREATE", CharacterIndexName, "ON", "HASH", "PREFIX", "1", "character:", "SCHEMA",
		"id", "NUMERIC",
		"name", "TEXT", "WEIGHT", "5.0", "SORTABLE",
		"movie_id", "TAG", "SORTABLE",
		"gender", "TAG", "SORTABLE",
		"height_cm", "NUMERIC", "SORTABLE",
	).
		Err()
	if err != nil {
		return err
	}

	return nil
}
