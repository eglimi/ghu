/*
Go Header Update (ghu) replaces the beginning of a file with the content of
another file.

This is useful e.g. for replacing the information in source code files.

Note that the file with the header content is expected to have no trailing
newlines (other than required ones). Beware of vim eol feature.
*/
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var hdrContent []byte

var path string
var hfile string
var ftype string

func main() {

	// Command line parsing
	flag.StringVar(&path, "path", "", "The path to process.")
	flag.StringVar(&hfile, "hfile", "", "The file with the header content.")
	flag.StringVar(&ftype, "ftype", "", "File type pattern (suffix).")
	flag.Parse()
	if path == "" {
		flag.PrintDefaults()
		log.Fatal("path argument not specified")
	}
	if hfile == "" {
		flag.PrintDefaults()
		log.Fatal("hfile argument not specified")
	}
	if ftype == "" {
		flag.PrintDefaults()
		log.Fatal("ftype argument not specified")
	}

	log.Printf("Starting to replace headers in %v", path)

	// Read header text
	hdr, err := ioutil.ReadFile(hfile)
	if err != nil {
		log.Fatalf("Could not read header content. %v", err)
	}
	hdrContent = hdr // ?

	// Walk path
	err = filepath.Walk(path, replaceHeader)
	if err != nil {
		log.Printf("Error occurred: %v", err)
	}

	log.Print("Job completed, Sir!")
}

func replaceHeader(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Fatal(err)
	}

	if info.IsDir() {
		return nil
	}

	if strings.HasSuffix(info.Name(), ftype) == false {
		return nil
	}

	log.Printf("Processing %s", path)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failed to process %v.", path)
		return err
	}

	start := true
	index := 0
	for count, c := range content {
		log.Print(c)
		if start == true {
			// Search for start
			if unicode.IsSpace(rune(c)) {
				continue // ignore
			}
			if len(content) >= count+1 && content[count] == '/' && content[count+1] == '*' {
				start = false
				continue
			}
			// No comment, use all text. Make sure that we have a newline after the header!
			hdrContent = append(hdrContent, byte('\n'))
			break
		} else {
			// Search for end	
			if len(content) >= count+2 {
				if content[count] == '*' && content[count+1] == '/' {
					index = count + 2
					break
				}
			} else {
				log.Print("Error in file. No end of format `*/` found")
				return nil // continue with other files
			}
		}
	}

	// Create results file
	resLength := len(hdrContent) + len(content) - index
	resFile := make([]byte, resLength)

	// Copy header
	copy(resFile, hdrContent)

	// Copy content
	copy(resFile[len(hdrContent):], content[index:])

	ioutil.WriteFile(path, resFile, info.Mode())

	return nil
}
