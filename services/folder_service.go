package services

import (
	"enterprise-architect-api/models"
	"enterprise-architect-api/repositories"
	"fmt"

	"github.com/google/uuid"
)

// FolderService handles business logic for folders
type FolderService struct {
	repo *repositories.FolderRepository
}

// NewFolderService creates a new FolderService
func NewFolderService(repo *repositories.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

// GetObjectTypeFolders retrieves folders and system repositories by library ID
func (s *FolderService) GetObjectTypeFolders(libraryID uuid.UUID) ([]models.ObjectTypeFolder, error) {
	folders, err := s.repo.GetObjectTypeFolders(libraryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object type folders: %w", err)
	}
	return folders, nil
}

// GetFoldersByLibrary retrieves folder contents by folder ID and profile ID
func (s *FolderService) GetFoldersByLibrary(folderID uuid.UUID, profileID int) ([]models.FolderContent, error) {
	if profileID == 0 {
		return nil, fmt.Errorf("profile ID is required")
	}

	contents, err := s.repo.GetFoldersByLibrary(folderID, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders by library: %w", err)
	}
	return contents, nil
}

