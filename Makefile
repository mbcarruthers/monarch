BIN_DIR = bin
PROTO_DIR = proto


ifeq ($(OS), Windows_NT)
	SHELL := powershell.exe
	.SHELLFLAGS := -NoProfile -Command
	SHELL_VERSION = $(shell (Get-Host | Select-Object Version | Format-Table -HideTableHeaders | Out-String).Trim())
	OS = $(shell "{0} {1}" -f "windows", (Get-ComputerInfo -Property OsVersion, OsArchitecture | Format-Table -HideTableHeaders | Out-String).Trim())
	PACKAGE = $(shell (Get-Content go.mod -head 1).Split(" ")[1])
	CHECK_DIR_CMD = if (!(Test-Path $@)) { $$e = [char]27; Write-Error "$$e[31mDirectory $@ doesn't exist$${e}[0m" }
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
	RM_F_CMD = Remove-Item -erroraction silentlycontinue -Force
	RM_RF_CMD = ${RM_F_CMD} -Recurse
	SERVER_BIN = ${SERVER_DIR}.exe
	CLIENT_BIN = ${CLIENT_DIR}.exe
else
	SHELL := bash
	SHELL_VERSION = $(shell echo $$BASH_VERSION)
	UNAME := $(shell uname -s)
	VERSION_AND_ARCH = $(shell uname -rm)
	ifeq ($(UNAME),Darwin)
		OS = macos ${VERSION_AND_ARCH}
	else ifeq ($(UNAME),Linux)
		OS = linux ${VERSION_AND_ARCH}
	else
    $(error OS not supported by this Makefile)
	endif
	PACKAGE = $(shell basename ${PWD})
	CHECK_DIR_CMD = test -d $@ || (echo "\033[31mDirectory $@ doesn't exist\033[0m" && false)
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	RM_F_CMD = rm -f
	RM_RF_CMD = ${RM_F_CMD} -r
	SERVER_BIN = ${SERVER_DIR}
	CLIENT_BIN = ${CLIENT_DIR}
endif

DEFAULT_GOAL := help
.PHONY: helio
project := helio

helio: $@help ## Keep this to display help when you just type 'make' in root dir

build_helio: ## build the database service
	cd helio && go build -o ${BIN_DIR}/helio ./cmd/

dBuild_helio: ## build database service for docker
	cd helio && CGO_ENABLED=0 GOOS=linux go build -o ${BIN_DIR}/helio ./cmd/

build_imageserver: ## build the image server normally
	cd imageserver && go build -o ${BIN_DIR}/imageserver ./cmd/

dBuild_imageserver: ## build the image server for docker
	cd imageserver && CGO_ENABLED=0 GOOS=linux go build -o ${BIN_DIR}/imageserver ./cmd/

dBuild_webServer: ## build the web-interface server for docker
	cd web-interface/webServer && GOOS=linux CGO_ENABLED=0 go build -o ${PWD}/web-interface/app/server/webServer ./cmd/

npm_build: ## build the web-interface(zebrafalter) for docker in app/server
	cd web-interface/zebrafalter && npm run build

npm_clean: ## clean the web-interface project from the docker location(app/client)
	cd web-interface/zebrafalter && npm run clean

npm_start: ## start the npm project independent of a server
	cd web-interface/zebrafalter && npm start

about: ## Display info related to the build
	@echo "Project: ${PACKAGE}"
	@echo "OS: ${OS}"
	@echo "Shell: ${SHELL} ${SHELL_VERSION}"
	@echo "Protoc version: $(shell protoc --version)"
	@echo "Go version: $(shell go version)"
	@echo "Go package: helio imageserver"
	@echo "Openssl version: $(shell openssl version)"

help: ## Show this help
	@${HELP_CMD}
