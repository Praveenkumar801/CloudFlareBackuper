# Release Guide

This document describes how to create releases for CloudFlare Backuper with automated binary distribution.

## Quick Start - Automated Releases

The easiest way to create a release is to update the VERSION file:

1. **Update the VERSION file**:
   ```bash
   echo "1.0.1" > VERSION
   git add VERSION
   git commit -m "Release version 1.0.1"
   git push origin main
   ```

2. **Automated Process**:
   - The `auto-release` workflow detects the VERSION file change
   - Automatically creates a git tag (e.g., `v1.0.1`)
   - Triggers the `release` workflow
   - Builds binaries for all platforms
   - Creates a GitHub release with all binaries attached

## Version File

The `VERSION` file in the repository root contains the current version number.

### Format

The version must follow [Semantic Versioning](https://semver.org/):

- `MAJOR.MINOR.PATCH` (e.g., `1.0.0`)
- Optional pre-release: `MAJOR.MINOR.PATCH-prerelease` (e.g., `1.0.0-beta.1`, `2.0.0-rc.2`)

### Examples

```
1.0.0           # Initial release
1.0.1           # Patch release (bug fixes)
1.1.0           # Minor release (new features, backward compatible)
2.0.0           # Major release (breaking changes)
1.0.0-beta.1    # Pre-release
```

## Semantic Versioning Guidelines

### When to Increment

- **MAJOR** version: Breaking changes (incompatible API changes, configuration format changes)
- **MINOR** version: New features (backward compatible)
- **PATCH** version: Bug fixes (backward compatible)

### Examples

- Configuration file format changes → Major version bump
- New notification provider support → Minor version bump
- Bug fix in backup logic → Patch version bump

## Release Process

### Automated Release (Recommended)

1. **Decide on the new version number** following semantic versioning

2. **Update VERSION file**:
   ```bash
   # For a patch release (bug fixes)
   echo "1.0.1" > VERSION
   
   # For a minor release (new features)
   echo "1.1.0" > VERSION
   
   # For a major release (breaking changes)
   echo "2.0.0" > VERSION
   ```

3. **Commit and push**:
   ```bash
   git add VERSION
   git commit -m "Release version X.Y.Z"
   git push origin main
   ```

4. **Monitor the workflows**:
   - Go to Actions tab in GitHub
   - Watch "Auto Release on Version Change" workflow
   - Once tag is created, "Release Build" workflow starts automatically
   - Release will appear in the Releases page with all binaries

### Manual Release (Advanced)

If you prefer manual control:

1. **Create and push a tag**:
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

2. **The release workflow triggers automatically** and builds binaries

## What Gets Released

Each release includes binaries for multiple platforms:

### Linux
- `cloudflare-backuper-linux-amd64-vX.Y.Z` - x86_64
- `cloudflare-backuper-linux-arm64-vX.Y.Z` - ARM64

### macOS
- `cloudflare-backuper-darwin-amd64-vX.Y.Z` - Intel Macs
- `cloudflare-backuper-darwin-arm64-vX.Y.Z` - Apple Silicon (M1/M2/M3)

### Windows
- `cloudflare-backuper-windows-amd64-vX.Y.Z.exe` - x86_64
- `cloudflare-backuper-windows-arm64-vX.Y.Z.exe` - ARM64

### Checksums
- Each binary has a corresponding `.sha256` file for integrity verification

## Binary Features

All release binaries include:

- **Version information**: Embedded in the binary
- **Build metadata**: Commit SHA and build date
- **Optimized**: Stripped symbols and compressed for smaller size
- **Static**: No external dependencies required

Check version:
```bash
./cloudflare-backuper-linux-amd64-v1.0.0 -version
```

Output:
```
CloudFlare Backuper
Version:    1.0.0
Commit:     abc1234
Build Date: 2024-01-15T10:30:00Z
```

## Pre-releases

For beta or release candidate versions:

1. Use a pre-release version in VERSION file:
   ```bash
   echo "2.0.0-beta.1" > VERSION
   ```

2. Commit and push:
   ```bash
   git add VERSION
   git commit -m "Release version 2.0.0-beta.1"
   git push origin main
   ```

The release will be marked as a pre-release automatically if the version contains a hyphen.

## Troubleshooting

### Tag Already Exists

If you try to release a version that already exists:
- The auto-release workflow will skip tag creation
- No new release will be created
- Increment the version number and try again

### Workflow Fails

Check the Actions tab for error details:
- Build failures: Usually dependency or compilation issues
- Tag creation failures: Permissions or duplicate tag issues

### Invalid Version Format

The VERSION file must contain only the version number:
- ✅ Good: `1.0.0`
- ✅ Good: `1.0.0-beta.1`
- ❌ Bad: `v1.0.0` (don't include 'v' prefix)
- ❌ Bad: `version 1.0.0` (no extra text)

## Best Practices

1. **Update VERSION file only when ready to release**
2. **Include meaningful commit message**: "Release version X.Y.Z" or "Bump version to X.Y.Z"
3. **Test before releasing**: Run tests and verify functionality
4. **Document changes**: Update README.md or CHANGELOG.md if needed
5. **Semantic versioning**: Follow semver guidelines consistently
6. **One version per commit**: Don't update VERSION file multiple times in one commit

## Manual Workflow Trigger

If needed, you can manually trigger the release workflow:

1. Go to Actions tab
2. Select "Release Build" workflow
3. Click "Run workflow"
4. Select the tag you want to build

This is useful for:
- Rebuilding a release
- Testing the workflow
- Creating a release without updating VERSION file

## CI/CD Pipeline Overview

```
VERSION file updated
       ↓
Auto-release workflow
       ↓
Tag created (vX.Y.Z)
       ↓
Release workflow triggered
       ↓
Build on 3 OS × 2 architectures
       ↓
Generate checksums
       ↓
Create GitHub Release
       ↓
Upload all binaries & checksums
```

## Questions?

- Check [README.md](README.md) for general information
- See [GitHub Actions](https://github.com/Praveenkumar801/CloudFlareBackuper/actions) for workflow status
- Open an issue for problems or questions
