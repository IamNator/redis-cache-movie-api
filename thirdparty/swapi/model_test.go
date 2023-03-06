package swapi

import "testing"

func TestFilm_Validate(t *testing.T) {
	tests := []struct {
		Film    Film
		WantErr bool
	}{
		{
			Film:    Film{},
			WantErr: true,
		},
	}

	for _, test := range tests {
		t.Run("test", func(t *testing.T) {
			if err := test.Film.Validate(); (err != nil) != test.WantErr {
				t.Errorf("expected error, got nil")
			}
		})
	}

}
