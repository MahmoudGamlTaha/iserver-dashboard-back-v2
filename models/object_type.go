package models

import (
	"time"

	"github.com/google/uuid"
)

// ObjectType represents the ObjectType table in the database
type ObjectType struct {
	ObjectTypeID                   int        `json:"objectTypeId" db:"ObjectTypeID"`
	ObjectTypeName                 *string    `json:"objectTypeName,omitempty" db:"ObjectTypeName"`
	ObjectTypeImage                []byte     `json:"objectTypeImage,omitempty" db:"ObjectTypeImage"`
	IsTemplateType                 bool       `json:"isTemplateType" db:"IsTemplateType"`
	GeneralType                    *int       `json:"generalType,omitempty" db:"GeneralType"`
	TemplateFileName               *string    `json:"templateFileName,omitempty" db:"TemplateFileName"`
	IsDefaultTemplate              bool       `json:"isDefaultTemplate" db:"IsDefaultTemplate"`
	ActiveType                     bool       `json:"activeType" db:"ActiveType"`
	EnforceUniqueNaming            bool       `json:"enforceUniqueNaming" db:"EnforceUniqueNaming"`
	CanHaveVisioAlias              bool       `json:"canHaveVisioAlias" db:"CanHaveVisioAlias"`
	IsConnector                    bool       `json:"isConnector" db:"IsConnector"`
	ImplicitlyAddObjectTypes       bool       `json:"implicitlyAddObjectTypes" db:"ImplicitlyAddObjectTypes"`
	CommitOverlapRelationships     bool       `json:"commitOverlapRelationships" db:"CommitOverlapRelationships"`
	DateCreated                    time.Time  `json:"dateCreated" db:"DateCreated"`
	CreatedBy                      int        `json:"createdBy" db:"CreatedBy"`
	DateModified                   time.Time  `json:"dateModified" db:"DateModified"`
	ModifiedBy                     int        `json:"modifiedBy" db:"ModifiedBy"`
	FileExtension                  *string    `json:"fileExtension,omitempty" db:"FileExtension"`
	HandlerToolId                  *uuid.UUID `json:"handlerToolId,omitempty" db:"HandlerToolId"`
	Color                          *int       `json:"color,omitempty" db:"Color"`
	Icon                           *int       `json:"icon,omitempty" db:"Icon"`
	IsExcludedFromBrokenConnectors bool       `json:"isExcludedFromBrokenConnectors" db:"IsExcludedFromBrokenConnectors"`
	Description                    *string    `json:"description,omitempty" db:"Description"`
	ExportShapeAttributes          bool       `json:"exportShapeAttributes" db:"ExportShapeAttributes"`
	ExportShapeSystemProperties    bool       `json:"exportShapeSystemProperties" db:"ExportShapeSystemProperties"`
	ImportShapeAttributes          bool       `json:"importShapeAttributes" db:"ImportShapeAttributes"`
	ExportDocumentAttributes       bool       `json:"exportDocumentAttributes" db:"ExportDocumentAttributes"`
	ExportDocumentSystemProperties bool       `json:"exportDocumentSystemProperties" db:"ExportDocumentSystemProperties"`
	DeleteNotSyncVisioShapeData    bool       `json:"deleteNotSyncVisioShapeData" db:"DeleteNotSyncVisioShapeData"`
	DeleteIfHasNoMaster            bool       `json:"deleteIfHasNoMaster" db:"DeleteIfHasNoMaster"`
}
type ObjectTypeHierarchy struct {
	ObjectTypeFolderId    int        `json:"objectTypeFolderId" db:"ObjectTypeFolderId"`
	ObjectTypeId          int        `json:"objectTypeId" db:"ObjectTypeId"`
	ObjectTypeName        *string    `json:"objectTypeName,omitempty" db:"ObjectTypeName"`
	ObjectTypeHierarchyId *uuid.UUID `json:"objectTypeHierarchyId" db:"objectTypeHierarchyId"`
	ObjectTypeParentId    *uuid.UUID `json:"objectTypeParentId" db:"objectTypeParentId"`
	Level                 *int       `json:"Level" db:"Level"`
	FullPath              *string    `json:"FullPath" db:"FullPath"`
}

// CreateObjectTypeRequest represents the request body for creating a new object type
type CreateObjectTypeRequest struct {
	ObjectTypeName *string   `json:"objectTypeName,omitempty"`
	Description    *string   `json:"description,omitempty"`
	FileExtension  *string   `json:"fileExtension,omitempty"`
	IsTemplateType bool      `json:"isTemplateType"`
	ActiveType     bool      `json:"activeType"`
	CreatedBy      int       `json:"createdBy" validate:"required"`
	DateCreated    time.Time `json:"dateCreated"`
	DateModified   time.Time `json:"dateModified"`
	ModifiedBy     int       `json:"modifiedBy" validate:"required"`
}

// UpdateObjectTypeRequest represents the request body for updating an object type
type UpdateObjectTypeRequest struct {
	ObjectTypeName *string `json:"objectTypeName,omitempty"`
	Description    *string `json:"description,omitempty"`
	FileExtension  *string `json:"fileExtension,omitempty"`
	IsTemplateType *bool   `json:"isTemplateType,omitempty"`
	ActiveType     *bool   `json:"activeType,omitempty"`
	ModifiedBy     int     `json:"modifiedBy" validate:"required"`
}

// AddFolderToTreeRequest represents the request body for adding a folder to the hierarchy tree
type AddFolderToTreeRequest struct {
	FolderObjectTypeId int        `json:"folderObjectTypeId"`
	ObjectTypeName     string     `json:"objectTypeName,omitempty"`
	ParentHierarchyId  *uuid.UUID `json:"parentHierarchyId,omitempty"`
}

type FolderObjectTypes struct {
	ObjectTypeID       int  `json:"objectTypeId"`
	FolderObjectTypeId int  `json:"folderObjectTypeId"`
	IsDocumentType     bool `json:"isDocumentType"`
}
type FolderObjectTypesNames struct {
	ObjectTypeID       int    `json:"objectTypeId"`
	FolderObjectTypeId int    `json:"folderObjectTypeId"`
	IsDocumentType     bool   `json:"isDocumentType"`
	ObjectTypeName     string `json:"objectTypeName"`
}
