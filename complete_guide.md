# 🚀 Complete REST API Implementation Guide

## 📋 Overview

This is a **beginner-friendly guide** for a Go REST API that manages:
- **Inventories** - Track items with stock levels
- **Users** - Register and authenticate users
- **Recipes** - Store cooking recipes (only superadmin can create/delete)

## 🔍 Critical Bugs Found

### 1. **Domain Models Issues**
- ❌ Inventory `Description` is `int` but should be `string`
- ❌ User model has duplicate ID field (gorm.Model already has ID)
- ❌ Missing validation tags

### 2. **Handler Problems**
- ❌ No error checking on `strconv.Atoi`
- ❌ No validation for request payloads
- ❌ Database errors exposed to clients
- ❌ Wrong HTTP status codes
- ❌ No Content-Type headers on responses

### 3. **Repository Issues**
- ❌ No error checking on `rows.Scan()`
- ❌ No `rows.Err()` check after iteration
- ❌ Update/Delete don't check if records exist
- ❌ Create doesn't return generated ID

### 4. **Service Layer**
- ❌ No business logic validation
- ❌ No logging
- ❌ Just passes through to repository

### 5. **Security Issues**
- ❌ JWT secret key hardcoded ("HIDUP_JOKOWI")
- ❌ No password strength validation in bcrypt
- ❌ Missing error handling in auth flow

### 6. **Main Application**
- ❌ No panic/404/405 handlers
- ❌ No database connection cleanup
- ❌ No structured logging

### 7. **SQL Migration**
- ❌ Syntax errors (curly braces instead of parentheses)
- ❌ Missing commas

## 📁 Fixed File Structure

```
avenger/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── domain/                  # Data models
│   │   ├── inventory.go
│   │   ├── user.go
│   │   └── recipe.go
│   ├── handler/                 # HTTP handlers (controllers)
│   │   ├── inventory_handler.go
│   │   ├── auth_handler.go
│   │   └── recipe_handler.go
│   ├── middleware/              # HTTP middleware
│   │   └── middleware.go
│   ├── repository/              # Database layer
│   │   ├── inventory_repository.go
│   │   ├── user_repository.go
│   │   └── recipe_repository.go
│   └── service/                 # Business logic layer
│       ├── inventory_service.go
│       ├── user_service.go
│       └── recipe_service.go
├── pkg/
│   ├── db/postgres.go          # Database connections
│   └── utils/jwt.go            # JWT utilities
├── migrations/                  # SQL migrations
│   └── 0001_init.sql
├── .env                        # Environment variables
├── .gitignore
├── go.mod
└── go.sum
```

## 🎯 Features Explained

### Feature 1: Inventory Management
**What it does:** Track items in your warehouse/store
- Create new items
- View all items or specific item
- Update item details
- Delete items
- Validate stock can't be negative

### Feature 2: User Authentication
**What it does:** Secure user registration and login
- Register new users with email/password
- Login returns JWT token
- Token used to access protected endpoints
- Two roles: admin and superadmin

### Feature 3: Recipe Management (Protected)
**What it does:** Store cooking recipes
- Anyone can view recipes
- Only superadmin can create recipes
- Only superadmin can delete recipes

### Feature 4: Middleware
**What it does:** Security and logging
- **AuthMiddleware**: Checks if user has valid token
- **LoggingMiddleware**: Logs all HTTP requests

## 🔧 Complete Implementation

### Setup Instructions

1. **Install Dependencies**
```bash
go get github.com/julienschmidt/httprouter
go get github.com/joho/godotenv
go get github.com/lib/pq
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/go-playground/validator/v10
go mod tidy
```

2. **Create .env file**
```env
PG_HOST=localhost
PG_PORT=5432
PG_USER=postgres
PG_PASSWORD=yourpassword
PG_DBNAME=avenger_db
JWT_SECRET=your-super-secret-key-change-this-in-production
```

3. **Run the application**
```bash
go run cmd/server/main.go
```

## 📡 API Endpoints

### Public Endpoints (No Auth Required)

#### 1. Register User
```bash
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepass123",
  "full_name": "John Doe",
  "age": 25,
  "occupation": "Software Engineer",
  "role": "admin"
}
```

#### 2. Login
```bash
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepass123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 3. Get All Inventories
```bash
GET /inventories

Response:
{
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Laptop Dell",
      "code": "LPT001",
      "stock": 10,
      "description": "Dell Inspiron 15",
      "status": "active"
    }
  ]
}
```

#### 4. Get Inventory by ID
```bash
GET /inventories/1
```

#### 5. Create Inventory
```bash
POST /inventories
Content-Type: application/json

