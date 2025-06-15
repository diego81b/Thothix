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

### 🏷️ **Version Management**

- **`version-bump.bat`** - Automatic version bumping (Windows)
  - Semantic versioning automation (major/minor/patch)
  - Automatic CHANGELOG.md updates
  - Git tag creation and release management

- **`version-bump.ps1`** - PowerShell version bumping (Cross-platform)
  - Enhanced error handling and validation
  - Interactive prompts for release descriptions
  - Git status verification before bumping

- **`version-bump.sh`** - Unix/Linux version bumping
  - Cross-platform compatibility for Unix environments
  - Consistent functionality across all platforms

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

### **Version Management & Releases**

```bash
# Patch release (bug fixes)
.\scripts\version-bump.bat patch "Fix authentication timeout bug"

# Minor release (new features)
.\scripts\version-bump.ps1 minor "Add user role management"

# Major release (breaking changes)
.\scripts\version-bump.sh major "Restructure API endpoints"

# Publish the release
git push origin main
git push origin v1.3.0
```

### **Release Process**

1. **Update CHANGELOG.md** - Scripts automatically move [Unreleased] content to new version
2. **Create Git commit** - Automatic commit with conventional message
3. **Create Git tag** - Semantic versioning tag (v1.2.3)
4. **Ready to publish** - Push to remote repository

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
