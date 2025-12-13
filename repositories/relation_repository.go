package repositories

import (
	"database/sql"
	"enterprise-architect-api/models"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type RelationRepository struct {
	db *sql.DB
}

func NewRelationRepository(db *sql.DB) *RelationRepository {
	return &RelationRepository{db: db}
}

// Create inserts a new relationship into the database
func (r *RelationRepository) Create(req models.CreateRelationRequest) (*models.Relation, error) {
	id := uuid.New()
	now := time.Now()

	// Transform UUIDs to SQL Server format
	relTypeID, _ := TransformUUID(req.RelationTypeId)
	fromID, _ := TransformUUID(req.FromObjectId)
	toID, _ := TransformUUID(req.ToObjectID)

	query := `
		INSERT INTO [Relation] (
			RelationshipId, RelationTypeId, RelationReason, FromObjectId, ToObjectId,
			DateCreated, CreatedBy, DateModified, ModifiedBy, RichTextDescription
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10
		)`

	var richTextDescriptionEncoded *string
	if req.RichTextDescription != nil {
		encoded := encodeRTFString(*req.RichTextDescription)
		richTextDescriptionEncoded = &encoded
	}

	_, err := r.db.Exec(query,
		id, relTypeID, req.RelationReason, fromID, toID,
		now, req.CreatedBy, now, req.CreatedBy, richTextDescriptionEncoded)

	if err != nil {
		return nil, fmt.Errorf("error creating relationship: %w", err)
	}

	return &models.Relation{
		RelationshipID:      id,
		RelationTypeID:      req.RelationTypeId,
		RelationReason:      req.RelationReason,
		FromObjectID:        req.FromObjectId,
		ToObjectID:          req.ToObjectID,
		DateCreated:         now,
		CreatedBy:           req.CreatedBy,
		DateModified:        now,
		ModifiedBy:          req.CreatedBy,
		RichTextDescription: req.RichTextDescription,
	}, nil
}

// GetByObjectID retrieves all relationships for a given object ID
// It includes details about the relationship type and the other object involved
func (r *RelationRepository) GetByObjectID(objectID uuid.UUID) ([]models.RelationWithDetails, error) {
	searchID, _ := TransformUUID(objectID)

	query := `
		SELECT 
			r.RelationshipId, r.RelationTypeId, r.RelationReason, 
			r.FromObjectId, r.ToObjectId, r.DateCreated, r.CreatedBy, 
			r.DateModified, r.ModifiedBy, r.RichTextDescription,
			rt.RelationTypeName,
			other.ObjectID, other.ObjectName, rt.FromToDescription, rt.ToFromDescription,
			CASE WHEN r.FromObjectId = @p1 THEN 'To' ELSE 'From' END as Direction
		FROM [Relation] r
		INNER JOIN [RelationType] rt ON r.RelationTypeId = rt.RelationTypeId
		INNER JOIN [Object] other ON other.ObjectID = CASE WHEN r.FromObjectId = @p1 THEN r.ToObjectId ELSE r.FromObjectId END
		WHERE r.FromObjectId = @p1 OR r.ToObjectId = @p1
		ORDER BY r.DateCreated DESC
	`

	rows, err := r.db.Query(query, searchID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving relationships: %w", err)
	}
	defer rows.Close()

	var relations []models.RelationWithDetails
	for rows.Next() {
		var rel models.RelationWithDetails
		var relIDBytes, typeIDBytes, fromIDBytes, toIDBytes, otherIDBytes []byte

		err := rows.Scan(
			&relIDBytes, &typeIDBytes, &rel.RelationReason,
			&fromIDBytes, &toIDBytes, &rel.DateCreated, &rel.CreatedBy,
			&rel.DateModified, &rel.ModifiedBy, &rel.RichTextDescription,
			&rel.RelationTypeName,
			&otherIDBytes, &rel.OtherObjectName,
			&rel.Direction,
			&rel.FromToDescription,
			&rel.ToFromDescription,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning relationship: %w", err)
		}

		// Parse UUIDs
		rel.RelationshipID, _ = parseSQLServerUUID(relIDBytes)
		rel.RelationTypeID, _ = parseSQLServerUUID(typeIDBytes)
		rel.FromObjectID, _ = parseSQLServerUUID(fromIDBytes)
		rel.ToObjectID, _ = parseSQLServerUUID(toIDBytes)
		rel.OtherObjectID, _ = parseSQLServerUUID(otherIDBytes)

		if rel.RichTextDescription != nil {
			decoded := decodeRTFString(*rel.RichTextDescription)
			rel.RichTextDescription = &decoded
		}

		relations = append(relations, rel)
	}

	return relations, nil
}

// GetAllRelationTypes retrieves all relation types with pagination
func (r *RelationRepository) GetAllRelationTypes(page, pageSize int) ([]models.RelationType, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM [RelationType]`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting relation types: %w", err)
	}

	query := `
		SELECT 
			RelationTypeId, RelationTypeName, RelationTypeDescription, 
			FromToDescription, ToFromDescription, isDirectionless, 
			DateCreated, CreatedBy, DateModified, ModifiedBy, 
			isHierarchical, usage
		FROM [RelationType]
		ORDER BY RelationTypeName
		OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
	`

	rows, err := r.db.Query(query, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving relation types: %w", err)
	}
	defer rows.Close()

	var types []models.RelationType
	for rows.Next() {
		var t models.RelationType
		var typeIDBytes []byte

		err := rows.Scan(
			&typeIDBytes, &t.RelationTypeName, &t.RelationTypeDescription,
			&t.FromToDescription, &t.ToFromDescription, &t.IsDirectionless,
			&t.DateCreated, &t.CreatedBy, &t.DateModified, &t.ModifiedBy,
			&t.IsHierarchical, &t.Usage,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning relation type: %w", err)
		}

		t.RelationTypeID, _ = parseSQLServerUUID(typeIDBytes)
		types = append(types, t)
	}

	return types, totalCount, nil
}

func encodeRTFString(s string) string {
	rtf := ""
	for _, r := range s {
		rtf += fmt.Sprintf("\\u%d?", r)
	}
	return rtf
}

func decodeRTFString(s string) string {
	re := regexp.MustCompile(`\\u(\d+)\?`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		val, err := strconv.Atoi(match[2 : len(match)-1])
		if err != nil {
			return match
		}
		return string(rune(val))
	})
}
