# Clerk Authentication Integration

## Overview

Thothix integra con Clerk per gestire l'autenticazione utenti in modo sicuro e scalabile. Questa guida spiega come configurare e utilizzare l'integrazione.

## Architettura di Autenticazione

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

## Configurazione Clerk

### 1. Setup Account Clerk

1. Crea un account su [clerk.com](https://clerk.com)
2. Crea una nuova applicazione
3. Configura i provider di autenticazione (Email, Google, etc.)

### 2. Ottieni le Chiavi API

```bash
# Da Clerk Dashboard → API Keys
CLERK_PUBLISHABLE_KEY=pk_test_...
CLERK_SECRET_KEY=sk_test_...
```

### 3. Configurazione Environment Variables

```bash
# .env file
CLERK_SECRET_KEY=sk_test_your_secret_key_here
```

## Configurazione Backend

### 1. Middleware di Autenticazione

Il middleware `ClerkAuth` verifica automaticamente i token JWT:

```go
// Configurazione automatica via environment
protected.Use(middleware.ClerkAuth(cfg.ClerkSecretKey))
```

#### Come Funziona il Middleware

1. **Estrazione Token**: Estrae il JWT dall'header `Authorization: Bearer <token>`
2. **Verifica con Clerk**: Chiama l'API Clerk `https://api.clerk.com/v1/users/me`
3. **Validazione**: Verifica che l'utente abbia ID ed email validi
4. **Context Setup**: Imposta nel context Gin:
   - `clerk_user_id` - ID utente Clerk
   - `clerk_email` - Email principale  
   - `clerk_first_name` - Nome
   - `clerk_last_name` - Cognome
   - `clerk_image_url` - Avatar URL
   - `user_id` - Per i hook BaseModel

#### Struttura Dati Clerk

```go
type ClerkUser struct {
    ID                   string                 `json:"id"`
    Username             *string                `json:"username"`
    FirstName            *string                `json:"first_name"`
    LastName             *string                `json:"last_name"`
    ImageURL             string                 `json:"image_url"`
    PrimaryEmailAddress  *ClerkEmailAddress     `json:"primary_email_address"`
    // Altri campi...
}
```

### 2. Endpoint Disponibili

#### Autenticazione
- `POST /api/v1/auth/sync` - Sincronizza utente da Clerk al DB locale
- `GET /api/v1/auth/me` - Ottieni utente corrente
- `POST /api/v1/auth/webhooks/clerk` - Webhook per sincronizzazione automatica

#### Flusso di Autenticazione
1. **Frontend**: Autentica con Clerk e ottiene JWT
2. **Frontend**: Invia richiesta con header `Authorization: Bearer <jwt>`
3. **Backend**: Middleware verifica JWT con Clerk
4. **Backend**: Estrae info utente e le mette nel context
5. **Backend**: Esegue business logic con utente autenticato

## Sincronizzazione Utenti

### Sincronizzazione Manuale

Il frontend deve chiamare `/api/v1/auth/sync` dopo il login:

```javascript
// Frontend (dopo login con Clerk)
const response = await fetch('/api/v1/auth/sync', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${clerkToken}`,
    'Content-Type': 'application/json'
  }
});
```

### Sincronizzazione Automatica (Webhook)

Configura webhook in Clerk Dashboard:

1. **URL**: `https://your-domain.com/api/v1/auth/webhooks/clerk`
2. **Eventi**: 
   - `user.created`
   - `user.updated`
   - `user.deleted`

Il sistema sincronizzerà automaticamente:
- ✅ Creazione nuovi utenti
- ✅ Aggiornamento dati esistenti (email, nome, avatar)
- ✅ Cancellazione utenti

## Mapping Dati Utente

### Da Clerk a Thothix

| Campo Clerk | Campo Thothix | Note |
|-------------|---------------|------|
| `id` | `id` | Chiave primaria |
| `primary_email_address.email_address` | `email` | Email principale |
| `first_name + last_name` | `name` | Nome completo |
| `image_url` | `avatar_url` | Avatar utente |
| `username` | `name` | Fallback se nome mancante |
| (default) | `system_role` | Sempre `user` per nuovi utenti |

