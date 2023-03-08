package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/iamnator/movie-api/model"
	"github.com/iamnator/movie-api/service/ports"
	"github.com/rs/zerolog/log"
	"sort"
	"time"
)

type IServices interface {
	GetMovies(page, pageSize int) ([]model.Movie, int64, error)
	SaveComment(movieID int, comment model.Comment) error
	GetComment(movieID int, page, pageSize int) ([]model.Comment, int64, error)
	GetCharactersByMovieID(arg model.GetCharactersByMovieIDArgs) (*model.CharacterList, int64, error)
}

type service struct {
	cache             ports.ICache
	commentRepository ports.ICommentRepository
	swapiClient       ports.ISwapi
}

func NewServices(cache ports.ICache, commentRepository ports.ICommentRepository, swapiClient ports.ISwapi) IServices {
	srv := service{
		cache:             cache,
		commentRepository: commentRepository,
		swapiClient:       swapiClient,
	}

	go func() {

		for {
			log.Info().Msg("running background job ...")
			if err := srv.backGroundJOB(); err == nil {
				log.Info().Msg("background job ran successfully")
				break
			} else {
				log.Error().Err(err).Msg("error running background job, retrying ...")
			}
		}

		//runs every 3 hours
		ticker := time.NewTicker(3 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if err := srv.backGroundJOB(); err != nil {
				log.Error().Err(err).Msg("error running background job")
			} else {
				log.Info().Msg("background job ran successfully")
			}
		}
	}()

	return srv
}

func (s service) GetMovies(page, pageSize int) ([]model.Movie, int64, error) {

	movies, count, err := s.cache.GetMovies(page, pageSize)
	if err != nil {
		log.Debug().Err(err).Msg("error getting movies from cache")
		return nil, 0, errors.New("error getting movies from cache")
	}

	sort.Slice(movies, func(i, j int) bool {
		return movies[i].ReleaseDate.Before(movies[j].ReleaseDate)
	})

	var movieList []model.Movie
	for _, movie := range movies {

		commentCount, err := s.commentRepository.GetCommentCountByMovieID(movie.ID)
		if err != nil {
			log.Error().Err(err).Msg("error getting comment count")
			return nil, 0, errors.New("error getting comment count")
		}

		movieList = append(movieList, model.Movie{
			SwapiMovieID: movie.ID,
			Name:         movie.Name,
			OpeningCrawl: movie.OpeningCrawl,
			CommentCount: commentCount,
			ReleaseDate:  movie.ReleaseDate,
		})
	}

	return movieList, count, nil
}

func (s service) GetCharactersByMovieID(arg model.GetCharactersByMovieIDArgs) (*model.CharacterList, int64, error) {
	//check if movie exists
	movie, err := s.cache.GetMovieByID(arg.MovieID)
	if err != nil {
		log.Debug().Err(err).Msg("movie not found")
		return nil, 0, errors.New("movie not found")
	}

	characters, count, err := s.cache.GetCharactersByMovieID(movie.ID, arg.Page, arg.PageSize, ports.GetCharacterFiler{
		SortKey:   arg.SortKey,
		SortOrder: arg.SortOrder,
		Gender:    arg.Gender,
	})

	if err != nil {
		log.Error().Err(err).Msg("error getting characters")
		return nil, 0, errors.New("error getting characters")
	}

	if len(characters) == 0 {
		return nil, 0, errors.New("no characters found")
	}

	var characterList model.CharacterList
	var inches float64
	var feets string
	for _, character := range characters {
		feets, inches = model.FeetsInches(character.HeightCm)
		characterList.Characters = append(characterList.Characters, model.CharacterList_Character{
			Name:     character.Name,
			Gender:   character.Gender,
			HeightCm: character.HeightCm,
			HeightFt: feets,
			HeightIn: inches,
		})

		characterList.TotalCm += character.HeightCm
	}

	characterList.TotalFt, characterList.TotalIn = model.FeetsInches(characterList.TotalCm)
	characterList.TotalCount = int(count)

	return &characterList, count, nil
}

func (s service) SaveComment(movieID int, comment model.Comment) error {

	//check if movie exists
	_, err := s.cache.GetMovieByID(movieID)
	if err != nil {
		log.Debug().Err(err).Msg("movie not found")
		return errors.New("movie not found")
	}

	comment = model.Comment{
		ID:           uuid.New(),
		SwapiMovieID: comment.SwapiMovieID,
		Message:      comment.Message,
		IPv4Addr:     comment.IPv4Addr,
		CreatedAt:    comment.CreatedAt,
	}

	if err := s.commentRepository.AddComment(comment); err != nil {
		log.Error().Err(err).Msg("error saving comment")
		return errors.New("error saving comment")
	}

	return err
}

func (s service) GetComment(movieID int, page, pageSize int) ([]model.Comment, int64, error) {
	//check if movie exists
	movie, err := s.cache.GetMovieByID(movieID)
	if err != nil {
		log.Debug().Err(err).Msg("movie not found")
		return nil, 0, errors.New("movie not found")
	}

	comments, count, err := s.commentRepository.GetCommentsByMovieID(movie.ID, page, pageSize)
	if err != nil {
		log.Error().Err(err).Msg("error getting comments")
		return nil, 0, err
	}

	return comments, count, nil
}
