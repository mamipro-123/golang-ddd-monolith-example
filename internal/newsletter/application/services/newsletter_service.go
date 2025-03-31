package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"monolith-domain/internal/newsletter/domain"
	"github.com/google/uuid"
)

type NewsletterService struct {
	repo domain.NewsletterRepository
}

func NewNewsletterService(repo domain.NewsletterRepository) *NewsletterService {
	return &NewsletterService{repo: repo}
}

func (s *NewsletterService) Subscribe(email string) (*domain.Newsletter, error) {
	// Check if already subscribed
	existing, err := s.repo.FindByEmail(email)
	if err == nil && existing != nil {
		return nil, errors.New("email already subscribed")
	}

	// Create new subscription
	newsletter := &domain.Newsletter{
		ID:    uuid.New(),
		Email: email,
		Token: s.GenerateUnsubscribeToken(),
	}

	if err := s.repo.Create(newsletter); err != nil {
		return nil, err
	}

	return newsletter, nil
}

func (s *NewsletterService) Unsubscribe(token string) error {
	newsletter, err := s.repo.FindByToken(token)
	if err != nil {
		return err
	}

	return s.repo.Delete(newsletter.ID)
}

func (s *NewsletterService) GenerateUnsubscribeToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
} 

func (s *NewsletterService) GetAllActiveSubscribers(page, size int) ([]*domain.Newsletter, int64, error) {
	return s.repo.FindAllActive(page, size)
}
