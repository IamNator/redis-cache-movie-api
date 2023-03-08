package ports

import (
	"context"
	swapi "github.com/iamnator/movie-api/thirdparty/swapi/lib"
)

//go:generate mockgen -source=thirdparty.go -destination=./mocks/thirdparty.go  -package=mocks github.com/iamnator/movie-api/service/ports ISwapi
type ISwapi interface {
	GetFilms(ctx context.Context, id ...int) ([]swapi.Film, error)
	GetCharacters(ctx context.Context, id ...int) ([]swapi.Person, error)
}
