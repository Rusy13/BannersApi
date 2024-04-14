package repository

import "context"

type BannerRepo interface {
	GetUserBanner(ctx context.Context, tagID, featureID int64, useLastRevision bool) (*Banner, error)
	GetBanners(ctx context.Context, featureID, tagID int64, limit, offset int) ([]*Banner, error)
	CreateBanner(ctx context.Context, banner *Banner) (int64, error)
	UpdateBanner(ctx context.Context, bannerID int64, banner *Banner) error
	DeleteBanner(ctx context.Context, bannerID int64) error
	UpdateFeatureTags(ctx context.Context, bannerID int64, featureID int, tagIDs []int) error
	DeleteByFeatureIDHandler(ctx context.Context, bannerID int64) error
	DeleteByTagIDHandler(ctx context.Context, bannerID int64) error
	GetBanner(ctx context.Context, bannerID int64) (*Banner, error)
	GetBannerVersionsCount(ctx context.Context, bannerID int64) (int, error)
	CreateBannerVersion(ctx context.Context, banner *Banner) error
	DeleteOldestBannerVersion(ctx context.Context, bannerID int64) error
	//GetVersionHandler(ctx context.Context, bannerID int64, bannerVersion int64) error
	//ApplyVersionHandler(ctx context.Context, bannerID int64, bannerVersion int64) error

}
