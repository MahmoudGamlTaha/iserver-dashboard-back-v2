package services

import (
	"enterprise-architect-api/models"
	"enterprise-architect-api/repositories"
	"fmt"
	"math"
)

// ObjectTypeService handles business logic for object types
type ObjectTypeService struct {
	repo *repositories.ObjectTypeRepository
}

// NewObjectTypeService creates a new ObjectTypeService
func NewObjectTypeService(repo *repositories.ObjectTypeRepository) *ObjectTypeService {
	return &ObjectTypeService{repo: repo}
}

// CreateObjectType creates a new object type
func (s *ObjectTypeService) CreateObjectType(req models.CreateObjectTypeRequest) (*models.ObjectType, error) {
	// Validate required fields
	if req.CreatedBy == 0 {
		return nil, fmt.Errorf("created by is required")
	}

	return s.repo.Create(req)
}

// GetObjectTypeByID retrieves an object type by its ID
func (s *ObjectTypeService) GetObjectTypeByID(id int) (*models.ObjectType, error) {
	return s.repo.GetByID(id)
}

// GetAllObjectTypes retrieves all object types with pagination
func (s *ObjectTypeService) GetAllObjectTypes(page, pageSize int) (*models.PaginatedResponse, error) {
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

	objectTypes, totalCount, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data:       objectTypes,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// SearchObjectTypesByName retrieves object types filtered by name with pagination
func (s *ObjectTypeService) SearchObjectTypesByName(name string, page, pageSize int) (*models.PaginatedResponse, error) {
	// reuse same pagination defaults
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	objectTypes, totalCount, err := s.repo.SearchByName(name, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data:       objectTypes,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// UpdateObjectType updates an existing object type
func (s *ObjectTypeService) UpdateObjectType(id int, req models.UpdateObjectTypeRequest) (*models.ObjectType, error) {
	// Validate that at least one field is being updated
	if req.ObjectTypeName == nil && req.Description == nil && req.FileExtension == nil &&
		req.IsTemplateType == nil && req.ActiveType == nil {
		return nil, fmt.Errorf("at least one field must be provided for update")
	}

	if req.ModifiedBy == 0 {
		return nil, fmt.Errorf("modified by is required")
	}

	return s.repo.Update(id, req)
}

// GetFolderRepositoryTree retrieves the hierarchy of object types
func (s *ObjectTypeService) GetFolderRepositoryTree() ([]models.ObjectTypeHierarchy, error) {
	return s.repo.GetFolderRepositoryTree()
}

// DeleteObjectType deletes an object type by its ID
func (s *ObjectTypeService) DeleteObjectType(id int) error {
	return s.repo.Delete(id)
}
