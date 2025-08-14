#!/bin/bash

# Tandem Install Script
# This script installs the tandem binary and sets up the required Docker environment
# Usage: ./install.sh [--dry-run]

set -e

# Parse command line arguments
DRY_RUN=false
if [ "$1" = "--dry-run" ]; then
    DRY_RUN=true
    echo "🔍 DRY RUN MODE - No actual changes will be made"
    echo
fi

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Repository information
REPO_OWNER="yyovil"
REPO_NAME="tandem"
DOCKER_IMAGE="kali:withtools"
CONTAINER_NAME="tandem-kali"

# Installation directory
INSTALL_DIR="/usr/local/bin"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${CYAN}[STEP]${NC} $1"
}

# Function to detect OS and architecture
detect_platform() {
    print_step "Detecting system platform..."
    
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $OS in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        msys*|mingw*|cygwin*)
            OS="windows"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
    
    case $ARCH in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="arm"
            ;;
        i386|i686)
            ARCH="386"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    print_status "Detected platform: ${OS}/${ARCH}"
}

# Function to get the latest release version
get_latest_version() {
    print_step "Fetching latest release information..."
    
    # Try to get the latest release version from GitHub API
    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -s "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' 2>/dev/null)
        if [ -z "$VERSION" ] || [ "$VERSION" = "null" ]; then
            print_warning "Could not fetch latest release. Trying tags endpoint..."
            VERSION=$(curl -s "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/tags" 2>/dev/null | grep '"name":' | head -1 | sed -E 's/.*"([^"]+)".*/\1/' 2>/dev/null)
        fi
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -qO- "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' 2>/dev/null)
        if [ -z "$VERSION" ] || [ "$VERSION" = "null" ]; then
            print_warning "Could not fetch latest release. Trying tags endpoint..."
            VERSION=$(wget -qO- "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/tags" 2>/dev/null | grep '"name":' | head -1 | sed -E 's/.*"([^"]+)".*/\1/' 2>/dev/null)
        fi
    else
        print_error "curl or wget is required to download the binary"
        exit 1
    fi
    
    # If we still don't have a version, offer fallback options
    if [ -z "$VERSION" ] || [ "$VERSION" = "null" ]; then
        print_warning "Could not determine the latest version from GitHub API"
        print_status "This might be due to network restrictions or API limits"
        
        # Offer to build from source as fallback
        if command -v go >/dev/null 2>&1; then
            print_status "Go is available. Offering to build from source..."
            read -p "Would you like to build from source instead? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                BUILD_FROM_SOURCE=true
                return
            fi
        fi
        
        # Ask user to provide version manually
        print_status "You can also specify a version manually (e.g., v0.1.0)"
        read -p "Enter version (or press Enter to exit): " VERSION
        if [ -z "$VERSION" ]; then
            print_error "No version specified. Exiting."
            exit 1
        fi
    fi
    
    print_status "Using version: $VERSION"
}

# Function to construct download URL
construct_download_url() {
    # Common archive naming patterns for Go binaries
    # Try different naming conventions
    local base_name="tandem"
    local file_extension=""
    
    if [ "$OS" = "windows" ]; then
        file_extension=".exe"
    fi
    
    # Try different naming patterns
    DOWNLOAD_URLS=(
        "https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${VERSION}/${base_name}_${OS}_${ARCH}${file_extension}"
        "https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${VERSION}/${base_name}-${OS}-${ARCH}${file_extension}"
        "https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${VERSION}/${base_name}_${VERSION}_${OS}_${ARCH}${file_extension}"
        "https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${VERSION}/${base_name}-${VERSION}-${OS}-${ARCH}${file_extension}"
    )
    
    print_status "Trying to find downloadable binary..."
}

