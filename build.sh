#!/bin/bash

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

function help() {
    echo 'Usage:'
    grep '^##' $0 | sed -n 's/^##//p' | column -t -s ':' | sed -e 's/^/ /'
}

function confirm() {
    read -p 'Are you sure? [y/N]' ans
    [ "${ans:-N}" == "y" ]
}

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: run the main application
function run() {
    air
}

## db/migrations/new name=$1: create a new database migration
function db_migrations_new() {
    if [ -z "$1" ]; then
        echo "Usage: db_migrations_new name"
        exit 1
    fi
    name=$1
    echo "Creating migration files for ${name}..."
    migrate create -seq -ext=.sql -dir=./migrations "${name}"
}

## db/migrations/up: apply all up database migrations
function db_migrations_up() {
    if confirm; then
        echo 'Running up migrations...'
        migrate -path ./migrations -database ${NAIT_DB_DSN} up
    fi
}

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
function audit() {
    vendor
    echo 'Formatting code...'
    go fmt ./...
    echo 'Vetting code...'
    go vet ./...
    staticcheck ./...
    echo 'Running tests...'
    go test -race -vet=off ./...
}

## vendor: tidy and vendor dependencies
function vendor() {
    echo 'Tidying and verifying module dependencies...'
    go mod tidy
    go mod verify
    echo 'Vendoring dependencies...'
    go mod vendor
}

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/css: build tailwind css sheet
function build_css() {
    tailwindcss -i ./views/index.css -o ./static/output.css --minify
}

current_time=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
git_description=$(git describe --always --dirty --tags --long)
linker_flags="-s -w -X main.buildTime=${current_time} -X main.version=${git_description} -linkmode external -extldflags '-static'"

## build: build the application
function build() {
    echo 'Building cmd/api...'
    go build -ldflags="${linker_flags}" -o=./dist/server .
    build_css
    cp -r ./static ./dist
    cp -r ./views ./dist
}

# ==================================================================================== #
# MAIN LOGIC
# ==================================================================================== #

case "$1" in
    help)
        help
        ;;
    run)
        run
        ;;
    db/migrations/new)
        db_migrations_new "$2"
        ;;
    db/migrations/up)
        db_migrations_up
        ;;
    audit)
        audit
        ;;
    vendor)
        vendor
        ;;
    build/css)
        build_css
        ;;
    build)
        build
        ;;
    *)
        echo "Unknown command: $1"
        help
        ;;
esac
