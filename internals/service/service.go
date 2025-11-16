package service

import (
	"bannerService/internals/storage"
	"log/slog"
)

type Service struct {
	SrvAnswer   *bannerService
	SrvQuestion *QuestionService
	storage     storage.Storage
}

func NewService(st storage.Storage, logger *slog.Logger) *Service {
	srvAnswer := NewbannerService(st)
	srvQuestion := NewQuestionService(st)

	return &Service{
		SrvAnswer:   srvAnswer,
		SrvQuestion: srvQuestion,
		storage:     st,
	}
}
