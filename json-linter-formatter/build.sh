#!/bin/bash
set -e

# Determine the host OS and Arch to name the binary correctly
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

PLUGIN_NAME="json-linter-formatter"
OUTPUT_NAME="${PLUGIN_NAME}"
RELEASE_DIR="release"

if [ "$GOOS" = "windows" ]; then
    OUTPUT_NAME="${PLUGIN_NAME}.exe"
fi

# 1. Build the frontend
echo "Building frontend..."
(cd frontend && npm install && npm run build)

# 2. Tidy the Go module to ensure dependencies are correct
echo "Tidying Go module..."
go mod tidy

# 3. Build the backend for the current host OS/Arch
echo "Building backend for ${GOOS}/${GOARCH}..."
go build -o "./${OUTPUT_NAME}" main.go

# 4. Add execute permission (for non-Windows builds)
if [ "$GOOS" != "windows" ]; then
    echo "Setting execute permission..."
    chmod +x "./${OUTPUT_NAME}"
fi

# 5. Assemble the package
echo "Assembling package..."
rm -rf "${RELEASE_DIR}"
mkdir -p "${RELEASE_DIR}/${PLUGIN_NAME}/frontend"

mv "./${OUTPUT_NAME}" "${RELEASE_DIR}/${PLUGIN_NAME}/"
cp "./frontend/dist/component.js" "${RELEASE_DIR}/${PLUGIN_NAME}/frontend/"

echo "Build complete! The packaged plugin is in the '${RELEASE_DIR}/${PLUGIN_NAME}' directory."
echo "Binary name: ${OUTPUT_NAME}"
