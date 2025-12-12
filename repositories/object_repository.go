package repositories

import (
	"database/sql"
	"enterprise-architect-api/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
)

// ObjectRepository handles database operations for objects
type ObjectRepository struct {
	db                  *sql.DB
	attributeRepository *AttributeRepository
}

// NewObjectRepository creates a new ObjectRepository
func NewObjectRepository(db *sql.DB, attrRepo *AttributeRepository) *ObjectRepository {
	return &ObjectRepository{db: db, attributeRepository: attrRepo}
}

func (r *ObjectRepository) ImportObjects(req models.ObjectImportRequest) (*models.ObjectImportResponse, error) {
	folderID, _ := TransformUUID(req.FolderId)
	libraryID, _ := TransformUUID(req.LibraryId)

	fmt.Println("transform: folder id", folderID)
	fmt.Println("transform: library id", libraryID)
	//checks
	checkFolderSql := `
	  SELECT ObjectTypeId,LibraryId 
	  FROM Object WHERE ObjectID = @p1
	`
	var objectTypeId *int64
	var libraryId uuid.UUID
	err := r.db.QueryRow(checkFolderSql, folderID).Scan(
		&objectTypeId,
		&libraryId,
	)
	fmt.Println("transform: library id after", libraryID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Check if object exists query
	checkExistsSql := `
		SELECT ObjectID,CurrentVersionId FROM [Object] 
		WHERE ObjectName = @p1 AND ExactObjectTypeID = @p2 and libraryid = @p3
	`

	// Update existing object query
	updateSql := `
		UPDATE [Object] SET 
			ObjectDescription = @p1, 
			ObjectTypeID = @p2, 
			Locked = @p3, 
			IsImported = @p4, 
			IsLibrary = @p5, 
			LibraryId = @p6, 
			FileExtension = @p7, 
			Prefix = @p8, 
			Suffix = @p9, 
			DateModified = @p10, 
			ModifiedBy = @p11, 
			RichTextDescription = @p12, 
			GeneralType = @p13, 
			HasVisioAlias = @p14
		WHERE ObjectID = @p15
	`

	var insertedObjectCount, insertedFailedObjectCount int
	insertedObjectCount = 0
	insertedFailedObjectCount = 0
	var attrs []models.AssignedAttribute
	tx, _ := r.db.Begin()
	for _, data := range req.Data {
		var objectName string
		var description string
		var objectId uuid.UUID
		var versionId uuid.UUID
		var currentAttrs []models.AssignedAttribute

		// First pass: collect data
		for key, entry := range data {
			if key == "Object Name" {
				if entry.AttributeValue != nil {
					objectName = *entry.AttributeValue
				}
			} else if key == "Description" {
				if entry.AttributeValue != nil {
					description = *entry.AttributeValue
				}
			} else {
				// Attribute processing
				var row models.AssignedAttribute
				if entry.AttributeId != nil {
					var attrUUID, _ = uuid.Parse(*entry.AttributeId)
					fmt.Print("attribute_id", attrUUID)
					row.AttributeID = attrUUID
				}
				if entry.AttributeName != nil {
					row.AttributeName = *entry.AttributeName
				}
				if entry.AttributeType != nil {
					row.AttributeType = *entry.AttributeType
				}
				if r.GetTypeId(row.AttributeType) == 4 {
					row.TextValue = entry.AttributeValue
				}
				if r.GetTypeId(row.AttributeType) == 1 {
					if entry.AttributeValue != nil {
						val, err := strconv.Atoi(*entry.AttributeValue)
						if err == nil {
							row.IntegerValue = &val
						}
					}
				}
				if row.AttributeID != uuid.Nil {
					currentAttrs = append(currentAttrs, row)
				}
			}
		}

		if objectName == "" {
			insertedFailedObjectCount++
			continue
		}

		// Check if object exists with the same ObjectName and ExactObjectTypeID
		var existingObjectId uuid.UUID
		var existingVersionId uuid.UUID
		err = tx.QueryRow(checkExistsSql, objectName, req.ObjectTypeId, libraryID).Scan(&existingObjectId, &existingVersionId)

		if err == sql.ErrNoRows {
			// Object doesn't exist, insert new one using CreateV2
			genType := int(r.GetTypeId("string"))
			createReq := models.CreateObjectRequest{
				ObjectName:          objectName,
				ObjectDescription:   description,
				ObjectTypeID:        int(r.GetTypeId("string")),
				ExactObjectTypeID:   req.ObjectTypeId,
				RichTextDescription: r.toRTFUnicode(description),
				IsLibrary:           false,
				IsImported:          true,
				LibraryId:           &libraryID,
				DirectParentId:      &folderID,
				CreatedBy:           62,
				GeneralType:         &genType,
			}

			createdObj, err := r.CreateV2(tx, createReq)
			if err != nil {
				insertedFailedObjectCount++
				fmt.Println("error insert Object", err)
				continue
			}
			objectId = createdObj.ObjectID
			versionId = *createdObj.CurrentVersionId
		} else if err != nil {
			// Error checking for existence
			insertedFailedObjectCount++
			fmt.Println("error checking object existence", err)
			continue
		} else {
			// Object exists, update it
			_, err = tx.Exec(updateSql, description, r.GetTypeId("string"), 0, 1, 0, libraryID, "", nil, nil,
				time.Now(), 62, r.toRTFUnicode(description), r.GetTypeId("string"), 0, existingObjectId)

			if err != nil {
				insertedFailedObjectCount++
				fmt.Println("error updating Object", err)
				continue
			}
			objectId = existingObjectId
			versionId = existingVersionId
		}

		// Assign IDs to attributes and add to main list
		for _, row := range currentAttrs {
			row.ObjectId = objectId
			row.VersionId = versionId
			fmt.Print("row data", row)
			attrs = append(attrs, row)
		}
		insertedObjectCount++
	}
	err = r.attributeRepository.UpdateAttributeValue(attrs)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	var response models.ObjectImportResponse
	response.Success = true
	response.FailedImportObjectCount = insertedFailedObjectCount
	response.SuccessImportedObjectCount = insertedObjectCount
	response.TotalImportedObjectCount = insertedFailedObjectCount + insertedObjectCount
	return &response, nil
}
func (r *ObjectRepository) GetTypeId(attributeType string) int64 {
	switch strings.ToLower(attributeType) {
	case "string":
		return 4
	case "integer":
		return 1
	case "float":
		return 3
	case "boolean":
		return 5
	case "date":
		return 2
	case "richtext":
		return 6
	default:
		return 4
	}
}
func (r *ObjectRepository) toRTFUnicode(s string) string {
	rtf := ""
	for _, r := range s {
		rtf += fmt.Sprintf("\\u%d?", r)
	}
	return rtf
}
func (r *ObjectRepository) CreateObjectVersion(objectId uuid.UUID, objectName string, objectDescription string) (*uuid.UUID, error) {
	var versionId uuid.UUID
	query := `INSERT INTO [VERSION] (ID, ObjectID,objectName,ObjectDescription, SystemVersionNo, userVersionNo,DateCreated,DateModified,ModifiedBy, CreatedBy) VALUES(
		@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10
	)`
	versionId = uuid.New()
	_, err := r.db.Exec(query,
		versionId, objectId, objectName, objectDescription, 1, "v1", time.Now(), time.Now(), 62, 62,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating object version: %w", err)
	}
	return &versionId, nil
}

func (r *ObjectRepository) CreateObjectVersionWithTx(tx *sql.Tx, objectId uuid.UUID, objectName string, objectDescription string) (*uuid.UUID, error) {
	var versionId uuid.UUID
	query := `INSERT INTO [VERSION] (ID, ObjectID,objectName,ObjectDescription, SystemVersionNo, userVersionNo,DateCreated,DateModified,ModifiedBy, CreatedBy) VALUES(
		@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10
	)`
	versionId = uuid.New()
	_, err := tx.Exec(query,
		versionId, objectId, objectName, objectDescription, 1, "v1", time.Now(), time.Now(), 62, 62,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating object version with tx: %w", err)
	}
	return &versionId, nil
}

// CreateV2 creates a new object in the database using a transaction
func (r *ObjectRepository) CreateV2(tx *sql.Tx, req models.CreateObjectRequest) (*models.Object, error) {
	objectID := uuid.New()
	now := time.Now()

	var libraryId *uuid.UUID
	if req.IsLibrary {
		libraryId = &objectID
	} else if req.LibraryId != nil {
		val, _ := TransformUUID(*req.LibraryId)
		libraryId = &val
	}

	query := ` SELECT GeneralType from ObjectType where ObjectTypeID =@p1`
	if req.GeneralType == nil {
		req.GeneralType = new(int)
		err := tx.QueryRow(query, req.ObjectTypeID).Scan(req.GeneralType)
		if err != nil {
			return nil, fmt.Errorf("error getting object type general type: %w", err)
		}
	}

	insertObjectSql := `
		INSERT INTO [Object] (
			ObjectID, ObjectName, ObjectDescription, ObjectTypeID, Locked, IsImported, 
			IsLibrary, LibraryId, FileExtension, Prefix, Suffix, DateCreated, CreatedBy, 
			DateModified, ModifiedBy, IsCheckedOut, ExactObjectTypeID, CurrentVersionId, CheckedInVersionId,DeleteFlag, RichTextDescription,generalType
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18, @p19,@p20,@p21,@p22
		)
	`
	versionId, err := r.CreateObjectVersionWithTx(tx, objectID, req.ObjectName, req.ObjectDescription)
	if err != nil {
		return nil, fmt.Errorf("error creating object version: %w", err)
	}

	_, err = tx.Exec(insertObjectSql,
		objectID, req.ObjectName, req.ObjectDescription, req.ObjectTypeID, false, req.IsImported,
		req.IsLibrary, libraryId, req.FileExtension, req.Prefix, req.Suffix, now, req.CreatedBy,
		now, req.CreatedBy, false, req.ExactObjectTypeID, versionId, versionId, 0, req.RichTextDescription, req.GeneralType,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating object: %w", err)
	}

	// Create ObjectContent if DirectParentId is provided
	if req.DirectParentId != nil {
		parentID, _ := TransformUUID(*req.DirectParentId)

		// Get Parent's CurrentVersionId
		var parentVersionID uuid.UUID
		err := tx.QueryRow("SELECT CurrentVersionId FROM [Object] WHERE ObjectID = @p1", parentID).Scan(&parentVersionID)
		if err != nil {
			return nil, fmt.Errorf("error getting parent version id: %w", err)
		}
		///
		insertContentSql := `
		INSERT INTO ObjectContents (
			DocumentObjectID, ContainerVersionID, ObjectID, Instances, IsShortCut, 
			ContainmentType, DateCreated, CreatedBy, DateModified, ModifiedBy
		) OUTPUT INSERTED.ID 
		SELECT 
			@p1,
			container.CurrentVersionId,
			@p2,
			@p3,
			@p4,
			@p5,
			@p6,
			@p7,
			@p8,
			@p9
		FROM [Object] AS container
		WHERE container.ObjectID = @p1
	`

		//	var id int
		_, err = tx.Exec(insertContentSql,
			parentID,      // p1
			objectID,      // p2
			1,             // p3
			false,         // p4
			4,             // p5
			now,           // p6
			req.CreatedBy, // p7
			now,           // p8
			req.CreatedBy, // p9
		)
		if err != nil {
			return nil, fmt.Errorf("error creating object content: %w", err)
		}
	}

	// Construct and return the object
	obj := &models.Object{
		ObjectID:           objectID,
		ObjectName:         req.ObjectName,
		ObjectDescription:  req.ObjectDescription,
		ObjectTypeID:       req.ObjectTypeID,
		CurrentVersionId:   versionId,
		CheckedInVersionId: versionId,
		// ... populate other fields as needed
	}
	return obj, nil
}

// Create creates a new object in the database
func (r *ObjectRepository) Create(req models.CreateObjectRequest) (*models.Object, error) {
	objectID := uuid.New()
	now := time.Now()

	var libraryId *uuid.UUID
	if req.IsLibrary {
		libraryId = &objectID
	} else if req.LibraryId != nil {
		val, _ := TransformUUID(*req.LibraryId)
		libraryId = &val
	}

	query := ` SELECT GeneralType from ObjectType where ObjectTypeID =@p1`
	if req.GeneralType == nil {
		req.GeneralType = new(int)
		err := r.db.QueryRow(query, req.ObjectTypeID).Scan(req.GeneralType)
		if err != nil {
			return nil, fmt.Errorf("error getting object type general type: %w", err)
		}
	}

	query = `
		INSERT INTO [Object] (
			ObjectID, ObjectName, ObjectDescription, ObjectTypeID, Locked, IsImported, 
			IsLibrary, LibraryId, FileExtension, Prefix, Suffix, DateCreated, CreatedBy, 
			DateModified, ModifiedBy, IsCheckedOut, ExactObjectTypeID, CurrentVersionId, CheckedInVersionId,DeleteFlag, RichTextDescription,generalType
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18, @p19,@p20,@p21,@p22
		)
	`
	versionId, err := r.CreateObjectVersion(objectID, req.ObjectName, req.ObjectDescription)
	if err != nil {
		return nil, fmt.Errorf("error creating object version: %w", err)
	}
	_, err = r.db.Exec(query,
		objectID, req.ObjectName, req.ObjectDescription, req.ObjectTypeID, false, false,
		req.IsLibrary, libraryId, req.FileExtension, req.Prefix, req.Suffix, now, req.CreatedBy,
		now, req.CreatedBy, false, req.ExactObjectTypeID, versionId, versionId, 0, req.RichTextDescription, req.GeneralType,
	)
	fmt.Println("version id", versionId)
	fmt.Println("object name", req.ObjectName)
	if err != nil {
		return nil, fmt.Errorf("error creating object: %w", err)
	}
	objectData, _ := r.GetByID(objectID)
	fmt.Println("get db version id", objectData.CurrentVersionId)
	return objectData, nil
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
	id, _ = TransformUUID(id)
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
	fmt.Println("current version get from db", obj.CurrentVersionId)
	fmt.Println("current name get from db", obj.ObjectName)
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
	fmt.Println("object Name", req.ObjectName)

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
	id, _ = TransformUUID(id)
	fmt.Println("obj ID: ", id)
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
		  , FolderId,
		  isFolder)
		  AS
	(
		SELECT  o.*, 
				CAST(CASE WHEN v.ApprovalStatus =  dbo.const_ApprovalStatus_PendingApproval() THEN 1 ELSE 0 END AS BIT) AS IsPendingApproval,
				o.CheckedOutUserId AS CheckedOutBy, 
				CAST(CASE WHEN v.SystemVersionNo = 1 AND o.IsCheckedOut = 1 THEN 1 ELSE 0 END AS BIT) AS IsFirstVersionCheckedOut,
		        do.FolderId, do.isFolder AS isFolder
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
		        do.FolderId, do.isFolder AS isFolder
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
	//sqlServerUUID := toSQLServerUUID(ObjectID)
	//sqlUUID, _ := uuid.FromBytes(sqlServerUUID)
	//fmt.Println("UUID:", sqlUUID.String())
	ObjectID, _ = TransformUUID(ObjectID)

	rows, err := r.db.Query(query, ObjectID, isFolder, profileID)
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
			&obj.IsFolder,
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

func (r *ObjectRepository) GetHierarchyFolderV2(ObjectID uuid.UUID, profileID int, isFolder bool) ([]models.ObjectTree, error) {
	query := `
			WITH recurse (
      [ObjectID]
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
    , [IsDeleted]
    , [VisioAlias]
    , [HasVisioAlias]
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
    , FolderId
    , isFolder
) AS (
    
    SELECT  
          o.ObjectID
        , do.FolderId AS ObjectParentID
        , o.ObjectName
        , o.ObjectDescription
        , o.CurrentVersionId
        , o.CheckedInVersionId
        , o.IsImported
        , o.IsLibrary
        , o.LibraryId
        , o.FileExtension
        , o.GeneralType
        , o.GeneralTypeName
        , o.TypeId
        , o.IsDeleted
        , o.VisioAlias
        , o.HasVisioAlias
        , o.SortOrder
        , o.AutoSort
        , o.Prefix
        , o.Suffix
        , o.ProvenanceId
        , o.ProvenanceVersionId
        , o.CreatedBy
        , o.DateCreated
        , o.ModifiedBy
        , o.DateModified
        , o.IsCheckedOut
        , o.CheckedOutUserId
        , o.RichTextDescription
        , CAST(CASE WHEN v.ApprovalStatus = dbo.const_ApprovalStatus_PendingApproval() 
                    THEN 1 ELSE 0 END AS BIT)
        , o.CheckedOutUserId
        , CAST(CASE WHEN v.SystemVersionNo = 1 AND o.IsCheckedOut = 1 
                    THEN 1 ELSE 0 END AS BIT)
        , do.FolderId
        , do.isFolder
    FROM vwFolderContents AS do
    INNER JOIN vwObjectSimple AS o ON o.ObjectID = do.ObjectId
    INNER JOIN dbo.[Version] AS v ON v.ID = o.CurrentVersionId
    WHERE do.FolderId = @p1
      AND o.IsDeleted = 0
      AND do.IsDeleted = 0
      AND do.isFolder IN (CAST(@p2 AS BIT), CAST(1 AS BIT))

    UNION ALL

    SELECT  
          o.ObjectID
        , do.FolderId AS ObjectParentID
        , o.ObjectName
        , o.ObjectDescription
        , o.CurrentVersionId
        , o.CheckedInVersionId
        , o.IsImported
        , o.IsLibrary
        , o.LibraryId
        , o.FileExtension
        , o.GeneralType
        , o.GeneralTypeName
        , o.TypeId
        , o.IsDeleted
        , o.VisioAlias
        , o.HasVisioAlias
        , o.SortOrder
        , o.AutoSort
        , o.Prefix
        , o.Suffix
        , o.ProvenanceId
        , o.ProvenanceVersionId
        , o.CreatedBy
        , o.DateCreated
        , o.ModifiedBy
        , o.DateModified
        , o.IsCheckedOut
        , o.CheckedOutUserId
        , o.RichTextDescription
        , CAST(CASE WHEN v.ApprovalStatus = dbo.const_ApprovalStatus_PendingApproval() 
                    THEN 1 ELSE 0 END AS BIT)
        , o.CheckedOutUserId
        , CAST(CASE WHEN v.SystemVersionNo = 1 AND o.IsCheckedOut = 1 
                    THEN 1 ELSE 0 END AS BIT)
        , do.FolderId
        , do.isFolder
    FROM vwFolderContents AS do
    INNER JOIN recurse ON do.FolderId = recurse.ObjectID
    INNER JOIN vwObjectSimple AS o ON o.ObjectID = do.ObjectId
    INNER JOIN dbo.[Version] AS v ON v.ID = o.CurrentVersionId
    WHERE o.IsDeleted = 0
      AND do.IsDeleted = 0
      AND do.isFolder IN (CAST(@p2 AS BIT), CAST(1 AS BIT))
)

SELECT  
     recurse.ObjectID
    ,recurse.ObjectParentID
    , recurse.ObjectName
    , recurse.ObjectDescription
    , recurse.CurrentVersionId
    , recurse.CheckedInVersionId
    , recurse.IsImported
    , recurse.IsLibrary
    , recurse.LibraryId
    , recurse.FileExtension
    , recurse.GeneralType
    , recurse.GeneralTypeName
    , recurse.TypeId
	, recurse.IsDeleted
    , recurse.VisioAlias
    , recurse.HasVisioAlias
    , recurse.SortOrder
    , recurse.AutoSort
    , recurse.Prefix
    , recurse.Suffix
    , recurse.ProvenanceID
    , recurse.ProvenanceVersionID
    , recurse.CreatedBy
    , recurse.DateCreated
    , recurse.ModifiedBy
    , recurse.DateModified
    , recurse.IsCheckedOut
    , recurse.CheckedOutUserId
    , recurse.RichTextDescription
    , recurse.IsPendingApproval
    , recurse.CheckedOutBy
    , recurse.IsFirstVersionCheckedOut
    , recurse.FolderId
    , recurse.isFolder
    , vchkin.ObjectName AS CheckedInName
    , vchkin.VisioAlias AS CheckedInVisioAlias
    , vchkin.HasVisioAlias AS CheckedInHasVisioAlias
    , ISNULL(op.HasRead, 0) AS HasReadPermission
    , ISNULL(op.HasModifyContents, 0) AS HasModifyContentsPermission
    , ISNULL(op.HasDelete, 0) AS HasDeletePermission
    , ISNULL(op.HasModify, 0) AS HasModifyPermission
    , ISNULL(op.HasModifyRelationships, 0) AS HasModifyRelationshipsPermission
FROM recurse
LEFT JOIN dbo.[Version] AS vchkin ON vchkin.ID = recurse.CheckedInVersionId
LEFT JOIN dbo.ObjectPermissions AS op 
       ON op.ObjectID = recurse.ObjectID 
      AND op.ProfileID = @p3;

	`

	// Convert UUID to SQL Server format
	//sqlServerUUID := toSQLServerUUID(ObjectID)
	//sqlUUID, _ := uuid.FromBytes(sqlServerUUID)
	//fmt.Println("UUID:", sqlUUID.String())
	ObjectID, _ = TransformUUID(ObjectID)

	rows, err := r.db.Query(query, ObjectID, isFolder, profileID)
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
			&obj.IsDeleted,
			&obj.VisioAlias,
			&obj.HasVisioAlias,
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
			&obj.IsFolder,
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
		fmt.Println("TyyyyypeId", obj.TypeId)
		objects = append(objects, obj)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return objects, nil
}

func (r *ObjectRepository) GetByObjectTypeIDAndLibraryID(objectTypeID int, libraryID uuid.UUID, page, pageSize int) ([]models.Object, int, error) {
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	// Transform UUID to SQL Server format before queries
	libraryID, _ = TransformUUID(libraryID)

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM [Object] WHERE ExactObjectTypeID = @p1 AND LibraryId = @p2`
	err := r.db.QueryRow(countQuery, objectTypeID, libraryID).Scan(&totalCount)
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
		WHERE ExactObjectTypeID = @p1 AND LibraryId = @p2
		ORDER BY DateCreated DESC
		OFFSET @p3 ROWS FETCH NEXT @p4 ROWS ONLY
	`
	fmt.Println(query)
	fmt.Println(objectTypeID)
	fmt.Println(libraryID)
	fmt.Println(offset)
	fmt.Println(pageSize)
	rows, err := r.db.Query(query, objectTypeID, libraryID, offset, pageSize)
	fmt.Println(rows)
	fmt.Println(err)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving objects by type: %w", err)
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
			return nil, 0, fmt.Errorf("error scanning object: %w", err)
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

/*
   select library repository list more efficent
   SELECT	vwObject.ObjectID,
        ObjectParentID,
        ObjectName,
        ObjectDescription,
        [CurrentVersionId],
        LastCheckedInVersionId,
        IsImported,
        IsLibrary,
        LibraryId,
        [LibraryName],
        FileExtension,
        CreatedBy,
        DateCreated,
        LastModifiedBy,
        LastModified,
        UserVersionNo,
        SystemVersionNo,
        IsFirstVersionCheckedOut,
        [GeneralType],
        [GeneralTypeName],
        [TypeId],
        [TypeName],
        VisioAlias,
        HasVisioAlias,
        IsLocked,
        CASE
            WHEN vwObject.ObjectId=dbo.const_GuidEmpty()
            THEN -1
            ELSE SortOrder
        END	AS SortOrder,
        RichTextDescription,
		AutoSort
    FROM	[vwObject]
    INNER JOIN vwObjectReadPermissions AS perm ON perm.ObjectID = vwObject.ObjectID AND perm.ProfileID = 1
    WHERE	(IsLibrary = CAST(1 AS BIT) OR vwObject.ObjectId=dbo.const_GuidEmpty())
    AND		IsDeleted = CAST(0 AS BIT)
    AND		perm.HasReadPermission = CAST(1 AS BIT)
    ORDER BY SortOrder,ObjectName */
