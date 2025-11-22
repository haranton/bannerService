package storage

import (
	"bannerService/internals/dto"
	"bannerService/internals/models"
	"context"
	"errors"
)

var (
	ErrBannerNotFound = errors.New("banner not found")
)

type DB interface {
	Close() error
}

type Storage interface {
	BannerStorage
	DB
}

type BannerStorage interface {
	Banner(ctx context.Context, params dto.BannerQuery) (*models.Banner, error)
	Banners(ctx context.Context, params dto.BannersQuery) ([]*models.Banner, error)
	CreateBanner(ctx context.Context, banner *models.Banner, featureTags *dto.FeatureTags) (*models.Banner, error)
	UpdateBanner(ctx context.Context, banner *models.Banner, featureTagBanner []*models.FeatureTagBanner) error
	DeleteBanner(ctx context.Context, idBanner int) error
}
