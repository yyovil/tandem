#!/usr/bin/env bash
set -euo pipefail
APP=tandem

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# XDG config home for storing Tandem config; set a default if not provided
XDG_CONFIG_HOME=${XDG_CONFIG_HOME:-$HOME/.config}

# Messaging helper used throughout the script
print_message() {
    local level=$1
    local message=$2
    local color=""

    case $level in
        info) color="${GREEN}" ;;
        warning) color="${YELLOW}" ;;
        error) color="${RED}" ;;
    esac

    echo -e "${color}${message}${NC}"
}

# Early Docker environment check: verify CLI presence and daemon availability
check_docker_installation() {
    # Exit codes:
    #   0   -> Docker CLI present and daemon reachable
    #   69  -> Service unavailable (daemon not running)
    #   127 -> Command not found (docker CLI missing)

    if ! command -v docker &>/dev/null; then
        print_message error "Docker CLI not found on this system. Please install Docker for your OS."
        # 127 is the conventional exit code for 'command not found'
        return 127
    fi

    if ! docker info &>/dev/null; then
        print_message warning "Docker daemon (dockerd) is not running or unreachable."
        # 69 corresponds to EX_UNAVAILABLE (service unavailable)
        return 69
    fi

    print_message info "Docker is installed and the daemon is running."
    return 0
}

# Ensure Docker environment is ready before any further work
check_docker_installation

requested_version=${VERSION:-}

os=$(uname -s | tr '[:upper:]' '[:lower:]')
if [[ "$os" == "darwin" ]]; then
    os="mac"
fi

arch=$(uname -m)
if [[ "$arch" == "aarch64" ]]; then
  arch="arm64"
fi

filename="$APP_$os_$arch.tar.gz"
case "$filename" in
    *"-linux-"*)
        [[ "$arch" == "x86_64" || "$arch" == "arm64" || "$arch" == "i386" ]] || exit 1
    ;;
    *"-mac-"*)
        [[ "$arch" == "x86_64" || "$arch" == "arm64" ]] || exit 1
    ;;
    *)
        echo "${RED}Unsupported OS/Arch: $os/$arch${NC}"
        exit 1
    ;;
esac

INSTALL_DIR=$HOME/.tandem/bin
mkdir -p "$INSTALL_DIR"

