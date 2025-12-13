package services

import (
	"enterprise-architect-api/models"
	"enterprise-architect-api/repositories"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type RelationService struct {
	repo *repositories.RelationRepository
}

func NewRelationService(repo *repositories.RelationRepository) *RelationService {
	return &RelationService{repo: repo}
}

func (s *RelationService) CreateRelation(req models.CreateRelationRequest) (*models.Relation, error) {
	return s.repo.Create(req)
}

func (s *RelationService) GetRelationsByObjectID(objectIDStr string) ([]models.RelationWithDetails, error) {
	objectID, err := uuid.Parse(objectIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid object id: %w", err)
	}
	return s.repo.GetByObjectID(objectID)
}

func (s *RelationService) GetAllRelationTypes(page, pageSize int) (*models.PaginatedResponse, error) {
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

	types, totalCount, err := s.repo.GetAllRelationTypes(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data:       types,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}
