# Paths to tools needed in dependencies
GO := $(shell which go)

# Build flags
BUILD_FLAGS = -ldflags "-s -w" 

# Output directory
BUILD_DIR := build

# Output binary name
OUTPUT_NAME := udpxy-go

# Targets
all: clean dependencies build build-openwrt


build: mkdir
	@echo "Building main.go with CGO enabled..."
	@CGO_ENABLED=1 ${GO} build ${BUILD_FLAGS} -o ${BUILD_DIR}/${OUTPUT_NAME} main.go


build-openwrt: mkdir
	@echo "Building main.go for OpenWrt..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 ${GO} build ${BUILD_FLAGS} -o ${BUILD_DIR}/openwrt_amd64_${OUTPUT_NAME} main.go

FORCE:

dependencies:
	@echo "Updating system and installing dependencies..."
	@sudo apt update
	@sudo apt install libavcodec-dev libavutil-dev libswscale-dev
	@test -f "${GO}" && test -x "${GO}"  || (echo "Missing go binary" && exit 1)

mkdir:
	@echo "Creating build directory..."
	mkdir -p ${BUILD_DIR}

clean:
	@echo "Cleaning up..."
	@${GO} mod tidy
	@${GO} clean
	@rm -fr $(BUILD_DIR)
