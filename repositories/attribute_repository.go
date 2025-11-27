package repositories

import (
	"database/sql"
	"enterprise-architect-api/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AttributeRepository struct {
	db *sql.DB
}

func NewAttributeRepository(db *sql.DB) *AttributeRepository {
	return &AttributeRepository{db: db}
}

func (r *AttributeRepository) GetAttributeForObject(objectID uuid.UUID, objectTypeId *int) (*models.ObjectInstanceAttribute, error) {
	sql := `SELECT attr.AttributeId,
			attr.objectId,
			attr.versionId,
			attr.DataType,
	        attr.textValue,
			attr.booleanValue,
			attr.dateValue,
			attr.floatValue,
			attr.intValue,
			attr.richTextValue,
			att.AttributeName, 
			att.AttributeType,
			att.IsMandatory
			FROM [vwAttributeValue] AS attr 
			INNER JOIN [AttributePermissions] AS attrPerm1 ON  attrPerm1.AttributeId = attr.AttributeId AND attrPerm1.ProfileId = 1 AND attrPerm1.HasRead = 1
			inner join Attribute att on att.AttributeId = attr.AttributeId
			AND attr.objectId = @p1
`
	fmt.Println("UUID:", objectID.String())
	objectID, _ = TransformUUID(objectID)
	fmt.Println("UUID:", objectID.String())
	rows, err := r.db.Query(sql, objectID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var attributes []models.AssignedAttribute
	var objectInstanceAttribute models.ObjectInstanceAttribute
	for rows.Next() {
		var attribute models.AssignedAttribute
		err := rows.Scan(
			&attribute.AttributeID,
			&attribute.ObjectId,
			&attribute.VersionId,
			&attribute.DataType,
			&attribute.TextValue,
			&attribute.BooleanValue,
			&attribute.DateValue,
			&attribute.FloatValue,
			&attribute.IntegerValue,
			&attribute.RichTextValue,
			&attribute.AttributeName,
			&attribute.AttributeType,
			&attribute.IsMandatory,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attribute: %w", err)
		}
		attributes = append(attributes, attribute)
	}
	var attributesRelateds []models.ObjectTypeAssignedAttribute
	if objectTypeId != nil && *objectTypeId > 0 {
		sql = `SELECT a.AttributeId,
            a.AttributeName,
            a.AttributeType,
            a.Description,
            a.tooltipText,
			a.TextDefaultValue,
			a.IntDefaultValue,
			a.BoolDefaultValue,
			a.listType,
			a.ListDefaultValue,
			a.ListValues,
            aa.AttributeGroupId,
            aa.ObjectTypeId, 
            aa.RelationTypeId,
            ot.GeneralType,
            aa.SequenceWithinGroup,
            ag.AttributeGroupName,
            ot.ObjectTypeName
            FROM vwAttribute a
        JOIN AttributeAssigned aa (NOLOCK)
            ON a.AttributeId = aa.AttributeId
        left JOIN AttributeGroup ag ON ag.AttributeGroupId = aa.AttributeGroupId
		  LEFT JOIN ObjectType ot
            ON aa.ObjectTypeId = ot.ObjectTypeId
        WHERE aa.ObjectTypeId = @p1`

		rows, err = r.db.Query(sql, *objectTypeId)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var attributesRelated models.ObjectTypeAssignedAttribute
			err := rows.Scan(
				&attributesRelated.AttributeId,
				&attributesRelated.AttributeName,
				&attributesRelated.AttributeType,
				&attributesRelated.Description,
				&attributesRelated.TooltipText,
				&attributesRelated.TextDefaultValue,
				&attributesRelated.IntDefaultValue,
				&attributesRelated.BoolDefaultValue,
				&attributesRelated.ListType,
				&attributesRelated.ListDefaultValue,
				&attributesRelated.ListValues,
				&attributesRelated.AttributeGroupId,
				&attributesRelated.ObjectTypeId,
				&attributesRelated.RelationTypeId,
				&attributesRelated.GeneralType,
				&attributesRelated.SequenceWithinGroup,
				&attributesRelated.AttributeGroupName,
				&attributesRelated.ObjectTypeName,
			)
			if err != nil {
				return nil, fmt.Errorf("error scanning attribute: %w", err)
			}
			attributesRelateds = append(attributesRelateds, attributesRelated)
		}
	}
	objectInstanceAttribute.AssignedAttribute = attributesRelateds
	objectInstanceAttribute.AssignedAttributesValues = attributes
	objectInstanceAttribute.Success = true
	objectInstanceAttribute.Message = "Attributes retrieved successfully"
	return &objectInstanceAttribute, nil
}
func (r *AttributeRepository) GetRelatedAttributeForObject(objectID uuid.UUID) {

}

