package ports

import "github.com/iamnator/movie-api/model"

//go:generate mockgen -destination=../mocks/cache.go -package=mocks github.com/iamnator/movie-api/adapter/cache ICache
type ICache interface {
	SetMovies([]model.MovieDetails) error
	SetMovieByID(id int, movie model.MovieDetails) error
	SetCharactersByMovieID(id int, characters []model.Character) error

	GetMovies(page, pageSize int) ([]model.MovieDetails, int64, error)
	GetMovieByID(id int) (*model.MovieDetails, error)
	GetCharactersByMovieID(id int, page, pageSize int) ([]model.Character, int64, error)
}
