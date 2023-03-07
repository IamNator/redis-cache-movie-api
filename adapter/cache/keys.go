package cache

import "strconv"

type (
	cacheTag string
)

func (tag cacheTag) String() string {
	return string(tag)
}

func (tag cacheTag) computeKey(id ...int) string {
	var key string

	for _, i := range id {
		key += ":" + strconv.Itoa(i)
	}

	return tag.String() + key
}

// computeParentKey returns the parent key for a given key
// e.g. let parentTag = "movie"
// then parentTag.computeParentKey(1) -> "movie:1"
func (parentTag cacheTag) computeParentKey(id int) cacheTag {
	return cacheTag(parentTag.computeKey(id))
}

// computeMovieKey returns the key for a movie
// e.g. movie:1 -> movie:<movie_id>
func computeMovieKey(id int) string {
	return movieTag.computeKey(id)
}

// computeCharacterKey returns the key for a character
// e.g. movie:1:character:2 -> movie:<movie_id>:character:<character_id>
func computeCharacterKey(movieID, characterID int) string {
	return movieTag.computeParentKey(movieID).computeKey(characterID)
}

const (
	movieTag     cacheTag = "movie"
	characterTag cacheTag = "character"

	// DefaultTTLSec is the default TTL for cache entries
	DefaultTTLSec = 0
)
