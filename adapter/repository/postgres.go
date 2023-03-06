package repository

import "github.com/iamnator/movie-api/model"

//go:generate mockgen -destination=../mocks/repository.go -package=mocks github.com/iamnator/movie-api/adapter/repository ICommentRepository
type ICommentRepository interface {
	AddComment(comment model.Comment) error
	GetComment(commentID int) (*model.Comment, error)
	GetCommentsByID(commentID ...int) ([]model.Comment, error)
	GetCommentsByIPAddr(ipAddr string) ([]model.Comment, error)
	GetCommentsByMovieID(movieID int, page, pageSize int) ([]model.Comment, error)
}

type (
	PgxCommentRepoImpl struct {
	}
)

func NewPostgresDB() PgxCommentRepoImpl {
	return PgxCommentRepoImpl{}
}

var _ ICommentRepository = PgxCommentRepoImpl{}

func (p PgxCommentRepoImpl) AddComment(comment model.Comment) error {
	return nil
}

func (p PgxCommentRepoImpl) GetComment(commentID int) (*model.Comment, error) {
	return nil, nil
}
func (p PgxCommentRepoImpl) GetCommentsByID(commentID ...int) ([]model.Comment, error) {
	return []model.Comment{}, nil
}
func (p PgxCommentRepoImpl) GetCommentsByIPAddr(ipAddr string) ([]model.Comment, error) {
	return []model.Comment{}, nil
}
func (p PgxCommentRepoImpl) GetCommentsByMovieID(movieID int, page, pageSize int) ([]model.Comment, error) {
	return []model.Comment{}, nil
}
