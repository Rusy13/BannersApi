package postgresql

import (
	errs "Avito/internal/config/errors"
	"Avito/internal/storage/db"
	"Avito/internal/storage/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"log"
	"time"
)

type BannerRepo struct {
	db db.PGX
}

func NewBannerRepo(database db.PGX) *BannerRepo {
	return &BannerRepo{db: database}
}

// ssssssssssssssssssssssssssssssssiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu
// "Banner not found", это означает, что сервер не смог найти баннер, который соответствует заданным критериям (тэг и фича) и который был актуальным в течение последних 5 минут.
func (r *BannerRepo) GetUserBanner(ctx context.Context, tagID, featureID int64, useLastRevision bool) (*repository.Banner, error) {
	var banner repository.Banner
	log.Println(tagID)
	log.Println(featureID)

	// Запрос для выбора баннера по tag_id и feature_id
	query := `
        SELECT b.id, ARRAY_AGG(ft.tag_id) AS tag_ids, ft.feature_id, b.content, b.is_active, b.created_at, b.updated_at
        FROM banners b
        INNER JOIN featuretag ft ON b.id = ft.banner_id
        WHERE ft.tag_id = $1 AND ft.feature_id = $2
        GROUP BY b.id, ft.feature_id, b.content, b.is_active, b.created_at, b.updated_at
        ORDER BY b.updated_at DESC
        LIMIT 1
    `

	err := r.db.Get(ctx, &banner, query, tagID, featureID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.ErrBannerNotFound
	} else if err != nil {
		return nil, err
	}

	return &banner, nil
}

