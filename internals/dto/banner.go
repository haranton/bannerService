package dto

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