# Function to build from source
build_from_source() {
    print_step "Building tandem from source..."
    
    # Check if we're in a git repository
    if [ ! -d ".git" ]; then
        print_status "Cloning repository..."
        TEMP_DIR=$(mktemp -d)
        trap "rm -rf $TEMP_DIR" EXIT
        
        git clone "https://github.com/${REPO_OWNER}/${REPO_NAME}.git" "$TEMP_DIR"
        cd "$TEMP_DIR"
    else
        print_status "Using current repository..."
    fi
    
    # Build the binary
    print_status "Building binary..."
    if go build -o tandem .; then
        print_success "Binary built successfully"
        
        # Install binary
        if [ -w "$INSTALL_DIR" ]; then
            mv tandem "$INSTALL_DIR/tandem"
        else
            print_status "Installing to $INSTALL_DIR requires sudo privileges..."
            sudo mv tandem "$INSTALL_DIR/tandem"
        fi
        
        print_success "tandem binary installed to $INSTALL_DIR/tandem"
    else
        print_error "Failed to build binary"
        exit 1
    fi
}

# Function to download and install binary
download_and_install() {
    print_step "Downloading and installing tandem binary..."
    
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    trap "rm -rf $TEMP_DIR" EXIT
    
    # Try each download URL
    DOWNLOAD_SUCCESS=false
    for url in "${DOWNLOAD_URLS[@]}"; do
        print_status "Trying: $url"
        
        if command -v curl >/dev/null 2>&1; then
            if curl -L -f -o "$TEMP_DIR/tandem" "$url" 2>/dev/null; then
                DOWNLOAD_SUCCESS=true
                break
            fi
        elif command -v wget >/dev/null 2>&1; then
            if wget -q -O "$TEMP_DIR/tandem" "$url" 2>/dev/null; then
                DOWNLOAD_SUCCESS=true
                break
            fi
        fi
    done
    
    if [ "$DOWNLOAD_SUCCESS" = false ]; then
        print_error "Failed to download binary. Please check if releases are available for ${OS}/${ARCH}"
        print_status "You can manually download from: https://github.com/${REPO_OWNER}/${REPO_NAME}/releases"
        exit 1
    fi
    
    print_success "Binary downloaded successfully"
    
    # Make binary executable
    chmod +x "$TEMP_DIR/tandem"
    
    # Install binary
    if [ "$DRY_RUN" = true ]; then
        print_status "[DRY RUN] Would install binary to $INSTALL_DIR/tandem"
    else
        if [ -w "$INSTALL_DIR" ]; then
            mv "$TEMP_DIR/tandem" "$INSTALL_DIR/tandem"
        else
            print_status "Installing to $INSTALL_DIR requires sudo privileges..."
            sudo mv "$TEMP_DIR/tandem" "$INSTALL_DIR/tandem"
        fi
    fi
    
    print_success "tandem binary installed to $INSTALL_DIR/tandem"
}

# Function to check if Docker is installed and running
check_docker() {
    print_step "Checking Docker installation and status..."
    
    if ! command -v docker >/dev/null 2>&1; then
        print_error "Docker is not installed. Please install Docker and run this script again."
        print_status "Installation instructions: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    print_success "Docker is installed"
    
    # Check if Docker daemon is running
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker daemon is not running. Please start Docker and run this script again."
        print_status "Try: sudo systemctl start docker (on Linux) or start Docker Desktop"
        exit 1
    fi
    
    print_success "Docker daemon is running"
}

