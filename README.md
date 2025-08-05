# elotus-challenge

A comprehensive solution for Elotus's backend challenge, featuring data structures & algorithms implementations and a complete file upload API with JWT authentication.

This project is an implementation for [elotus's challenge](https://github.com/elotusteam/challenges/blob/main/backend.md)

## I. Data Structures And Algorithms
Implementations placed under folder `data-structure-and-algorithms`
- Gray Code: `graycode.go`

- Sum of Distances in Tree: `findlength.go`

- Maximum Length of Repeated Subarray: `sumofdistanceintree.go`

## II. Hackathon

Implementation placed under folder `server`. The project implemented in golang and uses SQLite to persist user data.

This project follows a layered architecture with clear separation of concerns:

**Layers**
- Handlers → Services → Repository → Database
- Each layer has a single responsibility and depends only on lower layers
- Business logic is isolated in the service layer

**Request Flow**
1. **HTTP Request** → Middleware (auth, logging) → Handler
2. **Handler** → Service (business logic) → Repository (data access)
3. **Repository** → Database → Response back through layers

**Project structure:**
```
elotus-challenge/
├── data-structure-and-algorithms/     # Algorithm implementations
└── server/                           # Backend API implementation
    ├── database/                     # Database setup
    ├── internal/                     # Internal services
    │   ├── repository/               # Data access layer
    │   └── services/                 # Business logic layer
    ├── middleware/                   # HTTP middleware
    ├── handler/                      # HTTP handlers
    ├── static/                       # Web interface
    └── test/                         # Tests 
```
### Features
---

#### 1. JWT Authentication
- User registration with username/password
- User login with JWT token generation
- Token-based authentication for protected endpoints

#### 2. File Upload API
- Secure file upload endpoint at `/upload`
- Accepts only image files (JPEG, PNG, GIF, WebP, BMP, TIFF, SVG)
- Maximum file size: 8MB
- Files are saved to `/tmp` directory with unique names
- Stores file metadata in database with HTTP information


### Running the Application
---

#### 1. Prerequisites
- Go 1.24.4 or later
- SQLite (automatically initialized)

#### 2. Build and Run

```bash
# Navigate to server directory
cd server

# Install dependencies
go mod tidy

# Run the application
go run main.go
```

The server will start on port 8080 by default. You can set a custom port using the `PORT` environment variable:

```bash
PORT=3000 go run main.go
```

#### 3. Environment Variables

The application supports the following environment variables for configuration:

| Variable | Description | Default Value | Example |
|----------|-------------|---------------|---------|
| `PORT` | Server port | `8080` | `PORT=3000` |
| `JWT_SECRET` | Secret key for JWT token signing | `elotus-challenge-default` | `JWT_SECRET=my-super-secret-key` |
| `DB_PATH` | Path to SQLite database file | `./challenge.db` | `DB_PATH=/data/app.db` |
| `TEMP_DIR` | Directory for uploaded files | `./tmp` | `TEMP_DIR=/uploads` |
| `TOKEN_EXPIRATION_SECONDS` | JWT token expiration time in seconds | `86400` (24 hours) | `TOKEN_EXPIRATION_SECONDS=3600` |

**Example with environment variables:**

```bash
# Set environment variables
export JWT_SECRET="my-production-secret"
export DB_PATH="/data/production.db"
export TEMP_DIR="/var/uploads"
export TOKEN_EXPIRATION_SECONDS=7200
export PORT=3000

# Run the application
go run main.go
```

### Test uploading file
--- 

#### 1. Using the Web Forms

1. Start the server
2. Navigate to `http://localhost:8080/form/register` to register a new user
3. Navigate to `http://localhost:8080/form/login` to login and get a JWT token
4. Navigate to `http://localhost:8080/form/upload` to upload files using the token

#### 2. Using curl

```bash
# 1. Register a user
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# 2. Login to get token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# 3. Upload file (replace YOUR_JWT_TOKEN with actual token)
curl -X POST http://localhost:8080/api/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "data=@/path/to/your/image.jpg"
```


#### API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/api/register` | User registration | ❌ |
| `POST` | `/api/login` | User login | ❌ |
| `POST` | `/api/upload` | File upload | ✅ |
