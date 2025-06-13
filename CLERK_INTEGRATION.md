# Clerk Authentication Integration

## Overview

Thothix integrates with Clerk to manage user authentication securely and scalably. This guide explains how to configure and use the integration.

## Authentication Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Thothix API   │    │   Clerk.com     │
│   (Nuxt.js)     │    │     (Go)        │    │   (Auth SaaS)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                        │                        │
        │ 1. Login/Register      │                        │
        │───────────────────────▶│                        │
        │                        │ 2. Redirect to Clerk   │
        │                        │───────────────────────▶│
        │                        │                        │
        │ 3. Auth with Clerk     │                        │
        │◀───────────────────────────────────────────────│
        │                        │                        │
        │ 4. API calls with JWT  │                        │
        │───────────────────────▶│ 5. Verify JWT         │
        │                        │───────────────────────▶│
        │                        │ 6. User info           │
        │                        │◀───────────────────────│
        │ 7. Response            │                        │
        │◀───────────────────────│                        │
```

## Clerk Configuration

### 1. Clerk Account Setup

1. Create an account at [clerk.com](https://clerk.com)
2. Create a new application
3. Configure authentication providers (Email, Google, etc.)

### 2. Get API Keys

```bash
# From Clerk Dashboard → API Keys
CLERK_PUBLISHABLE_KEY=pk_test_...
CLERK_SECRET_KEY=sk_test_...
```

### 3. Environment Variables Configuration

```bash
# .env file
CLERK_SECRET_KEY=sk_test_your_secret_key_here
```

## Backend Configuration

### 1. Authentication Middleware

The `ClerkAuth` middleware automatically verifies JWT tokens:

```go
// Automatic configuration via environment
protected.Use(middleware.ClerkAuth(cfg.ClerkSecretKey))
```

#### How the Middleware Works

1. **Token Extraction**: Extracts JWT from `Authorization: Bearer <token>` header
2. **Clerk Verification**: Calls Clerk API `https://api.clerk.com/v1/users/me`
3. **Validation**: Verifies user has valid ID and email
4. **Context Setup**: Sets in Gin context:
   - `clerk_user_id` - Clerk user ID
   - `clerk_email` - Primary email  
   - `clerk_first_name` - First name
   - `clerk_last_name` - Last name
   - `clerk_image_url` - Avatar URL
   - `user_id` - For BaseModel hooks

#### Clerk Data Structure

```go
type ClerkUser struct {
    ID                   string                 `json:"id"`
    Username             *string                `json:"username"`
    FirstName            *string                `json:"first_name"`
    LastName             *string                `json:"last_name"`
    ImageURL             string                 `json:"image_url"`
    PrimaryEmailAddress  *ClerkEmailAddress     `json:"primary_email_address"`
    // Other fields...
}
```

### 2. Available Endpoints

#### Authentication

- `POST /api/v1/auth/sync` - Sync user from Clerk to local DB
- `GET /api/v1/auth/me` - Get current user
- `POST /api/v1/auth/webhooks/clerk` - Webhook for automatic sync

#### Authentication Flow

1. **Frontend**: Authenticates with Clerk and obtains JWT
2. **Frontend**: Sends request with `Authorization: Bearer <jwt>` header
3. **Backend**: Middleware verifies JWT with Clerk
4. **Backend**: Extracts user info and puts it in context
5. **Backend**: Executes business logic with authenticated user

## User Synchronization

### Manual Synchronization

The frontend must call `/api/v1/auth/sync` after login:

```javascript
// Frontend (after Clerk login)
const response = await fetch('/api/v1/auth/sync', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${clerkToken}`,
    'Content-Type': 'application/json'
  }
});
```

### Automatic Synchronization (Webhook)

Configure webhook in Clerk Dashboard:

1. **URL**: `https://your-domain.com/api/v1/auth/webhooks/clerk`
2. **Events**: 
   - `user.created`
   - `user.updated`
   - `user.deleted`

The system will automatically sync:

- ✅ New user creation
- ✅ Existing user data updates (email, name, avatar)
- ✅ User deletion

## User Data Mapping

### From Clerk to Thothix

| Clerk Field | Thothix Field | Notes |
|-------------|---------------|-------|
| `id` | `id` | Primary key |
| `primary_email_address.email_address` | `email` | Primary email |
| `first_name + last_name` | `name` | Full name |
| `image_url` | `avatar_url` | User avatar |
| `username` | `name` | Fallback if name missing |
| (default) | `system_role` | Always `user` for new users |

### Role Management

- **New users**: Default `user` role
- **Admin**: Must be assigned manually via API
- **Project roles**: Managed separately in `project_members`

## Integration Testing

### 1. Manual Testing

