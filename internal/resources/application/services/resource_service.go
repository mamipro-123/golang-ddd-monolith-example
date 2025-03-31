package services

import (
	"errors"
	"monolith-domain/internal/resources/domain"
	"github.com/google/uuid"
)

type ResourceService struct {
	repo domain.ResourceRepository
}

func NewResourceService(repo domain.ResourceRepository) *ResourceService {
	return &ResourceService{repo: repo}
}

func (s *ResourceService) CreateResource(key, value, langCode string) (*domain.Resource, error) {
	// Aynı key ve lang_code kombinasyonunun var olup olmadığını kontrol et
	existing, err := s.repo.FindByKeyAndLang(key, langCode)
	if err == nil && existing != nil {
		return nil, errors.New("resource with this key and language code already exists")
	}

	resource := &domain.Resource{
		Key:      key,
		Value:    value,
		LangCode: langCode,
	}

	if err := s.repo.Create(resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *ResourceService) UpdateResource(id uuid.UUID, value string) (*domain.Resource, error) {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resource.Value = value
	if err := s.repo.Update(resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *ResourceService) DeleteResource(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *ResourceService) GetResourceByID(id uuid.UUID) (*domain.Resource, error) {
	return s.repo.FindByID(id)
}

func (s *ResourceService) GetResourceByKeyAndLang(key, langCode string) (*domain.Resource, error) {
	return s.repo.FindByKeyAndLang(key, langCode)
}

func (s *ResourceService) GetAllResourcesByLang(langCode string) ([]*domain.Resource, error) {
	return s.repo.FindAllByLang(langCode)
}

func (s *ResourceService) GetAllResources(page, size int) ([]*domain.Resource, int64, error) {
	return s.repo.FindAll(page, size)
} 