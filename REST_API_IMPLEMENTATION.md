# REST API and Persistence Layer Implementation

## Overview

Successfully implemented a SQLite persistence layer and REST API endpoints for SBOM Sentinel, enabling users to submit and retrieve SBOM documents over HTTP. This establishes the foundation for a future web dashboard and makes the application network-accessible.

## Implementation Details

### 1. SQLite Persistence Layer (`internal/platform/database/sqlite.go`)

Created a complete SQLite repository implementation that:

- **Database Schema**: Stores SBOMs with ID, name, components (JSON), metadata (JSON), and timestamps
- **Interface Compliance**: Implements the `storage.Repository` interface
- **JSON Serialization**: Components and metadata are stored as JSON for flexibility
- **Upsert Logic**: Handles both insert and update operations based on SBOM ID
- **Proper Indexing**: Added indexes on name and created_at for performance
- **Error Handling**: Comprehensive error handling with descriptive messages

#### Key Features:
- Automatic schema initialization
- Transaction safety
- Efficient JSON storage and retrieval
- Connection management with proper cleanup

### 2. REST API Handlers (`internal/transport/rest/handlers.go`)

Implemented HTTP handlers for SBOM operations:

#### **POST /api/v1/sboms** - Submit SBOM
- Accepts `multipart/form-data` with SBOM file
- Validates file presence and size
- Parses CycloneDX format
- Stores in database
- Returns JSON response with SBOM ID

#### **GET /api/v1/sboms/get** - Retrieve SBOM
- Accepts SBOM ID as query parameter
- Returns full SBOM document as JSON
- Handles not found cases gracefully

#### Response Formats:
```json
// Success Response
{
  "id": "urn:uuid:...",
  "message": "SBOM submitted successfully"
}

// Error Response
{
  "error": "error_type",
  "message": "Human readable message"
}
```

### 3. Server Application (`cmd/sentinel-server/main.go`)

Enhanced the server with:

- **Database Initialization**: SQLite repository setup with configurable path
- **Route Registration**: RESTful endpoint configuration
- **Environment Configuration**: Database path and port configuration
- **Graceful Shutdown**: Proper resource cleanup
- **Informative Startup**: Clear endpoint documentation

#### Environment Variables:
- `DATABASE_PATH`: SQLite database file path (default: `./sentinel.db`)
- `PORT`: Server port (default: `8080`)

### 4. Database Schema

```sql
CREATE TABLE sboms (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    components TEXT NOT NULL, -- JSON-encoded components
    metadata TEXT NOT NULL,   -- JSON-encoded metadata
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX idx_sboms_name ON sboms(name);
CREATE INDEX idx_sboms_created_at ON sboms(created_at);
```

## API Usage Examples

### Submit an SBOM
```bash
curl -X POST \
  -F "sbom=@your-sbom.json" \
  http://localhost:8080/api/v1/sboms
```

### Retrieve an SBOM
```bash
curl "http://localhost:8080/api/v1/sboms/get?id=urn:uuid:12345678-1234-1234-1234-123456789012"
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Architecture Benefits

### **Hexagonal Architecture Compliance**
- Clear separation between core logic and infrastructure
- Repository pattern isolates persistence concerns
- Transport layer handles HTTP specifics
- Domain models remain database-agnostic

### **Scalability Considerations**
- JSON storage allows flexible schema evolution
- Indexing supports efficient queries
- Stateless HTTP handlers support horizontal scaling
- SQLite provides excellent read performance for moderate loads

### **Framework-Free Design**
- Uses only Go standard library HTTP server
- No external web framework dependencies
- Direct control over request/response handling
- Minimal dependency footprint

## Testing

### **Test Files Provided**
- `test-sbom.json`: Sample CycloneDX SBOM for testing
- `test-api.sh`: Comprehensive API testing script

### **Test Coverage**
- SBOM submission with valid file
- SBOM retrieval by ID
- Error handling for missing files
- Health endpoint validation

## Database Features

### **JSON Storage Benefits**
- Flexible component schema
- Easy metadata expansion
- Efficient querying with SQLite JSON functions
- No complex relational mapping required

### **Performance Optimizations**
- Primary key on SBOM ID for fast lookups
- Index on name for search capabilities
- Index on created_at for chronological queries
- Single table design minimizes joins

## Error Handling

### **Client Errors (4xx)**
- Invalid multipart forms
- Missing SBOM files
- Empty files
- Parse errors
- Missing required parameters

### **Server Errors (5xx)**
- Database connection failures
- Storage errors
- JSON serialization issues

### **Response Consistency**
- Standardized error format
- Proper HTTP status codes
- Descriptive error messages
- JSON content type headers

## Future Enhancements

### **Immediate Opportunities**
- URL path parameters instead of query parameters
- Pagination for SBOM listing
- Search and filtering capabilities
- File format validation
- Request size limits

### **Advanced Features**
- Authentication and authorization
- Rate limiting
- Request logging and metrics
- Database migrations
- Connection pooling
- Backup and restore functionality

### **Integration Possibilities**
- Analysis results storage
- Webhook notifications
- Batch processing APIs
- GraphQL endpoints
- OpenAPI documentation

## Security Considerations

### **Current Implementation**
- No authentication (suitable for internal/trusted networks)
- File size limits (32MB multipart form memory)
- Input validation on file presence
- SQL injection prevention through parameterized queries

### **Production Readiness Recommendations**
- Add authentication middleware
- Implement HTTPS/TLS
- Add request logging
- Configure CORS properly
- Add input sanitization
- Implement rate limiting

## Deployment

### **Local Development**
```bash
go build -o bin/sentinel-server ./cmd/sentinel-server/
./bin/sentinel-server
```

### **Production Considerations**
- Set `DATABASE_PATH` to persistent storage location
- Configure appropriate `PORT` for your environment
- Ensure SQLite file permissions are correct
- Consider database backup strategies
- Monitor disk usage for database growth

This implementation provides a solid foundation for web-based SBOM management and analysis, setting the stage for a comprehensive dashboard and advanced analytics features.