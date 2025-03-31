package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Resource struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Key       string         `json:"key" gorm:"uniqueIndex:idx_key_lang;not null"`
	Value     string         `json:"value" gorm:"type:text;not null"`
	LangCode  string         `json:"lang_code" gorm:"uniqueIndex:idx_key_lang;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName specifies the table name for GORM
func (Resource) TableName() string {
	return "resources"
}

// BeforeCreate hook for GORM to set UUID
func (r *Resource) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type ResourceRepository interface {
	Create(resource *Resource) error
	Update(resource *Resource) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*Resource, error)
	FindByKeyAndLang(key, langCode string) (*Resource, error)
	FindAllByLang(langCode string) ([]*Resource, error)
	FindAll(page, size int) ([]*Resource, int64, error)
}

type ResourceService interface {
	CreateResource(key, value, langCode string) (*Resource, error)
	UpdateResource(id uuid.UUID, value string) (*Resource, error)
	DeleteResource(id uuid.UUID) error
	GetResourceByID(id uuid.UUID) (*Resource, error)
	GetResourceByKeyAndLang(key, langCode string) (*Resource, error)
	GetAllResourcesByLang(langCode string) ([]*Resource, error)
	GetAllResources(page, size int) ([]*Resource, int64, error)
} 