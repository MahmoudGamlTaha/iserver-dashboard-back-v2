package repositories

import (
	"database/sql"
	"enterprise-architect-api/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ObjectTypeRepository handles database operations for object types
type ObjectTypeRepository struct {
	db *sql.DB
}

// NewObjectTypeRepository creates a new ObjectTypeRepository
func NewObjectTypeRepository(db *sql.DB) *ObjectTypeRepository {
	return &ObjectTypeRepository{db: db}
}

// Create creates a new object type in the database
func (r *ObjectTypeRepository) Create(req models.CreateObjectTypeRequest) (*models.ObjectType, error) {
	now := time.Now()

	query := `
		INSERT INTO ObjectType (
			ObjectTypeName, IsTemplateType, IsDefaultTemplate, ActiveType, 
			EnforceUniqueNaming, CanHaveVisioAlias, IsConnector, ImplicitlyAddObjectTypes, 
			CommitOverlapRelationships, DateCreated, CreatedBy, DateModified, ModifiedBy, 
			FileExtension, IsExcludedFromBrokenConnectors, Description, ExportShapeAttributes, 
			ExportShapeSystemProperties, ImportShapeAttributes, ExportDocumentAttributes, 
			ExportDocumentSystemProperties, DeleteNotSyncVisioShapeData, DeleteIfHasNoMaster
		) OUTPUT INSERTED.ObjectTypeID VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, 
			@p17, @p18, @p19, @p20, @p21, @p22, @p23
		)
	`

	var objectTypeID int
	err := r.db.QueryRow(query,
		req.ObjectTypeName, req.IsTemplateType, false, req.ActiveType,
		false, false, false, false,
		false, now, req.CreatedBy, now, req.CreatedBy,
		req.FileExtension, false, req.Description, false,
		false, false, false,
		false, false, false,
	).Scan(&objectTypeID)

	if err != nil {
		return nil, fmt.Errorf("error creating object type: %w", err)
	}

	return r.GetByID(objectTypeID)
}

// GetByID retrieves an object type by its ID
func (r *ObjectTypeRepository) GetByID(id int) (*models.ObjectType, error) {
	query := `
		SELECT ObjectTypeID, ObjectTypeName, ObjectTypeImage, IsTemplateType, GeneralType, 
			TemplateFileName, IsDefaultTemplate, ActiveType, EnforceUniqueNaming, CanHaveVisioAlias, 
			IsConnector, ImplicitlyAddObjectTypes, CommitOverlapRelationships, DateCreated, CreatedBy, 
			DateModified, ModifiedBy, FileExtension, HandlerToolId, Color, Icon, 
			IsExcludedFromBrokenConnectors, Description, ExportShapeAttributes, ExportShapeSystemProperties, 
			ImportShapeAttributes, ExportDocumentAttributes, ExportDocumentSystemProperties, 
			DeleteNotSyncVisioShapeData, DeleteIfHasNoMaster
		FROM ObjectType
		WHERE ObjectTypeID = @p1
	`

	objType := &models.ObjectType{}
	err := r.db.QueryRow(query, id).Scan(
		&objType.ObjectTypeID, &objType.ObjectTypeName, &objType.ObjectTypeImage, &objType.IsTemplateType,
		&objType.GeneralType, &objType.TemplateFileName, &objType.IsDefaultTemplate, &objType.ActiveType,
		&objType.EnforceUniqueNaming, &objType.CanHaveVisioAlias, &objType.IsConnector,
		&objType.ImplicitlyAddObjectTypes, &objType.CommitOverlapRelationships, &objType.DateCreated,
		&objType.CreatedBy, &objType.DateModified, &objType.ModifiedBy, &objType.FileExtension,
		&objType.HandlerToolId, &objType.Color, &objType.Icon, &objType.IsExcludedFromBrokenConnectors,
		&objType.Description, &objType.ExportShapeAttributes, &objType.ExportShapeSystemProperties,
		&objType.ImportShapeAttributes, &objType.ExportDocumentAttributes, &objType.ExportDocumentSystemProperties,
		&objType.DeleteNotSyncVisioShapeData, &objType.DeleteIfHasNoMaster,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("object type not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving object type: %w", err)
	}

	return objType, nil
}

// GetAll retrieves all object types with pagination
func (r *ObjectTypeRepository) GetAll(page, pageSize int) ([]models.ObjectType, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM ObjectType`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting object types: %w", err)
	}

	// Get paginated results
	query := `
		SELECT ObjectTypeID, ObjectTypeName, ObjectTypeImage, IsTemplateType, GeneralType, 
			TemplateFileName, IsDefaultTemplate, ActiveType, EnforceUniqueNaming, CanHaveVisioAlias, 
			IsConnector, ImplicitlyAddObjectTypes, CommitOverlapRelationships, DateCreated, CreatedBy, 
			DateModified, ModifiedBy, FileExtension, HandlerToolId, Color, Icon, 
			IsExcludedFromBrokenConnectors, Description, ExportShapeAttributes, ExportShapeSystemProperties, 
			ImportShapeAttributes, ExportDocumentAttributes, ExportDocumentSystemProperties, 
			DeleteNotSyncVisioShapeData, DeleteIfHasNoMaster
		FROM ObjectType
		ORDER BY DateCreated DESC
		OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
	`

	rows, err := r.db.Query(query, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving object types: %w", err)
	}
	defer rows.Close()

	var objectTypes []models.ObjectType
	for rows.Next() {
		var objType models.ObjectType
		err := rows.Scan(
			&objType.ObjectTypeID, &objType.ObjectTypeName, &objType.ObjectTypeImage, &objType.IsTemplateType,
			&objType.GeneralType, &objType.TemplateFileName, &objType.IsDefaultTemplate, &objType.ActiveType,
			&objType.EnforceUniqueNaming, &objType.CanHaveVisioAlias, &objType.IsConnector,
			&objType.ImplicitlyAddObjectTypes, &objType.CommitOverlapRelationships, &objType.DateCreated,
			&objType.CreatedBy, &objType.DateModified, &objType.ModifiedBy, &objType.FileExtension,
			&objType.HandlerToolId, &objType.Color, &objType.Icon, &objType.IsExcludedFromBrokenConnectors,
			&objType.Description, &objType.ExportShapeAttributes, &objType.ExportShapeSystemProperties,
			&objType.ImportShapeAttributes, &objType.ExportDocumentAttributes, &objType.ExportDocumentSystemProperties,
			&objType.DeleteNotSyncVisioShapeData, &objType.DeleteIfHasNoMaster,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning object type: %w", err)
		}
		objectTypes = append(objectTypes, objType)
	}

	return objectTypes, totalCount, nil
}

// Update updates an existing object type
func (r *ObjectTypeRepository) Update(id int, req models.UpdateObjectTypeRequest) (*models.ObjectType, error) {
	// Build dynamic update query
	var setClauses []string
	var args []interface{}
	argIndex := 1

	if req.ObjectTypeName != nil {
		setClauses = append(setClauses, fmt.Sprintf("ObjectTypeName = @p%d", argIndex))
		args = append(args, *req.ObjectTypeName)
		argIndex++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("Description = @p%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.FileExtension != nil {
		setClauses = append(setClauses, fmt.Sprintf("FileExtension = @p%d", argIndex))
		args = append(args, *req.FileExtension)
		argIndex++
	}
	if req.IsTemplateType != nil {
		setClauses = append(setClauses, fmt.Sprintf("IsTemplateType = @p%d", argIndex))
		args = append(args, *req.IsTemplateType)
		argIndex++
	}
	if req.ActiveType != nil {
		setClauses = append(setClauses, fmt.Sprintf("ActiveType = @p%d", argIndex))
		args = append(args, *req.ActiveType)
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
	query := fmt.Sprintf("UPDATE ObjectType SET %s WHERE ObjectTypeID = @p%d", strings.Join(setClauses, ", "), argIndex)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating object type: %w", err)
	}

	return r.GetByID(id)
}

// Delete deletes an object type by its ID
func (r *ObjectTypeRepository) Delete(id int) error {
	query := `DELETE FROM ObjectType WHERE ObjectTypeID = @p1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting object type: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("object type not found")
	}

	return nil
}

