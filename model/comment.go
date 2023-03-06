package model

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:id"`
	SwapiMovieID int            `gorm:"column:swapi_movie_id;not null;comment:'id of movie from www.swapi.dev'"`
	Message      string         `gorm:"column:message;not null;size:600;comment:'max length expected is 500; padding = 100 chars'"`
	IPv4Addr     string         `gorm:"column:ipv4_addr;not null;size:20;index;comment:'Ip address of the person commenting; max expected length is 15'"`
	CreatedAt    time.Time      `gorm:"column:created_at;not null;default:current_timestamp"`
	UpdatedAt    *time.Time     `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

// TableName specifies the table name for the Comment model.
func (Comment) TableName() string {
	return "comment"
}
