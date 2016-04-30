package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
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
	compress := flag.Bool("c", false, "Compress Output.")
	verbose := flag.Bool("v", false, "Verbose Logging")
	filename := flag.String("name", "output.json", "Where to output json")
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		log.Printf("%+v", args)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	dir := args[0]
	err := filepath.Walk(dir, registerFile)
	if err != nil {
		log.Fatal(err)
	}

	var enc *json.Encoder
	if *compress {
		fp, err := os.Create(fmt.Sprintf("%s.gz", *filename))
		if err != nil {
			log.Fatal(err)
		}
		f, err := gzip.NewWriterLevel(fp, gzip.BestCompression)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		enc = json.NewEncoder(f)
	} else {
		f, err := os.Create(*filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		enc = json.NewEncoder(f)
	}

	enc.Encode(files)

	if *verbose {
		for _, v := range files {
			log.Printf("%+v", v)
		}
	}
}
