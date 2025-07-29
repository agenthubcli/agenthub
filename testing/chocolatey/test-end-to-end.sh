#!/bin/bash

# End-to-End Chocolatey Testing
# This script provides COMPLETE validation:
# 1. Tests GitHub Actions workflow file generation 
# 2. Tests REAL Chocolatey installation in Windows Docker container
# 3. Verifies the installed binary actually works

set -e

VERSION="${1:-0.2.7}"
BINARY_NAME="agenthub"
ARCH="windows-amd64"

echo "🚀 COMPREHENSIVE END-TO-END CHOCOLATEY TESTING"
echo "=============================================="
echo "Version: $VERSION"
echo "This test will:"
echo "  ✅ Generate files using GitHub Actions scripts (zero duplication)"
echo "  ✅ Build Chocolatey package in Windows Docker container" 
echo "  ✅ Install package with Chocolatey in real Windows environment"
echo "  ✅ Verify installed binary is accessible and functional"
echo ""

# Clean up any previous test
rm -rf e2e-test
mkdir -p e2e-test/{release,chocolatey/tools}
cd e2e-test

echo "📦 Step 1: Creating mock release artifacts..."
# Create mock release files (simulating what GitHub Actions creates)
echo "Mock agenthub binary for testing" > "${BINARY_NAME}.exe"
zip "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip" "${BINARY_NAME}.exe"
shasum -a 256 "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip" > "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip.sha256"
rm "${BINARY_NAME}.exe"

echo "Generated release artifacts:"
ls -la release/
echo ""

echo "🔧 Step 2: Generating Chocolatey files using GitHub Actions scripts..."
echo "This ensures ZERO DUPLICATION with production workflow!"

# Use EXACT same scripts as GitHub Actions
../../../scripts/generate-nuspec.sh "$VERSION" "chocolatey/agenthub.nuspec"
../../../scripts/generate-chocolatey-install.sh "$VERSION" "release/agenthub-${VERSION}-${ARCH}.zip.sha256" "chocolatey/tools/chocolateyInstall.ps1"

echo ""
echo "📄 Generated files:"
echo "-------------------"
echo "Nuspec version:"
grep -A1 "<version>" chocolatey/agenthub.nuspec

echo ""
echo "PowerShell script:"
cat chocolatey/tools/chocolateyInstall.ps1

echo ""
echo "🐳 Step 3: Creating comprehensive Docker test environment..."

# Create a Dockerfile that does FULL end-to-end testing
cat > Dockerfile << 'EOF'
# escape=`
FROM mcr.microsoft.com/windows/servercore:ltsc2022

