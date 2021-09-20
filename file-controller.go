package main

import (
	"encoding/base64"
	"os"
)

func writeToFile(filePath string, based64Img string) {
	dec, err := base64.StdEncoding.DecodeString(based64Img)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}
