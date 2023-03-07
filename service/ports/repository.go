package ports

import (
	"github.com/google/uuid"
	"github.com/iamnator/movie-api/model"
)

//go:generate mockgen -destination=../mocks/repository.go -package=mocks github.com/iamnator/movie-api/adapter/repository ICommentRepository
type ICommentRepository interface {
	AddComment(comment model.Comment) error
	GetComment(commentID uuid.UUID) (*model.Comment, error)
	GetCommentsByID(commentID ...uuid.UUID) ([]model.Comment, error)
	GetCommentsByIPAddr(ipAddr string, page, pageSize int) ([]model.Comment, int64, error)
	GetCommentsByMovieID(movieID int, page, pageSize int) ([]model.Comment, int64, error)
	GetCommentCountByMovieID(movieID int) (int64, error)
}
