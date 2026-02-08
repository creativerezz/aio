#!/bin/bash
# Aio Installer Script
# Usage: curl -fsSL https://raw.githubusercontent.com/creativerezz/aio/main/scripts/installer/install.sh | bash
# Usage with custom directory: curl -fsSL https://raw.githubusercontent.com/creativerezz/aio/main/scripts/installer/install.sh | INSTALL_DIR=/usr/local/bin bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print functions
print_info() {
    printf "${BLUE}[INFO]${NC} %s\n" "$1"
}

print_success() {
    printf "${GREEN}[SUCCESS]${NC} %s\n" "$1"
}

print_warning() {
    printf "${YELLOW}[WARNING]${NC} %s\n" "$1"
}

print_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1" >&2
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            echo "Darwin"
            ;;
        Linux*)
            echo "Linux"
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            print_error "This installer only supports Darwin (macOS) and Linux"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            echo "x86_64"
            ;;
        arm64|aarch64)
            echo "arm64"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            print_error "This installer only supports x86_64 and arm64"
            exit 1
            ;;
    esac
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Download and extract aio
install_aio() {
    local os="$1"
    local arch="$2"
    local install_dir="$3"

    # Construct download URL
    local filename="aio_${os}_${arch}.tar.gz"
    local download_url="https://github.com/creativerezz/aio/releases/latest/download/${filename}"

    print_info "Downloading Aio for ${os} ${arch}..."
    print_info "URL: ${download_url}"

    # Create temporary directory
    local temp_dir
    temp_dir=$(mktemp -d)
    local temp_file="${temp_dir}/aio.tar.gz"

    # Download the archive
    if command_exists curl; then
        if ! curl -fsSL "${download_url}" -o "${temp_file}"; then
            print_error "Failed to download Aio"
            rm -rf "${temp_dir}"
            exit 1
        fi
    elif command_exists wget; then
        if ! wget -q "${download_url}" -O "${temp_file}"; then
            print_error "Failed to download Aio"
            rm -rf "${temp_dir}"
            exit 1
        fi
    else
        print_error "Neither curl nor wget found. Please install one of them and try again."
        exit 1
    fi

    print_info "Extracting Aio binary..."

    # Extract only the aio binary from the archive
    if ! tar -xzf "${temp_file}" -C "${temp_dir}" aio; then
        print_error "Failed to extract Aio binary"
        rm -rf "${temp_dir}"
        exit 1
    fi

    # Create install directory if it doesn't exist
    if [ ! -d "${install_dir}" ]; then
        print_info "Creating install directory: ${install_dir}"
        if ! mkdir -p "${install_dir}"; then
            print_error "Failed to create install directory: ${install_dir}"
            print_error "You may need to run with sudo or choose a different directory"
            rm -rf "${temp_dir}"
            exit 1
        fi
    fi

    # Move binary to install directory
    print_info "Installing Aio to ${install_dir}/aio..."
    if ! mv "${temp_dir}/aio" "${install_dir}/aio"; then
        print_error "Failed to install Aio to ${install_dir}"
        print_error "You may need to run with sudo or choose a different directory"
        rm -rf "${temp_dir}"
        exit 1
    fi

    # Make sure it's executable
    chmod +x "${install_dir}/aio"

    # Clean up
    rm -rf "${temp_dir}"

    print_success "Aio installed successfully to ${install_dir}/aio"
}

# Check PATH and provide instructions
check_path() {
    local install_dir="$1"

    if echo "$PATH" | grep -q "${install_dir}"; then
        print_success "✅ ${install_dir} is already in your PATH"
    else
        print_warning "⚠️  ${install_dir} is not in your PATH"
        print_info "To use aio from anywhere, add the following to your shell profile:"
        print_info "  export PATH=\"\$PATH:${install_dir}\""
        print_info ""
        print_info "For bash, add it to ~/.bashrc or ~/.bash_profile"
        print_info "For zsh, add it to ~/.zshrc"
        print_info "For fish, run: fish_add_path ${install_dir}"
    fi
}

# Verify installation
verify_installation() {
    local install_dir="$1"
    local aio_path="${install_dir}/aio"

    if [ -x "${aio_path}" ]; then
        print_info "Verifying installation..."
        local version
        if version=$("${aio_path}" --version 2>/dev/null); then
            print_success "Aio ${version} is working correctly!"
        else
            print_warning "Aio binary exists but --version failed"
        fi
    else
        print_error "Aio binary not found at ${aio_path}"
        exit 1
    fi
}

# Main installation function
main() {
    print_info "🚀 Starting Aio installation..."

    # Detect system
    local os
    local arch
    os=$(detect_os)
    arch=$(detect_arch)

    print_info "Detected system: ${os} ${arch}"

    # Determine install directory
    local install_dir="${INSTALL_DIR:-${HOME}/.local/bin}"

    print_info "Install directory: ${install_dir}"

    # Install aio
    install_aio "${os}" "${arch}" "${install_dir}"

    # Verify installation
    verify_installation "${install_dir}"

    # Check PATH
    check_path "${install_dir}"

    print_info ""
    print_success "🎉 Installation complete!"
    print_info ""
    print_info "Next steps:"
    print_info "  1. Run 'aio --setup' to configure Aio"
    print_info "  2. Add your API keys and preferences"
    print_info "  3. Start using Aio with 'aio --help'"
    print_info ""
    print_info "Documentation: https://github.com/creativerezz/aio"
}

# Run main function
main "$@"