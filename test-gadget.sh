#!/bin/bash

# Function to get the platform-specific binary name
get_platform_binary() {
  case "$(uname -s)" in
  Darwin)
    echo "test-gadget-client-macos"
    ;;
  Linux)
    echo "test-gadget-client-linux"
    ;;
  CYGWIN* | MINGW* | MSYS*)
    echo "test-gadget-client-windows.exe"
    ;;
  *)
    echo "Unsupported platform: $(uname -s)" >&2
    exit 1
    ;;
  esac
}

# Get the current script directory (similar to __file__ in Python)
script_dir="$(dirname "$(realpath "$0")")"
dist_dir="${script_dir}/.test-gadget"
binary="${dist_dir}/$(get_platform_binary)"

# Check if the binary exists
if [[ ! -f "$binary" ]]; then
  echo "Program not found: $binary" >&2
  exit 1
fi

# Execute the binary with all script arguments
"$binary" "$@"
