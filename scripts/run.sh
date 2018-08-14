#!/usr/bin/env bash

# Install the vcs server executable
go install ./cmd/gopx-vcs

# Run the server
gopx-vcs