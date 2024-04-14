package fixtures

import (
	"Avito/internal/storage/repository"
	"Avito/tests/states"
	"time"
)

type BannerBuilder struct {
	instance *repository.Banner
}

func Banner() *BannerBuilder {
	return &BannerBuilder{instance: &repository.Banner{}}
}

func (b *BannerBuilder) ID(v int) *BannerBuilder {
	b.instance.ID = v
	return b
}

func (b *BannerBuilder) TagIDs(v []int) *BannerBuilder {
	b.instance.TagIDs = v
	return b
}

func (b *BannerBuilder) FeatureID(v int) *BannerBuilder {
	b.instance.FeatureID = v
	return b
}

func (b *BannerBuilder) Content(v map[string]interface{}) *BannerBuilder {
	b.instance.Content = v
	return b
}

func (b *BannerBuilder) IsActive(v bool) *BannerBuilder {
	b.instance.IsActive = v
	return b
}

func (b *BannerBuilder) CreatedAt(v time.Time) *BannerBuilder {
	b.instance.CreatedAt = v
	return b
}

func (b *BannerBuilder) UpdatedAt(v time.Time) *BannerBuilder {
	b.instance.UpdatedAt = v
	return b
}

func (b *BannerBuilder) P() *repository.Banner {
	return b.instance
}

func (b *BannerBuilder) V() repository.Banner {
	return *b.instance
}

func (b *BannerBuilder) Valid() *BannerBuilder {
	return Banner().ID(states.Banner1ID).TagIDs(states.Banner1TagIDs).FeatureID(states.Banner1FeatureID).Content(states.Banner1Content).IsActive(states.Banner1IsActive).CreatedAt(states.Banner1CreatedAt).UpdatedAt(states.Banner1UpdatedAt)
}
