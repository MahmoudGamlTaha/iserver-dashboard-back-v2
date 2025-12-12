package services

import (
	"enterprise-architect-api/models"
	"enterprise-architect-api/repositories"
	"fmt"
	"math"

	"github.com/google/uuid"
)

// ObjectContentService handles business logic for object contents
type ObjectContentService struct {
	repo *repositories.ObjectContentRepository
}

// NewObjectContentService creates a new ObjectContentService
func NewObjectContentService(repo *repositories.ObjectContentRepository) *ObjectContentService {
	return &ObjectContentService{repo: repo}
}

// CreateObjectContent creates a new object content
func (s *ObjectContentService) CreateObjectContent(req models.CreateObjectContentRequest) (*models.ObjectContent, error) {
	// Validate required fields
	if req.CreatedBy == 0 {
		return nil, fmt.Errorf("created by is required")
	}

	return s.repo.Create(req)
}
func (s *ObjectContentService) CreateObjectContentV2(req models.CreateObjectContentRequest) (*models.ObjectContent, error) {
	// Validate required fields
	if req.CreatedBy == 0 {
		return nil, fmt.Errorf("created by is required")
	}

	return s.repo.CreateV2(req)
}

// GetObjectContentByID retrieves an object content by its ID
func (s *ObjectContentService) GetObjectContentByID(id int) (*models.ObjectContent, error) {
	return s.repo.GetByID(id)
}

// GetAllObjectContents retrieves all object contents with pagination
func (s *ObjectContentService) GetAllObjectContents(page, pageSize int) (*models.PaginatedResponse, error) {
	// Set default pagination values
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	objectContents, totalCount, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data:       objectContents,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// UpdateObjectContent updates an existing object content
func (s *ObjectContentService) UpdateObjectContent(id int, req models.UpdateObjectContentRequest) (*models.ObjectContent, error) {
	// Validate that at least one field is being updated
	if req.DocumentObjectID == nil && req.ContainerVersionID == nil && req.ObjectID == nil &&
		req.Instances == nil && req.IsShortCut == nil && req.ContainmentType == nil {
		return nil, fmt.Errorf("at least one field must be provided for update")
	}

	if req.ModifiedBy == 0 {
		return nil, fmt.Errorf("modified by is required")
	}

	return s.repo.Update(id, req)
}

// DeleteObjectContent deletes an object content by its ID
func (s *ObjectContentService) DeleteObjectContent(id int) error {
	return s.repo.Delete(id)
}

func (s *ObjectContentService) DashboardCount(libraryId uuid.UUID) ([]models.DashboardCount, error) {
	return s.repo.DashboardCount(libraryId)
}

// DashboardCountGrouped retrieves dashboard counts grouped by category with specified view type
func (s *ObjectContentService) DashboardCountGrouped(libraryId uuid.UUID, viewType string) (*models.GroupedDashboardResponse, error) {
	// Validate viewType
	if viewType == "" {
		viewType = "list"
	}
	if viewType != "list" && viewType != "cards" {
		return nil, fmt.Errorf("invalid view type: must be 'list' or 'cards'")
	}

	categories, err := s.repo.DashboardCountGrouped(libraryId)
	if err != nil {
		return nil, err
	}

	return &models.GroupedDashboardResponse{
		Categories: categories,
		ViewType:   viewType,
	}, nil
}