### Gestione Ruoli

- **Nuovi utenti**: Ruolo `user` di default
- **Admin**: Deve essere assegnato manualmente via API
- **Ruoli progetto**: Gestiti separatamente in `project_members`

## Testing dell'Integrazione

### 1. Test Manuale

```bash
# 1. Ottieni token da Clerk (frontend)
TOKEN="your_clerk_jwt_token"

# 2. Test sincronizzazione
curl -X POST "http://localhost:30000/api/v1/auth/sync" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

# 3. Test ottenimento utente corrente
curl -X GET "http://localhost:30000/api/v1/auth/me" \
  -H "Authorization: Bearer $TOKEN"
```

### 2. Test Webhook

```bash
# Simula webhook di creazione utente
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

## Debug e Troubleshooting

### Log di Debug

Per abilitare i log di debug dell'autenticazione:

```bash
# Nel docker-compose.yml o variabili ambiente
DEBUG=true
LOG_LEVEL=debug
GIN_MODE=debug
```

### Problemi Comuni

#### 1. Token JWT Invalido

**Errore**: `Invalid token: clerk verification failed with status: 401`

**Soluzioni**:
- Verifica che `CLERK_SECRET_KEY` sia configurata correttamente
- Controlla che il token non sia scaduto (il frontend deve fare refresh)
- Assicurati di usare il token corretto per l'environment (development/production)

#### 2. User ID Non Trovato

**Errore**: `Clerk user ID not found`

**Soluzioni**:
- Verifica che il middleware `ClerkAuth` sia configurato correttamente
- Controlla che l'endpoint sia protetto dal middleware
- Assicurati che il token sia valido

#### 3. Problemi di Sincronizzazione

**Errore**: Utente non creato nel database locale

**Soluzioni**:
- Chiama `/api/v1/auth/sync` dopo ogni login
- Verifica la connessione al database
- Controlla i log per errori di database

#### 4. Webhook Non Funziona

**Problemi**: Sincronizzazione automatica non avviene

**Soluzioni**:
- Verifica URL webhook in Clerk Dashboard
- Controlla che l'endpoint sia pubblico (non protetto da auth)
- Verifica eventi configurati: `user.created`, `user.updated`, `user.deleted`

### Testare l'Integrazione

```bash
# 1. Test del middleware di autenticazione
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:30000/api/v1/auth/me

# 2. Test sincronizzazione manuale
curl -X POST \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     http://localhost:30000/api/v1/auth/sync

# 3. Test webhook (simulazione)
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"type":"user.created","data":{"id":"user_123","email_addresses":[{"email_address":"test@example.com"}]}}' \
     http://localhost:30000/api/v1/auth/webhooks/clerk
```

### Monitoraggio

Per monitorare l'integrazione in produzione:

```bash
# Verifica utenti sincronizzati
docker-compose exec postgres psql -U postgres -d thothix-db \
  -c "SELECT id, email, name, created_at FROM users ORDER BY created_at DESC LIMIT 10;"

# Verifica log dell'applicazione
docker-compose logs thothix-api --tail=50 --follow
```

## Sicurezza

### Best Practices

- ✅ Mai esporre `CLERK_SECRET_KEY` nel frontend
- ✅ Usa HTTPS in produzione
- ✅ Valida sempre i token dal backend
- ✅ Implementa rate limiting
- ✅ Monitora webhook per attacchi

### Configurazione Produzione

```bash
# Variabili di produzione
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

## Conclusione

L'integrazione con Clerk fornisce:
- ✅ Autenticazione sicura e scalabile
- ✅ Sincronizzazione automatica utenti
- ✅ Gestione sessioni robusta
- ✅ Support multi-provider (Google, GitHub, etc.)
- ✅ UI components pre-built per il frontend

Per assistenza tecnica, consulta la [documentazione Clerk](https://clerk.com/docs) o il team di sviluppo.
