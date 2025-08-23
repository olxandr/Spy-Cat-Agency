#!/bin/sh

# Generate Swagger documentation
swag init --parseDependency --dir ./cmd/api,./internal

# Start the application
api
