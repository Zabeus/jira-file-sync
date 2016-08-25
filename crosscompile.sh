#!/bin/bash
gox -arch="386 amd64" -os="windows darwin linux" -output="output/{{.Dir}}_{{.OS}}_{{.Arch}}"
zip -r output/jira-file-sync.zip output/*
