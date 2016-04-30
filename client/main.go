package main

import (
	"crypto/sha512"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var files = make(map[[sha512.Size]byte]string)

func checkDuplicate(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}
	if info.IsDir() {
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Print(err)
		return nil
	}
	digest := sha512.Sum512(data)
	if v, ok := files[digest]; ok {
		fmt.Printf("%q is a duplicate of %q\n", path, v)
	} else {
		files[digest] = path
	}

	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	dir := os.Args[1]
	err := filepath.Walk(dir, checkDuplicate)
	if err != nil {
		log.Fatal(err)
	}
}
