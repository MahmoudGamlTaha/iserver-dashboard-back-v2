# API Documentation

## Base URL

```
http://localhost:8080/api
```

## Authentication

Currently, the API does not implement authentication. This should be added in production.

## Response Format

### Success Response

```json
{
  "data": { ... },
  "page": 1,
  "pageSize": 10,
  "totalCount": 100,
  "totalPages": 10
}
```

### Error Response

```json
{
  "error": "Error message",
  "message": "Detailed error description"
}
```

---

## Objects API

### 1. Get All Objects

**Endpoint:** `GET /api/objects`

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `pageSize` (optional, default: 10, max: 100) - Items per page

**Response:** `200 OK`

```json
{
  "data": [
    {
      "objectId": "123e4567-e89b-12d3-a456-426614174000",
      "objectName": "Sample Object",
      "objectDescription": "Description",
      "objectTypeId": 1,
      "isLibrary": false,
      "dateCreated": "2024-01-01T00:00:00Z",
      "dateModified": "2024-01-01T00:00:00Z",
      ...
    }
  ],
  "page": 1,
  "pageSize": 10,
  "totalCount": 50,
  "totalPages": 5
}
```

---

### 2. Get Object by ID

**Endpoint:** `GET /api/objects/{id}`

**Path Parameters:**
- `id` (UUID) - Object ID

**Response:** `200 OK`

```json
{
  "objectId": "123e4567-e89b-12d3-a456-426614174000",
  "objectName": "Sample Object",
  "objectDescription": "Description",
  "objectTypeId": 1,
  "exactObjectTypeId": 1,
  "richTextDescription": "Rich description",
  "isLibrary": false,
  "locked": false,
  "isImported": false,
  "isCheckedOut": false,
  "dateCreated": "2024-01-01T00:00:00Z",
  "createdBy": 1,
  "dateModified": "2024-01-01T00:00:00Z",
  "modifiedBy": 1
}
```

**Error Response:** `404 Not Found`

---

### 3. Create Object

**Endpoint:** `POST /api/objects`

**Request Body:**

```json
{
  "objectName": "New Object",
  "objectDescription": "Object description",
  "objectTypeId": 1,
  "exactObjectTypeId": 1,
  "richTextDescription": "Rich text description",
  "isLibrary": false,
  "fileExtension": ".vsdx",
  "prefix": "OBJ",
  "suffix": "001",
  "createdBy": 1
}
```

**Required Fields:**
- `objectName`
- `objectTypeId`
- `exactObjectTypeId`
- `createdBy`

**Response:** `201 Created`

```json
{
  "objectId": "123e4567-e89b-12d3-a456-426614174000",
  "objectName": "New Object",
  ...
}
```

---

### 4. Update Object

**Endpoint:** `PUT /api/objects/{id}`

**Path Parameters:**
- `id` (UUID) - Object ID

**Request Body:**

```json
{
  "objectName": "Updated Object Name",
  "objectDescription": "Updated description",
  "isLibrary": true,
  "modifiedBy": 1
}
```

**Required Fields:**
- `modifiedBy`

**Note:** All other fields are optional. Only provided fields will be updated.

**Response:** `200 OK`

```json
{
  "objectId": "123e4567-e89b-12d3-a456-426614174000",
  "objectName": "Updated Object Name",
  ...
}
```

---

### 5. Delete Object

**Endpoint:** `DELETE /api/objects/{id}`

**Path Parameters:**
- `id` (UUID) - Object ID

**Response:** `200 OK`

```json
{
  "message": "Object deleted successfully"
}
```

**Error Response:** `404 Not Found`

---

### 6. Get Libraries

**Endpoint:** `GET /api/objects/libraries`

**Description:** Returns all objects where `isLibrary = true`

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `pageSize` (optional, default: 10, max: 100) - Items per page

**Response:** `200 OK`

```json
{
  "data": [
    {
      "objectId": "123e4567-e89b-12d3-a456-426614174000",
      "objectName": "Library Object",
      "isLibrary": true,
      ...
    }
  ],
  "page": 1,
  "pageSize": 10,
  "totalCount": 25,
  "totalPages": 3
}
```

---

### 7. Get Objects by Type ID

**Endpoint:** `GET /api/objects/type/{typeId}`

**Path Parameters:**
- `typeId` (integer) - Object Type ID

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `pageSize` (optional, default: 10, max: 100) - Items per page

**Response:** `200 OK`

```json
{
  "data": [
    {
      "objectId": "123e4567-e89b-12d3-a456-426614174000",
      "objectName": "Object of Type 1",
      "objectTypeId": 1,
      ...
    }
  ],
  "page": 1,
  "pageSize": 10,
  "totalCount": 15,
  "totalPages": 2
}
```

---

## Object Types API

### 1. Get All Object Types

**Endpoint:** `GET /api/object-types`

