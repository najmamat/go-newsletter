# Authentication in Go Newsletter

## Overview

This project uses **client-side authentication with JWT validation** - the standard Supabase pattern. Authentication flows are handled by Supabase on the frontend, while the Go backend focuses on JWT validation and authorization.

## Architecture

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│   Frontend  │    │ Supabase     │    │ Go Backend  │
│ (React/Vue) │────│ Auth Service │    │             │
└─────────────┘    └──────────────┘    └─────────────┘
       │                   │                   │
       │ 1. Sign up/in     │                   │
       ├──────────────────▶│                   │
       │                   │                   │
       │ 2. JWT Token      │                   │
       │◀──────────────────┤                   │
       │                   │                   │
       │ 3. API Requests   │                   │
       │ (with JWT)        │                   │
       ├───────────────────┼──────────────────▶│
       │                   │                   │
       │                   │ 4. Validate JWT   │
       │                   │◀──────────────────┤
```

## How It Works

### 1. **User Registration/Login (Frontend)**
- Users sign up/login using Supabase client libraries
- Frontend calls: `supabase.auth.signUp({ email, password })`
- Supabase handles all auth complexity (passwords, email verification, etc.)

### 2. **JWT Token Management (Frontend)**
- Supabase returns JWT tokens automatically
- Frontend includes token in API requests: `Authorization: Bearer <jwt>`
- Supabase client handles token refresh automatically

### 3. **JWT Validation (Backend)**
- Go backend validates JWT signature using Supabase JWT secret
- Extracts user ID, email, and role from token claims
- No password handling or session management needed

### 4. **Profile Management (Backend)**
- When users sign up, a profile is created in the `profiles` table
- Backend manages additional user data (full_name, avatar_url, is_admin)
- Admin status is stored in database, not in JWT

## Implementation Details

### Environment Variables
```bash
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_JWT_SECRET=your-jwt-secret
```

### Middleware Usage
```go
// Require authentication
authMiddleware.RequireAuth(handler)

// Optional authentication (for public endpoints)
authMiddleware.OptionalAuth(handler)

// Require admin privileges
authMiddleware.RequireAdmin(handler)
```

### Getting User in Handlers
```go
func (s *Server) GetMe(w http.ResponseWriter, r *http.Request) {
    user, ok := services.GetUserFromContext(r.Context())
    if !ok {
        // Handle unauthorized
        return
    }
    // Use user.UserID, user.Email, etc.
}
```

## API Endpoints

### Documentation Endpoints
- `POST /auth/signup` - Returns instructions for frontend signup
- `POST /auth/signin` - Returns instructions for frontend signin  
- `POST /auth/password-reset-request` - Returns instructions for password reset

These endpoints don't handle actual authentication but provide documentation for frontend developers.

### Protected Endpoints
- `GET /me` - Get current user profile (requires auth)
- `PUT /me` - Update current user profile (requires auth)
- `GET /admin/*` - Admin endpoints (requires auth + admin check)

## Security Benefits

1. **Reduced Attack Surface**: Backend doesn't handle passwords or auth flows
2. **Automatic Security Updates**: Supabase handles security patches
3. **Built-in Features**: MFA, social auth, email verification come free
4. **Simple Validation**: Backend only validates JWT signatures
5. **Stateless**: No session management needed

## Frontend Integration Example

```javascript
// Sign up
const { user, error } = await supabase.auth.signUp({
  email: 'user@example.com',
  password: 'password'
})

// Sign in
const { user, error } = await supabase.auth.signInWithPassword({
  email: 'user@example.com', 
  password: 'password'
})

// Make authenticated API calls
const { data: profile } = await fetch('/api/v1/me', {
  headers: {
    'Authorization': `Bearer ${session.access_token}`
  }
})
```

## Database Schema

The `profiles` table extends Supabase's built-in `auth.users`:

```sql
CREATE TABLE profiles (
  id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  full_name TEXT,
  avatar_url TEXT,
  is_admin BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

## Benefits of This Approach

- ✅ **Standard Pattern**: How 99% of Supabase apps work
- ✅ **Simple Backend**: Focus on business logic, not auth complexity  
- ✅ **Better UX**: Real-time validation, built-in UI components
- ✅ **Secure**: Supabase handles sensitive operations
- ✅ **Maintainable**: Less code to maintain and test
- ✅ **Feature Rich**: MFA, social auth, email templates included 