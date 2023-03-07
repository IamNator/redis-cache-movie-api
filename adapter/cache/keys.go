package cache

import (
	"reflect"
	"strconv"
)

type (
	cacheTag string
)

func (tag cacheTag) String() string {
	return string(tag)
}

func (tag cacheTag) computeIntKey(id ...int) string {
	var key string

	for _, i := range id {
		key += ":" + strconv.Itoa(i)
	}

	return tag.String() + key
}

func (tag cacheTag) computeStringKey(id ...string) string {
	var key string

	for _, k := range id {
		key += ":" + k
	}

	return tag.String() + key
}

// computeParentKey returns the parent key for a given key
// e.g. let parentTag = "movie"
// then parentTag.computeParentKey(1) -> "movie:1"
func (parentTag cacheTag) computeParentKey(id int) cacheTag {
	return cacheTag(parentTag.computeIntKey(id))
}

// computeMovieKey returns the key for a movie
// e.g. movie:1 -> movie:<movie_id>
func computeMovieKey[k int | string](id k) string {

	typeOF := reflect.TypeOf(id)
	if typeOF.Kind() == reflect.Int {
		return movieTag.computeIntKey(int(reflect.ValueOf(id).Int()))
	}

	return movieTag.computeStringKey(reflect.ValueOf(id).String())
}

// computeCharacterKey returns the key for a character
// e.g. movie:1:character:2 -> movie:<movie_id>:character:<character_id>
func computeCharacterKey[k int | string](movieID, characterID k) string {

	typeOfMovieID := reflect.TypeOf(movieID)

	valueOfMovieID := reflect.ValueOf(movieID)
	valueOfCharacterID := reflect.ValueOf(characterID)

	if typeOfMovieID.Kind() == reflect.Int {
		return characterTag.computeIntKey(int(valueOfMovieID.Int()), int(valueOfCharacterID.Int()))
	} else {
		return characterTag.computeStringKey(valueOfMovieID.String(), valueOfCharacterID.String())
	}
}

const (
	movieTag     cacheTag = "movie"
	characterTag cacheTag = "character"

	// DefaultTTLSec is the default TTL for cache entries
	DefaultTTLSec = 0
)