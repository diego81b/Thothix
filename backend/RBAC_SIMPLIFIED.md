# Simplified RBAC System - Thothix

## Overview

The Thothix role and permission management system has been simplified to include four main roles, each with specific permissions and limitations.

## System Roles

### 1. Admin

- **Permissions**: Can manage the entire system
- **Description**: Full access to all functionalities
- **Limitations**: None

### 2. Manager

- **Permissions**: Can manage everything except user management
- **Description**: Can create and manage projects, channels, messages
- **Limitations**: Cannot assign roles or manage users

### 3. User

- **Permissions**: Can participate in projects and channels they've been added to
- **Description**: Can read assigned projects, participate in channels, create 1:1 chats
- **Limitations**: Cannot create projects, only participate in assigned ones

### 4. External

- **Permissions**: Can only participate in public conversations
- **Description**: Limited access to public channels only
- **Limitations**: Cannot access private channels or projects

## Public vs Private Channel Strategy

### Public Channels

- **Definition**: Channels with no explicit members in the `channel_members` table
- **Access**: Accessible to all authenticated users (except External who have project limitations)

### Private Channels

- **Definition**: Channels with at least one row in the `channel_members` table
- **Access**: Only to explicit members or users with Admin/Manager roles

### Computed `IsPrivate` Field

The `IsPrivate` field is calculated dynamically:

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

## Updated Data Models

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

## Access Controls

### Channel Logic

1. **External**: Only public channels
2. **User**: Public channels + private channels they're a member of (if they have project access)
3. **Manager/Admin**: All channels

### Project Logic

1. **External**: Only projects they're explicitly a member of
2. **User**: Only projects they're explicitly a member of  
3. **Manager/Admin**: All projects

### Message Logic

1. **Channel Messages**: Must have access to the channel
2. **Direct Messages**: All can create (except External)

## Main API Endpoints

### Channels

- `GET /api/v1/channels` - List accessible channels
- `POST /api/v1/channels` - Create new channel (Manager/Admin)
- `GET /api/v1/channels/{id}` - Channel details
- `POST /api/v1/channels/{id}/join` - Join public channel

### Messages

- `GET /api/v1/channels/{id}/messages` - Channel messages
- `POST /api/v1/channels/{id}/messages` - Send message to channel
- `POST /api/v1/messages/direct` - Send direct message

### Roles (Admin Only)

- `POST /api/v1/roles` - Assign role
- `DELETE /api/v1/roles/{roleId}` - Revoke role
- `GET /api/v1/users/{userId}/roles` - List user roles

## Security Middleware

The system uses middleware to control:

1. **RequirePermission**: Verifies specific permissions
2. **RequireSystemRole**: Verifies minimum required role  
3. **RequireProjectAccess**: Verifies project access
4. **RequireChannelAccess**: Verifies channel access

## Database Schema

The main tables are:

- `users` - with `system_role` field
- `channels` - without `type` field (calculated dynamically)
- `channel_members` - defines private channels
- `messages` - linked to channels or users for DMs
- `user_roles` - for future more granular roles (not used in simplified system)

## Migration

To apply changes to the database:

1. Remove the `type` field from the `channels` table
2. Ensure all users have a valid `system_role`
3. Existing relationships in `channel_members` will automatically define private channels