# Function to pull Docker image
pull_docker_image() {
    print_step "Setting up required Docker image: $DOCKER_IMAGE"
    
    # Check if the image already exists locally
    if docker images --filter "reference=$DOCKER_IMAGE" --format '{{.Repository}}:{{.Tag}}' | grep -q "$DOCKER_IMAGE"; then
        print_success "Docker image $DOCKER_IMAGE already exists locally"
        return
    fi
    
    # First, try multiple possible sources for the Docker image
    POSSIBLE_IMAGES=(
        "ghcr.io/yyovil/kali:withtools"
        "ghcr.io/yyovil/$DOCKER_IMAGE"
    )
    
    IMAGE_PULLED=false
    
    for img in "${POSSIBLE_IMAGES[@]}"; do
        print_status "Trying to pull: $img"
        
        if docker pull "$img" >/dev/null 2>&1; then
            print_success "Successfully pulled: $img"
            
            # Tag it with the expected name
            print_status "Tagging image as $DOCKER_IMAGE"
            docker tag "$img" "$DOCKER_IMAGE"
            
            IMAGE_PULLED=true
            break
        else
            print_warning "Failed to pull: $img"
        fi
    done
    
    if [ "$IMAGE_PULLED" = false ]; then
        print_warning "Could not pull pre-built Docker image from any remote source"
        print_status ""
        print_status "The tandem project uses a custom Kali Linux image with penetration testing tools."
        print_status "This image needs to be built locally using Nix or Docker."
        print_status ""
        
        # Check if we're in the tandem repository and have the Nix file
        if [ -f "kali.nix" ]; then
            print_status "Found kali.nix file in current directory."
            
            # Check if Nix is available
            if command -v nix >/dev/null 2>&1; then
                print_status "Nix is available. You can build the image with:"
                print_status "  nix-build kali.nix"
                print_status "  docker load < result"
                print_status ""
                
                read -p "Would you like to build the image now using Nix? (y/N): " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    print_step "Building Docker image with Nix..."
                    if nix-build kali.nix; then
                        print_status "Loading image into Docker..."
                        docker load < result
                        print_success "Docker image built and loaded successfully"
                        IMAGE_PULLED=true
                    else
                        print_error "Failed to build image with Nix"
                    fi
                fi
            else
                print_status "Nix is not installed. To build the image:"
                print_status "  1. Install Nix: https://nixos.org/download.html"
                print_status "  2. Run: nix-build kali.nix"
                print_status "  3. Run: docker load < result"
            fi
        else
            print_status "No kali.nix found. You may need to clone the tandem repository."
        fi
        
        print_status ""
        print_status "Alternatively, you can use a standard Kali image as a workaround:"
        print_status "  docker pull kalilinux/kali-rolling"
        print_status "  docker tag kalilinux/kali-rolling $DOCKER_IMAGE"
        print_status ""
        
        read -p "Do you want to use the standard Kali image as a fallback? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_step "Pulling standard Kali Linux image..."
            if docker pull kalilinux/kali-rolling; then
                print_status "Tagging as $DOCKER_IMAGE"
                docker tag kalilinux/kali-rolling "$DOCKER_IMAGE"
                print_success "Fallback image set up successfully"
                print_warning "Note: This image may not have all the tools expected by tandem"
                IMAGE_PULLED=true
            else
                print_error "Failed to pull fallback image"
            fi
        fi
        
        if [ "$IMAGE_PULLED" = false ]; then
            print_warning "Continuing without Docker image - you'll need to set it up manually"
            print_status "Tandem will still work, but terminal operations may fail until the image is available"
            SKIP_CONTAINER=true
        fi
    else
        print_success "Docker image ready"
    fi
}

# Function to create Docker container
create_docker_container() {
    # Skip if image wasn't pulled
    if [ "$SKIP_CONTAINER" = true ]; then
        print_warning "Skipping container creation due to missing image"
        return
    fi
    
    print_step "Creating Docker container: $CONTAINER_NAME"
    
    # Check if container already exists
    if docker ps -a --filter "name=$CONTAINER_NAME" --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        print_warning "Container $CONTAINER_NAME already exists"
        
        # Ask user if they want to remove and recreate
        read -p "Do you want to remove and recreate the container? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Removing existing container..."
            docker rm -f "$CONTAINER_NAME" || true
        else
            print_status "Using existing container"
            return
        fi
    fi
    
    # Create container with appropriate settings
    # Mount current directory to allow tandem to work with local files
    print_status "Creating container with current directory mounted..."
    
    if docker create \
        --name "$CONTAINER_NAME" \
        --workdir /workspace \
        -v "$(pwd):/workspace" \
        -it \
        "$DOCKER_IMAGE" \
        /bin/bash; then
        print_success "Docker container created successfully"
    else
        print_error "Failed to create Docker container"
        exit 1
    fi
}

