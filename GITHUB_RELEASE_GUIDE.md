# GitHub Release Creation Guide

## Manual Release Creation Steps

Since GitHub CLI requires interactive authentication, here are the steps to create a release manually:

### 1. Go to GitHub Repository

Visit: https://github.com/ankogit/ssh_keeper

### 2. Create New Release

1. Click on "Releases" tab
2. Click "Create a new release"
3. Click "Choose a tag" and select `v0.1.0`

### 3. Release Information

- **Tag version**: `v0.1.0`
- **Release title**: `SSH Keeper v0.1.0 - Initial Release`
- **Description**: Copy content from `RELEASE_NOTES.md`

### 4. Upload Release Assets

Upload the following files from `build/release/` directory:

- `ssh-keeper-0.1.0-darwin-amd64.tar.gz` (macOS Intel)
- `ssh-keeper-0.1.0-darwin-arm64.tar.gz` (macOS Apple Silicon)
- `ssh-keeper-0.1.0-linux-amd64.tar.gz` (Linux)
- `ssh-keeper-0.1.0-windows-amd64.zip` (Windows)

### 5. Release Settings

- ✅ Set as latest release
- ✅ Create a discussion for this release (optional)

### 6. Publish Release

Click "Publish release" to make it public.

## Alternative: Using GitHub CLI (if authenticated)

```bash
# Create release with assets
gh release create v0.1.0 \
  --title "SSH Keeper v0.1.0 - Initial Release" \
  --notes-file RELEASE_NOTES.md \
  build/release/ssh-keeper-0.1.0-darwin-amd64.tar.gz \
  build/release/ssh-keeper-0.1.0-darwin-arm64.tar.gz \
  build/release/ssh-keeper-0.1.0-linux-amd64.tar.gz \
  build/release/ssh-keeper-0.1.0-windows-amd64.zip
```

## Post-Release Tasks

1. **Update README**: Add actual download links once release is published
2. **Test Downloads**: Verify all platform downloads work correctly
3. **Community**: Share the release on social media/forums
4. **Documentation**: Add screenshots to complete the README

## Release Checklist

- [x] Code pushed to GitHub
- [x] Tag created and pushed
- [x] Release artifacts built
- [x] README updated with placeholders
- [x] LICENSE file added
- [x] Release notes prepared
- [ ] GitHub release created manually
- [ ] Download links tested
- [ ] Screenshots added to README
