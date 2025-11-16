package service

import (
	"bannerService/internals/models"
	"bannerService/internals/storage"
	"context"
	"errors"

	"gorm.io/gorm"
)

type QuestionService struct {
	storage storage.Storage
}

func NewQuestionService(storage storage.Storage) *QuestionService {
	return &QuestionService{storage: storage}
}

func (q *QuestionService) Questions(ctx context.Context) ([]models.Question, error) {
	return q.storage.Questions(ctx)
}

func (q *QuestionService) QuestionWithAnswers(ctx context.Context, id int) (*models.Question, error) {
	question, err := q.storage.QuestionWithAnswers(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}
	return question, nil
}

func (q *QuestionService) CreateQuestion(ctx context.Context, question *models.Question) (*models.Question, error) {
	return q.storage.CreateQuestion(ctx, question)
}

func (q *QuestionService) DeleteQuestion(ctx context.Context, id int) error {
	if _, err := q.storage.Question(ctx, id); err != nil {
		return err
	}

	return q.storage.DeleteQuestion(ctx, id)
}
