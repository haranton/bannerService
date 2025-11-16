package models

import "time"

type Banner struct {
	ID        int64     `db:"id"`
	Content   []byte    `db:"content"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type FeatureTagBanner struct {
	TagID     int64 `db:"tag_id"`
	FeatureID int64 `db:"feature_id"`
	BannerID  int64 `db:"banner_id"`
}
