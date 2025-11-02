package models

import (
	"time"

	"github.com/google/uuid"
)

// ObjectContent represents the ObjectContents table in the database
type ObjectContent struct {
	ID                              int        `json:"id" db:"ID"`
	DocumentObjectID                uuid.UUID  `json:"documentObjectId" db:"DocumentObjectID"`
	ContainerVersionID              uuid.UUID  `json:"containerVersionId" db:"ContainerVersionID"`
	ObjectID                        uuid.UUID  `json:"objectId" db:"ObjectID"`
	Instances                       int        `json:"instances" db:"Instances"`
	IsShortCut                      *bool      `json:"isShortCut,omitempty" db:"IsShortCut"`
	ShapeSheetKeysRequiringUpdateId *uuid.UUID `json:"shapeSheetKeysRequiringUpdateId,omitempty" db:"ShapeSheetKeysRequiringUpdateId"`
	ContainmentType                 int        `json:"containmentType" db:"ContainmentType"`
	DateCreated                     time.Time  `json:"dateCreated" db:"DateCreated"`
	CreatedBy                       int        `json:"createdBy" db:"CreatedBy"`
	DateModified                    time.Time  `json:"dateModified" db:"DateModified"`
	ModifiedBy                      int        `json:"modifiedBy" db:"ModifiedBy"`
}

// CreateObjectContentRequest represents the request body for creating a new object content
type CreateObjectContentRequest struct {
	DocumentObjectID   uuid.UUID `json:"documentObjectId" validate:"required"`
	ContainerVersionID uuid.UUID `json:"containerVersionId" validate:"required"`
	ObjectID           uuid.UUID `json:"objectId" validate:"required"`
	Instances          int       `json:"instances" validate:"required"`
	IsShortCut         *bool     `json:"isShortCut,omitempty"`
	ContainmentType    int       `json:"containmentType" validate:"required"`
	CreatedBy          int       `json:"createdBy" validate:"required"`
}

// UpdateObjectContentRequest represents the request body for updating an object content
type UpdateObjectContentRequest struct {
	DocumentObjectID   *uuid.UUID `json:"documentObjectId,omitempty"`
	ContainerVersionID *uuid.UUID `json:"containerVersionId,omitempty"`
	ObjectID           *uuid.UUID `json:"objectId,omitempty"`
	Instances          *int       `json:"instances,omitempty"`
	IsShortCut         *bool      `json:"isShortCut,omitempty"`
	ContainmentType    *int       `json:"containmentType,omitempty"`
	ModifiedBy         int        `json:"modifiedBy" validate:"required"`
}

type DashboardCount struct {
	ExactObjectTypeID *int64  `json:"ExactObjectTypeID" db:"ExactObjectTypeId"`
	ObjectTypeName    *string `json:"ObjectTypeName" db:"ObjectTypeName"`
	Count             *int64  `json:"count" db:"count"`
	Color             *int64  `json:"color" db:"color"`
	Icon              *int64  `json:"icon" db:"icon"`
}
