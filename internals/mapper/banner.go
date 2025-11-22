package mapper

import (
	"bannerService/internals/dto"
	"bannerService/internals/models"
)

func FeatureTagsBanner(request *dto.FeatureTags, bannerId int) []*models.FeatureTagBanner {

	featuresTags := make([]*models.FeatureTagBanner, len(request.TagIds))
	for _, tag := range request.TagIds {
		featureTagsBanner := models.FeatureTagBanner{
			TagID:     tag,
			FeatureID: request.FeatureId,
			BannerID:  bannerId,
		}
		featuresTags = append(featuresTags, &featureTagsBanner)
	}

	return featuresTags
}
