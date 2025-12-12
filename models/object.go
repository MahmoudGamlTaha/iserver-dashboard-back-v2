package models

import (
	"time"

	"github.com/google/uuid"
)

// Object represents the Object table in the database
type Object struct {
	ObjectID                 uuid.UUID  `json:"objectId" db:"ObjectID"`
	ObjectName               string     `json:"objectName" db:"ObjectName"`
	ObjectDescription        string     `json:"objectDescription" db:"ObjectDescription"`
	ObjectTypeID             int        `json:"objectTypeId" db:"ObjectTypeID"`
	CheckedInVersionId       *uuid.UUID `json:"checkedInVersionId,omitempty" db:"CheckedInVersionId"`
	DeleteFlag               *bool      `json:"deleteFlag,omitempty" db:"DeleteFlag"`
	Locked                   bool       `json:"locked" db:"Locked"`
	RequiresShapeSheetUpdate *bool      `json:"requiresShapeSheetUpdate,omitempty" db:"RequiresShapeSheetUpdate"`
	TemplateID               *int       `json:"templateId,omitempty" db:"TemplateID"`
	IsImported               bool       `json:"isImported" db:"IsImported"`
	IsLibrary                bool       `json:"isLibrary" db:"IsLibrary"`
	LibraryId                *uuid.UUID `json:"libraryId,omitempty" db:"LibraryId"`
	FileExtension            *string    `json:"fileExtension,omitempty" db:"FileExtension"`
	SortOrder                *int       `json:"sortOrder,omitempty" db:"SortOrder"`
	Prefix                   *string    `json:"prefix,omitempty" db:"Prefix"`
	Suffix                   *string    `json:"suffix,omitempty" db:"Suffix"`
	ProvenanceId             *uuid.UUID `json:"provenanceId,omitempty" db:"ProvenanceId"`
	ProvenanceVersionId      *uuid.UUID `json:"provenanceVersionId,omitempty" db:"ProvenanceVersionId"`
	GeneralType              *int       `json:"generalType,omitempty" db:"GeneralType"`
	CurrentVersionId         *uuid.UUID `json:"currentVersionId,omitempty" db:"CurrentVersionId"`
	VisioAlias               *string    `json:"visioAlias,omitempty" db:"VisioAlias"`
	HasVisioAlias            *bool      `json:"hasVisioAlias,omitempty" db:"HasVisioAlias"`
	DateCreated              time.Time  `json:"dateCreated" db:"DateCreated"`
	CreatedBy                int        `json:"createdBy" db:"CreatedBy"`
	DateModified             time.Time  `json:"dateModified" db:"DateModified"`
	ModifiedBy               int        `json:"modifiedBy" db:"ModifiedBy"`
	IsCheckedOut             bool       `json:"isCheckedOut" db:"IsCheckedOut"`
	CheckedOutUserId         *int       `json:"checkedOutUserId,omitempty" db:"CheckedOutUserId"`
	DeleteTransactionId      *uuid.UUID `json:"deleteTransactionId,omitempty" db:"DeleteTransactionId"`
	NameChecksum             *int       `json:"nameChecksum,omitempty" db:"NameChecksum"`
	ExactObjectTypeID        int        `json:"exactObjectTypeId" db:"ExactObjectTypeID"`
	RichTextDescription      string     `json:"richTextDescription" db:"RichTextDescription"`
	AutoSort                 *bool      `json:"autoSort,omitempty" db:"AutoSort"`
}

// CreateObjectRequest represents the request body for creating a new object
type CreateObjectRequest struct {
	ObjectName          string               `json:"objectName" validate:"required"`
	ObjectDescription   string               `json:"objectDescription"`
	ObjectTypeID        int                  `json:"objectTypeId" validate:"required"`
	ExactObjectTypeID   int                  `json:"exactObjectTypeId" validate:"required"`
	RichTextDescription string               `json:"richTextDescription"`
	IsLibrary           bool                 `json:"isLibrary"`
	IsImported          bool                 `json:"isImported"`
	LibraryId           *uuid.UUID           `json:"libraryId,omitempty"`
	FileExtension       *string              `json:"fileExtension,omitempty"`
	Prefix              *string              `json:"prefix,omitempty"`
	Suffix              *string              `json:"suffix,omitempty"`
	CreatedBy           int                  `json:"createdBy" validate:"required"`
	Attributes          *[]AssignedAttribute `json:"attributes,omitempty"`
	GeneralType         *int                 `json:"generalType,omitempty"`
	DirectParentId      *uuid.UUID           `json:"directParentId,omitempty"`
}

// UpdateObjectRequest represents the request body for updating an object
type UpdateObjectRequest struct {
	ObjectName          *string `json:"objectName,omitempty"`
	ObjectDescription   *string `json:"objectDescription,omitempty"`
	ObjectTypeID        *int    `json:"objectTypeId,omitempty"`
	ExactObjectTypeID   *int    `json:"exactObjectTypeId,omitempty"`
	RichTextDescription *string `json:"richTextDescription,omitempty"`
	IsLibrary           *bool   `json:"isLibrary,omitempty"`
	FileExtension       *string `json:"fileExtension,omitempty"`
	Prefix              *string `json:"prefix,omitempty"`
	Suffix              *string `json:"suffix,omitempty"`
	ModifiedBy          int     `json:"modifiedBy" validate:"required"`
}

type ObjectImportRequest struct {
	LibraryId    uuid.UUID                    `json:"libraryId"`
	FolderId     uuid.UUID                    `json:"folderId"`
	ObjectTypeId int                          `json:"objectTypeId"`
	Data         []map[string]ObjectImportRow `json:"data"`
	Mappings     []interface{}                `json:"mappings"`
}
type ObjectImportResponse struct {
	General
	SuccessImportedObjectCount int `json:"successImportedObjectCount"`
	FailedImportObjectCount    int `json:"failedImportObjectCount"`
	TotalImportedObjectCount   int `json:"totalImportedObjectCount"`
}
type ObjectImportRow struct {
	AttributeId    *string `json:"attributeId"`
	AttributeValue *string `json:"value"`
	AttributeType  *string `json:"attributeType"`
	AttributeName  *string `json:"attributeName"`
}
