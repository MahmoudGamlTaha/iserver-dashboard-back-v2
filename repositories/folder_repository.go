package repositories

import (
	"database/sql"
	"enterprise-architect-api/models"
	"fmt"

	"github.com/google/uuid"
)

// FolderRepository handles database operations for folders
type FolderRepository struct {
	db *sql.DB
}

// NewFolderRepository creates a new FolderRepository
func NewFolderRepository(db *sql.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

// GetObjectTypeFolders retrieves folders and system repositories by library ID
func (r *FolderRepository) GetObjectTypeFolders(libraryID uuid.UUID) ([]models.ObjectTypeFolder, error) {
	query := `
		SELECT o.ObjectID,
			o.GeneralType AS [GeneralType], 
			o.objectName,
			o.sortorder,
			CASE o.GeneralType 
				WHEN dbo.const_GeneralType_Folder() THEN 'Folder' 
				WHEN dbo.const_GeneralType_SystemRepository() THEN 'System Repository' 
			END AS [GeneralTypeName],
			ot.ObjectTypeID AS [TypeId],
			ot.ObjectTypeName AS [TypeName],
			ot.Color AS [Color],
			ot.icon AS [Icon],
			o.DeleteFlag AS [IsObjectDeleted],
			0 AS IsTemplateType,
			o.CurrentVersionId AS ObjectVersion,
			o.LibraryId
		FROM dbo.Object AS o
		INNER JOIN dbo.ObjectType AS ot ON ot.ObjectTypeID = o.ObjectTypeID
		WHERE o.GeneralType IN (dbo.const_GeneralType_Folder(), dbo.const_GeneralType_SystemRepository()) 
			AND o.LibraryId = @p1
	`

	rows, err := r.db.Query(query, libraryID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving object type folders: %w", err)
	}
	defer rows.Close()

	var folders []models.ObjectTypeFolder
	for rows.Next() {
		var folder models.ObjectTypeFolder
		var objectIDBytes []byte
		var objectVersionBytes, libraryIDBytes []byte
		
		err := rows.Scan(
			&objectIDBytes,
			&folder.GeneralType,
			&folder.ObjectName,
			&folder.SortOrder,
			&folder.GeneralTypeName,
			&folder.TypeId,
			&folder.TypeName,
			&folder.Color,
			&folder.Icon,
			&folder.IsObjectDeleted,
			&folder.IsTemplateType,
			&objectVersionBytes,
			&libraryIDBytes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning folder: %w", err)
		}
		
		// Parse UUID from bytes
		folder.ObjectID, err = parseSQLServerUUID(objectIDBytes)
		if err != nil {
			return nil, fmt.Errorf("error parsing ObjectID: %w", err)
		}
		
		if objectVersionBytes != nil {
			parsedUUID, err := parseSQLServerUUID(objectVersionBytes)
			if err != nil {
				return nil, fmt.Errorf("error parsing ObjectVersion: %w", err)
			}
			folder.ObjectVersion = &parsedUUID
		}
		
		if libraryIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(libraryIDBytes)
			if err != nil {
				return nil, fmt.Errorf("error parsing LibraryId: %w", err)
			}
			folder.LibraryId = &parsedUUID
		}
		
		folders = append(folders, folder)
	}

	return folders, nil
}

