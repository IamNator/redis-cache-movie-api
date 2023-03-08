package swapi

import (
	"context"
	"errors"
	swapi "github.com/iamnator/movie-api/thirdparty/swapi/lib"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
	"time"
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

func (s *Swapi) GetCharacters(ctx context.Context, ids ...int) ([]swapi.Person, error) {

	var characters []swapi.Person
	var errs []error
	if len(ids) > 0 {
		var character swapi.Person
		var err error

		charChan := make(chan swapi.Person, len(ids))
		erChan := make(chan error, len(ids))
		wait := sync.WaitGroup{}

		go func() {
			for p := range charChan {
				log.Info().Msgf("fetched character: %s", p.Name)
				characters = append(characters, p)
			}
		}()

		go func() {
			for e := range erChan {
				log.Error().Msgf("error fetching character: %s", e.Error())
				errs = append(errs, e)
			}
		}()

		idChan := make(chan int)

		workers := func() {

			defer wait.Done()
			for {
				select {
				case id, ok := <-idChan:
					if !ok {
						return
					}
					character, err = s.client.Person(ctx, id)
					if err != nil {
						erChan <- err
					} else {
						charChan <- character
					}

				}
			}

		}

		for i := 0; i < 5; i++ {
			wait.Add(1)
			go workers()
		}

		// push to workers
		for _, id := range ids {
			idChan <- id
		}

		close(idChan)

		wait.Wait()
		time.Sleep(time.Millisecond * 50) //allow time for workers to finish
		close(charChan)
		close(erChan)

	} else {
		_characters, err := s.client.AllPeople(ctx)
		if err != nil {
			return nil, err
		}
		characters = append(characters, _characters...)
	}

	if len(errs) > 0 {
		var ss string
		for _, err := range errs {
			ss += err.Error() + ", "
		}

		return characters, errors.New(ss)
	}

	return characters, nil
}
