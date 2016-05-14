// This client is for indexing a folder heirarchy and storing it on a remote
// server.
package main

import (
	"compress/gzip"
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var files = []EnkiFile{}
var db sql.DB

// EnkiFile is the data structure we use for storing data into "the cloud".
// This should get translated into whatever format we send over the wire.
type EnkiFile struct {
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	Path    string    `json:"path"`
	Host    string    `json:"host"`
}

// Our walk function. It takes in a path and file info and stores that data
// locally in a sqlite DB.
//
// The path argument contains the argument to Walk as a prefix; that is, if
// Walk is called with "dir", which is a directory containing the file "a", the
// walk function will be called with argument "dir/a". The info argument is the
// os.FileInfo for the named path.
//
// If there was a problem walking to the file or directory named by path, the
// incoming error will describe the problem and Walk will not descend into that
// directory. We just log the error and continue.
func walkFunction(path string, info os.FileInfo, err error) error {
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
	filename := flag.String("file", "output.gob", "Where to output gob")
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		log.Printf("%+v", args)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Open up database!
	// https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dir := args[0]
	err = filepath.Walk(dir, walkFunction)
	if err != nil {
		log.Fatal(err)
	}

	var enc *gob.Encoder
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
		enc = gob.NewEncoder(f)
	} else {
		f, err := os.Create(*filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		enc = gob.NewEncoder(f)
	}

	f, err := os.Create(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	enc = gob.NewEncoder(f)
	enc.Encode(files)

	if *verbose {
		for _, v := range files {
			log.Printf("%+v", v)
		}
	}
}
