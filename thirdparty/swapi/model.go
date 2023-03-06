package swapi

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type (
	Film struct {
		Title        string   `json:"title"`
		EpisodeID    int      `json:"episode_id"`
		OpeningCrawl string   `json:"opening_crawl"`
		Director     string   `json:"director"`
		Producer     string   `json:"producer"`
		ReleaseDate  string   `json:"release_date"`
		Characters   []string `json:"characters"`
		Planets      []string `json:"planets"`
		Starships    []string `json:"starships"`
		Vehicles     []string `json:"vehicles"`
		Species      []string `json:"species"`
		Created      string   `json:"created"`
		Edited       string   `json:"edited"`
		URL          string   `json:"url"`
	}

	Character struct {
		Name      string   `json:"name"`
		Height    string   `json:"height"`
		Mass      string   `json:"mass"`
		HairColor string   `json:"hair_color"`
		SkinColor string   `json:"skin_color"`
		EyeColor  string   `json:"eye_color"`
		BirthYear string   `json:"birth_year"`
		Gender    string   `json:"gender"`
		HomeWorld string   `json:"homeworld"`
		Films     []string `json:"films"`
		Species   []string `json:"species"`
		Vehicles  []string `json:"vehicles"`
		Starships []string `json:"starships"`
		Created   string   `json:"created"`
		Edited    string   `json:"edited"`
		URL       string   `json:"url"`
	}
)

func (f Film) Validate() error {
	validateDate := func(value interface{}) error {
		//check if time format is ISO 8601
		if _, err := time.Parse(time.RFC3339, value.(string)); err != nil {
			return errors.New("date is in ISO 8601 format")
		}
		return nil
	}

	return validation.ValidateStruct(&f,
		validation.Field(&f.Title, validation.Required),
		validation.Field(&f.EpisodeID, validation.Required),
		validation.Field(&f.OpeningCrawl, validation.Required),
		validation.Field(&f.Director, validation.Required),
		validation.Field(&f.Producer, validation.Required),
		validation.Field(&f.Characters, validation.Each(is.URL)),
		validation.Field(&f.Planets, validation.Each(is.URL)),
		validation.Field(&f.Starships, validation.Each(is.URL)),
		validation.Field(&f.Vehicles, validation.Each(is.URL)),
		validation.Field(&f.Species, validation.Each(is.URL)),
		validation.Field(&f.URL, validation.Required),
		validation.Field(&f.Created, validation.Required, validation.By(validateDate)),
		validation.Field(&f.Edited, validation.Required, validation.By(validateDate)),
		validation.Field(&f.ReleaseDate, validation.Required, validation.By(validateDate)),
	)
}
