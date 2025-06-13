# Backend di Thothix - Messaggistica Aziendale

Questo è il backend della piattaforma di messaggistica aziendale Thothix, sviluppato in Go con il framework Gin e integrazione con Clerk per l'autenticazione.

## Modelli di Dati

Il backend utilizza i seguenti modelli principali:

- **User**: Utenti della piattaforma (sincronizzati con Clerk)
- **Project**: Progetti aziendali 
- **ProjectMember**: Membri assegnati ai progetti con ruoli
- **Channel**: Canali di comunicazione all'interno dei progetti
- **ChannelMember**: Membri dei canali
- **Message**: Messaggi (in canali o diretti tra utenti)
- **File**: File condivisi nei progetti

## Struttura del progetto

```
backend/
├── internal/
│   ├── config/         # Configurazione applicazione
│   ├── database/       # Setup e migrazioni database
│   ├── handlers/       # Handler HTTP per le API
│   ├── middleware/     # Middleware personalizzati
│   ├── models/         # Modelli di dati
│   └── router/         # Setup delle rotte
├── main.go            # Entry point dell'applicazione
├── go.mod            # Dipendenze Go
├── Dockerfile        # Configurazione Docker
└── .env.example      # Template variabili d'ambiente
```

## Avvio veloce

1. **Configura le variabili d'ambiente**:
   ```bash
   cp .env.example .env
   # Modifica .env con le tue chiavi Clerk
   ```

2. **Avvia con Docker**:
   ```bash
   cd ..
   docker-compose up -d
   ```

3. **Verifica che funzioni**:
   ```bash
   curl http://localhost:30000/health
   ```

## API Endpoints

- **Health**: `GET /health`
- **Swagger**: `GET /swagger/index.html`
- **API**: `GET /api/v1/auth/sync`, `GET /api/v1/auth/me`
- **Users**: `GET /api/v1/users`, `PUT /api/v1/users/me`
- **Projects**: `GET /api/v1/projects` (TODO)
- **Channels**: `GET /api/v1/channels` (TODO)
- **Messages**: `GET /api/v1/channels/{id}/messages` (TODO)

## Sviluppo

Per sviluppare in locale:

```bash
# Installa le dipendenze
go mod tidy

# Genera documentazione Swagger
go install github.com/swaggo/swag/cmd/swag@latest
swag init

# Avvia il server
go run main.go
```

## Autenticazione con Clerk

Il backend utilizza Clerk per l'autenticazione. Gli utenti vengono sincronizzati automaticamente nel database locale al primo accesso.

Flow di autenticazione:
1. L'utente si autentica tramite Clerk nel frontend
2. Il frontend invia il token Clerk nelle richieste API
3. Il middleware `ClerkAuth` verifica il token con Clerk
4. L'endpoint `/auth/sync` crea/aggiorna l'utente nel database locale

## TODO

- [ ] Implementare handlers per progetti
- [ ] Implementare handlers per canali  
- [ ] Implementare handlers per messaggi
- [ ] Aggiungere WebSocket per real-time
- [ ] Aggiungere supporto per file upload (MinIO)
- [ ] Aggiungere test unitari
- [ ] Migliorare documentazione API
