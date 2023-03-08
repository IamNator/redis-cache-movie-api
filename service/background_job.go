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

	movCharMap := make(map[int][]int) // movieID, []characterID
	for msg := range chn {
		if _, ok := movCharMap[msg.movieID]; !ok {
			movCharMap[msg.movieID] = []int{msg.charID}
		} else {
			movCharMap[msg.movieID] = append(movCharMap[msg.movieID], msg.charID)
		}
	}

	var chxIDs []int
	for _, v := range movCharMap {
		chxIDs = append(chxIDs, v...)
	}

	characters, err := s.swapiClient.GetCharacters(context.Background(), chxIDs...)
	if err != nil {
		log.Error().Err(err).Msg("error getting character")
		return
	}

	movieCharacterMap := make(map[int][]model.Character) // movieID, []character

	for _, character := range characters {

		var heightCm int
		for movieID, charIDs := range movCharMap {

			for _, charID := range charIDs {
				id, _ := getCharacterIDFromURL(character.URL)
				if charID == id {
					if h, er := strconv.Atoi(character.Height); er == nil {
						heightCm = h
					}

					if _, ok := movieCharacterMap[movieID]; !ok {

						movieCharacterMap[movieID] = []model.Character{model.Character{
							ID:       id,
							MovieID:  movieID,
							Name:     character.Name,
							Gender:   character.Gender,
							HeightCm: heightCm,
						}}
					} else {
						movieCharacterMap[movieID] = append(movieCharacterMap[movieID], model.Character{
							ID:       id,
							MovieID:  movieID,
							Name:     character.Name,
							Gender:   character.Gender,
							HeightCm: heightCm,
						})
					}
				}
			}
		}
	}

	for msg, characterList := range movieCharacterMap {
		if err := s.cache.SetCharactersByMovieID(msg, characterList); err != nil {
			log.Error().Err(err).Msg("error saving character")
			return
		}
	}

	log.Info().Msgf("length of characters cached: %d", len(characters))
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
