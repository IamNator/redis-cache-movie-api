package service

import (
	"github.com/iamnator/movie-api/model"
	"github.com/iamnator/movie-api/service/ports"
	"sort"
)

type IServices interface {
	GetMovies() ([]model.Movie, error)
	SaveComment(movieID int, comment model.Comment) error
	GetComment(movieID int) ([]model.Comment, error)
	GetCharactersByMovieID(movieID int) ([]model.Character, error)
}

type service struct {
	cache             ports.ICache
	commentRepository ports.ICommentRepository
}

func NewServices(cache ports.ICache, commentRepository ports.ICommentRepository) service {
	return service{
		cache:             cache,
		commentRepository: commentRepository,
	}
}

func (h service) GetMovies() ([]model.Movie, error) {

	movies, err := h.cache.GetMovies()
	if err != nil {
		return nil, err
	}

	sort.Slice(movies, func(i, j int) bool {
		return movies[i].ReleaseDate.Before(movies[j].ReleaseDate)
	})

	var movieList []model.Movie
	for _, movie := range movies {

		comments, err := h.commentRepository.GetCommentsByMovieID(movie.ID, 1, 1)
		if err != nil {
			return nil, err
		}

		movieList = append(movieList, model.Movie{
			Title:        movie.Title,
			OpeningCrawl: movie.OpeningCrawl,
			CommentCount: len(comments),
		})
	}

	return movieList, nil
}

func (h service) SaveComment(movieID int, comment model.Comment) error {
	return h.commentRepository.AddComment(comment)
}

func (h service) GetComment(movieID int) ([]model.Comment, error) {
	return h.commentRepository.GetCommentsByMovieID(movieID, 1, 1)
}

func (h service) GetCharactersByMovieID(movieID int) ([]model.Character, error) {
	return h.cache.GetCharactersByMovieID(movieID)
}
