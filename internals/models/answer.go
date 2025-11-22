package models

import "time"

type Banner struct {
	ID        int       `db:"id"`
	Content   []byte    `db:"content"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type FeatureTagBanner struct {
	TagID     int `db:"tag_id"`
	FeatureID int `db:"feature_id"`
	BannerID  int `db:"banner_id"`
}