# Function to verify installation
verify_installation() {
    print_step "Verifying installation..."
    
    # Check if binary is accessible
    if command -v tandem >/dev/null 2>&1; then
        TANDEM_VERSION=$(tandem --version 2>/dev/null || echo "unknown")
        print_success "tandem binary is accessible: $TANDEM_VERSION"
    else
        print_warning "tandem binary is not in PATH. You may need to add $INSTALL_DIR to your PATH"
        print_status "Add this line to your ~/.bashrc or ~/.zshrc:"
        print_status "export PATH=\"$INSTALL_DIR:\$PATH\""
    fi
    
    # Check Docker container (only if not skipped)
    if [ "$SKIP_CONTAINER" != true ]; then
        if docker ps -a --filter "name=$CONTAINER_NAME" --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
            print_success "Docker container $CONTAINER_NAME is ready"
        else
            print_warning "Docker container was not created properly"
        fi
    else
        print_warning "Docker container was skipped - manual setup required"
    fi
    
    # Check Docker image
    if docker images --filter "reference=$DOCKER_IMAGE" --format '{{.Repository}}:{{.Tag}}' | grep -q "$DOCKER_IMAGE"; then
        print_success "Docker image $DOCKER_IMAGE is available"
    else
        print_warning "Docker image is not available - tandem may not work correctly"
    fi
}

# Function to print final instructions
print_final_instructions() {
    echo
    print_success "🎉 Tandem installation completed successfully!"
    echo
    print_status "Next steps:"
    echo "  1. Create a RoE.md file in your project directory with Rules of Engagement"
    echo "  2. Set up your AI provider API keys (see README.md for details)"
    echo "  3. Run 'tandem' to start the application"
    echo
    
    if [ "$SKIP_CONTAINER" = true ]; then
        print_warning "⚠️  Docker setup incomplete:"
        echo "  The Docker image '$DOCKER_IMAGE' still needs to be set up manually."
        echo "  Tandem's terminal operations will not work until this is resolved."
        echo
        echo "  To complete the setup:"
        if [ -f "kali.nix" ]; then
            echo "  - With Nix: nix-build kali.nix && docker load < result"
        fi
        echo "  - Or use fallback: docker pull kalilinux/kali-rolling && docker tag kalilinux/kali-rolling $DOCKER_IMAGE"
        echo
    else
        print_status "The Docker container '$CONTAINER_NAME' is ready for tandem's terminal operations."
        print_status "Your current directory is mounted to /workspace in the container."
    fi
    
    echo
    print_status "For more information, visit: https://github.com/${REPO_OWNER}/${REPO_NAME}"
}

# Main installation function
main() {
    echo "====================================="
    echo "    Tandem Installation Script"
    echo "====================================="
    echo
    
    print_status "Starting installation process..."
    echo
    
    # Check prerequisites
    print_step "Checking prerequisites..."
    
    # Check if we have necessary tools
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        print_error "Either curl or wget is required for installation"
        exit 1
    fi
    
    if ! command -v uname >/dev/null 2>&1; then
        print_error "uname command is required to detect system information"
        exit 1
    fi
    
    print_success "Prerequisites check passed"
    echo
    
    # Run installation steps
    detect_platform
    echo
    
    # Handle build from source or download binary
    BUILD_FROM_SOURCE=false
    SKIP_CONTAINER=false
    get_latest_version
    echo
    
    if [ "$BUILD_FROM_SOURCE" = true ]; then
        build_from_source
    else
        construct_download_url
        download_and_install
    fi
    echo
    
    check_docker
    echo
    
    pull_docker_image
    echo
    
    create_docker_container
    echo
    
    verify_installation
    echo
    
    print_final_instructions
}

# Handle Ctrl+C gracefully
trap 'echo; print_error "Installation interrupted by user"; exit 1' INT

# Run main function
main "$@"