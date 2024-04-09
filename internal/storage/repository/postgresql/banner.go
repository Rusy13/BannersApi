package postgresql

import (
	errs "Avito/internal/config/errors"
	"Avito/internal/storage/db"
	"Avito/internal/storage/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"time"
)

type BannerRepo struct {
	db db.DBops
}

func NewBannerRepo(database db.DBops) *BannerRepo {
	return &BannerRepo{db: database}
}

func (r *BannerRepo) GetUserBanner(ctx context.Context, tagID, featureID int64, useLastRevision bool) (*repository.Banner, error) {
	var banner repository.Banner

	var query string
	if useLastRevision {
		query = `SELECT b.id, b.tag_id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at
		         FROM banners b
		         INNER JOIN feature f ON b.id = f.banner_id
		         WHERE f.tag_id = $1 AND f.feature_id = $2 AND b.is_active = true
		         ORDER BY b.updated_at DESC
		         LIMIT 1`
	} else {
		query = `SELECT b.id, b.tag_id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at
		         FROM banners b
		         INNER JOIN feature f ON b.id = f.banner_id
		         WHERE f.tag_id = $1 AND f.feature_id = $2 AND b.is_active = true
		         ORDER BY b.updated_at DESC
		         LIMIT 1`
	}

	err := r.db.Get(ctx, &banner, query, tagID, featureID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.ErrBannerNotFound
	} else if err != nil {
		return nil, err
	}

	return &banner, nil
}

func (r *BannerRepo) GetBanners(ctx context.Context, featureID, tagID int64, limit, offset int) ([]*repository.Banner, error) {
	var banners []*repository.Banner

	query := `SELECT b.id, b.tag_id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at
	          FROM banners b
	          INNER JOIN feature f ON b.id = f.banner_id
	          WHERE f.feature_id = $1 AND f.tag_id = $2
	          ORDER BY b.updated_at DESC
	          LIMIT $3 OFFSET $4`

	err := r.db.Select(ctx, &banners, query, featureID, tagID, limit, offset)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (r *BannerRepo) CreateBanner(ctx context.Context, banner *repository.Banner) (int64, error) {
	// Вставляем запись в таблицу banners
	bannerQuery := `INSERT INTO banners (feature_id, content, is_active) VALUES ($1, $2, $3) RETURNING id`
	var bannerID int64
	err := r.db.Get(ctx, &bannerID, bannerQuery, banner.FeatureID, banner.Content, banner.IsActive)
	if err != nil {
		return 0, err
	}

	// Для каждого tag_id создаем запись в таблице feature с тем же feature_id
	for _, tagID := range banner.TagIDs {
		featureQuery := `INSERT INTO feature (banner_id, feature_id, tag_id) VALUES ($1, $2, $3)`
		_, err := r.db.Exec(ctx, featureQuery, bannerID, banner.FeatureID, tagID)
		if err != nil {
			return 0, err
		}
	}

	return bannerID, nil
}

func (r *BannerRepo) UpdateBanner(ctx context.Context, bannerID int64, banner *repository.Banner) error {
	query := `UPDATE banners SET tag_id = $1, feature_id = $2, content = $3, is_active = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.Exec(ctx, query, banner.FeatureID, banner.Content, banner.IsActive, time.Now(), bannerID)
	return err
}

func (r *BannerRepo) DeleteBanner(ctx context.Context, bannerID int64) error {
	query := `DELETE FROM banners WHERE id = $1`
	_, err := r.db.Exec(ctx, query, bannerID)
	return err
}