// ExistsByName checks if an attribute with the given name already exists
func (r *AttributeRepository) ExistsByName(name string) (bool, error) {
	query := `SELECT COUNT(*) FROM Attribute WHERE AttributeName = @p1`
	var count int
	err := r.db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking attribute name existence: %w", err)
	}
	return count > 0, nil
}

// Create creates a new attribute
func (r *AttributeRepository) Create(attribute *models.Attribute) error {
	query := `
		INSERT INTO Attribute (
			AttributeId, AttributeName, AttributeType, IsMandatory, IsSynchronised,
			VisioSyncName, Description, TooltipText, TextDefaultValue, TextRowCount,
			IntDefaultValue, IntLowerLimit, IntUpperLimit, FloatDefaultValue, FloatLowerLimit,
			FloatUpperLimit, DateDefaultValue, BoolDefaultValue, AutoIdPrefix, AutoIdSuffix,
			AutoIdPadding, AutoIdStartValue, AutoIdNextValue, ListDefaultValue, ListType,
			ListValues, IsCalculated
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10,
			@p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18, @p19, @p20,
			@p21, @p22, @p23, @p24, @p25, @p26, @p27
		)
	`

	_, err := r.db.Exec(query,
		attribute.AttributeId, attribute.AttributeName, attribute.AttributeType,
		attribute.IsMandatory, attribute.IsSynchronised, attribute.VisioSyncName,
		attribute.Description, attribute.TooltipText, attribute.TextDefaultValue,
		attribute.TextRowCount, attribute.IntDefaultValue, attribute.IntLowerLimit,
		attribute.IntUpperLimit, attribute.FloatDefaultValue, attribute.FloatLowerLimit,
		attribute.FloatUpperLimit, attribute.DateDefaultValue, attribute.BoolDefaultValue,
		attribute.AutoIdPrefix, attribute.AutoIdSuffix, attribute.AutoIdPadding,
		attribute.AutoIdStartValue, attribute.AutoIdNextValue, attribute.ListDefaultValue,
		attribute.ListType, attribute.ListValues, attribute.IsCalculated,
	)

	if err != nil {
		return fmt.Errorf("error creating attribute: %w", err)
	}

	return nil
}

