package main

import (
	"embed"
	"os"
	"path/filepath"
)

var audioFiles embed.FS

func writeEmbeddedFile(name string) string {
	data, err := audioFiles.ReadFile(name)
	if err != nil {
		panic(err)
	}
	tmpFile, err := os.CreateTemp("", filepath.Base(name))
	if err != nil {
		panic(err)
	}
	if _, err := tmpFile.Write(data); err != nil {
		panic(err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}
