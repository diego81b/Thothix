# Sistema RBAC Semplificato - Thothix

## Panoramica

Il sistema di gestione ruoli e permessi di Thothix è stato semplificato per includere quattro ruoli principali, ognuno con permessi specifici e limitazioni.

## Ruoli del Sistema

### 1. Admin

- **Permessi**: Può gestire tutto il sistema
- **Descrizione**: Accesso completo a tutte le funzionalità
- **Limitazioni**: Nessuna

### 2. Manager

- **Permessi**: Può gestire tutto tranne la gestione degli utenti
- **Descrizione**: Può creare e gestire progetti, canali, messaggi
- **Limitazioni**: Non può assegnare ruoli o gestire utenti

### 3. User

- **Permessi**: Può partecipare ai progetti e canali in cui è stato aggiunto
- **Descrizione**: Può leggere progetti assegnati, partecipare a canali, creare chat 1:1
- **Limitazioni**: Non può creare progetti, solo partecipare a quelli assegnati

### 4. External

- **Permessi**: Può solo partecipare alle conversazioni pubbliche
- **Descrizione**: Accesso limitato ai soli canali pubblici
- **Limitazioni**: Non può accedere a canali privati o progetti

## Strategia per Canali Pubblici vs Privati

### Canali Pubblici

- **Definizione**: Canali senza membri espliciti nella tabella `channel_members`
- **Accesso**: Accessibili a tutti gli utenti autenticati (tranne External che hanno limitazioni sui progetti)

### Canali Privati

- **Definizione**: Canali con almeno una riga nella tabella `channel_members`
- **Accesso**: Solo ai membri espliciti o agli utenti con ruoli Admin/Manager

### Campo Calcolato `IsPrivate`

Il campo `IsPrivate` viene calcolato dinamicamente:

```go
func (c *Channel) LoadIsPrivate(db *gorm.DB) error {
    var count int64
    err := db.Model(&ChannelMember{}).Where("channel_id = ?", c.ID).Count(&count).Error
    if err != nil {
        return err
    }
    c.IsPrivate = count > 0
    return nil
}
```

## Modelli Dati Aggiornati

### Channel

```go
type Channel struct {
    BaseModel
    Name      string      `json:"name"`
    ProjectID string      `json:"project_id"`
    Project   Project     `gorm:"foreignKey:ProjectID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"project,omitempty"`
    Members   []User      `gorm:"many2many:channel_members" json:"members,omitempty"`
    IsPrivate bool        `json:"is_private" gorm:"-"` // Computed field, not stored in DB
}
```

### Message

```go
type Message struct {
    BaseModel
    SenderID   string   `json:"sender_id"`
    Sender     User     `gorm:"foreignKey:SenderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"sender,omitempty"`
    ChannelID  *string  `json:"channel_id,omitempty"`
    Channel    *Channel `gorm:"foreignKey:ChannelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"channel,omitempty"`
    ReceiverID *string  `json:"receiver_id,omitempty"`
    Receiver   *User    `gorm:"foreignKey:ReceiverID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"receiver,omitempty"`
    Content    string   `json:"content"`
}
```

### User

```go
type User struct {
    BaseModel
    Email      string   `json:"email"`
    Name       string   `json:"name"`
    AvatarURL  string   `json:"avatar_url"`
    SystemRole RoleType `json:"system_role" gorm:"default:'user'"` // Default system role
}
```

## Controlli di Accesso

### Logica per Canali

1. **External**: Solo canali pubblici
2. **User**: Canali pubblici + canali privati di cui è membro (se ha accesso al progetto)
3. **Manager/Admin**: Tutti i canali

### Logica per Progetti

1. **External**: Solo progetti di cui è membro esplicito
2. **User**: Solo progetti di cui è membro esplicito  
3. **Manager/Admin**: Tutti i progetti

### Logica per Messaggi

1. **Messaggi in Canali**: Devono avere accesso al canale
2. **Messaggi Diretti**: Tutti possono creare (tranne External)

## API Endpoints Principali

### Canali

- `GET /api/v1/channels` - Lista canali accessibili
- `POST /api/v1/channels` - Crea nuovo canale (Manager/Admin)
- `GET /api/v1/channels/{id}` - Dettagli canale
- `POST /api/v1/channels/{id}/join` - Unisciti a canale pubblico

### Messaggi

- `GET /api/v1/channels/{id}/messages` - Messaggi del canale
- `POST /api/v1/channels/{id}/messages` - Invia messaggio al canale
- `POST /api/v1/messages/direct` - Invia messaggio diretto

### Ruoli (Solo Admin)

- `POST /api/v1/roles` - Assegna ruolo
- `DELETE /api/v1/roles/{roleId}` - Revoca ruolo
- `GET /api/v1/users/{userId}/roles` - Lista ruoli utente

## Middleware di Sicurezza

Il sistema utilizza middleware per controllare:

1. **RequirePermission**: Verifica permessi specifici
2. **RequireSystemRole**: Verifica ruolo minimo richiesto  
3. **RequireProjectAccess**: Verifica accesso al progetto
4. **RequireChannelAccess**: Verifica accesso al canale

## Database Schema

Le tabelle principali sono:

- `users` - con campo `system_role`
- `channels` - senza campo `type` (calcolato dinamicamente)
- `channel_members` - definisce i canali privati
- `messages` - collegati a canali o utenti per DM
- `user_roles` - per ruoli futuri più granulari (non utilizzato nel sistema semplificato)

## Migrazione

Per applicare i cambiamenti al database:

1. Rimuovere il campo `type` dalla tabella `channels`
2. Assicurarsi che tutti gli utenti abbiano un `system_role` valido
3. Le relazioni esistenti in `channel_members` definiranno automaticamente i canali privati
