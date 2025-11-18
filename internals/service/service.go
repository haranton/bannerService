package service

import (
	"bannerService/internals/storage"
	"log/slog"
)

type Service struct {
	SrvBanner *bannerService
	storage   storage.Storage
}

func NewService(st storage.Storage, logger *slog.Logger) *Service {
	SrvBanner := NewBannerService(st)

	return &Service{
		SrvBanner: SrvBanner,
		storage:   st,
	}
}
