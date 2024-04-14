package service

import (
	"Avito/internal/storage/repository"
	"context"
)

type Server1 struct {
	Repo repository.BannerRepo
}

type BannerResponse struct {
	BannerID  int64       `json:"banner_id"`
	TagIDs    []int64     `json:"tag_ids"`
	FeatureID int64       `json:"feature_id"`
	Content   interface{} `json:"content"`
	IsActive  bool        `json:"is_active"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

func (s *Server1) GetUserBanner(ctx context.Context, tagID, featureID int64, useLastRevision bool) (*repository.Banner, error) {
	banner, err := s.Repo.GetUserBanner(ctx, tagID, featureID, useLastRevision)
	return banner, err
}

func (s *Server1) GetBanners(ctx context.Context, featureID, tagID int64, limit, offset int) ([]*repository.Banner, error) {
	var banners []*repository.Banner
	banners, err := s.Repo.GetBanners(ctx, tagID, featureID, limit, offset)
	return banners, err
}
func (s *Server1) CreateBanner(ctx context.Context, banner *repository.Banner) (int64, error) {
	var bannerID int64
	bannerID, err := s.Repo.CreateBanner(ctx, banner)
	return bannerID, err
}

func (s *Server1) UpdateBanner(ctx context.Context, bannerID int64, banner *repository.Banner) error {
	err := s.Repo.UpdateBanner(ctx, bannerID, banner)
	return err
}

func (s *Server1) DeleteBanner(ctx context.Context, bannerID int64) error {
	err := s.Repo.DeleteBanner(ctx, bannerID)
	return err
}
func (s *Server1) UpdateFeatureTags(ctx context.Context, bannerID int64, featureID int, tagIDs []int) error {
	err := s.Repo.UpdateFeatureTags(ctx, bannerID, featureID, tagIDs)
	return err
}
func (s *Server1) DeleteByFeatureIDHandler(ctx context.Context, bannerID int64) error {
	err := s.Repo.DeleteByFeatureIDHandler(ctx, bannerID)
	return err
}
func (s *Server1) DeleteByTagIDHandler(ctx context.Context, bannerID int64) error {
	err := s.Repo.DeleteByTagIDHandler(ctx, bannerID)
	return err
}
func (s *Server1) GetBanner(ctx context.Context, bannerID int64) (*repository.Banner, error) {
	var banner *repository.Banner
	banner, err := s.Repo.GetBanner(ctx, bannerID)
	return banner, err
}
func (s *Server1) GetBannerVersionsCount(ctx context.Context, bannerID int64) (int, error) {
	var count int
	count, err := s.Repo.GetBannerVersionsCount(ctx, bannerID)
	return count, err
}
func (s *Server1) CreateBannerVersion(ctx context.Context, banner *repository.Banner) error {
	err := s.Repo.CreateBannerVersion(ctx, banner)
	return err
}
func (s *Server1) DeleteOldestBannerVersion(ctx context.Context, bannerID int64) error {
	err := s.Repo.DeleteOldestBannerVersion(ctx, bannerID)
	return err
}
func (s *Server1) GetVersionHandler(ctx context.Context, id int64) ([]*repository.BannerVersion, error) {
	var banners []*repository.BannerVersion
	banners, err := s.Repo.GetVersionHandler(ctx, id)
	return banners, err
}
func (s *Server1) ApplyVersionHandler(ctx context.Context, bannerID int64, bannerVersion int64) error {
	err := s.Repo.ApplyVersionHandler(ctx, bannerID, bannerVersion)
	return err
}
