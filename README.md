# 01Forum - Web Discussion Platform

## 🌐 Overview

01Forum is a dynamic web forum application that enables users to engage in meaningful discussions by allowing users to create, categorize, and interact with posts while maintaining a secure and efficient user authentication system. Built with Go and SQLite, Following RESTful principles and clean architecture patterns, 01Forum demonstrates how to build a scalable web application from scratch.

## ✨ Features

### 🔐 Authentication System
- **Secure Registration Flow**
  - Email uniqueness validation
  - Username availability checking
  - Password encryption using bcrypt
  - Custom validation rules for input fields
  
- **Session Management**
  - UUID-based session tracking
  - Cookie-based authentication
  - 24-hour session expiration
  - Protection against multiple concurrent sessions

### 📝 Post Management
- **Content Creation**
  - Multi-category association
  - Code block formatting
  - Content validation and sanitization

### 🔄 Interactive Elements
- **Real-time Features**
  - Like/dislike functionality
  - Comment threading
  - Post categories
  - Real-time updates

### 🔍 Advanced Filtering System
- **Category-based Filtering**
  - Subforum-style organization
  - Multi-category support
  - Dynamic management

### 🛡️ Security Features
- **Data Protection**
  - SQL injection prevention
  - XSS protection
  - Input sanitization

- **Rate Limiting**
  - Comment throttling
  - Post creation limits
  - API request limiting

### 🛠 Technical Implementation
- **Database Management**
  - SQLite integration
  - Efficient query optimization
  - Data integrity enforcement

- **Frontend Architecture**
  - Vanilla JavaScript implementation
  - Dynamic content loading
  - Responsive design
  - Browser compatibility

## 📚 What We Learn

Building this project from scratch provides invaluable experience in various aspects of web development:

### 🔧 Backend Development
1. **Go Web Programming**
   - HTTP server implementation
   - Router creation
   - Middleware development
   - Error handling patterns

2. **Database Management**
   - SQL database design
   - Query optimization
   - CRUD operations

3. **Authentication & Security**
   - Session management
   - Password hashing
   - Cookie handling
   - Security best practices

### 🎨 Frontend Development
1. **Pure JavaScript**
   - DOM manipulation
   - Event handling
   - AJAX requests

2. **HTML/CSS Skills**
   - Semantic HTML
   - CSS architecture
   - Responsive design
   - Cross-browser compatibility

### 📝 System Design
1. **Architecture Patterns**
   - MVC pattern implementation
   - Service layer design
   - Repository pattern
   - Clean architecture principles

2. **API Design**
   - RESTful API development
   - Status code handling
   - Request/Response formatting
   - API documentation

## 🚀 Installation Instructions

### 💾 Standard Installation
```bash
# Clone repository
git clone https://github.com/ANAS-OU/go_forum.git
cd go_forum

# Install dependencies
go mod download
go mod tidy

# Run application
go run main.go
```

### 🐳 Docker Installation
```bash
# Build image
docker build -t forum .

# Run container
docker run -p 8080:8080 forum
```

## 📁 Project Structure
```
GO_FORUM/
├── 🎮 controllers/          # controls database
├── 💾 database/            # Database initialization and schema
├── 🌐 handlers/            # HTTP request handlers
│   ├── 📡 api/            # API endpoints handlers
│   ├── 🔐 auth/           # Authentication handlers
│   ├── 🛡️ middleware/     # Request middleware
├── 📦 models/              # Data structures and types
├── 🛣️ routers/            # URL routing configuration
├── 🛠️ utils/              # Helper functions and utilities
└── 🎨 web/                # Frontend resources
    ├── 📂 static/         # Static assets (CSS, JS, images)
    └── 📝 templates/      # HTML templates
```

## 📝 API Documentation

### 🪪 Authentication Endpoints

- `POST /auth/register` - Register new user
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout

### 📨 Post Endpoints

- `GET /api/posts` - Get all posts
- `GET /api/posts/{id}` - Get specific post
- `POST /newpost` - Create new post
- `POST /newcomment` - Add comment to post
- `POST /reaction` - React to post/comment

### 📚 Category Endpoints

- `GET /api/categories` - Get all categories


## 🤝 Contributing Guidelines

**Pull Request Process**
   - Fork repository
   - Create feature branch
   - Submit PR with description
   - Ensure tests pass

## 👥 Team

Project developed by:
- [@ikazbat](https://github.com/kazbatdriss1)
- [@yelmach](https://github.com/yassinalmach)
- [@asaaoud](https://github.com/Saaoud99)
- [@amazighi](https://github.com/amazighii)
- [@oanass](https://github.com/ANAS-OU)