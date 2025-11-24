package models

import (
	"time"

	"github.com/google/uuid"
)

// ObjectTree represents a hierarchical object structure with permissions
type ObjectTree struct {
	// Core object fields from recurse CTE
	ObjectID                 uuid.UUID  `json:"objectId" db:"ObjectID"`
	ObjectParentID           *uuid.UUID `json:"objectParentId,omitempty" db:"ObjectParentID"`
	ObjectName               string     `json:"objectName" db:"ObjectName"`
	ObjectDescription        string     `json:"objectDescription" db:"ObjectDescription"`
	CurrentVersionId         *uuid.UUID `json:"currentVersionId,omitempty" db:"CurrentVersionId"`
	CheckedInVersionId       *uuid.UUID `json:"checkedInVersionId,omitempty" db:"CheckedInVersionId"`
	IsImported               bool       `json:"isImported" db:"IsImported"`
	IsLibrary                bool       `json:"isLibrary" db:"IsLibrary"`
	LibraryId                *uuid.UUID `json:"libraryId,omitempty" db:"LibraryId"`
	FileExtension            *string    `json:"fileExtension,omitempty" db:"FileExtension"`
	GeneralType              *int       `json:"generalType,omitempty" db:"GeneralType"`
	GeneralTypeName          *string    `json:"generalTypeName,omitempty" db:"GeneralTypeName"`
	TypeId                   *int       `json:"typeId,omitempty" db:"TypeId"`
	TypeName                 *string    `json:"typeName,omitempty" db:"TypeName"`
	IsDeleted                bool       `json:"isDeleted" db:"IsDeleted"`
	VisioAlias               *string    `json:"visioAlias,omitempty" db:"VisioAlias"`
	HasVisioAlias            *bool      `json:"hasVisioAlias,omitempty" db:"HasVisioAlias"`
	IsLocked                 bool       `json:"isLocked" db:"IsLocked"`
	SortOrder                *int       `json:"sortOrder,omitempty" db:"SortOrder"`
	AutoSort                 *bool      `json:"autoSort,omitempty" db:"AutoSort"`
	Prefix                   *string    `json:"prefix,omitempty" db:"Prefix"`
	Suffix                   *string    `json:"suffix,omitempty" db:"Suffix"`
	ProvenanceID             *uuid.UUID `json:"provenanceId,omitempty" db:"ProvenanceID"`
	ProvenanceVersionID      *uuid.UUID `json:"provenanceVersionId,omitempty" db:"ProvenanceVersionID"`
	CreatedBy                int        `json:"createdBy" db:"CreatedBy"`
	DateCreated              time.Time  `json:"dateCreated" db:"DateCreated"`
	ModifiedBy               int        `json:"modifiedBy" db:"ModifiedBy"`
	DateModified             time.Time  `json:"dateModified" db:"DateModified"`
	IsCheckedOut             bool       `json:"isCheckedOut" db:"IsCheckedOut"`
	CheckedOutUserId         *int       `json:"checkedOutUserId,omitempty" db:"CheckedOutUserId"`
	RichTextDescription      string     `json:"richTextDescription" db:"RichTextDescription"`
	IsPendingApproval        bool       `json:"isPendingApproval" db:"IsPendingApproval"`
	CheckedOutBy             *int       `json:"checkedOutBy,omitempty" db:"CheckedOutBy"`
	IsFirstVersionCheckedOut bool       `json:"isFirstVersionCheckedOut" db:"IsFirstVersionCheckedOut"`
	FolderId                 *uuid.UUID `json:"folderId,omitempty" db:"FolderId"`
	IsFolder                 bool       `json:"isFolder" db:"isFolder"`

	// Additional fields from joins
	CheckedInName                    *string `json:"checkedInName,omitempty" db:"CheckedInName"`
	CheckedInVisioAlias              *string `json:"checkedInVisioAlias,omitempty" db:"CheckedInVisioAlias"`
	CheckedInHasVisioAlias           *bool   `json:"checkedInHasVisioAlias,omitempty" db:"CheckedInHasVisioAlias"`
	HasReadPermission                bool    `json:"hasReadPermission" db:"HasReadPermission"`
	HasModifyContentsPermission      bool    `json:"hasModifyContentsPermission" db:"HasModifyContentsPermission"`
	HasDeletePermission              bool    `json:"hasDeletePermission" db:"HasDeletePermission"`
	HasModifyPermission              bool    `json:"hasModifyPermission" db:"HasModifyPermission"`
	HasModifyRelationshipsPermission bool    `json:"hasModifyRelationshipsPermission" db:"HasModifyRelationshipsPermission"`
}
