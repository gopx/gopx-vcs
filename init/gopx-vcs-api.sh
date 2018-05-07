#!/usr/bin/env bash

export $(cat init/.vcs-api.env | xargs) && go run cmd/gopx-vcs-api/*.go