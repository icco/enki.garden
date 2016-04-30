package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var files = []EnkiFile{}

type EnkiFile struct {
	Size    int64
	ModTime time.Time
	Path    string
}

func registerFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}

	var size int64
	if !info.IsDir() {
		size = info.Size()
	}

	files = append(files, EnkiFile{
		Path:    path,
		ModTime: info.ModTime(),
		Size:    size,
	})

	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	dir := os.Args[1]
	err := filepath.Walk(dir, registerFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range files {
		fmt.Printf("%+v\n", v)
	}
}
