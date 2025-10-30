package repositories

import (
	"database/sql"
	"enterprise-architect-api/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ObjectRepository handles database operations for objects
type ObjectRepository struct {
	db *sql.DB
}

// NewObjectRepository creates a new ObjectRepository
func NewObjectRepository(db *sql.DB) *ObjectRepository {
	return &ObjectRepository{db: db}
}

// Create creates a new object in the database
func (r *ObjectRepository) Create(req models.CreateObjectRequest) (*models.Object, error) {
	objectID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO [Object] (
			ObjectID, ObjectName, ObjectDescription, ObjectTypeID, Locked, IsImported, 
			IsLibrary, LibraryId, FileExtension, Prefix, Suffix, DateCreated, CreatedBy, 
			DateModified, ModifiedBy, IsCheckedOut, ExactObjectTypeID, RichTextDescription
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18
		)
	`

	_, err := r.db.Exec(query,
		objectID, req.ObjectName, req.ObjectDescription, req.ObjectTypeID, false, false,
		req.IsLibrary, req.LibraryId, req.FileExtension, req.Prefix, req.Suffix, now, req.CreatedBy,
		now, req.CreatedBy, false, req.ExactObjectTypeID, req.RichTextDescription,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating object: %w", err)
	}

	return r.GetByID(objectID)
}

// GetByID retrieves an object by its ID
func (r *ObjectRepository) GetByID(id uuid.UUID) (*models.Object, error) {
	query := `
		SELECT ObjectID, ObjectName, ObjectDescription, ObjectTypeID, CheckedInVersionId, 
			DeleteFlag, Locked, RequiresShapeSheetUpdate, TemplateID, IsImported, IsLibrary, 
			LibraryId, FileExtension, SortOrder, Prefix, Suffix, ProvenanceId, ProvenanceVersionId, 
			GeneralType, CurrentVersionId, VisioAlias, HasVisioAlias, DateCreated, CreatedBy, 
			DateModified, ModifiedBy, IsCheckedOut, CheckedOutUserId, DeleteTransactionId, 
			NameChecksum, ExactObjectTypeID, RichTextDescription, AutoSort
		FROM [Object]
		WHERE ObjectID = @p1
	`

	obj := &models.Object{}
	err := r.db.QueryRow(query, id).Scan(
		&obj.ObjectID, &obj.ObjectName, &obj.ObjectDescription, &obj.ObjectTypeID, &obj.CheckedInVersionId,
		&obj.DeleteFlag, &obj.Locked, &obj.RequiresShapeSheetUpdate, &obj.TemplateID, &obj.IsImported,
		&obj.IsLibrary, &obj.LibraryId, &obj.FileExtension, &obj.SortOrder, &obj.Prefix, &obj.Suffix,
		&obj.ProvenanceId, &obj.ProvenanceVersionId, &obj.GeneralType, &obj.CurrentVersionId,
		&obj.VisioAlias, &obj.HasVisioAlias, &obj.DateCreated, &obj.CreatedBy, &obj.DateModified,
		&obj.ModifiedBy, &obj.IsCheckedOut, &obj.CheckedOutUserId, &obj.DeleteTransactionId,
		&obj.NameChecksum, &obj.ExactObjectTypeID, &obj.RichTextDescription, &obj.AutoSort,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("object not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving object: %w", err)
	}

	return obj, nil
}

// GetAll retrieves all objects with pagination
func (r *ObjectRepository) GetAll(page, pageSize int) ([]models.Object, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM [Object]`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting objects: %w", err)
	}

	// Get paginated results
	query := `
		SELECT ObjectID, ObjectName, ObjectDescription, ObjectTypeID, CheckedInVersionId, 
			DeleteFlag, Locked, RequiresShapeSheetUpdate, TemplateID, IsImported, IsLibrary, 
			LibraryId, FileExtension, SortOrder, Prefix, Suffix, ProvenanceId, ProvenanceVersionId, 
			GeneralType, CurrentVersionId, VisioAlias, HasVisioAlias, DateCreated, CreatedBy, 
			DateModified, ModifiedBy, IsCheckedOut, CheckedOutUserId, DeleteTransactionId, 
			NameChecksum, ExactObjectTypeID, RichTextDescription, AutoSort
		FROM [Object]
		ORDER BY DateCreated DESC
		OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
	`

	rows, err := r.db.Query(query, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving objects: %w", err)
	}
	defer rows.Close()

	var objects []models.Object
	for rows.Next() {
		var obj models.Object
		err := rows.Scan(
			&obj.ObjectID, &obj.ObjectName, &obj.ObjectDescription, &obj.ObjectTypeID, &obj.CheckedInVersionId,
			&obj.DeleteFlag, &obj.Locked, &obj.RequiresShapeSheetUpdate, &obj.TemplateID, &obj.IsImported,
			&obj.IsLibrary, &obj.LibraryId, &obj.FileExtension, &obj.SortOrder, &obj.Prefix, &obj.Suffix,
			&obj.ProvenanceId, &obj.ProvenanceVersionId, &obj.GeneralType, &obj.CurrentVersionId,
			&obj.VisioAlias, &obj.HasVisioAlias, &obj.DateCreated, &obj.CreatedBy, &obj.DateModified,
			&obj.ModifiedBy, &obj.IsCheckedOut, &obj.CheckedOutUserId, &obj.DeleteTransactionId,
			&obj.NameChecksum, &obj.ExactObjectTypeID, &obj.RichTextDescription, &obj.AutoSort,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning object: %w", err)
		}
		objects = append(objects, obj)
	}

	return objects, totalCount, nil
}

