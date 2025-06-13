# Database Migration Notes

## Schema Alignment with Go Models

This document tracks the database schema alignment with the Go models in the Thothix application.

### Completed Migrations

#### 2025-06-13: Added system_role to users table

**Problem**: The `users` table was missing the `system_role` field required by the User model.

**Solution**: Added the `system_role` column to the users table:

```sql
ALTER TABLE users ADD COLUMN system_role text DEFAULT 'user';
```

**Verification**:

- The field was successfully added with the correct type and default value
- All existing users will have the default role of 'user'
- The backend starts without errors and can handle the new field

#### 2025-06-13: BaseModel alignment - Added updated_by field

**Problem**:

1. The `BaseModel` Go struct was missing the `updated_by` field for tracking who made updates
2. Several database tables were missing the `updated_by` column
3. The `channel_members` table was missing the `created_at` column

**Solution**:

1. Updated `BaseModel` Go struct to include `UpdatedBy string` field
2. Updated `BeforeUpdate` hook to automatically set `updated_by` from context
3. Added missing columns to all tables:

```sql
-- Add updated_by to all tables
ALTER TABLE users ADD COLUMN updated_by text DEFAULT '';
ALTER TABLE projects ADD COLUMN updated_by text DEFAULT '';
ALTER TABLE project_members ADD COLUMN updated_by text DEFAULT '';
ALTER TABLE channels ADD COLUMN updated_by text DEFAULT '';
ALTER TABLE messages ADD COLUMN updated_by text DEFAULT '';
ALTER TABLE files ADD COLUMN updated_by text DEFAULT '';
ALTER TABLE channel_members ADD COLUMN updated_by text DEFAULT '';

-- Add missing created_at to specific tables
ALTER TABLE channel_members ADD COLUMN created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE project_members ADD COLUMN created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP;
```

**Verification**:

- All tables now have the complete BaseModel fields: `id`, `created_by`, `created_at`, `updated_by`, `updated_at`
- The `project_members` table correctly has both `created_at` (from BaseModel) and `joined_at` (specific field)
- Backend starts without errors and can handle all new fields
- GORM hooks will automatically populate `updated_by` on updates

### Current Schema Status

All tables are now aligned with the Go models:

✅ **users** - Complete with system_role field  
✅ **projects** - Aligned  
✅ **project_members** - Aligned  
✅ **channels** - Aligned (type field removed as expected)  
✅ **channel_members** - Aligned  
✅ **messages** - Aligned  
✅ **files** - Aligned (has some extra fields but no conflicts)  

### Notes

- The `channels` table no longer has a `type` field, as public/private status is now determined by the presence of channel_members rows
- The `files` table has some extra fields (`uploaded_by`, `uploaded_at`) that don't conflict with the Go model
- All foreign key constraints are properly set up
- The simplified RBAC system (Admin, Manager, User, External) is fully implemented

### Future Migrations

For future schema changes, consider:

1. Creating proper migration scripts
2. Adding database migration automation to the startup process
3. Version tracking for schema changes
