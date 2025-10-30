package repositories

import (
	"database/sql"
	"enterprise-architect-api/models"
	"fmt"
	"strings"
	"time"
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

