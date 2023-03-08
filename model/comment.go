package model

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type AddCommentRequest struct {
	ID        uuid.UUID      `json:"id" swaggerignore:"true"`
	MovieID   int            `json:"movieID" swaggerignore:"true"  `
	Message   string         `json:"message" `
	IPv4Addr  string         `json:"ipv4_addr" swaggerignore:"true" `
	CreatedAt time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt *time.Time     `json:"updated_at" swaggerignore:"true"  gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt ` gorm:"column:deleted_at" swaggerignore:"true" json:"-"`
}

func (a AddCommentRequest) ToComment() Comment {
	return Comment{
		ID:           a.ID,
		SwapiMovieID: a.MovieID,
		Message:      a.Message,
		IPv4Addr:     a.IPv4Addr,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
		DeletedAt:    a.DeletedAt,
	}
}

type Comment struct {
	ID           uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:id"`
	SwapiMovieID int            `json:"swapi_movie_id" gorm:"column:swapi_movie_id;not null;comment:'id of movie from www.swapi.dev'"`
	Message      string         `json:"message" gorm:"column:message;not null;size:600;comment:'max length expected is 500; padding = 100 chars'"`
	IPv4Addr     string         `json:"ipv4_addr" swaggerignore:"true" gorm:"column:ipv4_addr;not null;size:20;index;comment:'Ip address of the person commenting; max expected length is 15'"`
	CreatedAt    time.Time      `json:"created_at" swaggerignore:"true"  gorm:"column:created_at;not null;default:current_timestamp"`
	UpdatedAt    *time.Time     `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt ` gorm:"column:deleted_at" swaggerignore:"true" json:"-"`
}

// TableName specifies the table name for the Comment model.
func (Comment) TableName() string {
	return "comment"
}
