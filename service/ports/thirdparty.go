package ports

import "github.com/peterhellberg/swapi"

//go:generate mockgen -destination=../mocks/thirdparty.go -package=mocks github.com/iamnator/movie-api/adapter/thirdparty/ports ISwapi
type ISwapi interface {
	GetFilms(id ...int) ([]swapi.Film, error)
	GetCharacters(id ...int) ([]swapi.Person, error)
}
