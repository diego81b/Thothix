# 🔐 Clerk Authentication Integration Guide

## Table of Contents

1. [Overview](#overview)
2. [Quick Start](#quick-start)
3. [Authentication Architecture](#authentication-architecture)
4. [Setup & Configuration](#setup--configuration)
5. [User Synchronization](#user-synchronization)
6. [Local Development with Webhooks](#local-development-with-webhooks)
7. [API Endpoints](#api-endpoints)
8. [Testing & Debugging](#testing--debugging)
9. [Frontend Integration](#frontend-integration)
10. [Security & Best Practices](#security--best-practices)
11. [Troubleshooting](#troubleshooting)

## Overview

Thothix integrates with Clerk to provide secure, scalable authentication for the enterprise messaging application. This guide covers everything from basic setup to advanced webhook configuration for local development.

### Key Features

- ✅ Secure JWT-based authentication
- ✅ Automatic user synchronization via webhooks
- ✅ Manual user sync endpoints
- ✅ Multi-provider support (Google, GitHub, email)
- ✅ Pre-built UI components
- ✅ Local development with ngrok tunneling
- ✅ Comprehensive testing scripts

## Quick Start

### 1. Prerequisites

- [Clerk account](https://clerk.com) with a configured application
- Go backend running on port 30000
- [Ngrok](https://ngrok.com/download) for local webhook testing (optional)

### 2. Environment Setup

```bash
# .env file
CLERK_SECRET_KEY=sk_test_your_secret_key_here
CLERK_WEBHOOK_SECRET=whsec_your_webhook_signing_secret
PORT=30000
```

### 3. Quick Start Script

```bash
# Windows - Start complete local development environment
scripts\start-local-dev.bat

# Or manually:
cd backend && go run main.go
```

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

### Authentication Flow

1. **Frontend**: User authenticates via Clerk UI components
2. **Frontend**: Obtains JWT token from Clerk
3. **Frontend**: Sends API requests with `Authorization: Bearer <jwt>` header
4. **Backend**: Middleware verifies JWT with Clerk API
5. **Backend**: Extracts user info and stores in request context
6. **Backend**: Executes business logic with authenticated user data

## Setup & Configuration

### 1. Clerk Dashboard Setup

1. **Create Account**: Sign up at [clerk.com](https://clerk.com)
2. **Create Application**: Set up a new application
3. **Configure Providers**: Enable authentication methods (Email, Google, etc.)
4. **Get API Keys**: Copy from Dashboard → API Keys

### 2. API Keys Configuration

```bash
# Development Environment (.env.dev)
CLERK_SECRET_KEY=sk_test_your_development_key
CLERK_PUBLISHABLE_KEY=pk_test_your_publishable_key

# Production Environment (.env)
CLERK_SECRET_KEY=sk_live_your_production_key
CLERK_PUBLISHABLE_KEY=pk_live_your_publishable_key
```

### 3. Backend Middleware

The authentication middleware is automatically configured:

```go
// Middleware automatically handles:
// - JWT extraction from Authorization header
// - Token verification with Clerk API
// - User data extraction and context setup
protected.Use(middleware.ClerkAuth(cfg.ClerkSecretKey))
```

#### Context Data Available

After authentication, the following data is available in request context:

- `clerk_user_id` - Clerk user ID
- `clerk_email` - Primary email address
- `clerk_first_name` - First name
- `clerk_last_name` - Last name
- `clerk_image_url` - Avatar URL
- `user_id` - Local database user ID (after sync)

## User Synchronization

### Manual Synchronization

After user login on the frontend, call the sync endpoint:

```javascript
// Frontend - after successful Clerk authentication
const response = await fetch('/api/v1/auth/sync', {
	method: 'POST',
	headers: {
		Authorization: `Bearer ${clerkToken}`,
		'Content-Type': 'application/json',
	},
});
```

### Automatic Synchronization (Webhooks)

Webhooks automatically sync user data when changes occur in Clerk.

**Supported Events**:

- `user.created` - New user registration
- `user.updated` - Profile updates (name, email, avatar)
- `user.deleted` - Account deletion

**User Data Mapping**:

| Clerk Field                           | Thothix Field | Notes                       |
| ------------------------------------- | ------------- | --------------------------- |
| `id`                                  | `id`          | Primary key                 |
| `primary_email_address.email_address` | `email`       | Primary email               |
| `first_name + last_name`              | `name`        | Full name                   |
| `image_url`                           | `avatar_url`  | User avatar                 |
| `username`                            | `name`        | Fallback if name missing    |
| (default)                             | `system_role` | Always `user` for new users |

### Manual User Import

For bulk operations or initial setup:

```bash
# Import all users from Clerk to local database
curl -X POST http://localhost:30000/api/v1/auth/import-users \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

## Local Development with Webhooks

### 1. Ngrok Setup

**Install Ngrok**:

```bash
# Download from https://ngrok.com/download
# Or via package managers:
choco install ngrok              # Windows (Chocolatey)
brew install ngrok               # macOS (Homebrew)
npm install -g @ngrok/ngrok      # Node.js
```

**Configure Ngrok**:

```bash
# Get auth token from ngrok.com dashboard
ngrok config add-authtoken YOUR_NGROK_TOKEN
```

### 2. Start Development Environment

**Option A: Automated Script (Windows)**:

```bash
# Starts backend + ngrok with pre-configured URL
scripts\start-local-dev.bat
```

**Option B: Manual Setup**:

```bash
# Terminal 1: Start backend
cd backend
go run main.go

# Terminal 2: Start ngrok tunnel
ngrok http --url=flying-mullet-socially.ngrok-free.app 30000
```

### 3. Webhook Configuration in Clerk Dashboard

1. **Navigate**: Go to [Clerk Dashboard](https://dashboard.clerk.com)
2. **Select Project**: Choose your application
3. **Configure Webhook**:

   - **URL**: `https://flying-mullet-socially.ngrok-free.app/api/v1/auth/webhooks/clerk`
   - **Events**: `user.created`, `user.updated`, `user.deleted`
   - **Version**: `v1`

4. **Copy Signing Secret**: Add to your `.env` file:
   ```bash
   CLERK_WEBHOOK_SECRET=whsec_your_webhook_signing_secret
   ```

### 4. Testing Setup

**Verify Installation**:

```bash
# Test backend health
curl http://localhost:30000/health

# Test ngrok tunnel (if running)
curl https://flying-mullet-socially.ngrok-free.app/health
```

**Verification URLs**:

- **Local Health**: http://localhost:30000/health
- **Ngrok Health**: https://flying-mullet-socially.ngrok-free.app/health
- **Swagger (Local)**: http://localhost:30000/swagger/index.html
- **Swagger (Ngrok)**: https://flying-mullet-socially.ngrok-free.app/swagger/index.html

## API Endpoints

### Authentication Endpoints

| Method | Endpoint                      | Description                  | Auth Required          |
| ------ | ----------------------------- | ---------------------------- | ---------------------- |
| `POST` | `/api/v1/auth/sync`           | Sync current user from Clerk | Yes (JWT)              |
| `GET`  | `/api/v1/auth/me`             | Get current user info        | Yes (JWT)              |
| `POST` | `/api/v1/auth/webhooks/clerk` | Webhook for auto-sync        | No (Webhook signature) |
| `POST` | `/api/v1/auth/import-users`   | Import all users from Clerk  | Yes (Admin)            |

### Protected Endpoints

All other API endpoints require JWT authentication:

```bash
# Example authenticated request
curl -X GET "http://localhost:30000/api/v1/users" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Webhook Payload Examples

**User Created**:

```json
{
	"type": "user.created",
	"data": {
		"id": "user_2abc123def456",
		"email_addresses": [{ "email_address": "user@example.com" }],
		"first_name": "John",
		"last_name": "Doe",
		"image_url": "https://img.clerk.com/avatar.jpg"
	}
}
```

## Testing & Debugging

### 1. API Testing

**Using Swagger UI** (recommended):

```bash
# Open in browser
http://localhost:30000/swagger/index.html
```

**Manual testing examples**:

# Health check

curl http://localhost:30000/health

# User sync (requires JWT)

curl -X POST http://localhost:30000/api/v1/auth/sync \
 -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Webhook simulation

curl -X POST http://localhost:30000/api/v1/auth/webhooks/clerk \
 -H "Content-Type: application/json" \
 -d '{"type":"user.created","data":{"id":"test_user","email_addresses":[{"email_address":"test@example.com"}]}}'

````

### 2. Debug Configuration

Enable detailed logging:

```bash
# Environment variables
DEBUG=true
LOG_LEVEL=debug
GIN_MODE=debug
````

### 3. Database Verification

```bash
# Check synchronized users
docker-compose exec postgres psql -U postgres -d thothix-db \
  -c "SELECT id, email, name, created_at FROM users ORDER BY created_at DESC LIMIT 10;"
```

### 4. Log Monitoring

```bash
# Application logs
docker-compose logs thothix-api --tail=50 --follow

# Real-time backend logs (development)
cd backend && go run main.go
```

## Frontend Integration

### Nuxt.js Setup

**Configuration**:

```javascript
// nuxt.config.js
export default {
	runtimeConfig: {
		public: {
			clerkPublishableKey: process.env.CLERK_PUBLISHABLE_KEY,
		},
	},
};
```

**Plugin Setup**:

```javascript
// plugins/clerk.client.js
import { ClerkProvider } from '@clerk/vue';

export default defineNuxtPlugin((nuxtApp) => {
	const config = useRuntimeConfig();

	nuxtApp.vueApp.use(ClerkProvider, {
		publishableKey: config.public.clerkPublishableKey,
	});
});
```

### Authentication Composable

```javascript
// composables/useAuth.js
export const useAuth = () => {
	const { isSignedIn, user, getToken } = useClerk();

	const syncUser = async () => {
		if (!isSignedIn.value) return;

		const token = await getToken();
		await $fetch('/api/v1/auth/sync', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${token}`,
			},
		});
	};

	const apiCall = async (url, options = {}) => {
		const token = await getToken();
		return $fetch(url, {
			...options,
			headers: {
				Authorization: `Bearer ${token}`,
				...options.headers,
			},
		});
	};

	return {
		isSignedIn,
		user,
		syncUser,
		apiCall,
	};
};
```

### Authentication Flow Example

```javascript
// pages/login.vue
<template>
  <div>
    <SignIn v-if="!isSignedIn" @success="onSignInSuccess" />
    <UserProfile v-else />
  </div>
</template>

<script setup>
const { isSignedIn, syncUser } = useAuth()

const onSignInSuccess = async () => {
  // Automatically sync user to backend after login
  await syncUser()
  await navigateTo('/dashboard')
}
</script>
```

## Security & Best Practices

### Production Security

- ✅ **Never expose secrets**: Keep `CLERK_SECRET_KEY` on backend only
- ✅ **Use HTTPS**: Always use secure connections in production
- ✅ **Validate tokens**: Backend must verify all JWT tokens with Clerk
- ✅ **Rate limiting**: Implement API rate limiting for protection
- ✅ **Monitor webhooks**: Log and monitor webhook calls for security

### Environment Configuration

```bash
# Production settings
CLERK_SECRET_KEY=sk_live_your_live_secret_key
ENVIRONMENT=production
GIN_MODE=release
```

### Best Practices

1. **Token Management**: Frontend should handle token refresh automatically
2. **Error Handling**: Implement proper error handling for auth failures
3. **Session Management**: Use Clerk's built-in session management
4. **User Roles**: Manage roles separately in your application database
5. **Webhook Security**: Always verify webhook signatures in production

## Troubleshooting

### Common Issues & Solutions

#### 1. Invalid JWT Token

**Error**: `Invalid token: clerk verification failed with status: 401`

**Solutions**:

- ✅ Verify `CLERK_SECRET_KEY` is correct for your environment
- ✅ Check token hasn't expired (frontend should refresh automatically)
- ✅ Ensure you're using the correct environment keys (dev vs prod)
- ✅ Verify token format: `Authorization: Bearer <token>`

#### 2. User ID Not Found

**Error**: `Clerk user ID not found`

**Solutions**:

- ✅ Verify `ClerkAuth` middleware is properly configured
- ✅ Check endpoint is protected by authentication middleware
- ✅ Ensure token is valid and not expired
- ✅ Call `/api/v1/auth/sync` after login

#### 3. Synchronization Issues

**Error**: User not created in local database

**Solutions**:

- ✅ Always call `/api/v1/auth/sync` after each login
- ✅ Verify database connection and configuration
- ✅ Check application logs for database errors
- ✅ Ensure user data mapping is correct

#### 4. Webhook Not Working

**Issues**: Automatic synchronization doesn't happen

**Solutions**:

- ✅ Verify webhook URL in Clerk Dashboard is correct
- ✅ Check webhook endpoint is public (not protected by auth)
- ✅ Confirm events are configured: `user.created`, `user.updated`, `user.deleted`
- ✅ Verify `CLERK_WEBHOOK_SECRET` is set correctly
- ✅ Test webhook endpoint manually with curl
- ✅ Check ngrok tunnel is active and accessible

#### 5. Local Development Issues

**Common Problems**:

**Ngrok not working**:

```bash
# Check if ngrok is installed
ngrok version

# Test tunnel manually
ngrok http 30000

# Use custom URL for consistency
ngrok http --url=flying-mullet-socially.ngrok-free.app 30000
```

**Backend not accessible**:

```bash
# Check if backend is running
curl http://localhost:30000/health

# Start backend manually
cd backend && go run main.go
```

**Webhook signature validation fails**:

```bash
# Verify webhook secret is correctly set
echo $CLERK_WEBHOOK_SECRET

# Check webhook endpoint logs
# Should show incoming webhook requests and validation results
```

### Debug Commands

```bash
# Test backend health
curl http://localhost:30000/health

# Check if backend is running
cd backend && go run main.go

# Test ngrok tunnel (if using)
curl https://flying-mullet-socially.ngrok-free.app/health
```

### Getting Help

1. **Clerk Documentation**: [clerk.com/docs](https://clerk.com/docs)
2. **API Reference**: Check Swagger UI at `/swagger/index.html`
3. **Development Team**: Contact internal development team
4. **Community**: Clerk Discord/GitHub for community support

## Conclusion

This integration provides a robust, secure authentication system for Thothix with:

- 🔐 **Enterprise-grade security** with JWT verification
- 🔄 **Automatic user synchronization** via webhooks
- 🛠️ **Developer-friendly** local development setup
- 🧪 **Comprehensive testing** tools and scripts
- 📚 **Complete documentation** and troubleshooting guides

The system is designed to handle both development and production environments seamlessly, with proper error handling, security measures, and monitoring capabilities.
