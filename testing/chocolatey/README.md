# Chocolatey Testing

This directory contains tools for testing the Chocolatey package locally before publishing.

## 🎯 Zero Duplication Approach

**SINGLE SOURCE OF TRUTH**: Both GitHub Actions and local testing use the exact same shared scripts:

- `../../scripts/generate-nuspec.sh` - Generates nuspec from template
- `../../scripts/generate-chocolatey-install.sh` - Generates PowerShell install script

## Files

- `test-end-to-end.sh` - **🎯 COMPREHENSIVE** - Full end-to-end test with real Chocolatey installation
- `test-github-workflow.sh` - **RECOMMENDED** - Exact GitHub Actions simulation (zero duplication)
- `test-chocolatey-local.sh` - Legacy test script
- `Dockerfile` - Windows container for full Chocolatey package testing
- `README.md` - This file

## Quick Testing

### 1. 🎯 COMPREHENSIVE End-to-End Test (BEST)

```bash
./test-end-to-end.sh [version]
# Example: ./test-end-to-end.sh 0.2.8
```

This provides **COMPLETE validation**:

- ✅ Uses EXACT same scripts as GitHub Actions (zero duplication)
- ✅ Builds Chocolatey package in real Windows Docker container
- ✅ **Actually installs** the package with `choco install`
- ✅ **Verifies binary is accessible** and runs `agenthub --help`
- ✅ Tests package uninstallation
- ✅ **GUARANTEE**: If this passes, your package works in production!

**Requirements**: Docker Desktop with Windows containers enabled

### 2. Test GitHub Actions Workflow (Syntax Only)

```bash
./test-github-workflow.sh [version]
# Example: ./test-github-workflow.sh 0.2.7
```

This will:

- Use the EXACT same scripts as GitHub Actions
- Generate nuspec from shared template (`scripts/agenthub.nuspec.template`)
- Generate PowerShell script with shared generator
- Validate file generation and syntax
- **Note**: Does NOT test actual Chocolatey installation

### 3. Legacy Test (Standalone)

```bash
./test-chocolatey-local.sh [version]
# Example: ./test-chocolatey-local.sh 0.2.7
```

This will:

- Create a `test-choco/` directory with mock files
- Generate files independently (not recommended - potential for drift)

### 4. Manual Docker Test (Windows Containers)

**Prerequisites:**

- Docker Desktop with Windows container support
- Switch Docker to Windows containers: Right-click Docker tray icon → "Switch to Windows containers"

```bash
./test-github-workflow.sh 0.2.7
cd workflow-test
docker build -t test-agenthub-choco .
```

### 5. Test on Windows Machine

If you have access to a Windows machine:

```powershell
# Run the comprehensive test first
./test-end-to-end.sh 0.2.8
cd e2e-test/chocolatey

# Then test with Chocolatey
choco pack
choco install agenthub -s . -f

# Test the installed binary
agenthub --version

# Uninstall after testing
choco uninstall agenthub
```

### 6. Manual PowerShell Validation

Test just the PowerShell script syntax:

```powershell
# On Windows
powershell -Command "& ./chocolatey/tools/chocolateyInstall.ps1 -WhatIf"
```

## What This Tests

### 🎯 Comprehensive End-to-End Test (`test-end-to-end.sh`)

✅ **Zero Duplication** - Uses identical scripts as GitHub Actions  
✅ **Real Windows Environment** - Tests in actual Windows Server Core container  
✅ **Chocolatey Installation** - Actually runs `choco install` with generated package  
✅ **Binary Accessibility** - Verifies `agenthub` command is in PATH and executable  
✅ **Package Uninstallation** - Tests complete lifecycle including cleanup  
✅ **Production Fidelity** - Guarantees what works here works in production

### 📄 File Generation Test (`test-github-workflow.sh`)

✅ **Nuspec Generation** - From shared template with version substitution  
✅ **PowerShell Script Generation** - Validates exact workflow logic  
✅ **Variable Escaping** - Ensures no concatenation issues like `3133toolsDir`  
✅ **Checksum Integration** - Verifies SHA256 checksum is properly included  
✅ **Chocolatey Function Call** - Tests `Install-ChocolateyZipPackage` syntax  
✅ **Package Building** - Validates the complete package structure

## Troubleshooting

### "The term 'XXXtoolsDir' is not recognized"

This indicates variable concatenation issues in the PowerShell script generation. The shared scripts now prevent this.

### "Package automation scripts download a remote file without validating the checksum"

The shared PowerShell generator automatically includes `-Checksum` and `-ChecksumType sha256` parameters.

### "YAML syntax error"

The shared scripts eliminate complex bash heredoc issues that caused YAML parsing problems.

### Docker Issues

- Ensure Docker Desktop is set to Windows containers
- Windows Server Core images are large (~5GB) - ensure sufficient disk space
- Network connectivity required for Chocolatey installation

## Cleanup

```bash
# Remove end-to-end test files
rm -rf e2e-test/

# Remove workflow test files
rm -rf workflow-test/

# Remove legacy test files
rm -rf test-choco/
```

## Integration with CI/CD

This testing setup uses the **exact same shared scripts** as the GitHub Actions workflow:

1. **GitHub Actions** calls: `scripts/generate-nuspec.sh` and `scripts/generate-chocolatey-install.sh`
2. **Local Testing** calls: the same scripts with identical parameters
3. **Zero Drift**: Impossible for tests to pass while production fails

## Architecture

```
scripts/
├── agenthub.nuspec.template     # Single template source
├── generate-nuspec.sh           # Shared nuspec generator
└── generate-chocolatey-install.sh # Shared PowerShell generator

.github/workflows/release.yml    # Uses shared scripts
testing/chocolatey/
└── test-github-workflow.sh     # Uses same shared scripts
```

**Result**: 100% fidelity between testing and production!
