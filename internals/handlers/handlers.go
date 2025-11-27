package handlers

import (
	"bannerService/internals/config"
	"bannerService/internals/middleware"
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
	base := &RouteGroup{mux: mux}
	auth := base.With(middleware.AuthMiddleware)
	admin := auth.With(middleware.AdminOnlyMiddleware)

	// Аутентифицированные маршруты
	auth.HandleFunc("GET /api/user_banner", h.Banner)

	// Административные маршруты
	admin.HandleFunc("GET /api/banner", h.Banners)
	admin.HandleFunc("POST /api/banner", h.CreateBanner)
	admin.HandleFunc("PATCH /api/banner/{id}", h.PatchBanner)
	admin.HandleFunc("DELETE /api/banner/{id}", h.DeleteBanner)
}