func (r *ObjectTypeRepository) GetFolderRepositoryTree() ([]models.ObjectTypeHierarchy, error) {
	query := `WITH FolderHierarchy AS (
				SELECT 
					ot.ObjectTypeName,
					ot.ObjectTypeID,
					fth.FolderTypeHierarchyId,
					fth.FolderObjectTypeId,
					fth.ParentHierarchyId,
					0 as Level,
					CAST(ot.ObjectTypeName AS NVARCHAR(MAX)) as FullPath
				FROM FolderTypeHierarchy fth
				INNER JOIN ObjectType ot ON ot.ObjectTypeID = fth.FolderObjectTypeId
				WHERE fth.ParentHierarchyId is null
				
				UNION ALL
				
				SELECT 
					ot.ObjectTypeName,
					ot.ObjectTypeID,
					fth.FolderTypeHierarchyId,
					fth.FolderObjectTypeId,
					fth.ParentHierarchyId,
					fh.Level + 1 as Level,
					CAST(fh.FullPath + ' > ' + ot.ObjectTypeName AS NVARCHAR(MAX)) as FullPath
				FROM FolderTypeHierarchy fth
				INNER JOIN ObjectType ot ON ot.ObjectTypeID = fth.FolderObjectTypeId
				INNER JOIN FolderHierarchy fh ON fth.ParentHierarchyId = fh.FolderTypeHierarchyId
			)
			SELECT 
			   
				fh.ObjectTypeName,
				fh.ObjectTypeID,
				fh.FolderObjectTypeId as objectTypeFolderId,
				fh.ParentHierarchyId as objectTypeParentId,
				fh.FolderTypeHierarchyId as objectTypeHierarchyId,
				fh.Level,
				fh.FullPath,
				CAST(CASE WHEN EXISTS (
					SELECT 1 FROM FolderObjectTypes fot 
					WHERE fot.FolderObjectTypeId = fh.FolderObjectTypeId 
					AND fot.IsDocumentType = 1
				) THEN 1 ELSE 0 END AS BIT) as IsDocumentType
			FROM FolderHierarchy fh
			ORDER BY fh.FullPath;
			`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving folder repository tree: %w", err)
	}
	defer rows.Close()

	var folderRepositoryTree []models.ObjectTypeHierarchy
	for rows.Next() {
		var objType models.ObjectTypeHierarchy
		err := rows.Scan(
			&objType.ObjectTypeName, &objType.ObjectTypeId, &objType.ObjectTypeFolderId, &objType.ObjectTypeParentId,
			&objType.ObjectTypeHierarchyId, &objType.Level, &objType.FullPath, &objType.IsDocumentType,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning folder repository tree: %w", err)
		}
		folderRepositoryTree = append(folderRepositoryTree, objType)
	}

	return folderRepositoryTree, nil
}