# Install Chocolatey
RUN powershell -Command `
    Set-ExecutionPolicy Bypass -Scope Process -Force; `
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; `
    iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Copy package files
COPY . /package
WORKDIR /package

# Step 1: Validate PowerShell script syntax
RUN powershell -Command `
    "Write-Host '=== STEP 1: Testing PowerShell Script Syntax ===' -ForegroundColor Cyan; `
    try { `
        . ./chocolatey/tools/chocolateyInstall.ps1 -WhatIf -ErrorAction Stop; `
        Write-Host 'PowerShell script syntax is VALID!' -ForegroundColor Green `
    } catch { `
        Write-Host \"PowerShell script syntax error: \$_\" -ForegroundColor Red; `
        exit 1 `
    }"

# Step 2: Build Chocolatey package
RUN powershell -Command `
    "Write-Host '=== STEP 2: Building Chocolatey Package ===' -ForegroundColor Cyan; `
    cd chocolatey; `
    choco pack --trace; `
    if (\$LASTEXITCODE -ne 0) { `
        Write-Host 'Package build FAILED!' -ForegroundColor Red; `
        exit \$LASTEXITCODE `
    }; `
    Write-Host 'Package built successfully!' -ForegroundColor Green; `
    Get-ChildItem *.nupkg"

# Step 3: Install package locally (this is the REAL test!)
RUN powershell -Command `
    "Write-Host '=== STEP 3: Installing Package with Chocolatey ===' -ForegroundColor Cyan; `
    cd chocolatey; `
    choco install agenthub -s . -f -y --debug; `
    if (\$LASTEXITCODE -ne 0) { `
        Write-Host 'Package installation FAILED!' -ForegroundColor Red; `
        exit \$LASTEXITCODE `
    }; `
    Write-Host 'Package installed successfully!' -ForegroundColor Green"

# Step 4: Verify binary is accessible (CRITICAL TEST)
RUN powershell -Command `
    "Write-Host '=== STEP 4: Verifying Installed Binary ===' -ForegroundColor Cyan; `
    Write-Host 'Checking if agenthub binary is in PATH...'; `
    \$found = Get-Command agenthub -ErrorAction SilentlyContinue; `
    if (\$found) { `
        Write-Host 'agenthub found in PATH!' -ForegroundColor Green; `
        Write-Host \"Location: \$(\$found.Source)\"; `
        Write-Host 'Testing binary execution...'; `
        try { `
            agenthub --help; `
            Write-Host 'Binary executed successfully!' -ForegroundColor Green `
        } catch { `
            Write-Host \"Binary execution failed: \$_\" -ForegroundColor Red; `
            exit 1 `
        } `
    } else { `
        Write-Host 'agenthub NOT found in PATH!' -ForegroundColor Red; `
        Write-Host 'Searching for installation location...'; `
        Get-ChildItem -Path C:\ -Name '*agenthub*' -Recurse -ErrorAction SilentlyContinue | Select-Object -First 10; `
        exit 1 `
    }"

# Step 5: Test uninstallation (cleanup test)
RUN powershell -Command `
    "Write-Host '=== STEP 5: Testing Package Uninstallation ===' -ForegroundColor Cyan; `
    choco uninstall agenthub -y; `
    if (\$LASTEXITCODE -ne 0) { `
        Write-Host 'Package uninstallation FAILED!' -ForegroundColor Red; `
        exit \$LASTEXITCODE `
    }; `
    Write-Host 'Package uninstalled successfully!' -ForegroundColor Green"

# Final success message
CMD ["powershell", "-Command", "Write-Host '🎉 ALL TESTS PASSED! Chocolatey package works perfectly!' -ForegroundColor Green"]
EOF

echo "✅ Comprehensive Docker test environment created!"
echo ""
echo "🚀 Step 4: Running FULL end-to-end test in Windows container..."
echo "This will take a few minutes as it downloads Windows Server Core (~5GB)"
echo ""

# Check if Docker is available and configured for Windows containers
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not available. Please install Docker Desktop."
    echo ""
    echo "Manual testing instructions:"
    echo "=============================="
    echo "1. Copy this directory to a Windows machine"
    echo "2. Run in PowerShell:"
    echo "   cd chocolatey"
    echo "   choco pack"
    echo "   choco install agenthub -s . -f"
    echo "   agenthub --help"
    echo "   choco uninstall agenthub"
    exit 1
fi

# Check Docker platform
docker_platform=$(docker version --format '{{.Server.Platform.Name}}' 2>/dev/null || echo "unknown")

if [[ "$docker_platform" != *"windows"* ]]; then
    echo "⚠️  Docker is not configured for Windows containers."
    echo ""
    echo "To switch to Windows containers:"
    echo "1. Right-click Docker Desktop system tray icon"
    echo "2. Select 'Switch to Windows containers'"
    echo "3. Wait for Docker to restart"
    echo "4. Re-run this test"
    echo ""
    echo "Alternative: Manual testing instructions above"
    echo ""
    echo "Building anyway to test Docker setup..."
fi

echo "Building Docker image (this will take several minutes)..."
docker build -t agenthub-e2e-test . 2>&1 | tee docker-build.log

if [ ${PIPESTATUS[0]} -eq 0 ]; then
    echo ""
    echo "🎉 DOCKER BUILD SUCCESSFUL!"
    echo ""
    echo "Running end-to-end test container..."
    docker run --rm agenthub-e2e-test
    
    if [ $? -eq 0 ]; then
        echo ""
        echo "🎊 SUCCESS! ALL TESTS PASSED!"
        echo "=============================================="
        echo "✅ GitHub Actions file generation: WORKS"
        echo "✅ Chocolatey package building: WORKS" 
        echo "✅ Chocolatey package installation: WORKS"
        echo "✅ Binary accessibility: WORKS"
        echo "✅ Package uninstallation: WORKS"
        echo ""
        echo "🚀 Your Chocolatey package is ready for production!"
    else
        echo ""
        echo "❌ Container execution failed. Check the logs above."
        exit 1
    fi
else
    echo ""
    echo "❌ Docker build failed. Possible issues:"
    echo "1. Docker not configured for Windows containers"
    echo "2. Network connectivity issues"
    echo "3. Insufficient disk space (~5GB needed)"
    echo ""
    echo "Check docker-build.log for details."
    exit 1
fi

cd ..
echo ""
echo "💡 To clean up: rm -rf e2e-test/" 