// SSSSSSSSSSSSSSSSSSSSSSSSSIIIIIIIIIIIIIIIIIIIIIIIIIIIIIUUUUUUUUUUUUUUUUUUUUUUUUUUU
// GetBanners retrieves banners with optional feature and tag filtering, with pagination.
func (r *BannerRepo) GetBanners(ctx context.Context, featureID, tagID int64, limit, offset int) ([]*repository.Banner, error) {
	var banners []*repository.Banner
	log.Println(featureID)
	log.Println(tagID)
	log.Println(limit)
	log.Println(offset)
	query := `
SELECT b.id, ARRAY_AGG(ft.tag_id) AS tag_ids, ft.feature_id, b.content, b.is_active, b.created_at, b.updated_at
FROM banners b
INNER JOIN featuretag ft ON b.id = ft.banner_id
WHERE (ft.feature_id = $1 OR $1 IS NULL) AND (ft.tag_id = $2 OR $2 IS NULL)
GROUP BY b.id, ft.feature_id, b.content, b.is_active, b.created_at, b.updated_at
ORDER BY b.updated_at DESC
LIMIT $3 OFFSET $4
`

	err := r.db.Select(ctx, &banners, query, featureID, tagID, limit, offset)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

// SSSSSSSSSSSSSSSSSSSIIIIIIIIIIIIIIIIIIIIIIIIIIIIIUUUUUUUUUUUUUUUUUUUUUUUUUUUUU
func (r *BannerRepo) CreateBanner(ctx context.Context, banner *repository.Banner) (int64, error) {
	// Вставляем запись в таблицу banners
	bannerQuery := `INSERT INTO banners (content, is_active) VALUES ($1, $2) RETURNING id`
	var bannerID int64
	err := r.db.Get(ctx, &bannerID, bannerQuery, banner.Content, banner.IsActive)
	if err != nil {
		return 0, err
	}

	log.Println(banner.TagIDs)

	// Для каждого tag_id создаем запись в таблице featuretag с тем же feature_id
	for _, tagID := range banner.TagIDs {
		featureQuery := `INSERT INTO featuretag (banner_id, feature_id, tag_id) VALUES ($1, $2, $3)`
		_, err := r.db.Exec(ctx, featureQuery, bannerID, banner.FeatureID, tagID)
		if err != nil {
			return 0, err
		}
	}

	return bannerID, nil
}

// SSSSSSSSSSSSSSSSSSSSSSIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIUUUUUUUUUUUUUUUUUUUUUU
func (r *BannerRepo) UpdateBanner(ctx context.Context, bannerID int64, banner *repository.Banner) error {
	query := `
		UPDATE banners
		SET content = $1, is_active = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query, banner.Content, banner.IsActive, time.Now(), bannerID)
	return err
}

func (r *BannerRepo) UpdateFeatureTags(ctx context.Context, bannerID int64, featureID int, tagIDs []int) error {
	// Удаляем существующие теги для данного баннера
	_, err := r.db.Exec(ctx, "DELETE FROM featuretag WHERE banner_id = $1", bannerID)
	if err != nil {
		return err
	}

	// Добавляем новые теги
	for _, tagID := range tagIDs {
		_, err := r.db.Exec(ctx, "INSERT INTO featuretag (feature_id, tag_id, banner_id) VALUES ($1, $2, $3)", featureID, tagID, bannerID)
		if err != nil {
			return err
		}
	}

	return err
}

// SSSSSSSSSSSSSSSSSSSSSSIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIUUUUUUUUUUUUUUUUUUUUUU
func (r *BannerRepo) DeleteBanner(ctx context.Context, bannerID int64) error {
	// Начнем транзакцию
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Удаляем связанные записи из таблицы featuretag
	_, err = tx.Exec(ctx, `DELETE FROM featuretag WHERE banner_id = $1`, bannerID)
	if err != nil {
		return err
	}

	// Удаляем баннер из таблицы banners
	_, err = tx.Exec(ctx, `DELETE FROM banners WHERE id = $1`, bannerID)
	if err != nil {
		return err
	}

	// Фиксируем транзакцию
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// ------------------------------------------
func (r *BannerRepo) DeleteByFeatureIDHandler(ctx context.Context, id int64) error {
	// Начинаем транзакцию
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	log.Println("DeleteByFeatureIDHandler")

	// Находим все banner_id, связанные с feature_id
	bannerIDs := []int64{}
	rows, err := tx.Query(ctx, `SELECT banner_id FROM featuretag WHERE feature_id = $1`, id)
	if err != nil {
		return err
	}
	for rows.Next() {
		var bannerID int64
		if err := rows.Scan(&bannerID); err != nil {
			rows.Close()
			return err
		}
		bannerIDs = append(bannerIDs, bannerID)
	}
	rows.Close()

	// Удаляем связанные записи из таблицы featuretag
	_, err = tx.Exec(ctx, `DELETE FROM featuretag WHERE feature_id = $1`, id)
	if err != nil {
		return err
	}

	// Удаляем связанные баннеры из таблицы banners
	for _, bannerID := range bannerIDs {
		_, err = tx.Exec(ctx, `DELETE FROM banners WHERE id = $1`, bannerID)
		if err != nil {
			return err
		}
	}

	// Фиксируем транзакцию
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------
func (r *BannerRepo) DeleteByTagIDHandler(ctx context.Context, id int64) error {
	// Начинаем транзакцию
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			log.Println("Rolling back transaction...")
			tx.Rollback(ctx)
		}
	}()

	log.Println("DeleteByTagIDHandler")

	// Находим все banner_id, связанные с feature_id
	tagIDs := []int64{}
	rows, err := tx.Query(ctx, `SELECT banner_id FROM featuretag WHERE tag_id = $1`, id)
	if err != nil {
		return err
	}
	defer rows.Close() // Убедимся, что строки будут закрыты после использования
	for rows.Next() {
		var bannerID int64
		if err := rows.Scan(&bannerID); err != nil {
			return err
		}
		tagIDs = append(tagIDs, bannerID)
	}

	// Удаляем связанные записи из таблицы featuretag
	if _, err := tx.Exec(ctx, `DELETE FROM featuretag WHERE tag_id = $1`, id); err != nil {
		return err
	}

	// Удаляем связанные баннеры из таблицы banners
	for _, tagID := range tagIDs {
		log.Println("Deleting banner with ID:", tagID)
		if _, err := tx.Exec(ctx, `DELETE FROM banners WHERE id = $1`, tagID); err != nil {
			return err
		}
	}

	// Фиксируем транзакцию
	log.Println("Committing transaction...")
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	log.Println("Transaction committed successfully!")
	return nil
}
