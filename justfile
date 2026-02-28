# kwelea — development tasks
# https://just.systems

default:
    @just --list

binary_name := "kwelea"
build_dir := "bin"
version_file := "VERSION"
version := `cat {{version_file}} 2>/dev/null || git describe --tags --always 2>/dev/null || echo 'dev'`
commit := `git rev-parse --short HEAD 2>/dev/null || echo 'none'`
date := `date -u +%Y-%m-%dT%H:%M:%SZ`
ldflags := "-s -w -X main.version=" + version + " -X main.commit=" + commit + " -X main.date=" + date

# Build the kwelea binary
build:
    @mkdir -p {{ build_dir }}
    go build -ldflags="{{ ldflags }}" -o {{ build_dir }}/{{ binary_name }} .

# 🌎 Build for all supported platforms
build-all:
    @mkdir -p {{ build_dir }}
    @for os in linux darwin; do \
        for arch in amd64 arm64; do \
            echo "🔨 Building {{ binary_name }}-${os}-${arch}..."; \
            GOOS=$os GOARCH=$arch go build -ldflags="{{ ldflags }}" -o {{ build_dir }}/{{ binary_name }}-${os}-${arch} .; \
        done; \
    done
    @echo "🔨 Building {{ binary_name }}-windows-amd64.exe..."
    @GOOS=windows GOARCH=amd64 go build -ldflags="{{ ldflags }}" -o {{ build_dir }}/{{ binary_name }}-windows-amd64.exe .
    @echo "✅ All platforms built in {{ build_dir }}/"

# Install kwelea globally
install:
    go install .

# Run all tests
test:
    @echo "🧪 Running tests..."
    gotestsum --format=pkgname-and-test-fails ./...

# Run tests with verbose output
test-v:
    go test -v ./...

# 📊 Run tests with coverage report
test-coverage:
    @echo "📊 Generating coverage report..."
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "✅ Coverage report generated: coverage.html"

# 🧹 Format all Go source files
fmt:
    @echo "🧹 Formatting code..."
    go fmt ./...

# 🔎 Vet the code for potential bugs
vet:
    @echo "🔎 Vetting code..."
    go vet ./...

# ✅ Run formatter and vet sequentially
check: fmt vet

# Build the kwelea docs with kwelea itself (dogfood)
docs: build
    {{ build_dir }}/{{ binary_name }} build

# Serve the kwelea docs locally (dogfood)
serve: build
    {{ build_dir }}/{{ binary_name }} serve

# Remove build artifacts
clean:
    rm -rf {{ build_dir }}/
    rm -rf site/

# Download/refresh self-hosted font files into assets/fonts/
download-fonts:
    mkdir -p assets/fonts
    curl -sL "https://fonts.gstatic.com/s/ibmplexmono/v20/-F63fjptAgt5VM-kVkqdyU8n1i8q1w.woff2" -o assets/fonts/ibm-plex-mono-400.woff2
    curl -sL "https://fonts.gstatic.com/s/ibmplexmono/v20/-F6qfjptAgt5VM-kVkqdyU8n3twJwlBFgg.woff2" -o assets/fonts/ibm-plex-mono-500.woff2
    curl -sL "https://fonts.gstatic.com/s/ibmplexsans/v23/zYXzKVElMYYaJe8bpLHnCwDKr932-G7dytD-Dmu1syxeKYY.woff2" -o assets/fonts/ibm-plex-sans.woff2
    curl -sL "https://fonts.gstatic.com/s/lora/v37/0QIvMX1D_JOuMwr7Iw.woff2" -o assets/fonts/lora.woff2
    curl -sL "https://fonts.gstatic.com/s/lora/v37/0QI8MX1D_JOuMw_hLdO6T2wV9KnW-MoFoq92nA.woff2" -o assets/fonts/lora-italic.woff2
    @echo "fonts downloaded to assets/fonts/"

# 📦 Create a new release with commit-and-tag-version (requires global installation: npm i -g commit-and-tag-version) (args optional, e.g. --release-as 0.1.0)
release *ARGS:
    #!/usr/bin/env sh
    if ! command -v commit-and-tag-version > /dev/null 2>&1; then
        echo "Error: commit-and-tag-version command not found."
        echo "Please install it globally with: npm i -g commit-and-tag-version"
        exit 1
    fi
    commit-and-tag-version {{ ARGS }}

# [🤖 CI task] extract content from CHANGELOG.md for use in Gitlab/Github Releases
release-notes:
    #!/usr/bin/env bash
    set -euo pipefail

    changelog_path="{{ invocation_directory() }}/CHANGELOG.md"
    release_notes_path="{{ invocation_directory() }}/../LATEST_RELEASE_NOTES.md"

    # Create header for release notes
    echo "## What's changed in this release" > "$release_notes_path"

    # Extract content between first and second level 2 heading
    awk '/^## /{
        if (count == 0) {
            count = 1
            next
        } else if (count == 1) {
            exit
        }
    }
    count == 1 { print }' "$changelog_path" >> "$release_notes_path"

    echo "Release notes extracted to $release_notes_path"
