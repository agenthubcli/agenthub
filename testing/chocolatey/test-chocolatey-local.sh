#!/bin/bash

# Local Chocolatey Testing Script
# This script simulates the GitHub Actions workflow and tests the Chocolatey package

set -e

VERSION="${1:-0.2.6}"
ARCH="windows-amd64"
BINARY_NAME="agenthub"

echo "🧪 Setting up local Chocolatey test environment..."

# Create test directory structure
mkdir -p test-choco/{tools,release}
cd test-choco

echo "📦 Creating mock release files..."
# Create a dummy ZIP file and SHA256 (simulating the actual release)
echo "Mock binary content for testing" > "${BINARY_NAME}.exe"
zip "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip" "${BINARY_NAME}.exe"
shasum -a 256 "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip" > "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip.sha256"
rm "${BINARY_NAME}.exe"

echo "📋 Generated SHA256:"
cat "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip.sha256"

echo "🔧 Generating Chocolatey package files..."

# Create nuspec file
cat > agenthub.nuspec << EOF
<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd">
  <metadata>
    <id>agenthub</id>
    <version>${VERSION}</version>
    <packageSourceUrl>https://github.com/agenthubcli/agenthub</packageSourceUrl>
    <owners>agenthubcli</owners>
    <title>AgentHub CLI</title>
    <authors>AgentHub Team</authors>
    <projectUrl>https://github.com/agenthubcli/agenthub</projectUrl>
    <copyright>2025 AgentHub</copyright>
    <licenseUrl>https://github.com/agenthubcli/agenthub/blob/main/LICENSE</licenseUrl>
    <requireLicenseAcceptance>false</requireLicenseAcceptance>
    <tags>cli agent ai tools</tags>
    <summary>Universal package manager and runtime for AI-native agents, tools, chains, and prompts</summary>
    <description>AgentHub is a universal package manager and runtime for AI-native agents, tools, chains, and prompts. It provides a CLI interface for discovering, installing, running, and publishing AI components.</description>
  </metadata>
  <files>
    <file src="tools\**" target="tools" />
  </files>
</package>
EOF

# Simulate the GitHub Actions workflow script generation
echo "🎯 Simulating GitHub Actions workflow..."
sha=$(cut -d' ' -f1 "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip.sha256")

cat > tools/chocolateyInstall.ps1 << EOF
\$toolsDir = "\$(Split-Path -parent \$MyInvocation.MyCommand.Definition)"
\$url = "https://github.com/agenthubcli/agenthub/releases/download/v${VERSION}/${BINARY_NAME}-${VERSION}-${ARCH}.zip"
\$packageName = "agenthub"
\$checksum = "$sha"
Install-ChocolateyZipPackage -PackageName \$packageName -Url \$url -UnzipLocation \$toolsDir -Checksum \$checksum -ChecksumType sha256
EOF

echo "✅ Generated PowerShell install script:"
echo "============================================="
cat tools/chocolateyInstall.ps1
echo "============================================="

# Basic syntax validation
echo "🔍 Validating PowerShell script..."
if grep -q '^\$toolsDir = ' tools/chocolateyInstall.ps1; then
    echo "✅ Variable declarations look correct"
else
    echo "❌ Variable declarations are malformed"
    exit 1
fi

if grep -q "^\$checksum = \"$sha\"" tools/chocolateyInstall.ps1; then
    echo "✅ Checksum variable is correctly set"
else
    echo "❌ Checksum variable is malformed"
    exit 1
fi

echo "📋 Creating Dockerfile for Windows testing..."

# Create Dockerfile for testing
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

# Test the chocolatey script syntax
RUN powershell -Command `
    "Write-Host 'Testing PowerShell script syntax...'; `
    try { `
        . ./tools/chocolateyInstall.ps1 -WhatIf -ErrorAction Stop; `
        Write-Host 'PowerShell script syntax is valid!' -ForegroundColor Green `
    } catch { `
        Write-Host \"PowerShell script error: \$_\" -ForegroundColor Red; `
        exit 1 `
    }"

# Build and test the package
RUN powershell -Command `
    "choco pack --trace; `
    if (\$LASTEXITCODE -ne 0) { exit \$LASTEXITCODE }; `
    Write-Host 'Package built successfully!' -ForegroundColor Green"

CMD ["powershell", "-Command", "Write-Host 'Chocolatey package test completed successfully!' -ForegroundColor Green"]
EOF

echo "🚀 Test commands:"
echo ""
echo "1. Build and test with Docker (requires Docker Desktop with Windows containers):"
echo "   docker build -t test-agenthub-choco ."
echo ""
echo "2. Test on Windows machine directly:"
echo "   choco pack"
echo "   choco install agenthub -s . -f"
echo ""
echo "3. Test PowerShell script syntax (on Windows):"
echo "   powershell -Command \"& ./tools/chocolateyInstall.ps1 -WhatIf\""
echo ""
echo "4. Manual inspection:"
echo "   cat tools/chocolateyInstall.ps1"
echo ""

cd ..
echo "✅ Test environment created in ./test-choco/"
echo "📂 Files created:"
ls -la test-choco/
echo ""
echo "🎯 The PowerShell script should now have proper syntax without variable concatenation issues!"
echo ""
echo "💡 To clean up after testing: rm -rf test-choco/" 