**Query Parameters:**
- `page` (optional, default: 1)
- `pageSize` (optional, default: 10, max: 100)

**Response:** `200 OK`

---

### 2. Get Object Type by ID

**Endpoint:** `GET /api/object-types/{id}`

**Path Parameters:**
- `id` (integer) - Object Type ID

**Response:** `200 OK`

```json
{
  "objectTypeId": 1,
  "objectTypeName": "Document",
  "description": "Document type",
  "isTemplateType": false,
  "activeType": true,
  "fileExtension": ".vsdx",
  "dateCreated": "2024-01-01T00:00:00Z",
  "createdBy": 1,
  "dateModified": "2024-01-01T00:00:00Z",
  "modifiedBy": 1
}
```

---

### 3. Create Object Type

**Endpoint:** `POST /api/object-types`

**Request Body:**

```json
{
  "objectTypeName": "New Type",
  "description": "Type description",
  "fileExtension": ".vsdx",
  "isTemplateType": false,
  "activeType": true,
  "createdBy": 1
}
```

**Required Fields:**
- `createdBy`

**Response:** `201 Created`

---

### 4. Update Object Type

**Endpoint:** `PUT /api/object-types/{id}`

**Request Body:**

```json
{
  "objectTypeName": "Updated Type Name",
  "description": "Updated description",
  "activeType": false,
  "modifiedBy": 1
}
```

**Required Fields:**
- `modifiedBy`

**Response:** `200 OK`

---

### 5. Delete Object Type

**Endpoint:** `DELETE /api/object-types/{id}`

**Response:** `200 OK`

```json
{
  "message": "Object type deleted successfully"
}
```

---

## Profiles API

### 1. Get All Profiles

**Endpoint:** `GET /api/profiles`

**Query Parameters:**
- `page` (optional, default: 1)
- `pageSize` (optional, default: 10, max: 100)

**Response:** `200 OK`

---

### 2. Get Profile by ID

**Endpoint:** `GET /api/profiles/{id}`

**Path Parameters:**
- `id` (integer) - Profile ID

**Response:** `200 OK`

```json
{
  "profileId": 1,
  "profileName": "Admin Profile",
  "profileDescription": "Administrator profile",
  "portalStartPageId": "123e4567-e89b-12d3-a456-426614174000",
  "dateCreated": "2024-01-01T00:00:00Z",
  "createdBy": 1,
  "dateModified": "2024-01-01T00:00:00Z",
  "modifiedBy": 1
}
```

---

### 3. Create Profile

**Endpoint:** `POST /api/profiles`

**Request Body:**

```json
{
  "profileName": "New Profile",
  "profileDescription": "Profile description",
  "portalStartPageId": "123e4567-e89b-12d3-a456-426614174000",
  "createdBy": 1
}
```

**Required Fields:**
- `profileName`
- `createdBy`

**Response:** `201 Created`

---

### 4. Update Profile

**Endpoint:** `PUT /api/profiles/{id}`

**Request Body:**

```json
{
  "profileName": "Updated Profile Name",
  "profileDescription": "Updated description",
  "modifiedBy": 1
}
```

**Required Fields:**
- `modifiedBy`

**Response:** `200 OK`

---

### 5. Delete Profile

**Endpoint:** `DELETE /api/profiles/{id}`

**Response:** `200 OK`

```json
{
  "message": "Profile deleted successfully"
}
```

---

## Object Contents API

### 1. Get All Object Contents

**Endpoint:** `GET /api/object-contents`

**Query Parameters:**
- `page` (optional, default: 1)
- `pageSize` (optional, default: 10, max: 100)

**Response:** `200 OK`

---

### 2. Get Object Content by ID

**Endpoint:** `GET /api/object-contents/{id}`

**Path Parameters:**
- `id` (integer) - Object Content ID

**Response:** `200 OK`

```json
{
  "id": 1,
  "documentObjectId": "123e4567-e89b-12d3-a456-426614174000",
  "containerVersionId": "123e4567-e89b-12d3-a456-426614174001",
  "objectId": "123e4567-e89b-12d3-a456-426614174002",
  "instances": 1,
  "isShortCut": false,
  "containmentType": 1,
  "dateCreated": "2024-01-01T00:00:00Z",
  "createdBy": 1,
  "dateModified": "2024-01-01T00:00:00Z",
  "modifiedBy": 1
}
```

---

### 3. Create Object Content

**Endpoint:** `POST /api/object-contents`

**Request Body:**

```json
{
  "documentObjectId": "123e4567-e89b-12d3-a456-426614174000",
  "containerVersionId": "123e4567-e89b-12d3-a456-426614174001",
  "objectId": "123e4567-e89b-12d3-a456-426614174002",
  "instances": 1,
  "isShortCut": false,
  "containmentType": 1,
  "createdBy": 1
}
```