// AddFolderToTree adds a new folder to the folder hierarchy tree
func (r *ObjectTypeRepository) AddFolderToTree(req models.AddFolderToTreeRequest) (*uuid.UUID, error) {
	var folderObjectTypeId int
	var err error

	// If FolderObjectTypeId is 0, create a new ObjectType
	if req.FolderObjectTypeId == 0 {
		// Validate that ObjectTypeName is provided
		if req.ObjectTypeName == "" {
			return nil, fmt.Errorf("object type name is required when creating a new object type")
		}

		// Insert new ObjectType
		insertObjectTypeQuery := `
			INSERT INTO ObjectType (
				ObjectTypeName, IsTemplateType, IsDefaultTemplate, ActiveType, 
				EnforceUniqueNaming, CanHaveVisioAlias, IsConnector, ImplicitlyAddObjectTypes, 
				CommitOverlapRelationships, DateCreated, CreatedBy, DateModified, ModifiedBy, 
				IsExcludedFromBrokenConnectors, ExportShapeAttributes, 
				ExportShapeSystemProperties, ImportShapeAttributes, ExportDocumentAttributes, 
				ExportDocumentSystemProperties, DeleteNotSyncVisioShapeData, DeleteIfHasNoMaster,GeneralType
			) OUTPUT INSERTED.ObjectTypeID VALUES (
				@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18, @p19, @p20, @p21,@p22
			)
		`

		now := time.Now()
		err := r.db.QueryRow(insertObjectTypeQuery,
			req.ObjectTypeName, false, false, true,
			true, false, false, false,
			false, now, 62, now, 62,
			false, false,
			false, false, false,
			false, false, false,
			2, // fixed for folder type
		).Scan(&folderObjectTypeId)

		if err != nil {
			return nil, fmt.Errorf("error creating new object type: %w", err)
		}
	} else {
		// Validate that the FolderObjectTypeId exists
		var exists bool
		checkQuery := `SELECT CASE WHEN EXISTS (SELECT 1 FROM ObjectType WHERE ObjectTypeID = @p1) THEN 1 ELSE 0 END`
		err := r.db.QueryRow(checkQuery, req.FolderObjectTypeId).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("error checking if object type exists: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("object type with ID %d does not exist", req.FolderObjectTypeId)
		}
		folderObjectTypeId = req.FolderObjectTypeId
	}
	var parentHierarchyId uuid.UUID
	parentHierarchyId, err = TransformUUID(*req.ParentHierarchyId)
	if err != nil {
		return nil, fmt.Errorf("error transforming parent hierarchy ID: %w", err)
	}

	// If ParentHierarchyId is provided, validate it exists
	if req.ParentHierarchyId != nil {
		var parentExists bool
		checkParentQuery := `SELECT CASE WHEN EXISTS (SELECT 1 FROM FolderTypeHierarchy WHERE FolderTypeHierarchyId = @p1) THEN 1 ELSE 0 END`
		err = r.db.QueryRow(checkParentQuery, parentHierarchyId).Scan(&parentExists)
		if err != nil {
			return nil, fmt.Errorf("error checking if parent hierarchy exists: %w", err)
		}
		if !parentExists {
			return nil, fmt.Errorf("parent hierarchy with ID %s does not exist", parentHierarchyId.String())
		}
	}

	// Insert new folder into hierarchy
	insertQuery := `
		INSERT INTO FolderTypeHierarchy (
			FolderTypeHierarchyId, 
			FolderObjectTypeId, 
			ParentHierarchyId
		) OUTPUT INSERTED.FolderTypeHierarchyId 
		VALUES (NEWID(), @p1, @p2)
	`

	var folderTypeHierarchyId uuid.UUID

	fmt.Println("Parent Hierarchy ID: ", parentHierarchyId)
	err = r.db.QueryRow(insertQuery, folderObjectTypeId, parentHierarchyId).Scan(&folderTypeHierarchyId)
	if err != nil {
		return nil, fmt.Errorf("error adding folder to tree: %w", err)
	}

	return &folderTypeHierarchyId, nil
}

