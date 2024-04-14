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

func (r *BannerRepo) GetBanners(ctx context.Context, featureID, tagID int64, limit, offset int) ([]*repository.Banner, error) {
	var banners []*repository.Banner
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

func (r *BannerRepo) CreateBanner(ctx context.Context, banner *repository.Banner) (int64, error) {
	bannerQuery := `INSERT INTO banners (content, is_active) VALUES ($1, $2) RETURNING id`
	var bannerID int64
	err := r.db.Get(ctx, &bannerID, bannerQuery, banner.Content, banner.IsActive)
	if err != nil {
		return 0, err
	}

	log.Println(banner.TagIDs)

	for _, tagID := range banner.TagIDs {
		featureQuery := `INSERT INTO featuretag (banner_id, feature_id, tag_id) VALUES ($1, $2, $3)`
		_, err := r.db.Exec(ctx, featureQuery, bannerID, banner.FeatureID, tagID)
		if err != nil {
			return 0, err
		}
	}

	return bannerID, nil
}

func (r *BannerRepo) UpdateBanner(ctx context.Context, bannerID int64, banner *repository.Banner) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	query := `
        UPDATE banners
        SET content = $1, is_active = $2, updated_at = $3
        WHERE id = $4
    `
	_, err = tx.Exec(ctx, query, banner.Content, banner.IsActive, time.Now(), bannerID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *BannerRepo) UpdateFeatureTags(ctx context.Context, bannerID int64, featureID int, tagIDs []int) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM featuretag WHERE banner_id = $1", bannerID)
	if err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		_, err := tx.Exec(ctx, "INSERT INTO featuretag (feature_id, tag_id, banner_id) VALUES ($1, $2, $3)", featureID, tagID, bannerID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *BannerRepo) DeleteBanner(ctx context.Context, bannerID int64) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `DELETE FROM featuretag WHERE banner_id = $1`, bannerID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM banners WHERE id = $1`, bannerID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *BannerRepo) DeleteByFeatureIDHandler(ctx context.Context, id int64) error {
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

	_, err = tx.Exec(ctx, `DELETE FROM featuretag WHERE feature_id = $1`, id)
	if err != nil {
		return err
	}

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

func (r *BannerRepo) DeleteByTagIDHandler(ctx context.Context, idtag int64) error {
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

	// Удаляем строки из таблицы featuretag по переданному idtag
	if _, err := tx.Exec(ctx, `DELETE FROM featuretag WHERE tag_id = $1`, idtag); err != nil {
		return err
	}

	// Получаем список id баннеров, связанных с данным тегом
	bannerIDs := []int64{}
	rows, err := tx.Query(ctx, `SELECT banner_id FROM featuretag WHERE tag_id = $1`, idtag)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var bannerID int64
		if err := rows.Scan(&bannerID); err != nil {
			return err
		}
		bannerIDs = append(bannerIDs, bannerID)
	}

	// Удаляем строки из таблицы banners, связанные с найденными id баннеров
	for _, bannerID := range bannerIDs {
		log.Println("Deleting banner with ID:", bannerID)
		if _, err := tx.Exec(ctx, `DELETE FROM banners WHERE id = $1`, bannerID); err != nil {
			return err
		}
	}

	log.Println("Committing transaction...")
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	log.Println("Transaction committed successfully!")
	return nil
}

func (r *BannerRepo) GetBanner(ctx context.Context, bannerID int64) (*repository.Banner, error) {
	var banner repository.Banner
	query := `
		SELECT id, content, is_active, created_at, updated_at
		FROM banners
		WHERE id = $1
		LIMIT 1
	`

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, query, bannerID)
	err = row.Scan(&banner.ID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &banner, nil
}

func (r *BannerRepo) GetBannerVersionsCount(ctx context.Context, bannerID int64) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM banner_versions
		WHERE banner_id = $1
	`

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Выполняем запрос и получаем результат
	row := tx.QueryRow(ctx, query, bannerID)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *BannerRepo) CreateBannerVersion(ctx context.Context, banner *repository.Banner) error {
	query := `
		INSERT INTO banner_versions (banner_id, version_number, content, is_active, created_at)
		VALUES ($1, (SELECT COALESCE(MAX(version_number), 0) + 1 FROM banner_versions WHERE banner_id = $1), $2, $3, $4)
	`

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, banner.ID, banner.Content, banner.IsActive, banner.CreatedAt)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *BannerRepo) DeleteOldestBannerVersion(ctx context.Context, bannerID int64) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Находим ID самой старой версии баннера
	var oldestVersionID int64
	query := `SELECT id FROM banner_versions WHERE banner_id = $1 ORDER BY created_at ASC LIMIT 1`
	err = tx.QueryRow(ctx, query, bannerID).Scan(&oldestVersionID)
	if err != nil {
		return err
	}

	// Удаляем самую старую версию баннера
	_, err = tx.Exec(ctx, "DELETE FROM banner_versions WHERE id = $1", oldestVersionID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// banner.go

func (r *BannerRepo) GetVersionHandler(ctx context.Context, id int64) ([]*repository.BannerVersion, error) {
	var banners []*repository.BannerVersion
	query := `
SELECT *
FROM banner_versions
WHERE banner_id = $1
`

	err := r.db.Select(ctx, &banners, query, id)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (r *BannerRepo) ApplyVersionHandler(ctx context.Context, id, version int64) error {
	// Получаем текущую версию баннера
	var currentBanner repository.Banner
	currentBannerQuery := `
SELECT content, is_active, created_at
FROM banners
WHERE id = $1
LIMIT 1
`
	err := r.db.Get(ctx, &currentBanner, currentBannerQuery, id)
	if err != nil {
		return err
	}

	// Получаем информацию о выбранной версии баннера
	var selectedBanner repository.BannerVersion
	selectedBannerQuery := `
SELECT banner_id, version_number, content, is_active, created_at
FROM banner_versions
WHERE banner_id = $1 AND version_number = $2
`
	err = r.db.Get(ctx, &selectedBanner, selectedBannerQuery, id, version)
	if err != nil {
		return err
	}

	// Обновляем текущую версию баннера
	updateCurrentBannerQuery := `
UPDATE banners
SET content = $1, is_active = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $3
`
	_, err = r.db.Exec(ctx, updateCurrentBannerQuery, selectedBanner.Content, selectedBanner.IsActive, selectedBanner.BannerID)
	if err != nil {
		return err
	}

	// Обновляем выбранную версию баннера, чтобы она стала текущей версией
	updateSelectedBannerQuery := `
UPDATE banner_versions
SET content = $1, is_active = $2, created_at = $3
WHERE banner_id = $4 AND version_number = $5
`
	log.Println("currentBanner.Content", currentBanner.Content)
	log.Println("currentBanner.IsActive", currentBanner.IsActive)
	log.Println("currentBanner.CreatedAt", currentBanner.CreatedAt)
	log.Println("selectedBanner.BannerID", selectedBanner.BannerID)
	log.Println("selectedBanner.Version", selectedBanner.Version)

	_, err = r.db.Exec(ctx, updateSelectedBannerQuery, currentBanner.Content, currentBanner.IsActive, currentBanner.CreatedAt, selectedBanner.BannerID, selectedBanner.Version)
	if err != nil {
		return err
	}

	return nil
}
