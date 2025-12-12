package services

import (
	"enterprise-architect-api/models"
	"enterprise-architect-api/repositories"
	"fmt"
	"math"

	"github.com/google/uuid"
)

// ObjectService handles business logic for objects
type ObjectService struct {
	repo *repositories.ObjectRepository
}

func (s *ObjectService) GetObjectsByObjectTypeIDAndLibraryID(objectTypeID int, libraryID uuid.UUID, page int, pageSize int) ([]models.Object, int, error) {
	return s.repo.GetByObjectTypeIDAndLibraryID(objectTypeID, libraryID, page, pageSize)
}

// NewObjectService creates a new ObjectService
func NewObjectService(repo *repositories.ObjectRepository) *ObjectService {
	return &ObjectService{repo: repo}
}

// CreateObject creates a new object
func (s *ObjectService) CreateObject(req models.CreateObjectRequest) (*models.Object, error) {
	// Validate required fields
	if req.ObjectName == "" {
		return nil, fmt.Errorf("object name is required")
	}
	if req.ObjectTypeID == 0 {
		return nil, fmt.Errorf("object type ID is required")
	}
	if req.ExactObjectTypeID == 0 {
		return nil, fmt.Errorf("exact object type ID is required")
	}
	if req.CreatedBy == 0 {
		return nil, fmt.Errorf("created by is required")
	}

	return s.repo.Create(req)
}

// GetObjectByID retrieves an object by its ID
func (s *ObjectService) GetObjectByID(id uuid.UUID) (*models.Object, error) {
	return s.repo.GetByID(id)
}

// GetAllObjects retrieves all objects with pagination
func (s *ObjectService) GetAllObjects(page, pageSize int) (*models.PaginatedResponse, error) {
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

	objects, totalCount, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data:       objects,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// UpdateObject updates an existing object
func (s *ObjectService) UpdateObject(id uuid.UUID, req models.UpdateObjectRequest) (*models.Object, error) {
	// Validate that at least one field is being updated
	if req.ObjectName == nil && req.ObjectDescription == nil && req.ObjectTypeID == nil &&
		req.ExactObjectTypeID == nil && req.RichTextDescription == nil && req.IsLibrary == nil &&
		req.FileExtension == nil && req.Prefix == nil && req.Suffix == nil {
		return nil, fmt.Errorf("at least one field must be provided for update")
	}

	if req.ModifiedBy == 0 {
		return nil, fmt.Errorf("modified by is required")
	}

	return s.repo.Update(id, req)
}

// DeleteObject deletes an object by its ID
func (s *ObjectService) DeleteObject(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// GetLibraries retrieves all objects where IsLibrary is true
func (s *ObjectService) GetLibraries(page, pageSize int) (*models.PaginatedResponse, error) {
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

	objects, totalCount, err := s.repo.GetLibraries(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data:       objects,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// GetObjectsByTypeID retrieves all objects by ObjectTypeID
func (s *ObjectService) GetObjectsByTypeID(objectTypeID, page, pageSize int) (*models.PaginatedResponse, error) {
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

	objects, totalCount, err := s.repo.GetByObjectTypeID(objectTypeID, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data:       objects,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (s *ObjectService) GetHierarchyFolder(ObjectID uuid.UUID, profileID int, isFolder bool) ([]models.ObjectTree, error) {
	return s.repo.GetHierarchyFolderV2(ObjectID, profileID, isFolder)
}
func (s *ObjectService) ImportObjects(req models.ObjectImportRequest) (*models.ObjectImportResponse, error) {
	return s.repo.ImportObjects(req)
}
