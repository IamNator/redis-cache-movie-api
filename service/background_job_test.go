package service_test

import (
	"github.com/iamnator/movie-api/service"
	"testing"
)

func Test_GetFilmIDFromURL(t *testing.T) {
	tests := []struct {
		URL     string
		ID      int
		WantErr bool
	}{
		{
			URL:     "https://swapi.dev/api/films/1/",
			ID:      1,
			WantErr: false,
		},
		{
			URL:     "https://swapi.dev/api/films/20933/",
			ID:      20933,
			WantErr: false,
		},
		{
			URL:     "https://swapi.dev/api/films/3erd/",
			ID:      3,
			WantErr: true,
		},
		{
			URL:     "https://swapi.dev/api/films/443/4344",
			ID:      443,
			WantErr: false,
		},
		{
			URL:     "https://swapi.dev/api/films/443?rye=4344",
			ID:      443,
			WantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.URL, func(t *testing.T) {
			id, err := service.GetFilmIDFromURL(tt.URL)
			if err != nil && !tt.WantErr {
				t.Errorf("url=%s | error=%s | id_gotten=%v | id_expected=%v", tt.URL, err.Error(), id, tt.ID)
				return
			}
			if id != tt.ID && !tt.WantErr {
				t.Errorf("url=%s | id_gotten=%v | id_expected=%v", tt.URL, id, tt.ID)
				return
			}
		})
	}
}

func Test_GetCharacterIDFromURL(t *testing.T) {

	tests := []struct {
		URL     string
		ID      int
		WantErr bool
	}{
		{
			URL:     "https://swapi.dev/api/people/1/",
			ID:      1,
			WantErr: false,
		},
		{
			URL:     "https://swapi.dev/api/people/20933/",
			ID:      20933,
			WantErr: false,
		},
		{
			URL:     "https://swapi.dev/api/people/3erd/",
			ID:      3,
			WantErr: true,
		},
		{
			URL:     "https://swapi.dev/api/people/443/4344",
			ID:      443,
			WantErr: false,
		},
		{
			URL:     "https://swapi.dev/api/people/443?rye=4344",
			ID:      443,
			WantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.URL, func(t *testing.T) {
			id, err := service.GetCharacterIDFromURL(tt.URL)
			if err != nil && !tt.WantErr {
				t.Errorf("url=%s | error=%s | id_gotten=%v | id_expected=%v", tt.URL, err.Error(), id, tt.ID)
				return
			}
			if id != tt.ID && !tt.WantErr {
				t.Errorf("url=%s | id_gotten=%v | id_expected=%v", tt.URL, id, tt.ID)
				return
			}
		})
	}
}