// GetByID retrieves an attribute by its ID
func (r *AttributeRepository) GetByID(id string) (*models.Attribute, error) {
	query := `
		SELECT AttributeId, AttributeName, AttributeType, IsMandatory, IsSynchronised,
			VisioSyncName, Description, TooltipText, TextDefaultValue, TextRowCount,
			IntDefaultValue, IntLowerLimit, IntUpperLimit, FloatDefaultValue, FloatLowerLimit,
			FloatUpperLimit, DateDefaultValue, BoolDefaultValue, AutoIdPrefix, AutoIdSuffix,
			AutoIdPadding, AutoIdStartValue, AutoIdNextValue, ListDefaultValue, ListType,
			ListValues, IsCalculated
		FROM Attribute
		WHERE AttributeId = @p1
	`

	attribute := &models.Attribute{}
	err := r.db.QueryRow(query, id).Scan(
		&attribute.AttributeId, &attribute.AttributeName, &attribute.AttributeType,
		&attribute.IsMandatory, &attribute.IsSynchronised, &attribute.VisioSyncName,
		&attribute.Description, &attribute.TooltipText, &attribute.TextDefaultValue,
		&attribute.TextRowCount, &attribute.IntDefaultValue, &attribute.IntLowerLimit,
		&attribute.IntUpperLimit, &attribute.FloatDefaultValue, &attribute.FloatLowerLimit,
		&attribute.FloatUpperLimit, &attribute.DateDefaultValue, &attribute.BoolDefaultValue,
		&attribute.AutoIdPrefix, &attribute.AutoIdSuffix, &attribute.AutoIdPadding,
		&attribute.AutoIdStartValue, &attribute.AutoIdNextValue, &attribute.ListDefaultValue,
		&attribute.ListType, &attribute.ListValues, &attribute.IsCalculated,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("attribute not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving attribute: %w", err)
	}

	return attribute, nil
}

// GetAll retrieves all attributes with pagination
func (r *AttributeRepository) GetAll(page, pageSize int) ([]models.Attribute, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM Attribute`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting attributes: %w", err)
	}

	// Get paginated results
	query := `
		SELECT AttributeId, AttributeName, AttributeType, IsMandatory, IsSynchronised,
			VisioSyncName, Description, TooltipText, TextDefaultValue, TextRowCount,
			IntDefaultValue, IntLowerLimit, IntUpperLimit, FloatDefaultValue, FloatLowerLimit,
			FloatUpperLimit, DateDefaultValue, BoolDefaultValue, AutoIdPrefix, AutoIdSuffix,
			AutoIdPadding, AutoIdStartValue, AutoIdNextValue, ListDefaultValue, ListType,
			ListValues, IsCalculated
		FROM Attribute
		ORDER BY AttributeName
		OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
	`

	rows, err := r.db.Query(query, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving attributes: %w", err)
	}
	defer rows.Close()

	var attributes []models.Attribute
	for rows.Next() {
		var attribute models.Attribute
		err := rows.Scan(
			&attribute.AttributeId, &attribute.AttributeName, &attribute.AttributeType,
			&attribute.IsMandatory, &attribute.IsSynchronised, &attribute.VisioSyncName,
			&attribute.Description, &attribute.TooltipText, &attribute.TextDefaultValue,
			&attribute.TextRowCount, &attribute.IntDefaultValue, &attribute.IntLowerLimit,
			&attribute.IntUpperLimit, &attribute.FloatDefaultValue, &attribute.FloatLowerLimit,
			&attribute.FloatUpperLimit, &attribute.DateDefaultValue, &attribute.BoolDefaultValue,
			&attribute.AutoIdPrefix, &attribute.AutoIdSuffix, &attribute.AutoIdPadding,
			&attribute.AutoIdStartValue, &attribute.AutoIdNextValue, &attribute.ListDefaultValue,
			&attribute.ListType, &attribute.ListValues, &attribute.IsCalculated,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning attribute: %w", err)
		}
		attribute.AttributeId, _ = TransformUUID(attribute.AttributeId)
		attributes = append(attributes, attribute)
	}

	return attributes, totalCount, nil
}

// Update updates an existing attribute
func (r *AttributeRepository) Update(id string, attribute *models.Attribute) error {
	// Build dynamic update query
	var setClauses []string
	var args []interface{}
	argIndex := 1

	setClauses = append(setClauses, fmt.Sprintf("AttributeName = @p%d", argIndex))
	args = append(args, attribute.AttributeName)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("AttributeType = @p%d", argIndex))
	args = append(args, attribute.AttributeType)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("IsMandatory = @p%d", argIndex))
	args = append(args, attribute.IsMandatory)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("IsSynchronised = @p%d", argIndex))
	args = append(args, attribute.IsSynchronised)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("VisioSyncName = @p%d", argIndex))
	args = append(args, attribute.VisioSyncName)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("Description = @p%d", argIndex))
	args = append(args, attribute.Description)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("TooltipText = @p%d", argIndex))
	args = append(args, attribute.TooltipText)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("TextDefaultValue = @p%d", argIndex))
	args = append(args, attribute.TextDefaultValue)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("TextRowCount = @p%d", argIndex))
	args = append(args, attribute.TextRowCount)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("IntDefaultValue = @p%d", argIndex))
	args = append(args, attribute.IntDefaultValue)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("IntLowerLimit = @p%d", argIndex))
	args = append(args, attribute.IntLowerLimit)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("IntUpperLimit = @p%d", argIndex))
	args = append(args, attribute.IntUpperLimit)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("FloatDefaultValue = @p%d", argIndex))
	args = append(args, attribute.FloatDefaultValue)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("FloatLowerLimit = @p%d", argIndex))
	args = append(args, attribute.FloatLowerLimit)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("FloatUpperLimit = @p%d", argIndex))
	args = append(args, attribute.FloatUpperLimit)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("DateDefaultValue = @p%d", argIndex))
	args = append(args, attribute.DateDefaultValue)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("BoolDefaultValue = @p%d", argIndex))
	args = append(args, attribute.BoolDefaultValue)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("AutoIdPrefix = @p%d", argIndex))
	args = append(args, attribute.AutoIdPrefix)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("AutoIdSuffix = @p%d", argIndex))
	args = append(args, attribute.AutoIdSuffix)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("AutoIdPadding = @p%d", argIndex))
	args = append(args, attribute.AutoIdPadding)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("AutoIdStartValue = @p%d", argIndex))
	args = append(args, attribute.AutoIdStartValue)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("AutoIdNextValue = @p%d", argIndex))
	args = append(args, attribute.AutoIdNextValue)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("ListDefaultValue = @p%d", argIndex))
	args = append(args, attribute.ListDefaultValue)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("ListType = @p%d", argIndex))
	args = append(args, attribute.ListType)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("ListValues = @p%d", argIndex))
	args = append(args, attribute.ListValues)
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("IsCalculated = @p%d", argIndex))
	args = append(args, attribute.IsCalculated)
	argIndex++

	args = append(args, id)
	query := fmt.Sprintf("UPDATE Attribute SET %s WHERE AttributeId = @p%d", strings.Join(setClauses, ", "), argIndex)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error updating attribute: %w", err)
	}

	return nil
}

// Delete deletes an attribute by its ID
func (r *AttributeRepository) Delete(id string) error {
	query := `DELETE FROM Attribute WHERE AttributeId = @p1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting attribute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("attribute not found")
	}

	return nil
}

