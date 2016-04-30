package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

var files = []EnkiFile{}

type EnkiFile struct {
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modified"`
	Path    string    `json:"path"`
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

	abs_path, err := filepath.Abs(path)
	if err != nil {
		log.Print(err)
		return nil
	}

	files = append(files, EnkiFile{
		Path:    abs_path,
		ModTime: info.ModTime(),
		Size:    size,
	})

	return nil
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Printf("Usage: %s path [filename.json]", os.Args[0])
		os.Exit(1)
	}

	dir := os.Args[1]
	err := filepath.Walk(dir, registerFile)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 3 {
		f, err := os.Create(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		enc := json.NewEncoder(f)
		enc.Encode(files)
	} else {
		for _, v := range files {
			log.Printf("%+v", v)
		}
	}
}
