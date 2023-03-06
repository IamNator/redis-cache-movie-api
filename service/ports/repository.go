package ports

import "github.com/iamnator/movie-api/model"

type (
	ICommentRepository interface {
		AddComment(comment model.Comment) error
		GetComment(commentID int) (*model.Comment, error)
		GetCommentsByID(commentID ...int) ([]model.Comment, error)
		GetCommentsByIPAddr(ipAddr string) ([]model.Comment, error)
		GetCommentsByMovieID(movieID int, page, pageSize int) ([]model.Comment, error)
	}
)
