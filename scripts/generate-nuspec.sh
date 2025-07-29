#!/bin/bash

# Nuspec File Generator from Template
# SINGLE SOURCE OF TRUTH for nuspec generation
# Used by both GitHub Actions workflow and local testing

set -e

if [ $# -ne 2 ]; then
    echo "Usage: $0 <version> <output_nuspec_path>"
    echo "Example: $0 0.2.6 chocolatey/agenthub.nuspec"
    exit 1
fi

VERSION="$1"
OUTPUT_FILE="$2"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_FILE="$SCRIPT_DIR/agenthub.nuspec.template"

echo "🔧 Generating nuspec file from template..."
echo "   Version: $VERSION"
echo "   Template: $TEMPLATE_FILE"
echo "   Output: $OUTPUT_FILE"

# Validate template exists
if [ ! -f "$TEMPLATE_FILE" ]; then
    echo "❌ Error: Template file not found: $TEMPLATE_FILE"
    exit 1
fi

# Create output directory if it doesn't exist
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Generate nuspec from template
sed "s/{{VERSION}}/$VERSION/g" "$TEMPLATE_FILE" > "$OUTPUT_FILE"

echo "✅ Nuspec file generated successfully: $OUTPUT_FILE"

# Validate the generated nuspec
echo "🔍 Validating generated nuspec..."

if grep -q "<version>$VERSION</version>" "$OUTPUT_FILE"; then
    echo "   ✅ Version is correctly set to: $VERSION"
else
    echo "   ❌ Version is incorrect in generated nuspec"
    exit 1
fi

if grep -q "<id>agenthub</id>" "$OUTPUT_FILE"; then
    echo "   ✅ Package ID is correct"
else
    echo "   ❌ Package ID is incorrect"
    exit 1
fi

echo "🎉 Nuspec validation passed!" 