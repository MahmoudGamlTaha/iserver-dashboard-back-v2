package services

import (
	"enterprise-architect-api/models"
	"enterprise-architect-api/repositories"
	"fmt"
	"math"
)

// ProfileService handles business logic for profiles
type ProfileService struct {
	repo *repositories.ProfileRepository
}

// NewProfileService creates a new ProfileService
func NewProfileService(repo *repositories.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

// CreateProfile creates a new profile
func (s *ProfileService) CreateProfile(req models.CreateProfileRequest) (*models.Profile, error) {
	// Validate required fields
	if req.ProfileName == "" {
		return nil, fmt.Errorf("profile name is required")
	}
	if req.CreatedBy == 0 {
		return nil, fmt.Errorf("created by is required")
	}

	return s.repo.Create(req)
}

// GetProfileByID retrieves a profile by its ID
func (s *ProfileService) GetProfileByID(id int) (*models.Profile, error) {
	return s.repo.GetByID(id)
}

// GetAllProfiles retrieves all profiles with pagination
func (s *ProfileService) GetAllProfiles(page, pageSize int) (*models.PaginatedResponse, error) {
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

	profiles, totalCount, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data:       profiles,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// UpdateProfile updates an existing profile
func (s *ProfileService) UpdateProfile(id int, req models.UpdateProfileRequest) (*models.Profile, error) {
	// Validate that at least one field is being updated
	if req.ProfileName == nil && req.ProfileDescription == nil && req.PortalStartPageId == nil {
		return nil, fmt.Errorf("at least one field must be provided for update")
	}

	if req.ModifiedBy == 0 {
		return nil, fmt.Errorf("modified by is required")
	}

	return s.repo.Update(id, req)
}

// DeleteProfile deletes a profile by its ID
func (s *ProfileService) DeleteProfile(id int) error {
	return s.repo.Delete(id)
}

