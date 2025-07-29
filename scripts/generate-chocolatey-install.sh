#!/bin/bash

# Chocolatey PowerShell Install Script Generator
# This script is the SINGLE SOURCE OF TRUTH for PowerShell script generation
# Used by both GitHub Actions workflow and local testing

set -e

# Input validation
if [ $# -ne 3 ]; then
    echo "Usage: $0 <version> <sha256_file_path> <output_ps1_path>"
    echo "Example: $0 0.2.6 ../release/agenthub-0.2.6-windows-amd64.zip.sha256 tools/chocolateyInstall.ps1"
    exit 1
fi

VERSION="$1"
SHA256_FILE="$2"
OUTPUT_FILE="$3"
ARCH="windows-amd64"

echo "🔧 Generating Chocolatey PowerShell install script..."
echo "   Version: $VERSION"
echo "   SHA256 file: $SHA256_FILE"
echo "   Output: $OUTPUT_FILE"

# Validate inputs
if [ ! -f "$SHA256_FILE" ]; then
    echo "❌ Error: SHA256 file not found: $SHA256_FILE"
    exit 1
fi

# Extract SHA256 hash (exactly as in GitHub Actions)
sha=$(cut -d' ' -f1 "$SHA256_FILE")

if [ -z "$sha" ]; then
    echo "❌ Error: Could not extract SHA256 from file: $SHA256_FILE"
    exit 1
fi

echo "   SHA256: $sha"

# Create output directory if it doesn't exist
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Generate PowerShell script (SINGLE SOURCE OF TRUTH)
# This is the EXACT same logic used in GitHub Actions
cat <<EOF > "$OUTPUT_FILE"
\$toolsDir = "\$(Split-Path -parent \$MyInvocation.MyCommand.Definition)"
\$url = "https://github.com/agenthubcli/agenthub/releases/download/v${VERSION}/agenthub-${VERSION}-${ARCH}.zip"
\$packageName = "agenthub"
\$checksum = "$sha"
Install-ChocolateyZipPackage -PackageName \$packageName -Url \$url -UnzipLocation \$toolsDir -Checksum \$checksum -ChecksumType sha256
EOF

echo "✅ PowerShell script generated successfully: $OUTPUT_FILE"

# Validate the generated script
echo "🔍 Validating generated script..."

if grep -q '^\$toolsDir = ' "$OUTPUT_FILE"; then
    echo "   ✅ toolsDir variable declaration is correct"
else
    echo "   ❌ toolsDir variable declaration is malformed"
    exit 1
fi

if grep -q "^\$checksum = \"$sha\"" "$OUTPUT_FILE"; then
    echo "   ✅ checksum variable is correctly set"
else
    echo "   ❌ checksum variable is malformed"
    exit 1
fi

if grep -q "Install-ChocolateyZipPackage.*-Checksum.*-ChecksumType sha256" "$OUTPUT_FILE"; then
    echo "   ✅ Chocolatey function call includes checksum validation"
else
    echo "   ❌ Chocolatey function call missing checksum validation"
    exit 1
fi

# Check for variable concatenation issues
if grep -q '^[a-f0-9]\+\$' "$OUTPUT_FILE"; then
    echo "   ❌ Found hash concatenated with variable name!"
    grep '^[a-f0-9]\+\$' "$OUTPUT_FILE"
    exit 1
else
    echo "   ✅ No variable concatenation issues found"
fi

echo "🎉 PowerShell script validation passed!"
echo ""
echo "📄 Generated script content:"
echo "----------------------------"
cat "$OUTPUT_FILE" 