if [ -z "$requested_version" ]; then
    url="https://github.com/yyovil/tandem/releases/latest/download/$filename"
    specific_version=$(curl -s https://api.github.com/repos/yyovil/tandem/releases/latest | awk -F'"' '/"tag_name": "/ {gsub(/^v/, "", $4); print $4}')

    if [[ $? -ne 0 ]]; then
        echo "${RED}Failed to fetch version information${NC}"
        exit 1
    fi
else
    url="https://github.com/yyovil/tandem/releases/download/v${requested_version}/$filename"
    specific_version=$requested_version
fi

 

check_version() {
    # NOTE: redirects both the stdout and stderr to /dev/null (nowhere). only concerned with the status code [0-255] here.
    if command -v tandem &>/dev/null; then
        installed_version=$(tandem --version | awk '{print $1}')
        if [[ "$installed_version" != "$specific_version" ]]; then
            print_message info "Installed version: ${YELLOW}$installed_version."
        else
            print_message info "Version ${YELLOW}$specific_version${GREEN} already installed"
            exit 0
        fi
    fi
}

download_and_install() {
    print_message info "Downloading ${BLUE}tandem ${GREEN}version: ${YELLOW}$specific_version ${GREEN}..."
    curl -# -L $url | tar xz
    cd tandem_tmp
    mv tandem $INSTALL_DIR
    # Put default config/docs into the user's XDG config directory without overwriting existing files.
    local app_cfg_dir="$XDG_CONFIG_HOME/$APP"
    mkdir -p "$app_cfg_dir"

    # Files expected in the archive root we want to place under config dir
    local files=("swarm.json" "LICENSE" "README.md" "tandem.example.env")
    for f in "${files[@]}"; do
        if [[ -f "$f" ]]; then
            local target="$app_cfg_dir/$f"
            if [[ -e "$target" ]]; then
                print_message warning "Skipping ${YELLOW}$f${NC}: already exists in $app_cfg_dir"
            else
                cp "$f" "$target"
                print_message info "Placed default ${BLUE}$f${GREEN} into $app_cfg_dir"
            fi
        else
            print_message warning "Expected ${YELLOW}$f${NC} not found in extracted archive"
        fi
    done
    cd .. && rm -rf tandem_tmp

    # Pull the Docker image and create/start a container (assumes Docker is available and running)
    local IMAGE="ghcr.io/yyovil/kali:headless"
    local CONTAINER_NAME="tandem"

    print_message info "Ensuring Docker image ${BLUE}$IMAGE${GREEN} is available..."
    if ! docker image inspect "$IMAGE" &>/dev/null; then
        docker pull "$IMAGE"
    fi

    # Create the container if it doesn't already exist; do not start it here
    if docker inspect -f '{{.State.Status}}' "$CONTAINER_NAME" &>/dev/null; then
        local status
        status=$(docker inspect -f '{{.State.Status}}' "$CONTAINER_NAME")
        if [[ "$status" == "running" ]]; then
            print_message info "Docker container ${BLUE}$CONTAINER_NAME${GREEN} already exists and is running."
        else
            print_message info "Docker container ${BLUE}$CONTAINER_NAME${GREEN} already exists (status: ${YELLOW}$status${GREEN})."
            print_message info "Start it when ready: docker start -ai $CONTAINER_NAME"
        fi
    else
        print_message info "Creating Docker container ${BLUE}$CONTAINER_NAME${GREEN} with hostname ${YELLOW}$CONTAINER_NAME${GREEN}..."
        docker create --privileged --hostname "$CONTAINER_NAME" --name "$CONTAINER_NAME" -it "$IMAGE" /bin/bash
        print_message info "Created container ${BLUE}$CONTAINER_NAME${GREEN}. Start it with: docker start -ai $CONTAINER_NAME"
    fi
}

check_version
download_and_install

add_to_path() {
    local config_file=$1
    local command=$2

    if [[ -w $config_file ]]; then
        echo -e "\n# tandem" >> "$config_file"
        echo "$command" >> "$config_file"
        print_message info "Successfully added ${BLUE}tandem ${GREEN}to \$PATH in $config_file"
    else
        print_message warning "Manually add the directory to $config_file (or similar):"
        print_message info "  $command"
    fi
}

# figuring out what config file should we use to put our app specific config into.
current_shell=$(basename "$SHELL")
case $current_shell in
    fish)
        config_files="$HOME/.config/fish/config.fish"
    ;;
    zsh)
        config_files="$HOME/.zshrc $HOME/.zshenv $XDG_CONFIG_HOME/zsh/.zshrc $XDG_CONFIG_HOME/zsh/.zshenv"
    ;;
    bash)
        config_files="$HOME/.bashrc $HOME/.bash_profile $HOME/.profile $XDG_CONFIG_HOME/bash/.bashrc $XDG_CONFIG_HOME/bash/.bash_profile"
    ;;
    ash)
        config_files="$HOME/.ashrc $HOME/.profile /etc/profile"
    ;;
    sh)
        config_files="$HOME/.ashrc $HOME/.profile /etc/profile"
    ;;
    *)
        # Default case if none of the above matches
        config_files="$HOME/.bashrc $HOME/.bash_profile $XDG_CONFIG_HOME/bash/.bashrc $XDG_CONFIG_HOME/bash/.bash_profile"
    ;;
esac

config_file=""
for file in $config_files; do
    if [[ -f $file ]]; then
        config_file=$file
        break
    fi
done

if [[ -z $config_file ]]; then
    print_message error "No config file found for $current_shell. Checked files: ${config_files[@]}"
    exit 1
fi

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    case $current_shell in
        fish)
            add_to_path "$config_file" "fish_add_path $INSTALL_DIR"
        ;;
        zsh)
            add_to_path "$config_file" "export PATH=$INSTALL_DIR:\$PATH"
        ;;
        bash)
            add_to_path "$config_file" "export PATH=$INSTALL_DIR:\$PATH"
        ;;
        ash)
            add_to_path "$config_file" "export PATH=$INSTALL_DIR:\$PATH"
        ;;
        sh)
            add_to_path "$config_file" "export PATH=$INSTALL_DIR:\$PATH"
        ;;
        *)
            print_message warning "Manually add the directory to $config_file (or similar):"
            print_message info "  export PATH=$INSTALL_DIR:\$PATH"
        ;;
    esac
fi

