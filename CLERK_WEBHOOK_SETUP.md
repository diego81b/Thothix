# 🔗 Configurazione Clerk Webhook per Sviluppo Locale

## 📋 **Setup Ngrok per Tunnel Locale**

### **Step 1: Installa Ngrok**

```bash
# Download da https://ngrok.com/download
# Oppure con Chocolatey:
choco install ngrok

# Oppure con npm:
npm install -g @ngrok/ngrok
```

### **Step 2: Configura Ngrok**

```bash
# Registrati su ngrok.com e ottieni il token
ngrok config add-authtoken YOUR_NGROK_TOKEN

# Avvia tunnel verso il tuo backend locale
ngrok http 30000
```

Output di esempio:

```
Session Status                online
Account                       Your Name (Plan: Free)
Version                       3.0.0
Region                        United States (us)
Latency                       150ms
Web Interface                 http://127.0.0.1:4040
Forwarding                    https://abc123.ngrok.io -> http://localhost:30000
```

## 🎯 **Configurazione Webhook su Clerk Dashboard**

### **Step 3: Configura Webhook su Clerk**

1. **Vai su [Clerk Dashboard](https://dashboard.clerk.com)**
2. **Seleziona il tuo progetto**
3. **Vai su "Webhooks" nel menu laterale**
4. **Clicca "Add Endpoint"**

**Configurazione:**

- **Endpoint URL**: `https://abc123.ngrok.io/api/v1/auth/webhooks/clerk`
- **Events to listen for**:
  - `user.created`
  - `user.updated`
  - `user.deleted`
- **Version**: `v1`

### **Step 4: Copia il Signing Secret**

Clerk genererà un **Webhook Signing Secret** - copialo!

## 🔧 **Aggiorna Configurazione Backend**

### **Step 5: Aggiungi il Webhook Secret al .env**

```env
# Clerk Authentication
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_your_key
CLERK_SECRET_KEY=sk_test_your_key
CLERK_WEBHOOK_SECRET=whsec_your_webhook_signing_secret

# Backend Configuration
PORT=30000
BACKEND_URL=https://abc123.ngrok.io
```

## 🧪 **Test del Webhook**

### **Step 6: Testa il Webhook**

1. **Avvia il backend**: `go run main.go`
2. **Avvia ngrok**: `ngrok http 30000`
3. **Nel Clerk Dashboard, vai su "Webhooks"**
4. **Clicca "Test" sul tuo webhook**
5. **Verifica i log del backend**

## 🔄 **Workflow Completo**

### **Sviluppo Locale:**

```bash
# Terminal 1: Avvia backend
cd backend
go run main.go

# Terminal 2: Avvia ngrok
ngrok http 30000

# Terminal 3: Avvia frontend (se hai)
cd frontend
npm run dev
```

### **Aggiornamenti Automatici:**

- Utente si registra su Clerk → Webhook → Backend crea utente nel DB
- Utente aggiorna profilo → Webhook → Backend aggiorna DB
- Utente cancella account → Webhook → Backend gestisce cancellazione

## 🔑 **Alternative per Sviluppo**

### **Opzione 1: Polling invece di Webhook**

Se ngrok è problematico, puoi usare polling:

```bash
# Chiama periodicamente l'API Clerk per sincronizzare
curl -X POST http://localhost:30000/api/v1/auth/sync \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### **Opzione 2: Import Manuale Utenti**

```bash
# Endpoint per importare tutti gli utenti da Clerk
curl -X POST http://localhost:30000/api/v1/auth/import-users \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

## 🎯 **Configurazione Specifica per Thothix**

### **URL Ngrok Configurato:**

- **Ngrok URL**: `https://flying-mullet-socially.ngrok-free.app`
- **Webhook URL**: `https://flying-mullet-socially.ngrok-free.app/api/v1/auth/webhooks/clerk`

### **Comando Ngrok:**

```bash
ngrok http --url=flying-mullet-socially.ngrok-free.app 30000
```

### **Link Rapidi:**

- **Health Check**: https://flying-mullet-socially.ngrok-free.app/health
- **Swagger UI**: https://flying-mullet-socially.ngrok-free.app/swagger/index.html
- **Webhook Endpoint**: https://flying-mullet-socially.ngrok-free.app/api/v1/auth/webhooks/clerk
