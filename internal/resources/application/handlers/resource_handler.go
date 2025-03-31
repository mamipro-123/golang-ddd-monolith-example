package handlers

import (
	"monolith-domain/internal/resources/application/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ResourceHandler struct {
	service *services.ResourceService
}

func NewResourceHandler(service *services.ResourceService) *ResourceHandler {
	return &ResourceHandler{service: service}
}

type CreateResourceRequest struct {
	Key      string `json:"key" validate:"required"`
	Value    string `json:"value" validate:"required"`
	LangCode string `json:"lang_code" validate:"required"`
}

type UpdateResourceRequest struct {
	Value string `json:"value" validate:"required"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Size       int        `json:"size"`
	TotalPages int        `json:"total_pages"`
}

func (h *ResourceHandler) CreateResource(c *fiber.Ctx) error {
	var req CreateResourceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resource, err := h.service.CreateResource(req.Key, req.Value, req.LangCode)
	if err != nil {
		if err.Error() == "resource with this key and language code already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Resource with this key and language code already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create resource",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Resource created successfully",
		"data":    resource,
	})
}

func (h *ResourceHandler) UpdateResource(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var req UpdateResourceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resource, err := h.service.UpdateResource(id, req.Value)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update resource",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Resource updated successfully",
		"data":    resource,
	})
}

func (h *ResourceHandler) DeleteResource(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	if err := h.service.DeleteResource(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete resource",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Resource deleted successfully",
	})
}

func (h *ResourceHandler) GetResourceByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	resource, err := h.service.GetResourceByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resource not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": resource,
	})
}

func (h *ResourceHandler) GetResourceByKeyAndLang(c *fiber.Ctx) error {
	key := c.Query("key")
	langCode := c.Query("lang_code")

	if key == "" || langCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Key and lang_code are required",
		})
	}

	resource, err := h.service.GetResourceByKeyAndLang(key, langCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resource not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": resource,
	})
}

func (h *ResourceHandler) GetAllResourcesByLang(c *fiber.Ctx) error {
	langCode := c.Params("lang_code")

	if langCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Language code is required",
		})
	}

	resources, err := h.service.GetAllResourcesByLang(langCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch resources",
		})
	}

	return c.JSON(fiber.Map{
		"data": resources,
	})
}

func (h *ResourceHandler) GetAllResources(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", 10)

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	resources, total, err := h.service.GetAllResources(page, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch resources",
		})
	}

	totalPages := int((total + int64(size) - 1) / int64(size))

	response := PaginationResponse{
		Data:       resources,
		Total:      total,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}

	return c.JSON(response)
} 