package models

import (
	"time"

	"github.com/google/uuid"
)

// ObjectTypeFolder represents a folder or system repository object
type ObjectTypeFolder struct {
	ObjectID        uuid.UUID  `json:"objectId" db:"ObjectID"`
	GeneralType     *int       `json:"generalType,omitempty" db:"GeneralType"`
	ObjectName      string     `json:"objectName" db:"objectName"`
	SortOrder       *int       `json:"sortOrder,omitempty" db:"sortorder"`
	GeneralTypeName *string    `json:"generalTypeName,omitempty" db:"GeneralTypeName"`
	TypeId          int        `json:"typeId" db:"TypeId"`
	TypeName        *string    `json:"typeName,omitempty" db:"TypeName"`
	IsObjectDeleted *bool      `json:"isObjectDeleted,omitempty" db:"IsObjectDeleted"`
	IsTemplateType  int        `json:"isTemplateType" db:"IsTemplateType"`
	ObjectVersion   *uuid.UUID `json:"objectVersion,omitempty" db:"ObjectVersion"`
	LibraryId       *uuid.UUID `json:"libraryId,omitempty" db:"LibraryId"`
	Color           *string    `json:"color,omitempty" db:"Color"`
	Icon            *string    `json:"icon,omitempty" db:"Icon"`
}

// FolderContent represents the content of a folder with permissions
type FolderContent struct {
	// Object fields
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

	// Additional fields
	IsPendingApproval                bool    `json:"isPendingApproval" db:"IsPendingApproval"`
	CheckedInName                    *string `json:"checkedInName,omitempty" db:"CheckedInName"`
	CheckedInVisioAlias              *string `json:"checkedInVisioAlias,omitempty" db:"CheckedInVisioAlias"`
	CheckedInHasVisioAlias           *bool   `json:"checkedInHasVisioAlias,omitempty" db:"CheckedInHasVisioAlias"`
	HasReadPermission                bool    `json:"hasReadPermission" db:"HasReadPermission"`
	HasModifyContentsPermission      bool    `json:"hasModifyContentsPermission" db:"HasModifyContentsPermission"`
	HasDeletePermission              bool    `json:"hasDeletePermission" db:"HasDeletePermission"`
	HasModifyPermission              bool    `json:"hasModifyPermission" db:"HasModifyPermission"`
	HasModifyRelationshipsPermission bool    `json:"hasModifyRelationshipsPermission" db:"HasModifyRelationshipsPermission"`
	CheckedOutBy                     *int    `json:"checkedOutBy,omitempty" db:"CheckedOutBy"`
	IsFirstVersionCheckedOut         bool    `json:"isFirstVersionCheckedOut" db:"IsFirstVersionCheckedOut"`
}
