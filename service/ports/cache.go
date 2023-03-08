package ports

import "github.com/iamnator/movie-api/model"

type GetCharacterFiler struct {
	SortKey   string
	SortOrder string
	Gender    string
}

//go:generate mockgen -source=cache.go -destination=../mocks/cache.go  -package=mocks github.com/iamnator/movie-api/service/ports ICache
type ICache interface {
	SetMovies([]model.MovieDetails) error
	SetMovieByID(id int, movie model.MovieDetails) error
	SetCharactersByMovieID(id int, characters []model.Character) error

	GetMovies(page, pageSize int) ([]model.MovieDetails, int64, error)
	GetMovieByID(id int) (*model.MovieDetails, error)
	GetCharactersByMovieID(id int, page, pageSize int, filter GetCharacterFiler) ([]model.Character, int64, error)
}