// GetFoldersByLibrary retrieves folder contents by folder ID and profile ID
func (r *FolderRepository) GetFoldersByLibrary(folderID uuid.UUID, profileID int) ([]models.FolderContent, error) {
	query := `
		SELECT  
			o.ObjectID,
			o.ObjectName,
			o.ObjectDescription,
			o.typeId,
			o.CheckedInVersionId,
			o.isDeleted,
			o.IsLocked,
			o.IsImported,
			o.IsLibrary,
			o.LibraryId,
			o.FileExtension,
			o.SortOrder,
			o.Prefix,
			o.Suffix,
			o.ProvenanceId,
			o.ProvenanceVersionId,
			o.GeneralType,
			o.CurrentVersionId,
			o.VisioAlias,
			o.HasVisioAlias,
			o.DateCreated,
			o.CreatedBy,
			o.DateModified,
			o.ModifiedBy,
			o.IsCheckedOut,
			o.CheckedOutUserId,
			o.RichTextDescription,
			o.AutoSort,
			CAST(CASE WHEN v.ApprovalStatus = dbo.const_ApprovalStatus_PendingApproval() THEN 1 ELSE 0 END AS BIT) AS IsPendingApproval,
			vchkin.ObjectName AS CheckedInName,
			vchkin.VisioAlias AS CheckedInVisioAlias,
			vchkin.HasVisioAlias AS CheckedInHasVisioAlias,
			ISNULL(op.HasRead, 0) AS HasReadPermission, 
			ISNULL(op.HasModifyContents, 0) AS HasModifyContentsPermission, 
			ISNULL(op.HasDelete, 0) AS HasDeletePermission, 
			ISNULL(op.HasModify, 0) AS HasModifyPermission, 
			ISNULL(op.HasModifyRelationships, 0) AS HasModifyRelationshipsPermission, 
			o.CheckedOutUserId AS CheckedOutBy, 
			CAST(CASE WHEN v.SystemVersionNo = 1 AND o.IsCheckedOut = 1 THEN 1 ELSE 0 END AS BIT) AS IsFirstVersionCheckedOut
		FROM vwFolderContents AS do
		INNER JOIN vwObjectSimple AS o ON o.ObjectID = do.ObjectId
		INNER JOIN [dbo].[Version] AS v on v.ID = o.CurrentVersionId
		LEFT JOIN [dbo].[Version] AS vchkin ON vchkin.ID = o.CheckedInVersionId
		LEFT JOIN [dbo].ObjectPermissions AS op ON op.ObjectID = o.ObjectID AND op.ProfileID = @p2
		WHERE do.FolderId = @p1
			AND o.IsDeleted = CAST(0 AS BIT)
			AND do.IsDeleted = CAST(0 AS BIT)
		ORDER BY SortOrder, CreatedBy
	`

	rows, err := r.db.Query(query, folderID, profileID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving folders by library: %w", err)
	}
	defer rows.Close()

	var contents []models.FolderContent
	for rows.Next() {
		var content models.FolderContent
		var objectIDBytes []byte
		var checkedInVersionIDBytes, libraryIDBytes, provenanceIDBytes, provenanceVersionIDBytes, currentVersionIDBytes []byte
		
		err := rows.Scan(
			&objectIDBytes,
			&content.ObjectName,
			&content.ObjectDescription,
			&content.ObjectTypeID,
			&checkedInVersionIDBytes,
			&content.DeleteFlag,
			&content.Locked,
			&content.IsImported,
			&content.IsLibrary,
			&libraryIDBytes,
			&content.FileExtension,
			&content.SortOrder,
			&content.Prefix,
			&content.Suffix,
			&provenanceIDBytes,
			&provenanceVersionIDBytes,
			&content.GeneralType,
			&currentVersionIDBytes,
			&content.VisioAlias,
			&content.HasVisioAlias,
			&content.DateCreated,
			&content.CreatedBy,
			&content.DateModified,
			&content.ModifiedBy,
			&content.IsCheckedOut,
			&content.CheckedOutUserId,
			&content.RichTextDescription,
			&content.AutoSort,
			&content.IsPendingApproval,
			&content.CheckedInName,
			&content.CheckedInVisioAlias,
			&content.CheckedInHasVisioAlias,
			&content.HasReadPermission,
			&content.HasModifyContentsPermission,
			&content.HasDeletePermission,
			&content.HasModifyPermission,
			&content.HasModifyRelationshipsPermission,
			&content.CheckedOutBy,
			&content.IsFirstVersionCheckedOut,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning folder content: %w", err)
		}
		
		// Parse UUID from bytes
		content.ObjectID, err = parseSQLServerUUID(objectIDBytes)
		if err != nil {
			return nil, fmt.Errorf("error parsing ObjectID: %w", err)
		}
		
		if checkedInVersionIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(checkedInVersionIDBytes)
			if err != nil {
				return nil, fmt.Errorf("error parsing CheckedInVersionId: %w", err)
			}
			content.CheckedInVersionId = &parsedUUID
		}
		
		if libraryIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(libraryIDBytes)
			if err != nil {
				return nil, fmt.Errorf("error parsing LibraryId: %w", err)
			}
			content.LibraryId = &parsedUUID
		}
		
		if provenanceIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(provenanceIDBytes)
			if err != nil {
				return nil, fmt.Errorf("error parsing ProvenanceId: %w", err)
			}
			content.ProvenanceId = &parsedUUID
		}
		
		if provenanceVersionIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(provenanceVersionIDBytes)
			if err != nil {
				return nil, fmt.Errorf("error parsing ProvenanceVersionId: %w", err)
			}
			content.ProvenanceVersionId = &parsedUUID
		}
		
		if currentVersionIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(currentVersionIDBytes)
			if err != nil {
				return nil, fmt.Errorf("error parsing CurrentVersionId: %w", err)
			}
			content.CurrentVersionId = &parsedUUID
		}
		
		contents = append(contents, content)
	}

	return contents, nil
}
