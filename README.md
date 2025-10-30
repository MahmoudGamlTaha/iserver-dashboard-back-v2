# Enterprise Architect API

A well-organized Go REST API application for managing Enterprise Architect data with MS SQL Server database.

## Project Structure

```
enterprise-architect-api/
├── config/              # Configuration management
├── handlers/            # HTTP request handlers
├── models/              # Data models and request/response structures
├── repositories/        # Database operations layer
├── services/            # Business logic layer
├── utils/               # Utility functions
├── main.go              # Application entry point
├── go.mod               # Go module dependencies
└── .env.example         # Environment variables template
```

## Features

### Modules

1. **Objects** - Manage enterprise architecture objects
2. **Object Types** - Manage object type definitions
3. **Profiles** - Manage user profiles
4. **Object Contents** - Manage object content relationships

### Operations

Each module supports:
- **Create** - Add new records
- **Read** - Get by ID or list all with pagination
- **Update** - Modify existing records
- **Delete** - Remove records

### Special Operations

- **Get Libraries** - Retrieve objects where `IsLibrary = 1`
- **Get Objects by Type ID** - Retrieve objects filtered by `ObjectTypeID`

## API Endpoints

### Objects

- `GET /api/objects` - List all objects (with pagination)
- `POST /api/objects` - Create a new object
- `GET /api/objects/{id}` - Get object by ID
- `PUT /api/objects/{id}` - Update object
- `DELETE /api/objects/{id}` - Delete object
- `GET /api/objects/libraries` - Get all library objects
- `GET /api/objects/type/{typeId}` - Get objects by type ID

### Object Types

- `GET /api/object-types` - List all object types (with pagination)
- `POST /api/object-types` - Create a new object type
- `GET /api/object-types/{id}` - Get object type by ID
- `PUT /api/object-types/{id}` - Update object type
- `DELETE /api/object-types/{id}` - Delete object type

### Profiles

- `GET /api/profiles` - List all profiles (with pagination)
- `POST /api/profiles` - Create a new profile
- `GET /api/profiles/{id}` - Get profile by ID
- `PUT /api/profiles/{id}` - Update profile
- `DELETE /api/profiles/{id}` - Delete profile

### Object Contents

- `GET /api/object-contents` - List all object contents (with pagination)
- `POST /api/object-contents` - Create a new object content
- `GET /api/object-contents/{id}` - Get object content by ID
- `PUT /api/object-contents/{id}` - Update object content
- `DELETE /api/object-contents/{id}` - Delete object content

### Folders

- `GET /api/folders/object-type/{libraryId}` - Get object type folders by library ID
- `GET /api/folders/{folderId}/contents?profileId={profileId}` - Get folder contents by folder ID and profile ID

### Health Check

- `GET /health` - Health check endpoint

## Pagination

All list endpoints support pagination with query parameters:
- `page` - Page number (default: 1)
- `pageSize` - Number of items per page (default: 10, max: 100)

Example: `GET /api/objects?page=1&pageSize=20`

## Setup and Installation

### Prerequisites

- Go 1.21 or higher
- MS SQL Server database
- Database tables created based on the provided schema

### Installation Steps

1. Clone the repository:
```bash
cd enterprise-architect-api
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. Set environment variables:
```bash
export DB_SERVER=your-server
export DB_PORT=1433
export DB_DATABASE=EnterpriseArchitect
export DB_USER=your-username
export DB_PASSWORD=your-password
export SERVER_PORT=8080
```

5. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Building for Production

```bash
go build -o enterprise-architect-api
./enterprise-architect-api
```

## Configuration

The application uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | Server host address | `0.0.0.0` |
| `SERVER_PORT` | Server port | `8080` |
| `DB_SERVER` | Database server address | `localhost` |
| `DB_PORT` | Database port | `1433` |
| `DB_DATABASE` | Database name | `EnterpriseArchitect` |
| `DB_USER` | Database username | `sa` |
| `DB_PASSWORD` | Database password | `` |

## Example API Requests

### Create an Object

```bash
curl -X POST http://localhost:8080/api/objects \
  -H "Content-Type: application/json" \
  -d '{
    "objectName": "My Object",
    "objectDescription": "Description of my object",
    "objectTypeId": 1,
    "exactObjectTypeId": 1,
    "richTextDescription": "Rich text description",
    "isLibrary": false,
    "createdBy": 1
  }'
```

### Get All Objects with Pagination

```bash
curl http://localhost:8080/api/objects?page=1&pageSize=10
```

### Get Libraries

```bash
curl http://localhost:8080/api/objects/libraries?page=1&pageSize=10
```

### Get Objects by Type ID

```bash
curl http://localhost:8080/api/objects/type/1?page=1&pageSize=10
```

### Update an Object

```bash
curl -X PUT http://localhost:8080/api/objects/{object-id} \
  -H "Content-Type: application/json" \
  -d '{
    "objectName": "Updated Object Name",
    "modifiedBy": 1
  }'
```

### Delete an Object

```bash
curl -X DELETE http://localhost:8080/api/objects/{object-id}
```

## Database Schema

The application works with the following MS SQL Server tables:

- **Object** - Main objects table
- **ObjectType** - Object type definitions
- **Profile** - User profiles
- **ObjectContents** - Object content relationships

Refer to the provided table schemas for detailed column definitions.

## Architecture

The application follows a layered architecture:

1. **Handlers Layer** - HTTP request handling and routing
2. **Services Layer** - Business logic and validation
3. **Repositories Layer** - Database operations
4. **Models Layer** - Data structures and DTOs

This separation ensures maintainability, testability, and scalability.

## Error Handling

All endpoints return appropriate HTTP status codes:
- `200 OK` - Successful operation
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

Error responses follow this format:
```json
{
  "error": "Error message",
  "message": "Detailed error description"
}
```

## License

This project is licensed under the MIT License.

