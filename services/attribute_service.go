package services

import (
	"enterprise-architect-api/models"
	"enterprise-architect-api/repositories"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type AttributeService struct {
	attributeRepository *repositories.AttributeRepository
}

func NewAttributeService(attributeRepository *repositories.AttributeRepository) *AttributeService {
	return &AttributeService{attributeRepository: attributeRepository}
}

func (as *AttributeService) GetAttributeForObject(objectID uuid.UUID) ([]models.AssignedAttribute, error) {
	return as.attributeRepository.GetAttributeForObject(objectID)
}

// CreateAttribute creates a new attribute
func (as *AttributeService) CreateAttribute(attribute *models.Attribute) error {
	// Validate required fields
	if attribute.AttributeName == "" {
		return fmt.Errorf("attribute name is required and must be unique")
	}
	if attribute.AttributeType == "" {
		return fmt.Errorf("attribute type is required")
	}

	// Check if attribute name already exists
	exists, err := as.attributeRepository.ExistsByName(attribute.AttributeName)
	if err != nil {
		return fmt.Errorf("error checking attribute name uniqueness: %w", err)
	}
	if exists {
		return fmt.Errorf("attribute name '%s' already exists, name must be unique", attribute.AttributeName)
	}

	return as.attributeRepository.Create(attribute)
}

// GetAttributeByID retrieves an attribute by its ID
func (as *AttributeService) GetAttributeByID(id string) (*models.Attribute, error) {
	return as.attributeRepository.GetByID(id)
}

// GetAllAttributes retrieves all attributes with pagination
func (as *AttributeService) GetAllAttributes(page, pageSize int) (*models.PaginatedResponse, error) {
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

	attributes, totalCount, err := as.attributeRepository.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data:       attributes,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// UpdateAttribute updates an existing attribute
func (as *AttributeService) UpdateAttribute(id string, attribute *models.Attribute) (*models.Attribute, error) {
	// Validate that attribute exists
	_, err := as.attributeRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if attribute.AttributeName == "" {
		return nil, fmt.Errorf("attribute name is required")
	}
	if attribute.AttributeType == "" {
		return nil, fmt.Errorf("attribute type is required")
	}

	err = as.attributeRepository.Update(id, attribute)
	if err != nil {
		return nil, err
	}

	return as.attributeRepository.GetByID(id)
}

// DeleteAttribute deletes an attribute by its ID
func (as *AttributeService) DeleteAttribute(id string) error {
	return as.attributeRepository.Delete(id)
}

// AssignAttributeToObjectType assigns an attribute to an object type
func (as *AttributeService) AssignAttributeToObjectType(req *models.AssignAttributeToObjectTypeRequest) error {
	// Validate required fields
	if req.AttributeGroupName == "" {
		return fmt.Errorf("attribute group name is required")
	}
	if req.AttributeId.String() == "00000000-0000-0000-0000-000000000000" {
		return fmt.Errorf("attribute ID is required")
	}
	if req.ObjectTypeId <= 0 && req.RelationTypeId.String() == "00000000-0000-0000-0000-000000000000" {
		return fmt.Errorf("either object type ID or relation type ID must be provided")
	}

	return as.attributeRepository.AssignAttributeToObjectType(req)
}
