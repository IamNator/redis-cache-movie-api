package swapi

import (
	"context"
	swapi "github.com/iamnator/movie-api/thirdparty/swapi/lib"
	"net/http"
)

type (
	ISwapi interface {
		GetFilms(ctx context.Context, id ...int) ([]swapi.Film, error)
		GetCharacters(ctx context.Context, id ...int) ([]swapi.Person, error)
	}

	Swapi struct {
		client *swapi.Client
	}
)

func NewSwapi(hc *http.Client) (ISwapi, error) {

	opts := swapi.HTTPClient(hc)

	c := swapi.NewClient(opts)

	return &Swapi{
		client: c,
	}, nil
}

func (s *Swapi) GetFilms(ctx context.Context, id ...int) ([]swapi.Film, error) {
	var films []swapi.Film

	if len(id) > 0 {
		for _, i := range id {
			film, err := s.client.Film(ctx, i)
			if err != nil {
				return nil, err
			}
			films = append(films, film)
		}
	} else {
		_films, err := s.client.AllFilms(ctx)
		if err != nil {
			return nil, err
		}
		
		films = append(films, _films...)
	}

	return films, nil
}

func (s *Swapi) GetCharacters(ctx context.Context, id ...int) ([]swapi.Person, error) {

	var characters []swapi.Person
	if len(id) > 0 {
		for _, i := range id {
			character, err := s.client.Person(ctx, i)
			if err != nil {
				return nil, err
			}
			characters = append(characters, character)
		}
	} else {
		_characters, err := s.client.AllPeople(ctx)
		if err != nil {
			return nil, err
		}
		characters = append(characters, _characters...)
	}

	return characters, nil
}
