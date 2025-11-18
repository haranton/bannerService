package service

import (
	"bannerService/internals/dto"
	"bannerService/internals/models"
	"bannerService/internals/storage"
	"context"
)

type bannerService struct {
	storage storage.Storage
}

func NewBannerService(storage storage.Storage) *bannerService {
	return &bannerService{storage: storage}
}

func (b *bannerService) Banner(ctx context.Context, paramsBanner dto.BannerQuery) (*models.Banner, error) {

	banner, err := b.storage.Banner(ctx, paramsBanner)
	if err != nil {
		return nil, err
	}

	return banner, nil
}

func (b *bannerService) Banners(ctx context.Context, params dto.BannersQuery) ([]*models.Banner, error) {

	banners, err := b.storage.Banners(ctx, params)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (b *bannerService) CreateBanner(
	ctx context.Context,
	banner *models.Banner,
	featureTagBanner []*models.FeatureTagBanner) (*models.Banner, error) {

	banner, err := b.storage.CreateBanner(ctx, banner, featureTagBanner)
	if err != nil {
		return nil, err
	}

	return banner, nil
}

func (b *bannerService) UpdateBanner(ctx context.Context,
	banner *models.Banner,
	featureTagBanner []*models.FeatureTagBanner) error {

	err := b.storage.UpdateBanner(ctx, banner, featureTagBanner)
	if err != nil {
		return err
	}

	return nil
}

func (b *bannerService) DeleteBanner(ctx context.Context, idBanner int) error {

	if err := b.storage.DeleteBanner(ctx, idBanner); err != nil {
		return err
	}

	return nil

}
