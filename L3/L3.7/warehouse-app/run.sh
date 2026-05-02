#!/bin/bash

set -e

echo "Cleaning..."
make down || true

echo "Installing deps..."
make deps

echo "Running tests..."
make test

echo "Starting app..."
make run