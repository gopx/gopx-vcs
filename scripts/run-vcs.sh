#!/usr/bin/env bash

# Install the server executable
go install ./cmd/gopx-vcs

# Run the server
gopx-vcs
