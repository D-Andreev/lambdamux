#!/bin/bash

# Ensure we're on the main branch
git checkout main
git pull origin main

# Get the latest tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo "Latest tag: $LATEST_TAG"

# Prompt for new version
read -p "Enter new version number (e.g., 1.0.1): " VERSION

# Create and push new tag
git tag -a v$VERSION -m "Release version $VERSION"
git push origin v$VERSION

# Build the project
go build -v -o lambdamux .

# Create GitHub release
gh release create v$VERSION ./lambdamux -t "Release v$VERSION" -n "Release notes for version $VERSION"

echo "Release v$VERSION created and published!"