{
  "name": "Laptop Dell",
  "code": "LPT001",
  "stock": 10,
  "description": "Dell Inspiron 15",
  "status": "active"
}
```

#### 6. Update Inventory
```bash
PUT /inventories/1
Content-Type: application/json

{
  "name": "Laptop HP",
  "code": "LPT001",
  "stock": 15,
  "description": "HP Pavilion",
  "status": "active"
}
```

#### 7. Delete Inventory
```bash
DELETE /inventories/1
```

#### 8. Get All Recipes (Public)
```bash
GET /recipes
```

### Protected Endpoints (Require Superadmin)

#### 9. Create Recipe (Superadmin Only)
```bash
POST /recipes
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
  "name": "Nasi Goreng",
  "description": "Indonesian fried rice",
  "cook_time": 20,
  "rating": 4.5
}
```

#### 10. Delete Recipe (Superadmin Only)
```bash
DELETE /recipes/1
Authorization: Bearer YOUR_JWT_TOKEN
```

## 🎓 Understanding the Code

### Layer Architecture Explained

```
Request → Handler → Service → Repository → Database
         ↓          ↓         ↓
      Validate   Business   SQL
       Input      Logic    Queries
```

#### Handler Layer
**Job:** Handle HTTP requests and responses
- Parse request body
- Validate input format
- Call service layer
- Return JSON response
- Set HTTP status codes

#### Service Layer
**Job:** Business logic and validation
- Validate business rules
- Log operations
- Transform data
- Handle errors gracefully

#### Repository Layer
**Job:** Database operations
- Execute SQL queries
- Map database rows to structs
- Handle database errors

## 🔒 Security Features

### Password Hashing
```go
// When registering
hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// When logging in
bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
```

### JWT Authentication
```go
// Generate token after login
token, _ := utils.GenerateJWT(userID, role)

// Validate token on protected routes
claims, err := utils.ValidateToken(token)
```

### Role-Based Access Control
```go
// Only superadmin can access
middleware.AuthMiddleware(handler, "superadmin")
```

## ✅ Validation Rules

### Inventory
- Name: required, 3-100 characters
- Code: required, 2-50 characters, unique
- Stock: required, >= 0
- Description: max 500 characters
- Status: required, must be "active" or "broken"

### User
- Email: required, valid email format
- Password: required, minimum 8 characters
- Full Name: required, 6-15 characters
- Age: required, >= 17
- Occupation: required
- Role: required, "admin" or "superadmin"

### Recipe
- Name: required
- Description: required
- Cook Time: required (minutes)
- Rating: required (0.0-5.0)

## 🐛 Error Responses

### 400 Bad Request
```json
{
  "message": "Validation failed",
  "errors": {
    "name": "name must be at least 3 characters",
    "stock": "stock must be greater than or equal to 0"
  }
}
```

### 401 Unauthorized
```json
{
  "message": "Unauthorized",
  "errors": "Missing or invalid authentication token"
}
```

### 403 Forbidden
```json
{
  "message": "Forbidden",
  "errors": "You don't have permission to access this resource"
}
```

### 404 Not Found
```json
{
  "message": "Resource not found"
}
```

### 409 Conflict
```json
{
  "message": "Inventory code already exists"
}
```

### 500 Internal Server Error
```json
{
  "message": "Internal server error"
}
```

## 📊 Testing the API

### Using cURL

```bash
# Register
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin12345","full_name":"Admin User","age":25,"occupation":"Administrator","role":"superadmin"}'

# Login
TOKEN=$(curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin12345"}' | jq -r '.token')

# Create Recipe (with auth)
curl -X POST http://localhost:8080/recipes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Rendang","description":"Spicy beef","cook_time":180,"rating":5.0}'
```

## 🎯 Key Improvements Made

1. ✅ Added validator/v10 for all models
2. ✅ Added comprehensive logging in services
3. ✅ Improved error handling (no DB errors exposed)
4. ✅ Added panic/404/405 handlers
5. ✅ Fixed authentication flow
6. ✅ Added proper HTTP status codes
7. ✅ Fixed domain models
8. ✅ Added input sanitization
9. ✅ Fixed SQL migration syntax
10. ✅ Added JWT secret from environment
11. ✅ Added database connection pooling
12. ✅ Added proper resource cleanup
13. ✅ Fixed duplicate code detection
14. ✅ Added comprehensive validation messages

## 🚦 Next Steps

1. Run the fixed application
2. Test all endpoints
3. Review logs in console
4. Check validation messages
5. Test authentication flow
6. Verify role-based access
