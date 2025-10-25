# 🚀 Avenger REST API - Complete Implementation

A production-ready REST API built with Go featuring inventory management, user authentication, and recipe management with role-based access control.

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Project Structure](#project-structure)
- [What Was Fixed](#what-was-fixed)

## ✨ Features

### 1. **Inventory Management** (Public)
- ✅ Create, Read, Update, Delete inventories
- ✅ Validate stock levels (cannot be negative)
- ✅ Track item status (active/broken)
- ✅ Unique inventory codes
- ✅ Full CRUD operations with validation

### 2. **User Authentication** (JWT-based)
- ✅ User registration with validation
- ✅ Secure password hashing (bcrypt)
- ✅ JWT token generation
- ✅ Role-based access control (admin/superadmin)
- ✅ Token expiration handling

### 3. **Recipe Management** (Protected)
- ✅ Public viewing of recipes
- ✅ Superadmin-only creation
- ✅ Superadmin-only deletion
- ✅ Rating validation (0-5)
- ✅ Cook time tracking

### 4. **Security & Middleware**
- ✅ JWT authentication middleware
- ✅ Request logging middleware
- ✅ Panic recovery
- ✅ Custom 404/405 handlers
- ✅ CORS-ready

### 5. **Data Validation**
- ✅ Request payload validation (validator/v10)
- ✅ Business logic validation in service layer
- ✅ User-friendly error messages
- ✅ SQL injection prevention

### 6. **Logging & Monitoring**
- ✅ Structured logging (slog)
- ✅ Debug logs in service layer
- ✅ Error logs with context
- ✅ HTTP request/response logging

## 🏗️ Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP Request
       ▼
┌─────────────────────────────────────┐
│         HTTP Handler Layer          │
│  - Parse & validate requests        │
│  - Set HTTP status codes            │
│  - Return JSON responses            │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Service Layer               │
│  - Business logic validation        │
│  - Logging (debug & error)          │
│  - Data transformation              │
│  - Error wrapping                   │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Repository Layer            │
│  - SQL queries                      │
│  - Database operations              │
│  - Error handling                   │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│          Database (PostgreSQL)      │
│  - Inventories (database/sql)       │
│  - Users & Recipes (GORM)           │
└─────────────────────────────────────┘
```

## 📦 Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Basic understanding of REST APIs
- Terminal/Command line access

## 🔧 Installation

### Step 1: Clone or create the project

```bash
mkdir avenger
cd avenger
```

### Step 2: Initialize Go module

```bash
go mod init avenger
```

### Step 3: Install dependencies

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

### Step 4: Setup PostgreSQL database

```bash
# Login to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE avenger_db;

# Exit psql
\q
```

### Step 5: Create .env file

```bash
cp .env.example .env
```

Edit `.env` with your database credentials:

```env
PG_HOST=localhost
PG_PORT=5432
PG_USER=postgres
PG_PASSWORD=yourpassword
PG_DBNAME=avenger_db
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

### Step 6: Run migrations (optional)

```bash
psql -U postgres -d avenger_db -f migrations/0001_init.sql
```

> **Note:** Tables will be auto-created when you run the application!

## 🚀 Running the Application

```bash
go run cmd/server/main.go
```

You should see:

```
🚀 Server starting on http://localhost:8080
Server running on localhost:8080
=====================================
📚 Available Endpoints:
  POST   /register          - Register new user
  POST   /login             - Login and get token
  GET    /inventories       - Get all inventories
  ...
=====================================
```

## 📡 