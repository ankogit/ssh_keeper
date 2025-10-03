#!/bin/bash

# Build Debian package script
set -e

echo "Building Debian package for ssh-keeper..."

# Clean previous builds
rm -rf debian/ssh-keeper
rm -f *.deb

# Build the package
dpkg-buildpackage -us -uc

echo "Debian package built successfully!"
echo "Package files:"
ls -la ../*.deb ../*.changes ../*.dsc 2>/dev/null || true