**Required Fields:**
- `documentObjectId`
- `containerVersionId`
- `objectId`
- `instances`
- `containmentType`
- `createdBy`

**Response:** `201 Created`

---

### 4. Update Object Content

**Endpoint:** `PUT /api/object-contents/{id}`

**Request Body:**

```json
{
  "instances": 2,
  "isShortCut": true,
  "modifiedBy": 1
}
```

**Required Fields:**
- `modifiedBy`

**Response:** `200 OK`

---

### 5. Delete Object Content

**Endpoint:** `DELETE /api/object-contents/{id}`

**Response:** `200 OK`

```json
{
  "message": "Object content deleted successfully"
}
```

---

## Health Check

### Health Check Endpoint

**Endpoint:** `GET /health`

**Response:** `200 OK`

```
OK
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request data |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

---

## Notes

1. All timestamps are in ISO 8601 format (UTC)
2. UUIDs are in standard UUID format (8-4-4-4-12)
3. Pagination is zero-indexed (page 1 is the first page)
4. Maximum page size is 100 items
5. All request bodies must be valid JSON
6. Content-Type header should be `application/json` for POST and PUT requests



---

## Folders API

### 1. Get Object Type Folders

**Endpoint:** `GET /api/folders/object-type/{libraryId}`

**Description:** Retrieves folders and system repositories filtered by library ID. Returns objects where GeneralType is either Folder or System Repository.

**Path Parameters:**
- `libraryId` (UUID) - Library ID to filter folders

**Response:** `200 OK`

```json
[
  {
    "objectId": "123e4567-e89b-12d3-a456-426614174000",
    "generalType": 1,
    "objectName": "My Folder",
    "sortOrder": 1,
    "generalTypeName": "Folder",
    "typeId": 5,
    "typeName": "Folder Type",
    "isObjectDeleted": false,
    "isTemplateType": 0,
    "objectVersion": "123e4567-e89b-12d3-a456-426614174001",
    "libraryId": "123e4567-e89b-12d3-a456-426614174002"
  }
]
```

**Error Response:** `400 Bad Request` - Invalid library ID

**Error Response:** `500 Internal Server Error` - Database error

---

### 2. Get Folders by Library

**Endpoint:** `GET /api/folders/{folderId}/contents`

**Description:** Retrieves the contents of a folder with permission information for a specific profile. Returns all non-deleted objects within the folder along with approval status, permissions, and version information.

**Path Parameters:**
- `folderId` (UUID) - Folder ID

**Query Parameters:**
- `profileId` (required, integer) - Profile ID for permission filtering

**Response:** `200 OK`

```json
[
  {
    "objectId": "123e4567-e89b-12d3-a456-426614174000",
    "objectName": "Document 1",
    "objectDescription": "Description of document",
    "objectTypeId": 1,
    "isLibrary": false,
    "sortOrder": 1,
    "dateCreated": "2024-01-01T00:00:00Z",
    "createdBy": 1,
    "dateModified": "2024-01-01T00:00:00Z",
    "modifiedBy": 1,
    "isPendingApproval": false,
    "checkedInName": "Document 1 v1",
    "hasReadPermission": true,
    "hasModifyContentsPermission": true,
    "hasDeletePermission": false,
    "hasModifyPermission": true,
    "hasModifyRelationshipsPermission": true,
    "isFirstVersionCheckedOut": false,
    ...
  }
]
```

**Required Query Parameters:**
- `profileId` - Must be provided as a query parameter

**Example Request:**
```
GET /api/folders/123e4567-e89b-12d3-a456-426614174000/contents?profileId=1
```

**Error Response:** `400 Bad Request` - Invalid folder ID or missing profileId

**Error Response:** `500 Internal Server Error` - Database error

---

## Notes on Folder Endpoints

1. **Object Type Folders** endpoint uses database functions:
   - `dbo.const_GeneralType_Folder()` - Returns the constant for folder type
   - `dbo.const_GeneralType_SystemRepository()` - Returns the constant for system repository type

2. **Folders by Library** endpoint uses views and joins:
   - `vwFolderContents` - View for folder content relationships
   - `vwObjectSimple` - Simplified object view
   - Joins with `Version` table for approval status
   - Joins with `ObjectPermissions` table for user permissions

3. **Permission Fields** in folder contents:
   - `hasReadPermission` - User can read the object
   - `hasModifyContentsPermission` - User can modify object contents
   - `hasDeletePermission` - User can delete the object
   - `hasModifyPermission` - User can modify object properties
   - `hasModifyRelationshipsPermission` - User can modify object relationships

4. **Version Information:**
   - `isPendingApproval` - Object version is pending approval
   - `isFirstVersionCheckedOut` - First version is currently checked out
   - `checkedInName` - Name of the checked-in version
   - `checkedOutBy` - User ID who has checked out the object

---

