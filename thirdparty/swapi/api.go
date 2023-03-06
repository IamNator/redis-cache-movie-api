package swapi

import (
	"context"
	"github.com/peterhellberg/swapi"
)

type (
	ISwapi interface {
		GetFilms(id ...int) ([]swapi.Film, error)
		GetCharacters(id ...int) ([]swapi.Person, error)
	}

	Swapi struct {
		client *swapi.Client
	}
)

func NewSwapi() (ISwapi, error) {
	c := swapi.NewClient(swapi.UserAgent("busha_swapi/1.0"))

	return &Swapi{
		client: c,
	}, nil
}

func (s *Swapi) GetFilms(id ...int) ([]swapi.Film, error) {
	var films []swapi.Film
	if len(id) > 0 {
		for _, i := range id {
			film, err := s.client.Film(context.Background(), i)
			if err != nil {
				return nil, err
			}
			films = append(films, film)
		}
	} else {
		_films, err := s.client.AllFilms(context.Background())
		if err != nil {
			return nil, err
		}
		films = append(films, _films...)
	}

	return nil, nil
}

func (s *Swapi) GetCharacters(id ...int) ([]swapi.Person, error) {

	var characters []swapi.Person
	if len(id) > 0 {
		for _, i := range id {
			character, err := s.client.Person(context.Background(), i)
			if err != nil {
				return nil, err
			}
			characters = append(characters, character)
		}
	} else {
		_characters, err := s.client.AllPeople(context.Background())
		if err != nil {
			return nil, err
		}
		characters = append(characters, _characters...)
	}

	return characters, nil
}
