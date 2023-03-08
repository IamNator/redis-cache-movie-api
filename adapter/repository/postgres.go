package repository

import (
	"github.com/google/uuid"
	"github.com/iamnator/movie-api/service/ports"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/iamnator/movie-api/model"
)

type PgxCommentRepository struct {
	db *gorm.DB
}

var _ ports.ICommentRepository = (*PgxCommentRepository)(nil)

func NewPgxCommentRepository(url string) (*PgxCommentRepository, error) {

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//ping the database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &PgxCommentRepository{
		db: db,
	}, nil
}

func (p PgxCommentRepository) AddComment(comment model.Comment) error {
	return p.db.Model(&model.Comment{}).Create(&comment).Error
}

func (p PgxCommentRepository) GetComment(commentID uuid.UUID) (comment *model.Comment, err error) {
	return comment, p.db.Model(&model.Comment{}).Where("id = ?", commentID.String()).First(&comment).Error
}

func (p PgxCommentRepository) GetCommentsByID(commentID ...uuid.UUID) (comments []model.Comment, err error) {
	return comments, p.db.Model(&model.Comment{}).Where("id IN ?", commentID).Order("created_at DESC").Find(&comments).Error
}

func (p PgxCommentRepository) GetCommentsByIPAddr(ipAddr string, page, pageSize int) (comments []model.Comment, count int64, err error) {
	if page <= 0 {
		page = 1
	}

	return comments, count, p.db.Model(&model.Comment{}).
		Where("ip_addr = ?", ipAddr).
		Count(&count).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Find(&comments).Error
}

func (p PgxCommentRepository) GetCommentsByMovieID(movieID int, page, pageSize int) (comments []model.Comment, count int64, err error) {
	if page <= 0 {
		page = 1
	}
	return comments, count, p.db.Model(&model.Comment{}).
		Where("swapi_movie_id = ?", movieID).
		Count(&count).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&comments).Error
}

func (p PgxCommentRepository) GetCommentCountByMovieID(movieID int) (count int64, err error) {
	return count, p.db.Model(&model.Comment{}).Where("swapi_movie_id = ?", movieID).Count(&count).Error
}
