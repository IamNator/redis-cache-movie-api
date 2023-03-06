package swapi

type (
	ISwapi interface {
		GetFilms(id ...int) ([]Film, error)
		GetCharacters(id ...int) ([]Character, error)
	}
)
