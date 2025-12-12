package services

import (
	"enterprise-architect-api/models"
	"enterprise-architect-api/repositories"
	"fmt"
	"math"

	"github.com/google/uuid"
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

// GetBaseLibrary retrieves the base library of object types

// AddFolderToTree adds a new folder to the folder hierarchy tree
func (s *ObjectTypeService) AddFolderToTree(req models.AddFolderToTreeRequest) (*uuid.UUID, error) {
	// Validate: if FolderObjectTypeId is 0, ObjectTypeName must be provided
	if req.FolderObjectTypeId == 0 && req.ObjectTypeName == "" {
		return nil, fmt.Errorf("object type name is required when creating a new object type")
	}

	return s.repo.AddFolderToTree(req)
}

// DeleteObjectType deletes an object type by its ID
func (s *ObjectTypeService) DeleteObjectType(id int) error {
	return s.repo.Delete(id)
}

// AssignObjectTypeToFolder assigns an object type to a folder type setting
func (s *ObjectTypeService) AssignObjectTypeToFolder(req models.FolderObjectTypes) error {
	// Validate required fields
	if req.FolderObjectTypeId == 0 {
		return fmt.Errorf("folder object type ID is required")
	}
	if req.ObjectTypeID == 0 {
		return fmt.Errorf("object type ID is required")
	}

	return s.repo.AssignObjectTypeToFolder(req)
}

// GetAvailableTypesForFolder retrieves available object types for a specific folder
func (s *ObjectTypeService) GetAvailableTypesForFolder(folderObjectTypeId int) ([]models.FolderObjectTypesNames, error) {
	if folderObjectTypeId == 0 {
		return nil, fmt.Errorf("folder object type ID is required")
	}

	return s.repo.GetAvailableTypesForFolder(folderObjectTypeId)
}

func (s *ObjectTypeService) GetAvailableTypesForLibsAndFolder(folderObjectTypeId int) ([]models.FolderObjectTypesNames, error) {
	if folderObjectTypeId == 0 {
		return nil, fmt.Errorf("folder object type ID is required")
	}

	return s.repo.GetAvailableTypesForLibsAndFolder(folderObjectTypeId)
}

// DeleteObjectTypeFromFolder removes an object type assignment from a folder
func (s *ObjectTypeService) DeleteObjectTypeFromFolder(folderObjectTypeId, objectTypeId int) error {
	if folderObjectTypeId == 0 {
		return fmt.Errorf("folder object type ID is required")
	}
	if objectTypeId == 0 {
		return fmt.Errorf("object type ID is required")
	}

	return s.repo.DeleteObjectTypeFromFolder(folderObjectTypeId, objectTypeId)
}
func (s *ObjectTypeService) GetBaseLibrary() ([]models.ObjectTypeHierarchy, error) {
	return s.repo.GetBaseLibrary()
}
