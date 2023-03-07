package model

import "fmt"

type (
	GetCharactersByMovieIDArgs struct {
		MovieID   int
		Page      int
		PageSize  int
		SortKey   string //name, gender or height
		SortOrder string //asc or desc
		Gender    string //  'female', 'male' or 'n/a' -> default
	}

	Character struct {
		ID       int    `json:"character_id"` //from swapi
		MovieID  int    `json:"movie_id"`     //from swapi
		Name     string `json:"name"`
		Gender   string `json:"gender"`
		HeightCm int    `json:"height_cm"`
	}

	CharacterList_Character struct {
		Name     string  `json:"name"`
		Gender   string  `json:"gender"`
		HeightCm int     `json:"height_cm"`
		HeightFt string  `json:"height_ft"`
		HeightIn float64 `json:"height_in"`
	}

	CharacterList struct {
		Characters []CharacterList_Character `json:"characters"`
		TotalCount int                       `json:"total_count"`
		TotalCm    int                       `json:"total_cm"`
		TotalFt    string                    `json:"total_ft"`
		TotalIn    float64                   `json:"total_in"`
	}
)

func (c Character) FeetsInches() (feets string, inches float64) {
	return FeetsInches(c.HeightCm)
}

func FeetsInches(heightCm int) (feets string, inches float64) {

	// 1 inch = 2.54 cm
	// 1 feet = 12 inches

	if heightCm == 0 {
		return "0", 0
	}

	inches = float64(heightCm) / 2.54
	feets = fmt.Sprintf("%d", int(inches/12))
	inches = inches - (float64(int(inches/12)) * 12)

	return feets, inches
}
