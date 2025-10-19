package repositories

import (
	"live-rescue/internal/models"

	"gorm.io/gorm"
)

type QuestionRepository struct {
	Db *gorm.DB
}

func NewQuestion(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{
		Db: db,
	}
}

func (r *QuestionRepository) Create(question *models.Question) error {
	return r.Db.Create(question).Error
}
