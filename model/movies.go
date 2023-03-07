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
	SwapiMovieID int    `json:"id" example:"1" swaggertype:"integer"`
	Name         string `json:"name" example:"A New Hope" swaggertype:"string"`
	OpeningCrawl string `json:"opening_crawl" example:"It is a period of civil war. Rebel spaceships, striking from a hidden base, have won their first victory against the evil Galactic Empire." swaggertype:"string"`
	CommentCount int64  `json:"comment_count" example:"10" swaggertype:"integer"`
}
