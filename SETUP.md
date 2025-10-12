# mvnenv-win Setup Guide

## Quick Setup for Your Environment

Since you already have Maven versions installed (3.9.11 and 3.8.9), follow these steps to enable mvnenv command interception:

### 1. Add Shims to PATH

The shims directory **must be first** in your PATH to intercept `mvn` commands:

#### PowerShell (Recommended)

```powershell
# Add to current session
$env:Path = "$env:USERPROFILE\.mvnenv\shims;$env:USERPROFILE\.mvnenv\bin;" + $env:Path

# Make permanent (User PATH)
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
$newPath = "$env:USERPROFILE\.mvnenv\shims;$env:USERPROFILE\.mvnenv\bin;$currentPath"
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")

# Restart your terminal for changes to take effect
```

#### Command Prompt

```cmd
# Add to current session
set PATH=%USERPROFILE%\.mvnenv\shims;%USERPROFILE%\.mvnenv\bin;%PATH%

# Make permanent using System Properties
# 1. Open: sysdm.cpl
# 2. Advanced > Environment Variables
# 3. Edit User PATH variable
# 4. Add to the beginning:
#    C:\Users\YourName\.mvnenv\shims
#    C:\Users\YourName\.mvnenv\bin
```

### 2. Verify Setup

After adding to PATH and restarting your terminal:

```bash
# Check which mvn is being used (should show shims directory)
where mvn
# Expected: C:\Users\YourName\.mvnenv\shims\mvn.exe

# Check current Maven version
mvn -version
# Should show the version set by mvnenv

# Verify mvnenv can see the version
mvnenv version
```

### 3. Test Version Switching

```bash
# Switch to 3.8.9
mvnenv global 3.8.9
mvn -version
# Should show: Apache Maven 3.8.9

# Switch to 3.9.11
mvnenv global 3.9.11
mvn -version
# Should show: Apache Maven 3.9.11
```

### 4. Common Issues

#### Issue: `mvn -version` still shows old Maven

**Solution:** Shims directory is not first in PATH

```bash
# Check current PATH
echo %PATH%

# The output should start with:
# C:\Users\YourName\.mvnenv\shims;...

# If not, adjust your PATH so shims comes first
```

#### Issue: "No Maven version is set"

**Solution:** Set a global version

```bash
mvnenv global 3.9.11
```

#### Issue: Shim not found

**Solution:** Regenerate shims

```bash
mvnenv rehash
```

## How It Works

1. When you type `mvn`, Windows finds `mvn.exe` in the shims directory first
2. The shim resolves which Maven version to use (shell > local > global)
3. The shim executes the correct Maven installation with all your arguments
4. Output is forwarded directly to your terminal

This allows seamless version switching without manual PATH manipulation!

## Using Project-Specific Versions

```bash
# In project directory
cd /path/to/my-project
mvnenv local 3.8.9

# Now this project always uses 3.8.9, regardless of global setting
mvn -version
# Shows: Apache Maven 3.8.9
```

The `.maven-version` file in your project tells mvnenv which version to use.
