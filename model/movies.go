package model

import "time"

type (
	MovieDetails struct {
		ID           int       `json:"id"`
		Title        string    `json:"title"`
		EpisodeID    int       `json:"episode_id"`
		OpeningCrawl string    `json:"opening_crawl"`
		Director     string    `json:"director"`
		Producer     string    `json:"producer"`
		ReleaseDate  time.Time `json:"release_date"`
	}
)

type Movie struct {
	SwapiMovieID int    `json:"id"`
	Name         string `json:"name"`
	OpeningCrawl string `json:"opening_crawl"`
	CommentCount int64  `json:"comment_count"`
}
