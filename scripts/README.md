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

- **`dev.bat`** - Main development script

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

### 🗄️ **Database & Verification**

- **`db-verify.bat`** - Database connectivity testing (Windows)
  - Database schema verification
  - Focused on Windows development environment

---

## 🎯 **Recommended Workflow**

### **Initial Setup** (once)

```bash
# Setup Git hooks
.\scripts\setup-hooks.ps1
```

### **Daily Development**

```bash
# Start backend for development
.\scripts\dev.bat

# For Clerk + webhook testing, also run ngrok in a separate terminal:
ngrok http --url=flying-mullet-socially.ngrok-free.app 30000

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

## Conclusion

The scripts directory is now clean, focused, and maintainable! 🎉
