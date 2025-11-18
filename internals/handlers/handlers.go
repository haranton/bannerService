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

	mux.HandleFunc("GET /api/user_banner/", h.Banner)
	mux.HandleFunc("GET /api/banner/", h.Banners)
	mux.HandleFunc("POST /api/banner/", h.CreateBanner)
	mux.HandleFunc("PATCH /api/banner/{id}", h.PatchBanner)
	mux.HandleFunc("DELETE /api/banner/{id}", h.DeleteBanner)

}
