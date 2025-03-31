package domain

import (
	"gorm.io/gorm"
	"time"

	"github.com/google/uuid"
)


type Newsletter struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Token     string         `json:"-" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName specifies the table name for GORM
func (Newsletter) TableName() string {
	return "newsletters"
}

// BeforeCreate hook for GORM to set UUID
func (n *Newsletter) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

type NewsletterRepository interface {
	Create(newsletter *Newsletter) error
	FindByEmail(email string) (*Newsletter, error)
	FindByToken(token string) (*Newsletter, error)
	Delete(id uuid.UUID) error
	FindAllActive(page, size int) ([]*Newsletter, int64, error)
}

type NewsletterService interface {
	Subscribe(email string) (*Newsletter, error)
	Unsubscribe(token string) error
	GenerateUnsubscribeToken() string
	GetAllActiveSubscribers() ([]*Newsletter, error)
}