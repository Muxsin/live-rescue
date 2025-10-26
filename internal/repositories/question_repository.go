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

func (r *QuestionRepository) GetAll() ([]models.Question, error) {
	var questions []models.Question

	result := r.Db.Find(&questions)
	if result.Error != nil {
		return nil, result.Error
	}

	return questions, nil
}

func (r *QuestionRepository) GetOne(id uint) (*models.Question, error) {
	var questions models.Question

	result := r.Db.Where("id = ?", id).First(&questions)
	if result.Error != nil {
		return nil, result.Error
	}

	return &questions, nil
}
