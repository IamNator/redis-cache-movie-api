package lib

import (
	"context"
	"fmt"
	"time"
)

// A Film is an single film.
type Film struct {
	Title         string   `json:"title"`
	EpisodeID     int      `json:"episode_id"`
	OpeningCrawl  string   `json:"opening_crawl"`
	Director      string   `json:"director"`
	Producer      string   `json:"producer"`
	CharacterURLs []string `json:"characters"`
	PlanetURLs    []string `json:"planets"`
	StarshipURLs  []string `json:"starships"`
	VehicleURLs   []string `json:"vehicles"`
	SpeciesURLs   []string `json:"species"`
	Created       string   `json:"created" example:"2014-12-20T18:49:38.403000Z"`
	Edited        string   `json:"edited"`
	URL           string   `json:"url" example:"https://swapi.dev/api/films/6/"`
	ReleaseDate   string   `json:"release_date" example:"2002-05-16"`
}

func (f Film) GetID() int {
	id, _ := getIDFromURL(f.URL)
	return id
}

func (f Film) GetReleaseDate() time.Time {
	t, _ := time.Parse("2006-01-02", f.ReleaseDate)
	return t
}

func (f Film) GetCreated() time.Time {
	t, _ := time.Parse(time.RFC3339, f.Created)
	return t
}

func (f Film) GetEdited() time.Time {
	t, _ := time.Parse(time.RFC3339, f.Edited)
	return t
}

// Film retrieves the film with the given id
func (c *Client) Film(ctx context.Context, id int) (Film, error) {
	req, err := c.newRequest(ctx, fmt.Sprintf("films/%d", id))
	if err != nil {
		return Film{}, err
	}

	var film Film

	if _, err = c.do(req, &film); err != nil {
		return Film{}, err
	}

	return film, nil
}

func (c *Client) AllFilms(ctx context.Context) ([]Film, error) {
	var films []Film

	req, err := c.newRequest(ctx, "films")
	if err != nil {
		return nil, err
	}

	for {
		var list List[Film]

		if _, err = c.do(req, &list); err != nil {
			return nil, err
		}

		films = append(films, list.Results...)

		if list.Next == nil {
			break
		}

		req, err = c.getRequest(ctx, *list.Next)
		if err != nil {
			return nil, err
		}
	}

	return films, nil
}
