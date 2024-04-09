package repository

import "context"

type BannerRepo interface {
	GetUserBanner(ctx context.Context, tagID, featureID int64, useLastRevision bool) (*Banner, error)
	GetBanners(ctx context.Context, featureID, tagID int64, limit, offset int) ([]*Banner, error)
	CreateBanner(ctx context.Context, banner *Banner) (int64, error)
	UpdateBanner(ctx context.Context, bannerID int64, banner *Banner) error
	DeleteBanner(ctx context.Context, bannerID int64) error
}