// AssignObjectTypeToFolder assigns an object type to a folder type
func (r *ObjectTypeRepository) AssignObjectTypeToFolder(req models.FolderObjectTypes) error {
	query := `
		INSERT INTO FolderObjectTypes (
			FolderObjectTypeId, 
			ObjectTypeId, 
			IsDocumentType
		) VALUES (@p1, @p2, @p3)
	`

	_, err := r.db.Exec(query, req.FolderObjectTypeId, req.ObjectTypeID, req.IsDocumentType)
	if err != nil {
		return fmt.Errorf("error assigning object type to folder: %w", err)
	}

	return nil
}

// GetAvailableTypesForFolder retrieves available object types for a specific folder
func (r *ObjectTypeRepository) GetAvailableTypesForFolder(folderObjectTypeId int) ([]models.FolderObjectTypesNames, error) {
	query := `
		SELECT ot.ObjectTypeName, fo.FolderObjectTypeId, fo.ObjectTypeId, fo.IsDocumentType 
		FROM FolderObjectTypes fo 
		INNER JOIN ObjectType ot ON fo.ObjectTypeId = ot.ObjectTypeID
		WHERE FolderObjectTypeId = @p1
	`

	rows, err := r.db.Query(query, folderObjectTypeId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving available types for folder: %w", err)
	}
	defer rows.Close()

	var folderObjectTypes []models.FolderObjectTypesNames
	for rows.Next() {
		var fot models.FolderObjectTypesNames
		err := rows.Scan(
			&fot.ObjectTypeName,
			&fot.FolderObjectTypeId,
			&fot.ObjectTypeID,
			&fot.IsDocumentType,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning folder object type: %w", err)
		}
		folderObjectTypes = append(folderObjectTypes, fot)
	}

	return folderObjectTypes, nil
}

