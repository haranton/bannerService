package service

import (
	"bannerService/internals/models"
	"bannerService/internals/storage"
	"context"
	"errors"
)

var (
	ErrQuestionNotFound = errors.New("question not found")
	ErrAnswerNotFound   = errors.New("answer not found")
)

type bannerService struct {
	storage storage.Storage
}

func NewbannerService(storage storage.Storage) *bannerService {
	return &bannerService{storage: storage}
}

func (a *bannerService) Answer(ctx context.Context, id int) (*models.Answer, error) {

	answer, err := a.storage.Answer(ctx, id)
	if err != nil {
		return nil, err
	}
	if answer == nil {
		return nil, ErrAnswerNotFound
	}
	return answer, nil
}

func (a *bannerService) CreateAnswer(ctx context.Context, answer *models.Answer) (*models.Answer, error) {
	question, err := a.storage.Question(ctx, answer.QuestionID)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return nil, ErrQuestionNotFound
	}
	return a.storage.CreateAnswer(ctx, answer)
}

func (a *bannerService) DeleteAnswer(ctx context.Context, id int) error {
	question, err := a.storage.Answer(ctx, id)
	if err != nil {
		return err
	}
	if question == nil {
		return ErrAnswerNotFound
	}
	return a.storage.DeleteAnswer(ctx, id)
}
