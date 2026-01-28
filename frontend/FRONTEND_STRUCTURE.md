# Frontend Structure

This frontend is built with Vue 3 + Vite and is fully integrated with the Go backend.

## Directory Structure

```
src/
├── layouts/           # Page layouts (PublicLayout, MainLayout)
├── views/             # Page components (Feed, Profile, Settings, etc.)
├── stores/            # Pinia state management
│   ├── auth.js       # Authentication store (login, register, logout)
│   ├── posts.js      # Posts management
│   ├── groups.js     # Groups management
│   └── events.js     # Events management
├── services/          # API communication layer
│   └── api.js        # All backend endpoints
├── utils/             # Helper functions
│   └── helpers.js    # Date, string, validation utilities
├── router/            # Vue Router configuration
│   └── index.js      # Routes with auth guards
├── components/        # Reusable UI components
├── App.vue           # Root component with dynamic layouts
└── main.js           # App entry point
```

## Backend Integration

All API endpoints are configured in `src/services/api.js` and point to `http://localhost:8080`.

### Available Endpoints

#### Authentication
- `POST /register` - User registration
- `POST /login` - User login
- `POST /logout` - User logout

#### Posts
- `GET /posts` - Get all posts
- `POST /posts` - Create a new post
- `POST /comments` - Create a comment

#### Groups
- `GET /groups` - Get all groups
- `POST /groups` - Create a new group
- `POST /groups/invite` - Invite user to group
- `POST /groups/invitations/accept` - Accept group invitation
- `POST /groups/invitations/reject` - Reject group invitation
- `GET /groups/invitations/pending` - Get pending invitations
- `POST /groups/requests` - Request to join group
- `POST /groups/requests/accept` - Accept join request
- `POST /groups/requests/reject` - Reject join request
- `GET /groups/requests/pending` - Get pending requests

#### Events
- `GET /events` - Get events for a group
- `POST /events` - Create an event
- `POST /events/response` - Submit event response

#### Upload
- `POST /upload/image` - Upload image

## State Management (Pinia)

### Auth Store (`stores/auth.js`)
- `user` - Current logged-in user
- `isLoading` - Loading state
- `error` - Error message
- `isAuthenticated` - Boolean, whether user is logged in
- Methods: `login()`, `register()`, `logout()`

### Post Store (`stores/posts.js`)
- `posts` - Array of posts
- Methods: `fetchPosts()`, `createPost()`

### Group Store (`stores/groups.js`)
- `groups` - Array of user's groups
- `pendingInvites` - Pending group invitations
- `pendingRequests` - Pending join requests
- Methods: `fetchGroups()`, `createGroup()`, `inviteUser()`, `acceptInvitation()`, etc.

### Event Store (`stores/events.js`)
- `events` - Array of events
- Methods: `fetchEvents()`, `createEvent()`, `submitResponse()`

## Layouts

### PublicLayout.vue
Used for login and register pages. Shows minimal navigation.

### MainLayout.vue
Used for authenticated pages. Includes:
- Top navigation bar with logout
- Sidebar with navigation links
- Main content area

## Router Configuration

Routes are defined in `src/router/index.js` with authentication guards:
- Public routes: `/login`, `/register`
- Protected routes: `/feed`, `/profile`, `/settings`
- Redirects: Unauthenticated users → `/login`, Authenticated users trying to access auth pages → `/feed`

## Running the Frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend will run on `http://localhost:5173` by default.

## Environment Setup

Ensure the backend is running on `http://localhost:8080`. Modify `API_BASE_URL` in `src/services/api.js` if needed.

## Future Enhancements

- Add more reusable components in `components/` folder
- Implement search functionality
- Add notifications system
- Implement private and group chat
- Add more event and group management features
