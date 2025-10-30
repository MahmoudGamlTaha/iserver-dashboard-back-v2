package repositories

import (
	"database/sql"
	"enterprise-architect-api/models"
	"fmt"
	"strings"
	"time"
)

// ObjectContentRepository handles database operations for object contents
type ObjectContentRepository struct {
	db *sql.DB
}

// NewObjectContentRepository creates a new ObjectContentRepository
func NewObjectContentRepository(db *sql.DB) *ObjectContentRepository {
	return &ObjectContentRepository{db: db}
}

// Create creates a new object content in the database
func (r *ObjectContentRepository) Create(req models.CreateObjectContentRequest) (*models.ObjectContent, error) {
	now := time.Now()

	query := `
		INSERT INTO ObjectContents (
			DocumentObjectID, ContainerVersionID, ObjectID, Instances, IsShortCut, 
			ContainmentType, DateCreated, CreatedBy, DateModified, ModifiedBy
		) OUTPUT INSERTED.ID VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10
		)
	`

	var id int
	err := r.db.QueryRow(query,
		req.DocumentObjectID, req.ContainerVersionID, req.ObjectID, req.Instances, req.IsShortCut,
		req.ContainmentType, now, req.CreatedBy, now, req.CreatedBy,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("error creating object content: %w", err)
	}

	return r.GetByID(id)
}

// GetByID retrieves an object content by its ID
func (r *ObjectContentRepository) GetByID(id int) (*models.ObjectContent, error) {
	query := `
		SELECT ID, DocumentObjectID, ContainerVersionID, ObjectID, Instances, IsShortCut, 
			ShapeSheetKeysRequiringUpdateId, ContainmentType, DateCreated, CreatedBy, DateModified, ModifiedBy
		FROM ObjectContents
		WHERE ID = @p1
	`

	objContent := &models.ObjectContent{}
	err := r.db.QueryRow(query, id).Scan(
		&objContent.ID, &objContent.DocumentObjectID, &objContent.ContainerVersionID, &objContent.ObjectID,
		&objContent.Instances, &objContent.IsShortCut, &objContent.ShapeSheetKeysRequiringUpdateId,
		&objContent.ContainmentType, &objContent.DateCreated, &objContent.CreatedBy, &objContent.DateModified,
		&objContent.ModifiedBy,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("object content not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving object content: %w", err)
	}

	return objContent, nil
}

// GetAll retrieves all object contents with pagination
func (r *ObjectContentRepository) GetAll(page, pageSize int) ([]models.ObjectContent, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM ObjectContents`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting object contents: %w", err)
	}

	// Get paginated results
	query := `
		SELECT ID, DocumentObjectID, ContainerVersionID, ObjectID, Instances, IsShortCut, 
			ShapeSheetKeysRequiringUpdateId, ContainmentType, DateCreated, CreatedBy, DateModified, ModifiedBy
		FROM ObjectContents
		ORDER BY DateCreated DESC
		OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
	`

	rows, err := r.db.Query(query, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving object contents: %w", err)
	}
	defer rows.Close()

	var objectContents []models.ObjectContent
	for rows.Next() {
		var objContent models.ObjectContent
		err := rows.Scan(
			&objContent.ID, &objContent.DocumentObjectID, &objContent.ContainerVersionID, &objContent.ObjectID,
			&objContent.Instances, &objContent.IsShortCut, &objContent.ShapeSheetKeysRequiringUpdateId,
			&objContent.ContainmentType, &objContent.DateCreated, &objContent.CreatedBy, &objContent.DateModified,
			&objContent.ModifiedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning object content: %w", err)
		}
		objectContents = append(objectContents, objContent)
	}

	return objectContents, totalCount, nil
}

// Update updates an existing object content
func (r *ObjectContentRepository) Update(id int, req models.UpdateObjectContentRequest) (*models.ObjectContent, error) {
	// Build dynamic update query
	var setClauses []string
	var args []interface{}
	argIndex := 1

	if req.DocumentObjectID != nil {
		setClauses = append(setClauses, fmt.Sprintf("DocumentObjectID = @p%d", argIndex))
		args = append(args, *req.DocumentObjectID)
		argIndex++
	}
	if req.ContainerVersionID != nil {
		setClauses = append(setClauses, fmt.Sprintf("ContainerVersionID = @p%d", argIndex))
		args = append(args, *req.ContainerVersionID)
		argIndex++
	}
	if req.ObjectID != nil {
		setClauses = append(setClauses, fmt.Sprintf("ObjectID = @p%d", argIndex))
		args = append(args, *req.ObjectID)
		argIndex++
	}
	if req.Instances != nil {
		setClauses = append(setClauses, fmt.Sprintf("Instances = @p%d", argIndex))
		args = append(args, *req.Instances)
		argIndex++
	}
	if req.IsShortCut != nil {
		setClauses = append(setClauses, fmt.Sprintf("IsShortCut = @p%d", argIndex))
		args = append(args, *req.IsShortCut)
		argIndex++
	}
	if req.ContainmentType != nil {
		setClauses = append(setClauses, fmt.Sprintf("ContainmentType = @p%d", argIndex))
		args = append(args, *req.ContainmentType)
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
	query := fmt.Sprintf("UPDATE ObjectContents SET %s WHERE ID = @p%d", strings.Join(setClauses, ", "), argIndex)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating object content: %w", err)
	}

	return r.GetByID(id)
}

// Delete deletes an object content by its ID
func (r *ObjectContentRepository) Delete(id int) error {
	query := `DELETE FROM ObjectContents WHERE ID = @p1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting object content: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("object content not found")
	}

	return nil
}

