#!/bin/bash

# GitHub Actions Workflow Simulation Script
# This script exactly simulates the GitHub Actions workflow for Chocolatey package generation

set -e

VERSION="${1:-0.2.6}"
BINARY_NAME="agenthub"
ARCH="windows-amd64"

echo "🎯 Testing GitHub Actions Workflow for Chocolatey Package Generation"
echo "====================================================================="
echo "Version: $VERSION"
echo ""

# Clean up any previous test
rm -rf workflow-test
mkdir -p workflow-test/{release,chocolatey/tools}
cd workflow-test

echo "📦 Step 1: Simulating release artifact generation..."
# Create mock release files (simulating what cross-compile and package-release would create)
echo "Mock agenthub binary for testing" > "${BINARY_NAME}.exe"
zip "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip" "${BINARY_NAME}.exe"
shasum -a 256 "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip" > "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip.sha256"
rm "${BINARY_NAME}.exe"

echo "Generated release files:"
ls -la release/
echo ""
echo "SHA256 file content:"
cat "release/${BINARY_NAME}-${VERSION}-${ARCH}.zip.sha256"
echo ""

echo "📋 Step 2: Executing GitHub Actions workflow logic (ZERO DUPLICATION)..."

echo "Generating nuspec using shared script..."
# Use the EXACT same script as GitHub Actions (SINGLE SOURCE OF TRUTH):
../../../scripts/generate-nuspec.sh "$VERSION" "chocolatey/agenthub.nuspec"

echo "Generating PowerShell install script using shared generator..."
# Use the EXACT same script as GitHub Actions (SINGLE SOURCE OF TRUTH):
../../../scripts/generate-chocolatey-install.sh "$VERSION" "release/agenthub-${VERSION}-${ARCH}.zip.sha256" "chocolatey/tools/chocolateyInstall.ps1"

echo ""
echo "✅ Step 3: Validating generated files..."
echo "========================================"

echo ""
echo "📄 Generated nuspec file:"
echo "-------------------------"
grep -A2 -B2 "<version>" chocolatey/agenthub.nuspec

echo ""
echo "📄 Generated PowerShell script:"
echo "--------------------------------"
cat chocolatey/tools/chocolateyInstall.ps1

echo ""
echo "🔍 Step 4: Comprehensive Validation..."
echo "======================================"

# Extract SHA for validation (but don't use it for generation)
sha=$(cut -d' ' -f1 release/agenthub-${VERSION}-${ARCH}.zip.sha256)

# Test 1: Basic syntax validation (shared scripts already validate, but double-check)
echo "1. Testing PowerShell variable syntax..."
if grep -q '^\$toolsDir = ' chocolatey/tools/chocolateyInstall.ps1; then
    echo "   ✅ toolsDir variable declaration is correct"
else
    echo "   ❌ toolsDir variable declaration is malformed"
    exit 1
fi

if grep -q '^\$url = ' chocolatey/tools/chocolateyInstall.ps1; then
    echo "   ✅ url variable declaration is correct"
else
    echo "   ❌ url variable declaration is malformed"
    exit 1
fi

if grep -q '^\$packageName = ' chocolatey/tools/chocolateyInstall.ps1; then
    echo "   ✅ packageName variable declaration is correct"
else
    echo "   ❌ packageName variable declaration is malformed"
    exit 1
fi

if grep -q "^\$checksum = \"$sha\"" chocolatey/tools/chocolateyInstall.ps1; then
    echo "   ✅ checksum variable is correctly set to: $sha"
else
    echo "   ❌ checksum variable is malformed"
    echo "   Expected: \$checksum = \"$sha\""
    echo "   Actual:"
    grep "checksum" chocolatey/tools/chocolateyInstall.ps1 || echo "   (no checksum found)"
    exit 1
fi

# Test 2: Check for variable concatenation issues
echo ""
echo "2. Testing for variable concatenation issues..."
if grep -q '^[a-f0-9]\+\$' chocolatey/tools/chocolateyInstall.ps1; then
    echo "   ❌ Found hash concatenated with variable name!"
    grep '^[a-f0-9]\+\$' chocolatey/tools/chocolateyInstall.ps1
    exit 1
else
    echo "   ✅ No variable concatenation issues found"
fi

# Test 3: Check Chocolatey function call
echo ""
echo "3. Testing Chocolatey function call..."
if grep -q "Install-ChocolateyZipPackage.*-Checksum.*-ChecksumType sha256" chocolatey/tools/chocolateyInstall.ps1; then
    echo "   ✅ Install-ChocolateyZipPackage call includes checksum validation"
else
    echo "   ❌ Install-ChocolateyZipPackage call missing checksum validation"
    exit 1
fi

# Test 4: URL validation
echo ""
echo "4. Testing download URL format..."
expected_url="https://github.com/agenthubcli/agenthub/releases/download/v${VERSION}/agenthub-${VERSION}-${ARCH}.zip"
if grep -q "$(printf '%s\n' "$expected_url" | sed 's/[[\.*^$()+?{|]/\\&/g')" chocolatey/tools/chocolateyInstall.ps1; then
    echo "   ✅ Download URL format is correct"
    echo "   URL: $expected_url"
else
    echo "   ❌ Download URL format is incorrect"
    echo "   Expected: $expected_url"
    echo "   Actual:"
    grep '\$url = ' chocolatey/tools/chocolateyInstall.ps1
    exit 1
fi

# Test 5: Version consistency
echo ""
echo "5. Testing version consistency..."
nuspec_version=$(grep -o '<version>[^<]*</version>' chocolatey/agenthub.nuspec | sed 's/<[^>]*>//g')
if [ "$nuspec_version" = "$VERSION" ]; then
    echo "   ✅ Nuspec version updated correctly: $nuspec_version"
else
    echo "   ❌ Nuspec version incorrect. Expected: $VERSION, Got: $nuspec_version"
    exit 1
fi

echo ""
echo "🎉 All tests passed!"
echo "==================="
echo ""
echo "📊 Summary:"
echo "• Nuspec generation: ✅ Uses shared template"
echo "• PowerShell script syntax: ✅ Valid"
echo "• Variable declarations: ✅ Proper format"
echo "• Checksum integration: ✅ Included ($sha)"
echo "• Chocolatey function: ✅ Proper parameters"
echo "• URL format: ✅ Correct"
echo "• Version consistency: ✅ Matches"
echo ""
echo "🚀 This matches EXACTLY what GitHub Actions will generate!"
echo "🎯 ZERO DUPLICATION - Both use the same shared scripts"
echo ""
echo "📂 Test files created in: $(pwd)/workflow-test"
echo "💡 To clean up: rm -rf workflow-test"

cd .. 