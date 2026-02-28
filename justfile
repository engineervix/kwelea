# kwelea — development tasks
# https://just.systems

default:
    @just --list

# Build the kwelea binary
build:
    go build -o kwelea .

# Install kwelea globally
install:
    go install .

# Run all tests
test:
    go test ./...

# Run tests with verbose output
test-v:
    go test -v ./...

# Run go vet
vet:
    go vet ./...

# Build the kwelea docs with kwelea itself (dogfood)
docs: build
    ./kwelea build

# Serve the kwelea docs locally (dogfood)
serve: build
    ./kwelea serve

# Remove build artifacts
clean:
    rm -f kwelea
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