```bash
# 1. Get token from Clerk (frontend)
TOKEN="your_clerk_jwt_token"

# 2. Test synchronization
curl -X POST "http://localhost:30000/api/v1/auth/sync" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

# 3. Test get current user
curl -X GET "http://localhost:30000/api/v1/auth/me" \
  -H "Authorization: Bearer $TOKEN"
```

### 2. Webhook Testing

```bash
# Simulate user creation webhook
curl -X POST "http://localhost:30000/api/v1/auth/webhooks/clerk" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user.created",
    "data": {
      "id": "user_test123",
      "email_addresses": [
        {"email_address": "test@example.com"}
      ],
      "first_name": "Test",
      "last_name": "User",
      "image_url": "https://example.com/avatar.jpg"
    }
  }'
```

## Debug and Troubleshooting

### Debug Logs

To enable authentication debug logs:

```bash
# In docker-compose.yml or environment variables
DEBUG=true
LOG_LEVEL=debug
GIN_MODE=debug
```

### Common Issues

#### 1. Invalid JWT Token

**Error**: `Invalid token: clerk verification failed with status: 401`

**Solutions**:

- Verify that `CLERK_SECRET_KEY` is configured correctly
- Check that the token hasn't expired (frontend must refresh)
- Make sure you're using the correct token for the environment (development/production)

#### 2. User ID Not Found

**Error**: `Clerk user ID not found`

**Solutions**:

- Verify that the `ClerkAuth` middleware is configured correctly
- Check that the endpoint is protected by the middleware
- Make sure the token is valid

#### 3. Synchronization Issues

**Error**: User not created in local database

**Solutions**:

- Call `/api/v1/auth/sync` after each login
- Verify database connection
- Check logs for database errors

#### 4. Webhook Not Working

**Issues**: Automatic synchronization doesn't happen

**Solutions**:

- Verify webhook URL in Clerk Dashboard
- Check that the endpoint is public (not protected by auth)
- Verify configured events: `user.created`, `user.updated`, `user.deleted`

### Testing the Integration

```bash
# 1. Test authentication middleware
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:30000/api/v1/auth/me

# 2. Test manual synchronization
curl -X POST \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     http://localhost:30000/api/v1/auth/sync

# 3. Test webhook (simulation)
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"type":"user.created","data":{"id":"user_123","email_addresses":[{"email_address":"test@example.com"}]}}' \
     http://localhost:30000/api/v1/auth/webhooks/clerk
```

### Monitoring

To monitor the integration in production:

```bash
# Verify synchronized users
docker-compose exec postgres psql -U postgres -d thothix-db \
  -c "SELECT id, email, name, created_at FROM users ORDER BY created_at DESC LIMIT 10;"

# Verify application logs
docker-compose logs thothix-api --tail=50 --follow
```

## Security

### Best Practices

- ✅ Never expose `CLERK_SECRET_KEY` in frontend
- ✅ Use HTTPS in production
- ✅ Always validate tokens from backend
- ✅ Implement rate limiting
- ✅ Monitor webhooks for attacks

### Production Configuration

```bash
# Production variables
CLERK_SECRET_KEY=sk_live_your_live_secret_key
ENVIRONMENT=production
GIN_MODE=release
```

## Frontend Integration

### Clerk Setup (Nuxt.js)

```javascript
// nuxt.config.js
export default {
  runtimeConfig: {
    public: {
      clerkPublishableKey: process.env.CLERK_PUBLISHABLE_KEY
    }
  }
}

// plugins/clerk.client.js
import { ClerkProvider } from '@clerk/vue'

export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig()
  
  nuxtApp.vueApp.use(ClerkProvider, {
    publishableKey: config.public.clerkPublishableKey
  })
})
```

### Authentication Hook

```javascript
// composables/useAuth.js
export const useAuth = () => {
  const { isSignedIn, user, getToken } = useClerk()
  
  const syncUser = async () => {
    if (!isSignedIn.value) return
    
    const token = await getToken()
    await $fetch('/api/v1/auth/sync', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
  }
  
  const apiCall = async (url, options = {}) => {
    const token = await getToken()
    return $fetch(url, {
      ...options,
      headers: {
        'Authorization': `Bearer ${token}`,
        ...options.headers
      }
    })
  }
  
  return {
    isSignedIn,
    user,
    syncUser,
    apiCall
  }
}
```

## Conclusion

The Clerk integration provides:
- ✅ Secure and scalable authentication
- ✅ Automatic user synchronization
- ✅ Robust session management
- ✅ Multi-provider support (Google, GitHub, etc.)
- ✅ Pre-built UI components for frontend

For technical support, see the [Clerk documentation](https://clerk.com/docs) or contact the development team.
