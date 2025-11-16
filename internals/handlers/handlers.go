package handlers

import (
	"bannerService/internals/config"
	"bannerService/internals/service"
	"log/slog"
	"net/http"
)

var (
	UserIdErrIsRequired = "user id is required"
)

type Handler struct {
	service *service.Service
	logger  *slog.Logger
	cfg     *config.Config
}

func NewHandler(service *service.Service, logger *slog.Logger, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
		cfg:     cfg,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /api/questions/", h.GetQuestions)
	mux.HandleFunc("POST /api/questions/", h.CreateQuestion)
	mux.HandleFunc("GET /api/questions/{id}", h.GetQuestionWithAnswers)
	mux.HandleFunc("DELETE /api/questions/{id}", h.DeleteQuestion)

	mux.HandleFunc("POST /api/questions/{id}/answers/", h.CreateAnswer)
	mux.HandleFunc("GET /api/answers/{id}", h.GetAnswer)
	mux.HandleFunc("DELETE /api/answers/{id}", h.DeleteAnswer)
}
