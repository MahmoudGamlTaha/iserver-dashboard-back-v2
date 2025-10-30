# Folder Features Documentation

## Overview

Two new endpoints have been added to support folder operations in the Enterprise Architect API.

## New Endpoints

### 1. Get Object Type Folders

**Endpoint:** `GET /api/folders/object-type/{libraryId}`

**Purpose:** Retrieves all folders and system repositories associated with a specific library.

**Query Used:**
```sql
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
    o.DeleteFlag AS [IsObjectDeleted],
    0 AS IsTemplateType,
    o.CurrentVersionId AS ObjectVersion,
    o.LibraryId
FROM dbo.Object AS o
INNER JOIN dbo.ObjectType AS ot ON ot.ObjectTypeID = o.ObjectTypeID
WHERE o.GeneralType IN (dbo.const_GeneralType_Folder(), dbo.const_GeneralType_SystemRepository()) 
    AND o.LibraryId = @libraryId
```

**Parameters:**
- `libraryId` (UUID, path parameter) - The library ID to filter folders

**Response Model:** `ObjectTypeFolder`

**Fields Returned:**
- `objectId` - Unique identifier of the folder
- `generalType` - Type indicator (Folder or System Repository)
- `objectName` - Name of the folder
- `sortOrder` - Display order
- `generalTypeName` - Human-readable type name ("Folder" or "System Repository")
- `typeId` - Object type ID
- `typeName` - Object type name
- `isObjectDeleted` - Deletion flag
- `isTemplateType` - Template indicator (always 0)
- `objectVersion` - Current version ID
- `libraryId` - Associated library ID

**Example Request:**
```bash
curl http://localhost:8080/api/folders/object-type/123e4567-e89b-12d3-a456-426614174000
```

---

### 2. Get Folders by Library

**Endpoint:** `GET /api/folders/{folderId}/contents?profileId={profileId}`

**Purpose:** Retrieves the contents of a specific folder with permission information based on the user's profile.

**Query Used:**
```sql
SELECT  
    o.*,
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
LEFT JOIN [dbo].ObjectPermissions AS op ON op.ObjectID = o.ObjectID AND op.ProfileID = @profileId
WHERE do.FolderId = @folderId
    AND o.IsDeleted = CAST(0 AS BIT)
    AND do.IsDeleted = CAST(0 AS BIT)
ORDER BY SortOrder, CreatedBy
```

**Parameters:**
- `folderId` (UUID, path parameter) - The folder ID
- `profileId` (integer, query parameter, **required**) - The profile ID for permission filtering

**Response Model:** `FolderContent`

**Fields Returned:**

**Object Information:**
- All standard object fields (objectId, objectName, objectDescription, etc.)

**Version Information:**
- `isPendingApproval` - Whether the version is pending approval
- `checkedInName` - Name of the checked-in version
- `checkedInVisioAlias` - Visio alias of checked-in version
- `checkedInHasVisioAlias` - Whether checked-in version has Visio alias
- `isFirstVersionCheckedOut` - Whether the first version is checked out

**Permission Information:**
- `hasReadPermission` - User can read the object
- `hasModifyContentsPermission` - User can modify object contents
- `hasDeletePermission` - User can delete the object
- `hasModifyPermission` - User can modify object properties
- `hasModifyRelationshipsPermission` - User can modify object relationships

**Checkout Information:**
- `checkedOutBy` - User ID who has checked out the object

**Example Request:**
```bash
curl http://localhost:8080/api/folders/123e4567-e89b-12d3-a456-426614174000/contents?profileId=1
```

---

## Database Dependencies

### Views Used:
- `vwFolderContents` - Contains folder-object relationships
- `vwObjectSimple` - Simplified view of objects

### Tables Used:
- `Object` - Main objects table
- `ObjectType` - Object type definitions
- `Version` - Version information and approval status
- `ObjectPermissions` - User permissions per object

### Functions Used:
- `dbo.const_GeneralType_Folder()` - Returns the constant value for folder type
- `dbo.const_GeneralType_SystemRepository()` - Returns the constant value for system repository type
- `dbo.const_ApprovalStatus_PendingApproval()` - Returns the constant value for pending approval status

---

## Implementation Details

### Files Added:

1. **models/folder.go**
   - `ObjectTypeFolder` - Model for folder/system repository objects
   - `FolderContent` - Model for folder contents with permissions

2. **repositories/folder_repository.go**
   - `GetObjectTypeFolders()` - Repository method for fetching folders by library
   - `GetFoldersByLibrary()` - Repository method for fetching folder contents

3. **services/folder_service.go**
   - `GetObjectTypeFolders()` - Service layer with business logic
   - `GetFoldersByLibrary()` - Service layer with validation

4. **handlers/folder_handler.go**
   - `GetObjectTypeFolders()` - HTTP handler for folder endpoint
   - `GetFoldersByLibrary()` - HTTP handler for folder contents endpoint

### Routes Added to main.go:

```go
// Folder routes
api.HandleFunc("/folders/object-type/{libraryId}", folderHandler.GetObjectTypeFolders).Methods("GET")
api.HandleFunc("/folders/{folderId}/contents", folderHandler.GetFoldersByLibrary).Methods("GET")
```

---

## Testing with Postman

The Postman collection has been updated with two new requests under the "Folders" folder:

1. **Get Object Type Folders**
   - Method: GET
   - URL: `{{baseUrl}}/folders/object-type/:libraryId`
   - Variable: `libraryId` (UUID)

2. **Get Folders by Library**
   - Method: GET
   - URL: `{{baseUrl}}/folders/:folderId/contents?profileId=1`
   - Variables: `folderId` (UUID)
   - Query Parameter: `profileId` (integer)

---

## Error Handling

### Common Errors:

**400 Bad Request:**
- Invalid UUID format for libraryId or folderId
- Missing profileId query parameter
- Invalid profileId format

**500 Internal Server Error:**
- Database connection issues
- Missing database views or functions
- SQL query execution errors

---

## Notes

1. The `profileId` parameter is **required** for the folder contents endpoint to properly filter permissions.

2. Both endpoints rely on database functions (`const_GeneralType_Folder()`, etc.) that must exist in your database schema.

3. The folder contents query uses views (`vwFolderContents`, `vwObjectSimple`) that must be present in the database.

4. Permission fields default to `false` (0) if no permissions are found for the given profile.

5. Only non-deleted objects and folder relationships are returned.

6. Results are ordered by `SortOrder` and `CreatedBy` for consistent display.

