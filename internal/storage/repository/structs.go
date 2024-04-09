package repository

import (
	"errors"
	"time"
)

var ErrObjectNotFound = errors.New("not found")

// BannerDB представляет собой структуру данных для баннера в базе данных
type Banner struct {
	ID        int   `db:"id"`         // Идентификатор баннера
	FeatureID int   `db:"feature_id"` // Идентификатор фичи
	TagIDs    []int // Идентификаторы тэгов
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		URL   string `json:"url"`
	} `json:"content"` // Содержимое баннера
	IsActive  bool      `db:"is_active"`  // Флаг активности баннера
	CreatedAt time.Time `db:"created_at"` // Дата создания баннера
	UpdatedAt time.Time `db:"updated_at"` // Дата обновления баннера
}
