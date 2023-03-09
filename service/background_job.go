package service

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"net/url"
	"strconv"
	"strings"
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

		filmID, err = GetFilmIDFromURL(film.URL)
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

			characterID, err := GetCharacterIDFromURL(characterURL)
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

func chunkSlice(slice []int, chunkSize int) [][]int {
	var chunks [][]int
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func (s service) refreshCharacterCache(chn chan struct {
	charID  int
	movieID int
}) {

	movCharMap := make(map[int][]int) // movieID, []characterID
	var chxIDs []int

	for msg := range chn {
		if _, ok := movCharMap[msg.movieID]; !ok {
			movCharMap[msg.movieID] = []int{msg.charID}
		} else {
			movCharMap[msg.movieID] = append(movCharMap[msg.movieID], msg.charID)
		}
	}

	var chxMap = make(map[int]bool)
	for _, v := range movCharMap {
		for _, chxID := range v {
			if _, ok := chxMap[chxID]; !ok {
				chxMap[chxID] = true
				chxIDs = append(chxIDs, chxID)
			}
		}
	}

	log.Info().Msgf("length of movie characters to fetch: %d", len(chxIDs))

	// chunk the slice of character ids
	//
	// the reason for this is that it takes a while to fetch all the characters
	var steps [][]int
	if len(chxIDs) > 10 {
		steps = chunkSlice(chxIDs, 10)
	} else {
		steps = [][]int{chxIDs}
	}

	log.Info().Msgf("length of steps: %d", len(steps))

	for _, stepIds := range steps {

		characters, err := s.swapiClient.GetCharacters(context.Background(), stepIds...)
		if err != nil {
			log.Error().Err(err).Msg("error getting character")
		}

		if len(characters) == 0 {
			log.Info().Msg("no characters to cache")
			return
		}

		log.Info().Msgf("length of fetched characters: %d", len(characters))

		movieCharacterMap := make(map[int][]model.Character) // movieID, []character

		for _, character := range characters {

			var heightCm int
			for movieID, charIDs := range movCharMap {

				for _, charID := range charIDs {
					id, _ := GetCharacterIDFromURL(character.URL)
					if charID == id {
						if h, er := strconv.Atoi(character.Height); er == nil {
							heightCm = h
						}

						if _, ok := movieCharacterMap[movieID]; !ok { //

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

}

func GetFilmIDFromURL(urL string) (int, error) {

	// url format: https://swapi.dev/api/films/1/
	var id int

	parseURL, err := url.Parse(urL)
	if err != nil {
		return 0, err
	}

	path := parseURL.Path

	parts := strings.TrimSuffix(path, "/")
	split := strings.Split(parts, "/")
	if len(split) < 2 {
		return 0, errors.New("error getting id from url")
	}

	for i, v := range split {
		if v == "films" {
			if i+1 >= len(split) {
				return 0, errors.New("error getting id from url")
			}
			id, err = strconv.Atoi(split[i+1])
			if err != nil {
				return 0, err
			}
			break
		}
	}

	return id, nil
}

func GetCharacterIDFromURL(urL string) (int, error) {

	// url format: https://swapi.dev/api/people/%d/
	var id int

	parseURL, err := url.Parse(urL)
	if err != nil {
		return 0, err
	}

	path := parseURL.Path

	parts := strings.TrimSuffix(path, "/")
	split := strings.Split(parts, "/")
	if len(split) < 2 {
		return 0, errors.New("error getting id from url")
	}

	for i, v := range split {
		if v == "people" {
			if i+1 >= len(split) {
				return 0, errors.New("error getting id from url")
			}
			id, err = strconv.Atoi(split[i+1])
			if err != nil {
				return 0, err
			}
			break
		}
	}

	return id, nil
}
