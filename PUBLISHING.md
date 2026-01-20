# Publishing Rasedi Go SDK

This guide explains how to release new versions of the Rasedi Go SDK and make them available on [pkg.go.dev](https://pkg.go.dev/).

## Prerequisites

1.  **Code Check**: Ensure your code is thoroughly tested and ready for release.
2.  **Clean Dependencies**: Run `go mod tidy` to ensure `go.mod` and `go.sum` are up to date.
3.  **Commit**: Ensure all changes are committed and pushed to the `main` branch.

## Step-by-Step Publishing Guide

Go uses **Semantic Versioning** (vMajor.Minor.Patch). "Publishing" simply means pushing a git tag.

### 1. Check Existing Versions
Check what versions have already been released:
```bash
git tag
```

### 2. Tag the New Version
Create a new tag for your release (e.g., `v1.0.1`, `v1.1.0`, or `v2.0.0`).
```bash
git tag v1.0.1
```
*Note: If releasing a v2+ major version, update your `go.mod` module path to verify (e.g., `.../rasedi-go-sdk/v2`).*

### 3. Push the Tag
Push the tag to GitHub. This effectively "publishes" the package.
```bash
git push origin v1.0.1
```

### 4. Trigger Indexing on pkg.go.dev
The Go proxy caches versions. To make your new version show up immediately:

1.  Visit this URL in your browser (replace `v1.0.1` with your version):
    ```
    https://pkg.go.dev/github.com/MirotechTeam/rasedi-go-sdk@v1.0.1
    ```
2.  Assuming the page hasn't been cached yet, click the **Request** button if you see a "not found" or "not yet indexed" message.

Alternatively, you can trigger it via the command line:
```bash
go list -m github.com/MirotechTeam/rasedi-go-sdk@v1.0.1
```

## Verification
To verify users can fetch the new version:
```bash
go get github.com/MirotechTeam/rasedi-go-sdk@v1.0.1
```
