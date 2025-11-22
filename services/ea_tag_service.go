package services

import (
	"enterprise-architect-api/models"
	"enterprise-architect-api/repositories"
	"fmt"
	"math"
)

// EATagService handles business logic for EA tags
type EATagService struct {
	repo *repositories.ReportConfigRepository
}

func (s *EATagService) GetEAObjectTypesAssignedToDimension(objectTypeId int64) (any, error) {
	return s.repo.GetEAObjectTypesAssignedToDimension(objectTypeId)
}

// NewEATagService creates a new EATagService
func NewEATagService(repo *repositories.ReportConfigRepository) *EATagService {
	return &EATagService{repo: repo}
}

// CreateEATag creates a new EA tag
func (s *EATagService) CreateEATag(req models.CreateEATagRequest) (*models.EATag, error) {
	// Validate required fields
	if req.NameAr == "" {
		return nil, fmt.Errorf("name_ar is required")
	}
	if req.NameEn == "" {
		return nil, fmt.Errorf("name_en is required")
	}

	return s.repo.CreateEATag(req)
}

// GetEATagByID retrieves an EA tag by its ID
func (s *EATagService) GetEATagByID(id int) (*models.EATag, error) {
	return s.repo.GetEATagByID(id)
}

// GetAllEATags retrieves all EA tags with pagination
func (s *EATagService) GetAllEATags(page, pageSize int) (*models.PaginatedResponse, error) {
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

	tags, totalCount, err := s.repo.GetAllEATags(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data:       tags,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// UpdateEATag updates an existing EA tag
func (s *EATagService) UpdateEATag(id int, req models.UpdateEATagRequest) (*models.EATag, error) {
	// Validate that at least one field is being updated
	if req.NameAr == nil && req.NameEn == nil {
		return nil, fmt.Errorf("at least one field must be provided for update")
	}

	return s.repo.UpdateEATag(id, req)
}

// DeleteEATag deletes an EA tag by its ID
func (s *EATagService) DeleteEATag(id int) error {
	return s.repo.DeleteEATag(id)
}

// AssignObjectTypeToDimention assigns an object type to a dimension
func (s *EATagService) AssignObjectTypeToDimention(req models.AssignObjectTypeToDimentionRequest) (*models.EATagDimention, error) {
	// Validate required fields
	if req.ObjectTypeID <= 0 {
		return nil, fmt.Errorf("object_type_id is required and must be greater than 0")
	}
	if req.EAID <= 0 {
		return nil, fmt.Errorf("ea_tag_id is required and must be greater than 0")
	}

	return s.repo.AssignObjectTypeToDimention(req)
}
