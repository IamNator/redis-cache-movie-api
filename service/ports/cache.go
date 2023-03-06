package ports

import "github.com/iamnator/movie-api/model"

type ICache interface {
	SetMovies([]model.MovieDetails) error
	SetMovieByID(id int, movie model.MovieDetails) error
	SetCharactersByMovieID(id int, characters []model.Character) error

	GetMovies() ([]model.MovieDetails, error)
	GetMovieByID(id int) (*model.MovieDetails, error)
	GetCharactersByMovieID(id int) ([]model.Character, error)
}