// DeleteObjectTypeFromFolder removes an object type assignment from a folder
func (r *ObjectTypeRepository) DeleteObjectTypeFromFolder(folderObjectTypeId, objectTypeId int) error {
	query := `
		DELETE FROM FolderObjectTypes 
		WHERE FolderObjectTypeId = @p1 AND ObjectTypeId = @p2
	`

	result, err := r.db.Exec(query, folderObjectTypeId, objectTypeId)
	if err != nil {
		return fmt.Errorf("error deleting object type from folder: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("object type assignment not found")
	}

	return nil
}

// SearchByName retrieves object types filtered by name with pagination
func (r *ObjectTypeRepository) SearchByName(name string, page, pageSize int) ([]models.ObjectType, int, error) {
	offset := (page - 1) * pageSize

	// Total count with filter
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM ObjectType WHERE ObjectTypeName LIKE '%' + @p1 + '%'`
	if err := r.db.QueryRow(countQuery, name).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("error counting object types by name: %w", err)
	}

	query := `
		SELECT ObjectTypeID, ObjectTypeName, ObjectTypeImage, IsTemplateType, GeneralType,
			TemplateFileName, IsDefaultTemplate, ActiveType, EnforceUniqueNaming, CanHaveVisioAlias,
			IsConnector, ImplicitlyAddObjectTypes, CommitOverlapRelationships, DateCreated, CreatedBy,
			DateModified, ModifiedBy, FileExtension, HandlerToolId, Color, Icon,
			IsExcludedFromBrokenConnectors, Description, ExportShapeAttributes, ExportShapeSystemProperties,
			ImportShapeAttributes, ExportDocumentAttributes, ExportDocumentSystemProperties,
			DeleteNotSyncVisioShapeData, DeleteIfHasNoMaster
		FROM ObjectType
		WHERE ObjectTypeName LIKE '%' + @p1 + '%'
		ORDER BY DateCreated DESC
		OFFSET @p2 ROWS FETCH NEXT @p3 ROWS ONLY
	`

	rows, err := r.db.Query(query, name, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving object types by name: %w", err)
	}
	defer rows.Close()

	var objectTypes []models.ObjectType
	for rows.Next() {
		var objType models.ObjectType
		if err := rows.Scan(
			&objType.ObjectTypeID, &objType.ObjectTypeName, &objType.ObjectTypeImage, &objType.IsTemplateType,
			&objType.GeneralType, &objType.TemplateFileName, &objType.IsDefaultTemplate, &objType.ActiveType,
			&objType.EnforceUniqueNaming, &objType.CanHaveVisioAlias, &objType.IsConnector,
			&objType.ImplicitlyAddObjectTypes, &objType.CommitOverlapRelationships, &objType.DateCreated,
			&objType.CreatedBy, &objType.DateModified, &objType.ModifiedBy, &objType.FileExtension,
			&objType.HandlerToolId, &objType.Color, &objType.Icon, &objType.IsExcludedFromBrokenConnectors,
			&objType.Description, &objType.ExportShapeAttributes, &objType.ExportShapeSystemProperties,
			&objType.ImportShapeAttributes, &objType.ExportDocumentAttributes, &objType.ExportDocumentSystemProperties,
			&objType.DeleteNotSyncVisioShapeData, &objType.DeleteIfHasNoMaster,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning object type by name: %w", err)
		}
		objectTypes = append(objectTypes, objType)
	}

	return objectTypes, totalCount, nil
}

// GetBaseLibrary retrieves the base library of object types
func (r *ObjectTypeRepository) GetBaseLibrary() ([]models.ObjectTypeHierarchy, error) {
	sql := `SELECT 
					ot.ObjectTypeName,
					fth.FolderTypeHierarchyId AS ObjectTypeHierarchyId,
					fth.FolderObjectTypeId AS ObjectTypeId
				FROM FolderTypeHierarchy fth
				INNER JOIN ObjectType ot ON ot.ObjectTypeID = fth.FolderObjectTypeId
				WHERE fth.ParentHierarchyId is null`

	var baseLib models.ObjectTypeHierarchy

	err := r.db.QueryRow(sql).Scan(&baseLib.ObjectTypeName, &baseLib.ObjectTypeHierarchyId, &baseLib.ObjectTypeId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving base library: %w", err)
	}

	var baseLibraryList []models.ObjectTypeHierarchy

	sql = `SELECT 
					ot.ObjectTypeName,
					fth.FolderTypeHierarchyId AS ObjectTypeHierarchyId,
					fth.FolderObjectTypeId AS ObjectTypeId
					FROM FolderTypeHierarchy fth
					INNER JOIN ObjectType ot ON ot.ObjectTypeID = fth.FolderObjectTypeId
					WHERE fth.ParentHierarchyId = @p1`
	fmt.Println("Base Library ID: ", baseLib.ObjectTypeHierarchyId)
	if baseLib.ObjectTypeHierarchyId != nil {
		transformedUUID, _ := TransformUUID(*baseLib.ObjectTypeHierarchyId)
		baseLib.ObjectTypeHierarchyId = &transformedUUID
	}
	rows, err := r.db.Query(sql, baseLib.ObjectTypeHierarchyId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving base library: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var baseLib models.ObjectTypeHierarchy
		if err := rows.Scan(&baseLib.ObjectTypeName, &baseLib.ObjectTypeHierarchyId, &baseLib.ObjectTypeId); err != nil {
			return nil, fmt.Errorf("error scanning base library: %w", err)
		}
		baseLibraryList = append(baseLibraryList, baseLib)
	}

	return baseLibraryList, nil
}
func (r *ObjectTypeRepository) GetAvailableTypesForLibsAndFolder(folderObjectTypeId int) ([]models.FolderObjectTypesNames, error) {

	query := `
		SELECT 
			ot.ObjectTypeName,
			fth.FolderTypeHierarchyId AS objectTypeHierarchyId,
			fth.FolderObjectTypeId,
			CASE 
				WHEN ot.GeneralType = 2 THEN CAST(1 AS bit) 
				ELSE CAST(0 AS bit)
			END AS isDocumentType
			FROM FolderTypeHierarchy fth
			INNER JOIN ObjectType ot 
				ON ot.ObjectTypeID = fth.FolderObjectTypeId
			WHERE ot.ObjectTypeID = @p1`

	var fot models.FolderObjectTypesNames
	err := r.db.QueryRow(query, folderObjectTypeId).Scan(
		&fot.ObjectTypeName,
		&fot.ObjectTypeHierarchyId,
		&fot.FolderObjectTypeId,
		&fot.IsDocumentType,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tree of available types for folder: %w", err)
	}

	var folderObjectTypes []models.FolderObjectTypesNames
	*fot.ObjectTypeHierarchyId, _ = TransformUUID(*fot.ObjectTypeHierarchyId)
	query = `
		SELECT 
			distinct ot.ObjectTypeName,
			fth.FolderTypeHierarchyId AS ObjectTypeHierarchyId,
			fth.FolderObjectTypeId,
			ot.ObjectTypeID,
			fth.ParentHierarchyId,
			CASE 
				WHEN ot.GeneralType = 2 THEN CAST(0 AS bit) 
				ELSE CAST(1 AS bit)
			END AS isDocumentType
		    FROM FolderTypeHierarchy fth
		    INNER JOIN ObjectType ot 
			ON ot.ObjectTypeID = fth.FolderObjectTypeId
		   WHERE fth.ParentHierarchyId = @p1`
	rows, err := r.db.Query(query, *fot.ObjectTypeHierarchyId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving list of available types for folder: %w", err)
	}
	for rows.Next() {
		var fot models.FolderObjectTypesNames
		err := rows.Scan(
			&fot.ObjectTypeName,
			&fot.ObjectTypeHierarchyId,
			&fot.FolderObjectTypeId,
			&fot.ObjectTypeID,
			&fot.ParentHierarchyId,
			&fot.IsDocumentType,
		)
		if err != nil {
			return nil, fmt.Errorf("error- scanning folder object type: %w", err)
		}
		folderObjectTypes = append(folderObjectTypes, fot)
	}

	return folderObjectTypes, nil
}
