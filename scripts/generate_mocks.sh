#!/bin/bash

# Generate mocks from all interfaces in internal/interfaces
# Output goes to internal/testrunners/mocks

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INTERFACES_DIR="${PROJECT_ROOT}/internal/interfaces"
MOCKS_DIR="${PROJECT_ROOT}/internal/testrunners/mocks"

# Check if mockgen is installed
if ! command -v mockgen &> /dev/null; then
    echo -e "${RED}Error: mockgen is not installed${NC}"
    echo "Install with: go install github.com/golang/mock/cmd/mockgen@latest"
    exit 1
fi

mkdir -p "${MOCKS_DIR}"

echo -e "${YELLOW}Generating mocks from interfaces...${NC}"
echo "Interfaces: ${INTERFACES_DIR}"
echo "Output: ${MOCKS_DIR}"
echo ""

# Process each package: repository, handler, usecase, middleware
for package in repository handler usecase middleware; do
    package_dir="${INTERFACES_DIR}/${package}"
    output_dir="${MOCKS_DIR}/${package}"

    if [ ! -d "${package_dir}" ]; then
        continue
    fi

    mkdir -p "${output_dir}"

    echo -e "${YELLOW}Processing ${package}...${NC}"

        # For other packages, generate individual mocks
        find "${package_dir}" -maxdepth 1 -name "*.go" -type f | sort | while read file; do
            filename=$(basename "$file")

            if [[ "$filename" == ".gitkeep" ]]; then
                continue
            fi

            output_file="${output_dir}/mock_${filename}"
            echo -n "  ${filename}... "

            if mockgen -source="$file" -destination="$output_file" -package="mock${package}"; then
                echo -e "${GREEN}✓${NC}"
            else
                echo -e "${RED}✗${NC}"
            fi
        done


    echo ""
done

echo -e "${YELLOW}=== Done ===${NC}"
