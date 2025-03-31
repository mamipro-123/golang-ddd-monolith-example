package infrastructure

import (
	"monolith-domain/internal/newsletter/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(newsletter *domain.Newsletter) error {
	return r.db.Create(newsletter).Error
}

func (r *PostgresRepository) FindByEmail(email string) (*domain.Newsletter, error) {
	var newsletter domain.Newsletter
	err := r.db.Where("email = ?", email).First(&newsletter).Error
	if err != nil {
		return nil, err
	}
	return &newsletter, nil
}

func (r *PostgresRepository) FindByToken(token string) (*domain.Newsletter, error) {
	var newsletter domain.Newsletter
	err := r.db.Where("token = ?", token).First(&newsletter).Error
	if err != nil {
		return nil, err
	}
	return &newsletter, nil
}

func (r *PostgresRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Newsletter{}, id).Error
} 

func (r *PostgresRepository) FindAllActive(page, size int) ([]*domain.Newsletter, int64, error) {
	var newsletters []*domain.Newsletter
	var total int64

	// Toplam kayıt sayısını al
	if err := r.db.Model(&domain.Newsletter{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sayfalama ile kayıtları al
	offset := (page - 1) * size
	err := r.db.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&newsletters).Error

	return newsletters, total, err
}