// AssignAttributeToObjectType assigns an attribute to an object type
func (r *AttributeRepository) AssignAttributeToObjectType(req *models.AssignAttributeToObjectTypeRequest) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	emptyGuid := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	attributeGroupId := req.AttributeGroupId

	// Check if AttributeGroup exists, if not create it
	checkGroupQuery := `
		SELECT ag.AttributeGroupId 
		FROM dbo.AttributeGroup AS ag (NOLOCK) 
		INNER JOIN dbo.AttributeGroupAssigned AS aga ON aga.AttributeGroupId = ag.AttributeGroupId 
		WHERE ag.AttributeGroupName = @p1 
		AND (aga.ObjectTypeId > 0 AND aga.ObjectTypeId = @p2) OR (aga.RelationTypeId <>  dbo.const_GuidEmpty() AND aga.RelationTypeId = @p3)
	`

	var existingGroupId uuid.UUID
	if attributeGroupId != uuid.MustParse("00000000-0000-0000-0000-000000000001") {
		err = tx.QueryRow(checkGroupQuery, req.AttributeGroupName, req.ObjectTypeId, req.RelationTypeId).Scan(&existingGroupId)
	} else {
		err = sql.ErrNoRows
	}
	if err == sql.ErrNoRows {
		// Insert new AttributeGroup
		insertGroupQuery := `
		INSERT INTO dbo.AttributeGroup (AttributeGroupId, AttributeGroupName)
		OUTPUT INSERTED.AttributeGroupId
		VALUES (NEWID(), @p1)`

		err = tx.QueryRow(insertGroupQuery, req.AttributeGroupName).Scan(&attributeGroupId)
		if err != nil {
			return fmt.Errorf("error inserting attribute group: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking attribute group: %w", err)
	} else {
		// Use existing group ID
		attributeGroupId = existingGroupId
	}

	// Check if AttributeGroupAssigned exists
	checkGroupAssignedQuery := `
		SELECT AttributeGroupId 
		FROM dbo.AttributeGroupAssigned (NOLOCK) 
		WHERE AttributeGroupId = @p1 
		AND ((@p2 > 0 AND ObjectTypeId = @p2) OR (@p3 <> @p4 AND RelationTypeId = @p3)) 
		AND AttributeGroupId <> @p5
	`

	var existingAssignment uuid.UUID
	err = tx.QueryRow(checkGroupAssignedQuery, attributeGroupId, req.ObjectTypeId, req.RelationTypeId, emptyGuid, req.AttributeGroupId).Scan(&existingAssignment)

	if err == sql.ErrNoRows {
		// Get max GroupSequence
		var groupSequence sql.NullInt32
		getSequenceQuery := `
			SELECT MAX(GroupSequence) 
			FROM dbo.AttributeGroupAssigned (NOLOCK) 
			WHERE ((@p1 > 0 AND ObjectTypeId = @p1) OR (@p2 <> @p3 AND RelationTypeId = @p2)) 
			AND AttributeGroupId <> @p4
		`
		err = tx.QueryRow(getSequenceQuery, req.ObjectTypeId, req.RelationTypeId, emptyGuid, req.AttributeGroupId).Scan(&groupSequence)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error getting group sequence: %w", err)
		}

		sequence := 1
		if groupSequence.Valid {
			sequence = int(groupSequence.Int32) + 1
		}

		// Insert AttributeGroupAssigned
		insertGroupAssignedQuery := `
			INSERT INTO dbo.AttributeGroupAssigned (ObjectTypeId, RelationTypeId, AttributeGroupId, GroupSequence)
			VALUES (@p1, @p2, @p3, @p4)
		`
		_, err = tx.Exec(insertGroupAssignedQuery, req.ObjectTypeId, req.RelationTypeId, attributeGroupId, sequence)
		if err != nil {
			return fmt.Errorf("error inserting attribute group assigned: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking attribute group assigned: %w", err)
	}

	// Get max SequenceWithinGroup
	var sequenceWithinGroup sql.NullInt32
	getAttributeSequenceQuery := `
		SELECT MAX(SequenceWithinGroup) 
		FROM dbo.AttributeAssigned (NOLOCK) 
		WHERE AttributeGroupId = @p1 
		AND ((@p2 > 0 AND ObjectTypeId = @p2) OR (@p3 <> @p4 AND RelationTypeId = @p3))
		AND AttributeGroupId <> @p5
	`
	err = tx.QueryRow(getAttributeSequenceQuery, attributeGroupId, req.ObjectTypeId, req.RelationTypeId, emptyGuid, req.AttributeGroupId).Scan(&sequenceWithinGroup)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error getting attribute sequence: %w", err)
	}

	attrSequence := 1
	if sequenceWithinGroup.Valid {
		attrSequence = int(sequenceWithinGroup.Int32) + 1
	}
	fmt.Println("attributeGroupId : ", attributeGroupId)
	// Insert AttributeAssigned
	insertAttributeAssignedQuery := `
		INSERT INTO dbo.AttributeAssigned (ObjectTypeId, RelationTypeId, AttributeId, AttributeGroupId, SequenceWithinGroup)
		VALUES (@p1, @p2, @p3, @p4, @p5)
	`
	attributeGroupId, _ = TransformUUID(attributeGroupId)

	fmt.Println("attributeGroupId v4: ", attributeGroupId)

	_, err = tx.Exec(insertAttributeAssignedQuery, req.ObjectTypeId, req.RelationTypeId, req.AttributeId, attributeGroupId, attrSequence)
	if err != nil {
		return fmt.Errorf("error inserting attribute assigned: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
func (r *AttributeRepository) GetAttributeAssignments(objectTypeId int, relationTypeId uuid.UUID) ([]models.AttributeAssignment, error) {
	query := `
        SELECT a.AttributeId,
            a.AttributeName,
            a.AttributeType,
            a.Description,
            a.tooltipText,
            aa.AttributeGroupId,
            aa.ObjectTypeId, 
            aa.RelationTypeId,
            ot.GeneralType,
            aa.SequenceWithinGroup,
            ag.AttributeGroupName,
            ot.ObjectTypeName
        FROM vwAttribute a
        JOIN AttributeAssigned aa (NOLOCK)
            ON a.AttributeId = aa.AttributeId
        LEFT JOIN ObjectType ot
            ON aa.ObjectTypeId = ot.ObjectTypeId
        LEFT JOIN AttributeGroup ag ON ag.AttributeGroupId = aa.AttributeGroupId
        WHERE aa.ObjectTypeId = @p1
    `

	// Add relationTypeId condition if not empty
	var rows *sql.Rows
	var err error

	if relationTypeId != uuid.Nil {
		query += " AND aa.RelationTypeId = @p2"
		rows, err = r.db.Query(query, objectTypeId, relationTypeId)
	} else {
		rows, err = r.db.Query(query, objectTypeId)
	}

	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var assignments []models.AttributeAssignment
	for rows.Next() {
		var assignment models.AttributeAssignment
		err := rows.Scan(
			&assignment.AttributeId,
			&assignment.AttributeName,
			&assignment.AttributeType,
			&assignment.Description,
			&assignment.TooltipText,
			&assignment.AttributeGroupId,
			&assignment.ObjectTypeId,
			&assignment.RelationTypeId,
			&assignment.GeneralType,
			&assignment.SequenceWithinGroup,
			&assignment.AttributeGroupName,
			&assignment.ObjectTypeName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attribute assignment: %w", err)
		}
		assignment.AttributeGroupId, _ = TransformUUID(assignment.AttributeGroupId)
		assignment.RelationTypeId, _ = TransformUUID(assignment.RelationTypeId)
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

// UnassignAttributeFromObjectType removes an attribute assignment from an object type
func (r *AttributeRepository) UnassignAttributeFromObjectType(req *models.UnassignAttributeFromObjectTypeRequest) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	//emptyGuid := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	fmt.Println("req.AttributeId : ", req.AttributeId)
	fmt.Println("req.AttributeGroupId : ", req.AttributeGroupId)
	fmt.Println("req.ObjectTypeId : ", req.ObjectTypeId)
	fmt.Println("req.RelationTypeId : ", req.RelationTypeId)
	req.AttributeId, _ = TransformUUID(req.AttributeId)
	req.AttributeGroupId, _ = TransformUUID(req.AttributeGroupId)
	// Delete from AttributeAssigned
	deleteAttributeAssignedQuery := `
		DELETE FROM dbo.AttributeAssigned 
		WHERE AttributeId = @p1 
		AND AttributeGroupId = @p2
		AND ((@p3 > 0 AND ObjectTypeId = @p3))
	`

	result, err := tx.Exec(deleteAttributeAssignedQuery, req.AttributeId, req.AttributeGroupId, req.ObjectTypeId)
	if err != nil {
		return fmt.Errorf("error deleting attribute assignment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("attribute assignment not found")
	}

	// Check if AttributeGroup still has any attributes assigned for this object type
	checkGroupHasAttributesQuery := `
		SELECT COUNT(*) 
		FROM dbo.AttributeAssigned (NOLOCK) 
		WHERE AttributeGroupId = @p1
		AND ((@p2 > 0 AND ObjectTypeId = @p2))
	`
	//OR (@p3 <> @p4 AND RelationTypeId = @p3))
	var attributeCount int
	err = tx.QueryRow(checkGroupHasAttributesQuery, req.AttributeGroupId, req.ObjectTypeId).Scan(&attributeCount)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking attribute group: %w", err)
	}

	// If no attributes left in the group for this object type, delete AttributeGroupAssigned
	if attributeCount == 0 {
		// Delete AttributeGroupAssigned
		deleteGroupAssignedQuery := `
			DELETE FROM dbo.AttributeGroupAssigned 
			WHERE AttributeGroupId = @p1 
			AND ((@p2 > 0 AND ObjectTypeId = @p2)) 
		`
		_, err = tx.Exec(deleteGroupAssignedQuery, req.AttributeGroupId, req.ObjectTypeId)
		if err != nil {
			return fmt.Errorf("error deleting attribute group assignment: %w", err)
		}

		// Delete AttributeGroup if not used elsewhere
		deleteGroupQuery := `
			DELETE FROM dbo.AttributeGroup 
			WHERE AttributeGroupId = @p1 
			AND NOT EXISTS (
				SELECT 1 FROM dbo.AttributeGroupAssigned 
				WHERE AttributeGroupId = @p1
			)
		`
		_, err = tx.Exec(deleteGroupQuery, req.AttributeGroupId)
		if err != nil {
			return fmt.Errorf("error deleting attribute group: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// UpdateAttributeValue updates the values of multiple attributes
func (r *AttributeRepository) UpdateAttributeValue(attrs []models.AssignedAttribute) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE AttributeValue
		SET 
			ValueText =
				CASE WHEN DataType = 4 THEN @p1 ELSE ValueText END,

			ValueBigInt =
				CASE 
					WHEN DataType = 1 THEN @p2      -- int
					WHEN DataType = 5 THEN CASE WHEN @p3 = 1 THEN 1 ELSE 0 END  -- boolean
					ELSE ValueBigInt
				END,

			ValueFloat =
				CASE WHEN DataType = 3 THEN @p4 ELSE ValueFloat END,

			ValueDate =
				CASE WHEN DataType = 2 THEN @p5 ELSE ValueDate END,

			ValueRichText =
				CASE WHEN DataType = 6 THEN @p6 ELSE ValueRichText END
		WHERE 
			AttributeId = @p7
			AND ObjectId = @p8
			AND VersionId = @p9
	`
	fmt.Println("attrs : ", *attrs[0].TextValue)
	for _, attr := range attrs {
		// Transform UUIDs
		attributeID, _ := TransformUUID(attr.AttributeID)
		objectID, _ := TransformUUID(attr.ObjectId)
		versionID, _ := TransformUUID(attr.VersionId)

		var boolVal interface{}
		if attr.BooleanValue != nil {
			if *attr.BooleanValue {
				boolVal = 1
			} else {
				boolVal = 0
			}
		}
		var textValue string
		if attr.TextValue != nil {
			textValue = *attr.TextValue
		}
		var richTextValue string
		if attr.RichTextValue != nil {
			richTextValue = *attr.RichTextValue
		}
		var integerValue int64
		if attr.IntegerValue != nil {
			integerValue = int64(*attr.IntegerValue)
		}
		var floatValue float64
		if attr.FloatValue != nil {
			floatValue = *attr.FloatValue
		}
		var dateValue time.Time
		if attr.DateValue != nil {
			dateValue = *attr.DateValue
		}

		_, err := tx.Exec(query,
			textValue,
			integerValue,
			boolVal,
			floatValue,
			dateValue,
			richTextValue,
			attributeID,
			objectID,
			versionID,
		)

		if err != nil {
			return fmt.Errorf("error updating attribute value for AttributeId %s: %w", attr.AttributeID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
