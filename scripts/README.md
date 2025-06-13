# 📁### 🚀 **Development & Build**

- **`dev.bat`** - Main development script (Windows)

  - Handles formatting, linting, and pre-commit workflow
  - Uses `gofmt` for simple, reliable formatting
  - Integrated with Windows development environment

- **`pre-commit.bat`** - Manual pre-commit hook backup
  - Fallback for direct script execution
  - Integrates with Git hooksirectory - Simplified & Clean

## ✅ **Current Scripts Overview**

### � **Development & Build**

- **`dev.bat`** / **`dev.sh`** - Main development script

  - Handles formatting, linting, and pre-commit workflow
  - Uses `gofmt` for simple, reliable formatting
  - Cross-platform compatibility

- **`pre-commit.bat`** - Manual pre-commit hook backup
  - Fallback for direct script execution
  - Integrates with Git hooks

### � **Setup & Configuration**

- **`setup-hooks.ps1`** - One-time Git hooks installation

  - Configures pre-commit automation
  - PowerShell script for Windows

- **`start-local-dev.bat`** - Complete Clerk development environment
  - Starts backend server
  - Launches ngrok tunnel with pre-configured URL
  - Provides testing URLs and next steps

### 🗄️ **Database & Verification**

- **`db-verify.bat`** - Database connectivity testing (Windows)
  - Database schema verification
  - Focused on Windows development environment

---

## 🗑️ **Removed Scripts**

The following scripts have been removed to eliminate redundancy:

❌ `test-api.bat` - Replaced by Swagger UI (`/swagger/index.html`)
❌ `test-ngrok.bat` - Functionality merged into `start-local-dev.bat`
❌ `diagnostic-format.bat` - No longer needed with simplified `gofmt` workflow
❌ `fix-formatting.bat` - Empty file, removed
❌ `fix-imports.bat` - Obsolete with `gofmt`-only approach
❌ `dev.sh` - Removed to focus on Windows development environment
❌ `db-verify.sh` - Removed to simplify maintenance

---

## 🎯 **Recommended Workflow**

### **Initial Setup** (once)

```bash
# Setup Git hooks
.\scripts\setup-hooks.ps1
```

### **Daily Development**

```bash
# Start complete development environment (with Clerk + ngrok)
.\scripts\start-local-dev.bat

# Quick formatting and linting
.\scripts\dev.bat format
.\scripts\dev.bat lint
.\scripts\dev.bat pre-commit
```

### **API Testing**

```bash
# Use Swagger UI for comprehensive API testing
http://localhost:30000/swagger/index.html

# Or quick health checks
curl http://localhost:30000/health
```

---

## ✅ **Benefits of Simplification**

✅ **Reduced complexity** - Fewer scripts to maintain
✅ **Clear purpose** - Each script has a specific, non-overlapping function
✅ **Better developer experience** - Less confusion about which script to use
✅ **Consistent workflow** - Standardized on `gofmt` and Swagger UI
✅ **Documentation accuracy** - README matches actual available scripts
✅ **Easier onboarding** - New developers have fewer tools to learn

**The scripts directory is now clean, focused, and maintainable! 🎉**
