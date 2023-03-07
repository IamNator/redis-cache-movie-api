package service

import (
	"errors"
	"fmt"
	
	"github.com/rs/zerolog/log"

	"github.com/iamnator/movie-api/model"
)

func getFilmIDFromURL(url string) (int, error) {
	var id int
	if _, er := fmt.Sscanf(url, "https://swapi.dev/api/films/%d/", &id); er != nil {
		return 0, errors.New("error getting id from url")
	}
	return id, nil
}
func (s service) backGroundJOB() error {

	//get all movies
	go s.updateMovies()
	return nil
}

func (s service) updateMovies() error {

	films, err := s.swapiClient.GetFilms()
	if err != nil {
		log.Error().Err(err).Msg("error getting movies")
		return errors.New("error getting movies")
	}
	var movies []model.MovieDetails
	var filmID int

	characterIDChan := make(chan int, 5)

	go s.updateCharacters(characterIDChan)

	for _, film := range films {
		filmID, err = getFilmIDFromURL(film.URL)
		if err != nil {
			log.Error().Err(err).Msg("error getting film id")
			return errors.New("error getting film id")
		}

		movies = append(movies, model.MovieDetails{
			ID:           filmID,
			Name:         film.Title,
			EpisodeID:    film.EpisodeID,
			OpeningCrawl: film.OpeningCrawl,
			Director:     film.Director,
			Producer:     film.Producer,
		})

		for _, characterURL := range film.CharacterURLs {

			characterID, err := getCharacterIDFromURL(characterURL)
			if err != nil {
				log.Error().Err(err).Msg("error getting character id")
				return errors.New("error getting character id")
			}

			characterIDChan <- characterID
		}
	}

	//save movies to cache
	if err := s.cache.SetMovies(movies); err != nil {
		log.Error().Err(err).Msg("error saving movies")
		return errors.New("error saving movies")
	}

	return err
}

func getCharacterIDFromURL(url string) (int, error) {
	var id int
	if _, er := fmt.Sscanf(url, "https://swapi.dev/api/people/%d/", &id); er != nil {
		return 0, errors.New("error getting id from url")
	}
	return id, nil
}

func (s service) updateCharacters(chn chan int) {

	for characterID := range chn {
		characters, err := s.swapiClient.GetCharacters(characterID)
		if err != nil {
			log.Error().Err(err).Msg("error getting character")
			continue
		}
		var characterList []model.Character

		for _, character := range characters {
			characterList = append(characterList, model.Character{
				Name: character.Name,
			})
		}

		if err := s.cache.SetCharactersByMovieID(0, characterList); err != nil {
			log.Error().Err(err).Msg("error saving character")
			continue
		}
	}
}
