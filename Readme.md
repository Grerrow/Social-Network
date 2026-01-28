## Full-Stack Social Network (Short Overview)

A full-stack social networking platform built from scratch with **Go** and **Vue.js**, focusing on **real-time communication**, **networking fundamentals**, and **secure backend design**.

The application supports user authentication, content sharing, private and group messaging, notifications, and group management. Real-time features are implemented using a custom **WebSocket hub** designed to handle concurrent connections, multi-device users, and targeted message delivery.

### Key Highlights
- **Backend**: Go (`net/http`) with session-based authentication using HTTP-only cookies
- **Real-Time**: Custom WebSocket hub for private messaging, group chat, presence, and notifications
- **Frontend**: Vue 3 (Composition API) with Pinia state management and Vue Router
- **Data**: SQLite with versioned migrations and enforced foreign key constraints
- **Security**: Bcrypt password hashing, centralized auth middleware, prepared SQL statements
- **Infrastructure**: Dockerized frontend and backend with multi-stage builds and Docker Compose

### What This Project Demonstrates
- Designing and implementing a complete client–server system
- Managing concurrent WebSocket connections safely using Go routines and channels
- Building RESTful APIs and integrating them with real-time protocols
- Practical understanding of networking concepts (HTTP, WebSocket, CORS, sessions)
- Applying security fundamentals in a real application

This project was developed as part of my training at **Zone01**, with an emphasis on systems thinking, networking, and production-oriented backend development.


# Full-Stack Social Network (Detailed OverView)

A real-time social networking platform built from scratch with Go and Vue.js. This project demonstrates end-to-end implementation of a production-grade web application, including WebSocket-based real-time communication, session-based authentication, and a containerized deployment architecture.

## Overview

This application replicates core social networking functionality: user authentication, content sharing, real-time messaging, group management, and notification systems. The project focuses on clean separation of concerns, scalable architecture patterns, and practical security implementations.

The system handles concurrent real-time connections, manages complex many-to-many relationships in the database, and coordinates state across multiple clients using a custom WebSocket hub.

## Core Features

### Authentication & Authorization
- Session-based authentication using HTTP-only cookies
- Bcrypt password hashing with configurable cost
- Server-side session management stored in SQLite
- Middleware-based route protection for all authenticated endpoints
- Session validation on every protected request

### Real-Time Communication
- **WebSocket Hub Architecture**: Custom hub manages concurrent client connections with user indexing for efficient message routing
- **Private Messaging**: Direct one-on-one chat with message history retrieval
- **Group Chat**: Multi-user conversations within groups with role-based access
- **Presence System**: Real-time online/offline status broadcast to all connected clients
- **Notification Delivery**: Push notifications for follow requests, group invitations, event updates, and comments delivered via WebSocket

### Social Features
- **Privacy Controls**: Public, private, and selective post visibility (almost-private allows targeting specific followers)
- **Follow System**: Follow requests with accept/decline workflows for private accounts
- **Groups**: Create groups, invite members, request to join, and manage membership
- **Events**: Group events with RSVP functionality (going/not going)
- **Comments**: Nested commenting on posts and group posts
- **Search**: Search for users and groups by name

### Data Management
- SQLite database with 17+ migration files managing schema evolution
- Foreign key constraints enforced for referential integrity
- Separate tables for users, sessions, posts, comments, messages, groups, group_members, group_invitations, group_requests, events, event_responses, notifications, followers, and post_visibility
- Database connection pooling and transaction support

## Architecture & Design

### Backend (Go)
- **HTTP Server**: Standard library `net/http` with custom multiplexer routing
- **WebSocket Management**: Custom hub implementation using Gorilla WebSocket library
  - Maintains two indexes: global client map and per-user client map for multi-device support
  - Uses channels for thread-safe registration/deregistration
  - Selective message delivery (private messages go only to intended recipient)
- **Middleware Pattern**: Authentication middleware wraps protected endpoints and injects user context
- **Structured Logging**: Request logging and error tracking throughout handlers
- **RESTful API Design**: Method-based routing (GET, POST) with JSON request/response bodies

### Frontend (Vue.js 3)
- **Composition API**: Modern Vue 3 patterns with `<script setup>` syntax
- **Pinia State Management**: Centralized stores for auth, chat, notifications, posts, groups, and search
- **Vue Router**: Client-side routing with route guards for authentication checks
- **Composables**: Reusable logic for API calls (`useApi`), chat functionality (`useChat`), and follow status (`useFollowStatus`)
- **Layout System**: Separate layouts for authenticated (MainLayout) and public (PublicLayout) views
- **WebSocket Integration**: Persistent connection managed through Vue component lifecycle

### API Communication
- REST endpoints for CRUD operations (posts, comments, groups, users)
- WebSocket for bidirectional real-time events (messages, notifications, presence)
- JSON serialization for all data transfer
- CORS headers configured for frontend-backend separation

### WebSocket Message Flow
1. Client authenticates via REST endpoint, receives session cookie
2. Client opens WebSocket connection (`/ws`), session validated from cookie
3. Server registers client in hub, broadcasts presence
4. Clients send typed messages (`{"type": "chat_message", "data": {...}}`)
5. Hub routes messages based on type:
   - `presence`: broadcast to all
   - `chat_message`: send only to recipient
   - `notification`: send only to target user
6. On disconnect, server cleans up and broadcasts offline presence

