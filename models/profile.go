package models

import (
	"time"

	"github.com/google/uuid"
)

// Profile represents the Profile table in the database
type Profile struct {
	ProfileID          int        `json:"profileId" db:"ProfileID"`
	ProfileName        string     `json:"profileName" db:"ProfileName"`
	ProfileDescription *string    `json:"profileDescription,omitempty" db:"ProfileDescription"`
	PortalStartPageId  *uuid.UUID `json:"portalStartPageId,omitempty" db:"PortalStartPageId"`
	DateCreated        time.Time  `json:"dateCreated" db:"DateCreated"`
	CreatedBy          int        `json:"createdBy" db:"CreatedBy"`
	DateModified       time.Time  `json:"dateModified" db:"DateModified"`
	ModifiedBy         int        `json:"modifiedBy" db:"ModifiedBy"`
}

// CreateProfileRequest represents the request body for creating a new profile
type CreateProfileRequest struct {
	ProfileName        string     `json:"profileName" validate:"required"`
	ProfileDescription *string    `json:"profileDescription,omitempty"`
	PortalStartPageId  *uuid.UUID `json:"portalStartPageId,omitempty"`
	CreatedBy          int        `json:"createdBy" validate:"required"`
}

// UpdateProfileRequest represents the request body for updating a profile
type UpdateProfileRequest struct {
	ProfileName        *string    `json:"profileName,omitempty"`
	ProfileDescription *string    `json:"profileDescription,omitempty"`
	PortalStartPageId  *uuid.UUID `json:"portalStartPageId,omitempty"`
	ModifiedBy         int        `json:"modifiedBy" validate:"required"`
}

