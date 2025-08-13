#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Installing Tandem...${NC}"

# Function to check if running as root
check_root() {
    if [ "$EUID" -ne 0 ]; then 
        echo -e "${YELLOW}Some installation steps might require sudo privileges${NC}"
    fi
}

# Function to install Git based on the system
install_git() {
    if command -v apt-get &> /dev/null; then
        sudo apt-get update
        sudo apt-get install -y git
    elif command -v yum &> /dev/null; then
        sudo yum install -y git
    elif command -v pacman &> /dev/null; then
        sudo pacman -S --noconfirm git
    else
        echo -e "${RED}Could not install git. Please install it manually.${NC}"
        exit 1
    fi
}

# Function to install Go based on the system
install_go() {
    local GO_VERSION="1.24.2"
    echo -e "${YELLOW}Installing Go ${GO_VERSION}...${NC}"
    
    if command -v apt-get &> /dev/null; then
        sudo apt-get update
        sudo apt-get install -y wget
        wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
        rm go${GO_VERSION}.linux-amd64.tar.gz
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        source ~/.bashrc
    elif command -v yum &> /dev/null; then
        sudo yum install -y wget
        wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
        rm go${GO_VERSION}.linux-amd64.tar.gz
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        source ~/.bashrc
    else
        echo -e "${RED}Could not install Go. Please install it manually.${NC}"
        exit 1
    fi
}

# Function to install Nix
install_nix() {
    if ! command -v nix &> /dev/null; then
        echo -e "${YELLOW}Installing Nix...${NC}"
        curl -L https://nixos.org/nix/install | sh
        # Source nix
        . ~/.nix-profile/etc/profile.d/nix.sh
    fi
}

# Check and install git if needed
if ! command -v git &> /dev/null; then
    echo -e "${YELLOW}Git not found. Installing...${NC}"
    install_git
fi

# Function to handle repository setup
setup_repository() {
    # If we're already in a directory with main.go, we're good
    if [ -f "main.go" ]; then
        echo -e "${GREEN}Repository already set up in current directory${NC}"
        return 0
    fi

    # If tandem directory exists but we're not in it
    if [ -d "tandem" ]; then
        echo -e "${YELLOW}Found existing tandem directory${NC}"
        read -p "Would you like to remove it and clone again? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf tandem
        else
            echo -e "${YELLOW}Trying to use existing directory...${NC}"
            cd tandem || {
                echo -e "${RED}Failed to enter existing project directory${NC}"
                exit 1
            }
            if [ -f "main.go" ]; then
                echo -e "${GREEN}Successfully entered existing repository${NC}"
                return 0
            else
                echo -e "${RED}Existing directory seems invalid${NC}"
                read -p "Would you like to remove it and clone again? (y/n) " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    cd ..
                    rm -rf tandem
                else
                    echo -e "${RED}Cannot proceed without valid repository${NC}"
                    exit 1
                fi
            fi
        fi
    fi

    # At this point, either there was no tandem directory or we removed it
    echo -e "${YELLOW}Cloning repository...${NC}"
    if git clone --depth 1 https://github.com/yyovil/tandem.git; then
        echo -e "${GREEN}Successfully cloned repository${NC}"
        cd tandem || {
            echo -e "${RED}Failed to enter project directory${NC}"
            exit 1
        }
        return 0
    else
        echo -e "${RED}Failed to clone repository. Please check your internet connection and try again.${NC}"
        echo -e "${YELLOW}You can also manually download the project from:${NC}"
        echo -e "${GREEN}https://github.com/yyovil/tandem/archive/refs/heads/main.zip${NC}"
        exit 1
    fi
}

# Setup repository
setup_repository

# Check for Go installation
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}Go is not installed. Would you like to install it? (y/n)${NC}"
    read -p "" -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        install_go
    else
        echo -e "${RED}Cannot proceed without Go installation.${NC}"
        exit 1
    fi
fi

# Check Go version
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}Go is not installed. Would you like to install it? (y/n)${NC}"
    read -p "" -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        install_go
    else
        echo -e "${RED}Cannot proceed without Go installation.${NC}"
        exit 1
    fi