// Update updates an existing object
func (r *ObjectRepository) Update(id uuid.UUID, req models.UpdateObjectRequest) (*models.Object, error) {
	// Build dynamic update query
	var setClauses []string
	var args []interface{}
	argIndex := 1

	if req.ObjectName != nil {
		setClauses = append(setClauses, fmt.Sprintf("ObjectName = @p%d", argIndex))
		args = append(args, *req.ObjectName)
		argIndex++
	}
	if req.ObjectDescription != nil {
		setClauses = append(setClauses, fmt.Sprintf("ObjectDescription = @p%d", argIndex))
		args = append(args, *req.ObjectDescription)
		argIndex++
	}
	if req.ObjectTypeID != nil {
		setClauses = append(setClauses, fmt.Sprintf("ObjectTypeID = @p%d", argIndex))
		args = append(args, *req.ObjectTypeID)
		argIndex++
	}
	if req.ExactObjectTypeID != nil {
		setClauses = append(setClauses, fmt.Sprintf("ExactObjectTypeID = @p%d", argIndex))
		args = append(args, *req.ExactObjectTypeID)
		argIndex++
	}
	if req.RichTextDescription != nil {
		setClauses = append(setClauses, fmt.Sprintf("RichTextDescription = @p%d", argIndex))
		args = append(args, *req.RichTextDescription)
		argIndex++
	}
	if req.IsLibrary != nil {
		setClauses = append(setClauses, fmt.Sprintf("IsLibrary = @p%d", argIndex))
		args = append(args, *req.IsLibrary)
		argIndex++
	}
	if req.FileExtension != nil {
		setClauses = append(setClauses, fmt.Sprintf("FileExtension = @p%d", argIndex))
		args = append(args, *req.FileExtension)
		argIndex++
	}
	if req.Prefix != nil {
		setClauses = append(setClauses, fmt.Sprintf("Prefix = @p%d", argIndex))
		args = append(args, *req.Prefix)
		argIndex++
	}
	if req.Suffix != nil {
		setClauses = append(setClauses, fmt.Sprintf("Suffix = @p%d", argIndex))
		args = append(args, *req.Suffix)
		argIndex++
	}

	// Always update DateModified and ModifiedBy
	setClauses = append(setClauses, fmt.Sprintf("DateModified = @p%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("ModifiedBy = @p%d", argIndex))
	args = append(args, req.ModifiedBy)
	argIndex++

	if len(setClauses) == 2 { // Only DateModified and ModifiedBy
		return nil, fmt.Errorf("no fields to update")
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE [Object] SET %s WHERE ObjectID = @p%d", strings.Join(setClauses, ", "), argIndex)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating object: %w", err)
	}

	return r.GetByID(id)
}

// Delete deletes an object by its ID
func (r *ObjectRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM [Object] WHERE ObjectID = @p1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting object: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("object not found")
	}

	return nil
}

