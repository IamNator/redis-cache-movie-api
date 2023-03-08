package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"

	"github.com/iamnator/movie-api/model"
)

func (s service) backGroundJOB() error {
	//get all movies and characters
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return s.refreshMovieCache(ctx)
}

func (s service) refreshMovieCache(ctx context.Context) error {

	films, err := s.swapiClient.GetFilms(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error getting movies")
		return errors.New("error getting movies")
	}

	log.Info().Msgf("length of films fetched: %d", len(films))

	characterIDChan := make(chan struct {
		charID  int
		movieID int
	}, 5)
	defer close(characterIDChan)
	go s.refreshCharacterCache(characterIDChan)

	var movies []model.MovieDetails
	var filmID int
	var movie model.MovieDetails

	for _, film := range films {

		filmID, err = getFilmIDFromURL(film.URL)
		if err != nil {
			log.Error().Err(err).Msg("error getting film id")
			return errors.New("error getting film id")
		}

		movie = model.MovieDetails{
			ID:           filmID,
			Name:         film.Title,
			EpisodeID:    film.EpisodeID,
			OpeningCrawl: film.OpeningCrawl,
			Director:     film.Director,
			Producer:     film.Producer,
			ReleaseDate:  film.GetReleaseDate(),
			CreatedAt:    film.GetCreated(),
			UpdatedAt:    film.GetEdited(),
		}

		movies = append(movies, movie)

		for _, characterURL := range film.CharacterURLs {

			characterID, err := getCharacterIDFromURL(characterURL)
			if err != nil {
				log.Error().Err(err).Msg("error getting character id")
				return errors.New("error getting character id")
			}

			characterIDChan <- struct {
				charID  int
				movieID int
			}{
				charID:  characterID,
				movieID: filmID,
			}
		}
	}

	//save movies to cache
	if err := s.cache.SetMovies(movies); err != nil {
		log.Error().Err(err).Msg("error saving movies")
		return errors.New("error saving movies")
	}

	log.Info().Msgf("length of movies cached: %d", len(movies))

	return nil
}

func (s service) refreshCharacterCache(chn chan struct {
	charID  int
	movieID int
}) {

	for msg := range chn {
		characters, err := s.swapiClient.GetCharacters(context.Background(), msg.charID)
		if err != nil {
			log.Error().Err(err).Msg("error getting character")
			continue
		}
		var characterList []model.Character
		var heightCm int
		var characterID int

		for _, character := range characters {

			if h, er := strconv.Atoi(character.Height); er == nil {
				heightCm = h
			}

			characterID, err = getCharacterIDFromURL(character.URL)
			if err != nil {
				log.Error().Err(err).Msg("error getting character id")
				continue
			}

			characterList = append(characterList, model.Character{
				ID:       characterID,
				MovieID:  msg.movieID,
				Name:     character.Name,
				Gender:   character.Gender,
				HeightCm: heightCm,
			})
		}

		if err := s.cache.SetCharactersByMovieID(msg.movieID, characterList); err != nil {
			log.Error().Err(err).Msg("error saving character")
			continue
		}
	}
}

func getFilmIDFromURL(url string) (int, error) {
	var id int
	if _, er := fmt.Sscanf(url, "https://swapi.dev/api/films/%d/", &id); er != nil {
		return 0, errors.New("error getting id from url")
	}
	return id, nil
}

func getCharacterIDFromURL(url string) (int, error) {
	var id int
	if _, er := fmt.Sscanf(url, "https://swapi.dev/api/people/%d/", &id); er != nil {
		return 0, errors.New("error getting id from url")
	}
	return id, nil
}
