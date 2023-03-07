package lib

import (
	"fmt"
	"strconv"
	"strings"
)

func getIDFromURL(url string) (int, error) {
	// https://swapi.dev/api/films/6/

	//stripe the last slash
	url = strings.TrimSuffix(url, "/")

	// 1. Split the url by "/"
	ss := strings.Split(url, "/")
	// 2. Get the last element
	last := ss[len(ss)-1]
	// 3. Convert the last element to an int
	id, err := strconv.Atoi(last)
	if err != nil {
		return 0, fmt.Errorf("error converting %s to int: %w", last, err)
	}
	// 4. Return the int
	return id, nil
}
