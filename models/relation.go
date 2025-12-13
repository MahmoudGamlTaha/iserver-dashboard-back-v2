package models

import (
	"time"

	"github.com/google/uuid"
)

// Relation represents the Relationship table in the database
type Relation struct {
	RelationshipID      uuid.UUID `json:"relationshipId" db:"RelationshipId"`
	RelationTypeID      uuid.UUID `json:"relationTypeId" db:"RelationTypeId"`
	RelationReason      *string   `json:"relationReason" db:"RelationReason"`
	FromObjectID        uuid.UUID `json:"fromObjectId" db:"FromObjectId"`
	ToObjectID          uuid.UUID `json:"toObjectId" db:"ToObjectId"`
	DateCreated         time.Time `json:"dateCreated" db:"DateCreated"`
	CreatedBy           int       `json:"createdBy" db:"CreatedBy"`
	DateModified        time.Time `json:"dateModified" db:"DateModified"`
	ModifiedBy          int       `json:"modifiedBy" db:"ModifiedBy"`
	RichTextDescription *string   `json:"richTextDescription" db:"RichTextDescription"`
}

// RelationType represents the RelationType table in the database
type RelationType struct {
	RelationTypeID          uuid.UUID `json:"relationTypeId" db:"RelationTypeId"`
	RelationTypeName        string    `json:"relationTypeName" db:"RelationTypeName"`
	RelationTypeDescription string    `json:"relationTypeDescription" db:"RelationTypeDescription"`
	FromToDescription       string    `json:"fromToDescription" db:"FromToDescription"`
	ToFromDescription       string    `json:"toFromDescription" db:"ToFromDescription"`
	IsDirectionless         bool      `json:"isDirectionless" db:"isDirectionless"`
	DateCreated             time.Time `json:"dateCreated" db:"DateCreated"`
	CreatedBy               int       `json:"createdBy" db:"CreatedBy"`
	DateModified            time.Time `json:"dateModified" db:"DateModified"`
	ModifiedBy              int       `json:"modifiedBy" db:"ModifiedBy"`
	IsHierarchical          bool      `json:"isHierarchical" db:"isHierarchical"`
	Usage                   string    `json:"usage" db:"usage"`
}

// CreateRelationRequest represents the request body for creating a new relationship
type CreateRelationRequest struct {
	RelationTypeId      uuid.UUID `json:"relationTypeId" validate:"required"`
	RelationReason      *string   `json:"relationReason"`
	FromObjectId        uuid.UUID `json:"fromObjectId" validate:"required"`
	ToObjectID          uuid.UUID `json:"toObjectId" validate:"required"`
	CreatedBy           int       `json:"createdBy" validate:"required"`
	RichTextDescription *string   `json:"richTextDescription"`
}

// RelationWithDetails extends Relation to include details about the related object and type
type RelationWithDetails struct {
	Relation
	RelationTypeName string    `json:"relationTypeName"`
	OtherObjectID    uuid.UUID `json:"otherObjectId"`
	OtherObjectName  string    `json:"otherObjectName"`
	Direction        string    `json:"direction"` // "From" or "To"
}