// GetLibraries retrieves all objects where IsLibrary is true
func (r *ObjectRepository) GetLibraries(page, pageSize int) ([]models.Object, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM [Object] WHERE IsLibrary = 1`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting libraries: %w", err)
	}

	// Get paginated results
	query := `
		SELECT ObjectID, ObjectName, ObjectDescription, ObjectTypeID, CheckedInVersionId, 
			DeleteFlag, Locked, RequiresShapeSheetUpdate, TemplateID, IsImported, IsLibrary, 
			LibraryId, FileExtension, SortOrder, Prefix, Suffix, ProvenanceId, ProvenanceVersionId, 
			GeneralType, CurrentVersionId, VisioAlias, HasVisioAlias, DateCreated, CreatedBy, 
			DateModified, ModifiedBy, IsCheckedOut, CheckedOutUserId, DeleteTransactionId, 
			NameChecksum, ExactObjectTypeID, RichTextDescription, AutoSort
		FROM [Object]
		WHERE IsLibrary = 1
		ORDER BY DateCreated DESC
		OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
	`

	rows, err := r.db.Query(query, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving libraries: %w", err)
	}
	defer rows.Close()

	var objects []models.Object
	for rows.Next() {
		var obj models.Object
		var objectIDBytes []byte
		var checkedInVersionIDBytes, libraryIDBytes, provenanceIDBytes, provenanceVersionIDBytes, currentVersionIDBytes, deleteTransactionIDBytes []byte

		err := rows.Scan(
			&objectIDBytes, &obj.ObjectName, &obj.ObjectDescription, &obj.ObjectTypeID, &checkedInVersionIDBytes,
			&obj.DeleteFlag, &obj.Locked, &obj.RequiresShapeSheetUpdate, &obj.TemplateID, &obj.IsImported,
			&obj.IsLibrary, &libraryIDBytes, &obj.FileExtension, &obj.SortOrder, &obj.Prefix, &obj.Suffix,
			&provenanceIDBytes, &provenanceVersionIDBytes, &obj.GeneralType, &currentVersionIDBytes,
			&obj.VisioAlias, &obj.HasVisioAlias, &obj.DateCreated, &obj.CreatedBy, &obj.DateModified,
			&obj.ModifiedBy, &obj.IsCheckedOut, &obj.CheckedOutUserId, &deleteTransactionIDBytes,
			&obj.NameChecksum, &obj.ExactObjectTypeID, &obj.RichTextDescription, &obj.AutoSort,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning library: %w", err)
		}

		// Parse UUID from bytes
		obj.ObjectID, err = parseSQLServerUUID(objectIDBytes)
		if err != nil {
			return nil, 0, fmt.Errorf("error parsing ObjectID: %w", err)
		}

		if checkedInVersionIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(checkedInVersionIDBytes)
			if err != nil {
				return nil, 0, fmt.Errorf("error parsing CheckedInVersionId: %w", err)
			}
			obj.CheckedInVersionId = &parsedUUID
		}

		if libraryIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(libraryIDBytes)
			if err != nil {
				return nil, 0, fmt.Errorf("error parsing LibraryId: %w", err)
			}
			obj.LibraryId = &parsedUUID
		}

		if provenanceIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(provenanceIDBytes)
			if err != nil {
				return nil, 0, fmt.Errorf("error parsing ProvenanceId: %w", err)
			}
			obj.ProvenanceId = &parsedUUID
		}

		if provenanceVersionIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(provenanceVersionIDBytes)
			if err != nil {
				return nil, 0, fmt.Errorf("error parsing ProvenanceVersionId: %w", err)
			}
			obj.ProvenanceVersionId = &parsedUUID
		}

		if currentVersionIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(currentVersionIDBytes)
			if err != nil {
				return nil, 0, fmt.Errorf("error parsing CurrentVersionId: %w", err)
			}
			obj.CurrentVersionId = &parsedUUID
		}

		if deleteTransactionIDBytes != nil {
			parsedUUID, err := parseSQLServerUUID(deleteTransactionIDBytes)
			if err != nil {
				return nil, 0, fmt.Errorf("error parsing DeleteTransactionId: %w", err)
			}
			obj.DeleteTransactionId = &parsedUUID
		}

		objects = append(objects, obj)
	}

	return objects, totalCount, nil
}

// GetByObjectTypeID retrieves all objects by ObjectTypeID with pagination
func (r *ObjectRepository) GetByObjectTypeID(objectTypeID, page, pageSize int) ([]models.Object, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM [Object] WHERE ObjectTypeID = @p1`
	err := r.db.QueryRow(countQuery, objectTypeID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting objects by type: %w", err)
	}

	// Get paginated results
	query := `
		SELECT ObjectID, ObjectName, ObjectDescription, ObjectTypeID, CheckedInVersionId, 
			DeleteFlag, Locked, RequiresShapeSheetUpdate, TemplateID, IsImported, IsLibrary, 
			LibraryId, FileExtension, SortOrder, Prefix, Suffix, ProvenanceId, ProvenanceVersionId, 
			GeneralType, CurrentVersionId, VisioAlias, HasVisioAlias, DateCreated, CreatedBy, 
			DateModified, ModifiedBy, IsCheckedOut, CheckedOutUserId, DeleteTransactionId, 
			NameChecksum, ExactObjectTypeID, RichTextDescription, AutoSort
		FROM [Object]
		WHERE ObjectTypeID = @p1
		ORDER BY DateCreated DESC
		OFFSET @p2 ROWS FETCH NEXT @p3 ROWS ONLY
	`

	rows, err := r.db.Query(query, objectTypeID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving objects by type: %w", err)
	}
	defer rows.Close()

	var objects []models.Object
	for rows.Next() {
		var obj models.Object
		err := rows.Scan(
			&obj.ObjectID, &obj.ObjectName, &obj.ObjectDescription, &obj.ObjectTypeID, &obj.CheckedInVersionId,
			&obj.DeleteFlag, &obj.Locked, &obj.RequiresShapeSheetUpdate, &obj.TemplateID, &obj.IsImported,
			&obj.IsLibrary, &obj.LibraryId, &obj.FileExtension, &obj.SortOrder, &obj.Prefix, &obj.Suffix,
			&obj.ProvenanceId, &obj.ProvenanceVersionId, &obj.GeneralType, &obj.CurrentVersionId,
			&obj.VisioAlias, &obj.HasVisioAlias, &obj.DateCreated, &obj.CreatedBy, &obj.DateModified,
			&obj.ModifiedBy, &obj.IsCheckedOut, &obj.CheckedOutUserId, &obj.DeleteTransactionId,
			&obj.NameChecksum, &obj.ExactObjectTypeID, &obj.RichTextDescription, &obj.AutoSort,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning object: %w", err)
		}
		objects = append(objects, obj)
	}

	return objects, totalCount, nil
}