else
    GO_VERSION=$(go version | grep -oP "go\K[0-9]+\.[0-9]+\.[0-9]+")
    echo -e "${GREEN}Found Go version ${GO_VERSION}${NC}"
    if [ "$(printf '%s\n' "1.24.2" "$GO_VERSION" | sort -V | head -n1)" = "1.24.2" ]; then
        echo -e "${GREEN}Go version is sufficient for building Tandem${NC}"
    else
        echo -e "${RED}Go version 1.24.2 or higher is required. Current version: ${GO_VERSION}${NC}"
        exit 1
    fi
fi

# Function to build tandem
build_tandem() {
    echo -e "${YELLOW}Building Tandem...${NC}"
    
    # Create necessary directories
    mkdir -p .tandem
    mkdir -p .go
    
    # Install dependencies
    echo -e "${YELLOW}Installing Go dependencies...${NC}"
    go mod download || {
        echo -e "${RED}Failed to download Go dependencies${NC}"
        exit 1
    }
    go mod tidy || {
        echo -e "${RED}Failed to tidy Go modules${NC}"
        exit 1
    }

    # Build the binary
    echo -e "${YELLOW}Compiling Tandem...${NC}"
    go build -o tandem main.go || {
        echo -e "${RED}Failed to build Tandem${NC}"
        exit 1
    }
    chmod +x tandem || {
        echo -e "${RED}Failed to make binary executable${NC}"
        exit 1
    }
    echo -e "${GREEN}Successfully built Tandem binary${NC}"
}

# Function to build with system Go
build_with_system_go() {
    # Check Go version
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}Go is not installed. Would you like to install it? (y/n)${NC}"
        read -p "" -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            install_go
        else
            echo -e "${RED}Cannot proceed without Go installation.${NC}"
            exit 1
        fi
    else
        GO_VERSION=$(go version | grep -oP "go\K[0-9]+\.[0-9]+\.[0-9]+")
        echo -e "${GREEN}Found Go version ${GO_VERSION}${NC}"
        if [ "$(printf '%s\n' "1.24.2" "$GO_VERSION" | sort -V | head -n1)" = "1.24.2" ]; then
            echo -e "${GREEN}Go version is sufficient for building Tandem${NC}"
        else
            echo -e "${RED}Go version 1.24.2 or higher is required. Current version: ${GO_VERSION}${NC}"
            exit 1
        fi
    fi
    
    build_tandem
}

# Ask about Nix installation
echo -e "${YELLOW}Would you like to set up the Nix development environment? (y/n)${NC}"
read -p "" -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    install_nix
    echo -e "${GREEN}Running nix develop...${NC}"
    # Start a new nix-shell and execute the build without Go version check
    nix develop --command bash -c "$(declare -f build_tandem); echo -e '${GREEN}Using Nix-provided Go...${NC}'; build_tandem"
else
    # Build with system Go
    build_with_system_go
fi

# Create necessary directories
mkdir -p .tandem
mkdir -p .go

# Setup environment file if it doesn't exist
if [ ! -f .env ]; then
    cp .example.env .env
    echo "Created .env file. Please edit it to add your API keys"
fi

# Install dependencies
echo "Installing Go dependencies..."
go mod download
go mod tidy

# Build the binary
echo "Building Tandem..."
go build -o tandem main.go

# Make the binary executable
chmod +x tandem



# Verify the build was successful
if [ -f "tandem" ] && [ -x "tandem" ]; then
    echo -e "${GREEN}Installation complete!${NC}"
    echo -e "${YELLOW}Next steps:${NC}"
    echo -e "1. Edit ${GREEN}.env${NC} file to add your API keys"
    echo -e "2. Run ${GREEN}./tandem${NC} to start the application"
    
    # Show current directory and binary
    echo -e "\n${YELLOW}Current directory:${NC}"
    pwd
    echo -e "\n${YELLOW}Tandem binary:${NC}"
    ls -l tandem
    
    echo -e "\n${YELLOW}For development:${NC}"
    echo -e "- Use ${GREEN}nix develop${NC} to enter the development environment"
    echo -e "- All required tools are installed and configured"
else
    echo -e "${RED}Installation seems to have failed - tandem binary not found or not executable${NC}"
    echo -e "${YELLOW}Please check the error messages above and try again${NC}"
    exit 1
fi
