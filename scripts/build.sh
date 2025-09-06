#!/bin/bash

# Build the Go application
go build -o event-booking-api ./cmd/api/

# Check if the build was successful
if [ $? -eq 0 ]; then
  echo "Build successful. Starting API..."
  # Run the executable
  ./event-booking-api
else
  echo "Build failed. Please check the errors."
  exit 1
fi