func (r *ObjectRepository) GetHierarchyFolder(ObjectID uuid.UUID, profileID int, isFolder bool) ([]models.ObjectTree, error) {
	query := `
			WITH recurse ([ObjectID]
		  , [ObjectParentID]
		  , [ObjectName]
		  , [ObjectDescription]
		  , [CurrentVersionId]
		  , [CheckedInVersionId]
		  , [IsImported]
		  , [IsLibrary]
		  , [LibraryId]
		  , [FileExtension]
		  , [GeneralType]
		  , [GeneralTypeName]
		  , [TypeId]
		  , [TypeName]
		  , [IsDeleted]
		  , [VisioAlias]
		  , [HasVisioAlias]
		  , [IsLocked]
		  , [SortOrder]
		  , [AutoSort]
		  , [Prefix]
		  , [Suffix]
		  , [ProvenanceID]
		  , [ProvenanceVersionID]
		  , [CreatedBy]
		  , [DateCreated]
		  , [ModifiedBy]
		  , [DateModified]
		  , [IsCheckedOut]
		  , [CheckedOutUserId]
		  , [RichTextDescription]
		  , IsPendingApproval
		  , CheckedOutBy
		  , IsFirstVersionCheckedOut
		  , FolderId)
		  AS
	(
		SELECT  o.*, 
				CAST(CASE WHEN v.ApprovalStatus =  dbo.const_ApprovalStatus_PendingApproval() THEN 1 ELSE 0 END AS BIT) AS IsPendingApproval,
				o.CheckedOutUserId AS CheckedOutBy, 
				CAST(CASE WHEN v.SystemVersionNo = 1 AND o.IsCheckedOut = 1 THEN 1 ELSE 0 END AS BIT) AS IsFirstVersionCheckedOut,
		        do.FolderId
		FROM vwFolderContents AS do
		INNER JOIN vwObjectSimple AS o ON o.ObjectID = do.ObjectId
		INNER JOIN [dbo].[Version] AS v on v.ID = o.CurrentVersionId
		WHERE do.FolderId = @p1
			AND o.IsDeleted = CAST(0 AS BIT)
			AND do.IsDeleted = CAST(0 AS BIT)
			AND do.isFolder in (CAST(@p2 AS BIT), CAST(1 AS BIT))
		UNION ALL
		SELECT  o.*, 
				CAST(CASE WHEN v.ApprovalStatus =  dbo.const_ApprovalStatus_PendingApproval() THEN 1 ELSE 0 END AS BIT) AS IsPendingApproval,
				o.CheckedOutUserId AS CheckedOutBy, 
				CAST(CASE WHEN v.SystemVersionNo = 1 AND o.IsCheckedOut = 1 THEN 1 ELSE 0 END AS BIT) AS IsFirstVersionCheckedOut,
		        do.FolderId
		FROM vwFolderContents AS do
		INNER JOIN recurse ON do.FolderId = recurse.ObjectID
		INNER JOIN vwObjectSimple AS o ON o.ObjectID = do.ObjectId
		INNER JOIN [dbo].[Version] AS v on v.ID = o.CurrentVersionId
		WHERE o.IsDeleted = CAST(0 AS BIT)
			AND do.IsDeleted = CAST(0 AS BIT)
			AND do.isFolder in (CAST(@p2 AS BIT), CAST(1 AS BIT))
	)
	SELECT recurse.* ,
		vchkin.ObjectName AS CheckedInName,
		vchkin.VisioAlias AS CheckedInVisioAlias,
		vchkin.HasVisioAlias AS CheckedInHasVisioAlias,
		ISNULL(op.HasRead, 0) AS HasReadPermission, 
		ISNULL(op.HasModifyContents, 0) AS HasModifyContentsPermission, 
		ISNULL(op.HasDelete, 0) AS HasDeletePermission, 
		ISNULL(op.HasModify, 0) AS HasModifyPermission, 
		ISNULL(op.HasModifyRelationships, 0) AS HasModifyRelationshipsPermission
	FROM recurse
	LEFT JOIN [dbo].[Version] AS vchkin ON vchkin.ID = CheckedInVersionId
	LEFT JOIN [dbo].ObjectPermissions AS op ON op.ObjectID = recurse.ObjectID AND op.ProfileID =@p3
	`

	// Convert UUID to SQL Server format
	sqlServerUUID := toSQLServerUUID(ObjectID)
	sqlUUID, _ := uuid.FromBytes(sqlServerUUID)
	fmt.Println("UUID:", sqlUUID.String())
	rows, err := r.db.Query(query, sqlUUID, isFolder, profileID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var objects []models.ObjectTree
	for rows.Next() {
		var obj models.ObjectTree
		err := rows.Scan(
			// Fields from recurse CTE
			&obj.ObjectID,
			&obj.ObjectParentID,
			&obj.ObjectName,
			&obj.ObjectDescription,
			&obj.CurrentVersionId,
			&obj.CheckedInVersionId,
			&obj.IsImported,
			&obj.IsLibrary,
			&obj.LibraryId,
			&obj.FileExtension,
			&obj.GeneralType,
			&obj.GeneralTypeName,
			&obj.TypeId,
			&obj.TypeName,
			&obj.IsDeleted,
			&obj.VisioAlias,
			&obj.HasVisioAlias,
			&obj.IsLocked,
			&obj.SortOrder,
			&obj.AutoSort,
			&obj.Prefix,
			&obj.Suffix,
			&obj.ProvenanceID,
			&obj.ProvenanceVersionID,
			&obj.CreatedBy,
			&obj.DateCreated,
			&obj.ModifiedBy,
			&obj.DateModified,
			&obj.IsCheckedOut,
			&obj.CheckedOutUserId,
			&obj.RichTextDescription,
			&obj.IsPendingApproval,
			&obj.CheckedOutBy,
			&obj.IsFirstVersionCheckedOut,
			&obj.FolderId,
			// Additional fields from joins
			&obj.CheckedInName,
			&obj.CheckedInVisioAlias,
			&obj.CheckedInHasVisioAlias,
			&obj.HasReadPermission,
			&obj.HasModifyContentsPermission,
			&obj.HasDeletePermission,
			&obj.HasModifyPermission,
			&obj.HasModifyRelationshipsPermission,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		objects = append(objects, obj)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return objects, nil
}
