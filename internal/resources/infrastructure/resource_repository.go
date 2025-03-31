package infrastructure

import (
	"monolith-domain/internal/resources/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(resource *domain.Resource) error {
	return r.db.Create(resource).Error
}

func (r *PostgresRepository) Update(resource *domain.Resource) error {
	return r.db.Save(resource).Error
}

func (r *PostgresRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Resource{}, id).Error
}

func (r *PostgresRepository) FindByID(id uuid.UUID) (*domain.Resource, error) {
	var resource domain.Resource
	err := r.db.First(&resource, id).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *PostgresRepository) FindByKeyAndLang(key, langCode string) (*domain.Resource, error) {
	var resource domain.Resource
	err := r.db.Where("key = ? AND lang_code = ?", key, langCode).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *PostgresRepository) FindAllByLang(langCode string) ([]*domain.Resource, error) {
	var resources []*domain.Resource
	err := r.db.Where("lang_code = ?", langCode).Find(&resources).Error
	return resources, err
}

func (r *PostgresRepository) FindAll(page, size int) ([]*domain.Resource, int64, error) {
	var resources []*domain.Resource
	var total int64

	// Toplam kayıt sayısını al
	if err := r.db.Model(&domain.Resource{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sayfalama ile kayıtları al
	offset := (page - 1) * size
	err := r.db.Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&resources).Error

	return resources, total, err
} 