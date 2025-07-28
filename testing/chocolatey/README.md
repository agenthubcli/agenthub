# Chocolatey Testing

This directory contains tools for testing the Chocolatey package locally before publishing.

## Files

- `test-chocolatey-local.sh` - Script to generate test package and validate PowerShell syntax
- `Dockerfile` - Windows container for full Chocolatey package testing
- `README.md` - This file

## Quick Testing

### 1. Generate Test Package

```bash
./test-chocolatey-local.sh [version]
# Example: ./test-chocolatey-local.sh 0.2.7
```

This will:

- Create a `test-choco/` directory with mock files
- Generate the PowerShell install script using the same logic as GitHub Actions
- Validate PowerShell syntax
- Create a Dockerfile for Windows testing

### 2. Test with Docker (Windows Containers)

**Prerequisites:**

- Docker Desktop with Windows container support
- Switch Docker to Windows containers: Right-click Docker tray icon → "Switch to Windows containers"

```bash
cd test-choco
docker build -t test-agenthub-choco .
```

### 3. Test on Windows Machine

If you have access to a Windows machine:

```powershell
# In the test-choco directory
choco pack
choco install agenthub -s . -f

# Test the installed binary
agenthub --version

# Uninstall after testing
choco uninstall agenthub
```

### 4. Manual PowerShell Validation

Test just the PowerShell script syntax:

```powershell
# On Windows
powershell -Command "& ./tools/chocolateyInstall.ps1 -WhatIf"
```

## What This Tests

✅ **PowerShell Script Generation** - Validates the GitHub Actions workflow produces valid PowerShell  
✅ **Variable Escaping** - Ensures no concatenation issues like `3133toolsDir`  
✅ **Checksum Integration** - Verifies SHA256 checksum is properly included  
✅ **Chocolatey Function Call** - Tests `Install-ChocolateyZipPackage` syntax  
✅ **Package Building** - Validates the `.nuspec` and package structure

## Troubleshooting

### "The term 'XXXtoolsDir' is not recognized"

This indicates variable concatenation issues in the PowerShell script generation. Check the bash heredoc escaping in `.github/workflows/release.yml`.

### "Package automation scripts download a remote file without validating the checksum"

Ensure the PowerShell script includes `-Checksum` and `-ChecksumType sha256` parameters.

### Docker Issues

- Ensure Docker Desktop is set to Windows containers
- Windows Server Core images are large (~5GB) - ensure sufficient disk space
- Network connectivity required for Chocolatey installation

## Cleanup

```bash
# Remove test files after testing
rm -rf test-choco/
```

## Integration with CI/CD

This testing setup simulates the exact same PowerShell script generation process used in the GitHub Actions workflow, ensuring consistency between local testing and production releases.
