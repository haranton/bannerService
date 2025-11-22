package dto

import "encoding/json"

type BannerQuery struct {
	UseLastRevision bool
	Feature_id      int
	Tag_id          int
}

type BannersQuery struct {
	Feature_id int
	Tag_id     int
	Limit      int
	Offset     int
}

type BannerCreateUpdateRequest struct {
	TagIds    []int           `json:"tag_ids"`
	FeatureId int             `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
}

type FeatureTags struct {
	TagIds    []int
	FeatureId int
}
