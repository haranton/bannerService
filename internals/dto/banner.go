package dto

type BannerQuery struct {
	useLastRevision bool
	feature_id      int
	tag_id          int
}

type BannersQuery struct {
	Feature_id int
	Tag_id     int
	Limit      int
	Offset     int
}
