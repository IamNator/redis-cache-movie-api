package service

import (
	"github.com/iamnator/movie-api/model"
	"github.com/iamnator/movie-api/service/ports"
	"sort"
)

type IServices interface {
	GetMovies(page, pageSize int) ([]model.Movie, int64, error)
	SaveComment(movieID int, comment model.Comment) error
	GetComment(movieID int, page, pageSize int) ([]model.Comment, int64, error)
	GetCharactersByMovieID(movieID int, page, pageSize int) ([]model.Character, int64, error)
}

type service struct {
	cache             ports.ICache
	commentRepository ports.ICommentRepository
	swapiClient       ports.ISwapi
}

func NewServices(cache ports.ICache, commentRepository ports.ICommentRepository, swapiClient ports.ISwapi) service {
	return service{
		cache:             cache,
		commentRepository: commentRepository,
		swapiClient:       swapiClient,
	}
}

func (s service) GetMovies(page, pageSize int) ([]model.Movie, int64, error) {

	movies, _, err := s.cache.GetMovies(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	sort.Slice(movies, func(i, j int) bool {
		return movies[i].ReleaseDate.Before(movies[j].ReleaseDate)
	})

	var movieList []model.Movie
	for _, movie := range movies {

		_, count, err := s.commentRepository.GetCommentsByMovieID(movie.ID, 1, 1)
		if err != nil {
			return nil, 0, err
		}

		movieList = append(movieList, model.Movie{
			Name:         movie.Title,
			OpeningCrawl: movie.OpeningCrawl,
			CommentCount: count,
		})
	}

	return movieList, 0, nil
}

func (s service) SaveComment(movieID int, comment model.Comment) error {
	return s.commentRepository.AddComment(comment)
}

func (s service) GetComment(movieID int, page, pageSize int) ([]model.Comment, int64, error) {
	return s.commentRepository.GetCommentsByMovieID(movieID, page, pageSize)
}

func (s service) GetCharactersByMovieID(movieID int, page, pageSize int) ([]model.Character, int64, error) {
	return s.cache.GetCharactersByMovieID(movieID, page, pageSize)
}