### Database Schema
The database uses normalized tables with foreign key relationships:
- `users` ← `sessions` (1:many)
- `users` ← `posts` ← `comments` (1:many:many)
- `users` ← `groups` (creator), `group_members` (many:many)
- `users` ← `followers` (many:many, self-referential)
- `users` ← `messages` (many:many via sender/receiver)
- `groups` ← `group_messages`, `group_posts`, `events`

Migrations are versioned and applied automatically on startup using golang-migrate library.

### Containerization
- **Multi-stage Docker builds** for both frontend and backend to minimize image size
- **Frontend**: Vite build → Nginx static serving
- **Backend**: Go compilation → scratch-based runtime image
- **Docker Compose**: Orchestrates both services with shared network, volume persistence for database and uploads
- **Health checks**: Defined for both services to ensure readiness
- Separate production configuration (`docker-compose.prod.yml`) available

## Networking & Security Considerations

### Network Communication
- **HTTP/HTTPS**: All REST API calls use standard HTTP protocol (production would use HTTPS with TLS termination)
- **WebSocket Protocol**: Upgrades HTTP connection to persistent bidirectional TCP connection
- **Same-Origin Policy**: CORS headers explicitly allow frontend origin
- **Session Cookies**: HttpOnly flag prevents XSS access, SameSite attribute prevents CSRF

### Security Implementations
- **Password Storage**: Bcrypt hashing (never stores plaintext)
- **Session Tokens**: UUIDs stored in database, validated on every authenticated request
- **Input Validation**: JSON decoding with type safety, explicit checks for required fields
- **SQL Injection Prevention**: Prepared statements used throughout (sql.DB QueryRow/Exec methods)
- **Authentication Middleware**: Centralized auth logic prevents bypassing on protected routes
- **Connection Security**: WebSocket hijacking prevented by requiring valid session before upgrade

### Known Security Limitations
- No rate limiting on login attempts (brute force vulnerability)
- No CSRF tokens (relies on SameSite cookies)
- No input sanitization for XSS (trusts frontend validation)
- Sessions don't expire automatically (no TTL implementation)
- No HTTPS enforcement (production deployment needed)
- WebSocket connections not encrypted separately (relies on WSS in production)

## Technology Stack

**Backend**
- Go 1.24 (standard library for HTTP server, contexts, JSON handling)
- Gorilla WebSocket (RFC 6455 WebSocket implementation)
- SQLite3 with mattn/go-sqlite3 driver
- golang-migrate for schema versioning
- bcrypt for password hashing
- google/uuid for session token generation

**Frontend**
- Vue.js 3.5 (Composition API)
- Vue Router 4 (client-side routing)
- Pinia 3 (state management)
- Vite 7 (build tool and dev server)
- Nginx (production static file serving)

**Infrastructure**
- Docker & Docker Compose
- SQLite (file-based database)
- Nginx (reverse proxy and static hosting)

## Local Setup & Running

### Prerequisites
- Docker (24.x or later)
- Docker Compose (2.x or later)

### Running with Docker
```bash
# Clone the repository
git clone <repository-url>
cd social-network

# Start both services
docker compose up --build

# Frontend available at: http://localhost
# Backend API at: http://localhost:8080
```

The application will:
1. Build the Go backend and Vue frontend from source
2. Apply database migrations automatically
3. Create persistent volumes for the database and user uploads
4. Start the backend on port 8080 and frontend on port 80

### Development Mode
```bash
# Backend (requires Go 1.24+)
cd backend
go run main.go

# Frontend (requires Node.js 20+)
cd frontend
npm install
npm run dev
```

### Useful Commands
```bash
# View logs
docker compose logs -f

# Stop services
docker compose down

# Clean restart (removes all data)
docker compose down -v
docker compose up --build
```

### Testing
Create a user account, log in, and:
- Create posts with different privacy settings
- Follow another user and test follow request flow
- Create a group and invite members
- Send private messages and group messages
- Observe real-time notifications and presence updates

## What I Learned

### Full-Stack System Design
- Architected a complete client-server system from database to UI
- Made pragmatic trade-offs between feature complexity and maintainability
- Designed RESTful APIs with consistent patterns and error handling
- Separated concerns across layers (handlers, middleware, database, business logic)

### Real-Time Communication
- Implemented a custom WebSocket hub with concurrent connection management
- Used Go channels and mutexes for thread-safe state management across goroutines
- Understood the difference between broadcast and directed messaging patterns
- Handled edge cases like reconnections, multi-device scenarios, and ungraceful disconnects

### Networking Fundamentals
- HTTP request/response lifecycle and statelessness
- WebSocket protocol upgrade handshake and persistent connections
- How CORS works and why it's necessary for cross-origin requests
- Session management using cookies and the HTTP-only security attribute
- Difference between TCP-level persistence (WebSocket) and application-level state

### Security in Practice
- Session-based authentication vs token-based (chose sessions for simplicity)
- Defense against common vulnerabilities: SQL injection, XSS, CSRF, password attacks
- Middleware pattern for centralized authorization checks
- Secure password storage with salting and adaptive hashing (bcrypt)
- Realized that security is layered: no single mechanism is sufficient

### Working with Containers
- Multi-stage Docker builds to reduce production image size
- Orchestrated multi-container applications with Docker Compose
- Used volumes for data persistence across container restarts
- Configured health checks for service readiness verification
- Understood the difference between development and production container configurations

This project taught me that building production-grade software requires thinking beyond features—it's about reliability, security, maintainability, and understanding the underlying protocols that make everything work.

## Author

- Christos M
