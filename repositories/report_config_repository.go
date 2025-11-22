package repositories

import (
	"database/sql"
	"enterprise-architect-api/models"
	"fmt"

	"github.com/google/uuid"
)

// ReportConfigRepository handles database operations for report configurations
type ReportConfigRepository struct {
	db *sql.DB
}

func (r *ReportConfigRepository) GetEAObjectTypesAssignedToDimension(param any) (models.AssignObjectTypeToDimentionResponse, error) {
	query := `select ea_tag_id, object_type_id from EA_Tags_Dimentions where object_type_id = @p1`

	var assignObjectTypeToDimentionResponse models.AssignObjectTypeToDimentionResponse = models.AssignObjectTypeToDimentionResponse{}
	err := r.db.QueryRow(query, param).Scan(&assignObjectTypeToDimentionResponse.EAID, &assignObjectTypeToDimentionResponse.ObjectTypeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.AssignObjectTypeToDimentionResponse{}, nil
		}
		return models.AssignObjectTypeToDimentionResponse{}, err
	}

	return assignObjectTypeToDimentionResponse, nil
}

// NewReportConfigRepository creates a new ReportConfigRepository
func NewReportConfigRepository(db *sql.DB) *ReportConfigRepository {
	return &ReportConfigRepository{db: db}
}

// Create creates a new object content in the database

func (r *ReportConfigRepository) DashboardCount(libraryID uuid.UUID) ([]models.DashboardCount, error) {

	query := ` select o.ExactObjectTypeID,COUNT(ObjectID) [count], ot.ObjectTypeName, ot.color, ot.icon from Object o
			inner join objecttype ot on ot.ObjectTypeID = o.ExactObjectTypeID
			where o.LibraryId = @p1
			and o.GeneralType <> dbo.const_GeneralType_Folder()

			group by o.ExactObjectTypeID, ot.ObjectTypeName, ot.color, ot.icon`
	fmt.Println("libraryID:", libraryID)
	//sqlServerUUID := toSQLServerUUID(libraryID)
	//sqlUUID, _ := uuid.FromBytes(sqlServerUUID)
	//fmt.Println("UUID:", sqlUUID.String())

	resultSet, err := r.db.Query(query, &libraryID)

	var dashboardCounts []models.DashboardCount

	if err != nil {
		return nil, err
	}

	for resultSet.Next() {
		var dashboardCount models.DashboardCount
		err := resultSet.Scan(&dashboardCount.ExactObjectTypeID, &dashboardCount.Count,
			&dashboardCount.ObjectTypeName, &dashboardCount.Color, &dashboardCount.Icon)

		if err != nil {
			return nil, fmt.Errorf("error scanning folder: %w", err)
		}
		dashboardCounts = append(dashboardCounts, dashboardCount)
	}
	return dashboardCounts, nil
}

// ========== EA_Tags CRUD Operations ==========

// CreateEATag creates a new EA tag in the database
func (r *ReportConfigRepository) CreateEATag(req models.CreateEATagRequest) (*models.EATag, error) {
	query := `INSERT INTO EA_Tags (name_ar, name_en) VALUES (@p1, @p2); SELECT SCOPE_IDENTITY()`

	var id int
	err := r.db.QueryRow(query, req.NameAr, req.NameEn).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error creating EA tag: %w", err)
	}

	return &models.EATag{
		ID:     id,
		NameAr: req.NameAr,
		NameEn: req.NameEn,
	}, nil
}

// GetEATagByID retrieves an EA tag by its ID
func (r *ReportConfigRepository) GetEATagByID(id int) (*models.EATag, error) {
	query := `SELECT id, name_ar, name_en FROM EA_Tags WHERE id = @p1`

	var tag models.EATag
	err := r.db.QueryRow(query, id).Scan(&tag.ID, &tag.NameAr, &tag.NameEn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("EA tag not found with ID: %d", id)
		}
		return nil, fmt.Errorf("error getting EA tag: %w", err)
	}

	return &tag, nil
}

// GetAllEATags retrieves all EA tags with pagination
func (r *ReportConfigRepository) GetAllEATags(page, pageSize int) ([]models.EATag, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	countQuery := `SELECT COUNT(*) FROM EA_Tags`
	var totalCount int
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting EA tags: %w", err)
	}

	// Get paginated data
	query := `SELECT id, name_ar, name_en FROM EA_Tags ORDER BY id OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY`

	rows, err := r.db.Query(query, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting EA tags: %w", err)
	}
	defer rows.Close()

	var tags []models.EATag
	for rows.Next() {
		var tag models.EATag
		if err := rows.Scan(&tag.ID, &tag.NameAr, &tag.NameEn); err != nil {
			return nil, 0, fmt.Errorf("error scanning EA tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, totalCount, nil
}

// UpdateEATag updates an existing EA tag
func (r *ReportConfigRepository) UpdateEATag(id int, req models.UpdateEATagRequest) (*models.EATag, error) {
	// Check if the tag exists
	existing, err := r.GetEATagByID(id)
	if err != nil {
		return nil, err
	}

	// Build dynamic update query
	query := `UPDATE EA_Tags SET `
	params := []interface{}{}
	paramIndex := 1

	if req.NameAr != nil {
		query += fmt.Sprintf("name_ar = @p%d, ", paramIndex)
		params = append(params, *req.NameAr)
		paramIndex++
		existing.NameAr = *req.NameAr
	}

	if req.NameEn != nil {
		query += fmt.Sprintf("name_en = @p%d, ", paramIndex)
		params = append(params, *req.NameEn)
		paramIndex++
		existing.NameEn = *req.NameEn
	}

	// Remove trailing comma and space
	query = query[:len(query)-2]

	// Add WHERE clause
	query += fmt.Sprintf(" WHERE id = @p%d", paramIndex)
	params = append(params, id)

	_, err = r.db.Exec(query, params...)
	if err != nil {
		return nil, fmt.Errorf("error updating EA tag: %w", err)
	}

	return existing, nil
}

// DeleteEATag deletes an EA tag by its ID
func (r *ReportConfigRepository) DeleteEATag(id int) error {
	query := `DELETE FROM EA_Tags WHERE id = @p1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting EA tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("EA tag not found with ID: %d", id)
	}

	return nil
}

// ========== EA_Tags_Dimentions Operations ==========

// AssignObjectTypeToDimention assigns an object type to a dimension
func (r *ReportConfigRepository) AssignObjectTypeToDimention(req models.AssignObjectTypeToDimentionRequest) (*models.EATagDimention, error) {
	query := `delete from EA_Tags_Dimentions where object_type_id = @p1`
	_, err := r.db.Exec(query, req.ObjectTypeID)
	if err != nil {
		return nil, fmt.Errorf("error deleting object type from dimension: %w", err)
	}

	query = `INSERT INTO EA_Tags_Dimentions (ea_tag_id, object_type_id) VALUES (@p1, @p2); SELECT SCOPE_IDENTITY()`

	var id int
	err = r.db.QueryRow(query, req.EAID, req.ObjectTypeID).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error assigning object type to dimension: %w", err)
	}

	return &models.EATagDimention{
		ID:           id,
		EATagID:      req.EAID,
		ObjectTypeID: req.ObjectTypeID,
	}, nil
}
