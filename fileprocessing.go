package main

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func recurseThroughDirs(
	dirPath string,
	filter *regexp.Regexp,
	outChan chan<- string,
) error {
	entries, err := os.ReadDir(dirPath)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err = recurseThroughDirs(filepath.Join(dirPath, entry.Name()), filter, outChan)

			if err != nil {
				return err
			}
		} else if filter.MatchString(entry.Name()) {
			outChan <- filepath.Join(dirPath, entry.Name())
		}

		// else ignore
	}

	return nil
}

// Read open file from current seek position to end
func readFileContents(file *os.File) (string, error) {
	var buffer bytes.Buffer

	_, err := io.Copy(&buffer, file)

	if err != nil && err != io.EOF {
		return "", err
	}

	return buffer.String(), nil
}

func processFile(filepath string, fileModifier func(string) string) error {
	// open for reading and writing
	file, err := os.OpenFile(filepath, os.O_RDWR, fs.ModeExclusive)

	if err != nil {
		return err
	}

	defer file.Close()

	fileContents, err := readFileContents(file)

	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)

	if err != nil {
		return err
	}

	modifiedFileContents := fileModifier(fileContents)

	err = file.Truncate(0)

	if err != nil {
		return err
	}

	_, err = file.WriteString(modifiedFileContents)

	if err != nil {
		return err
	}

	return nil
}

func processFiles(
	fileModifier func(input string) string,
	filepathChannel <-chan string,
) {
	for filepath := range filepathChannel {
		if err := processFile(filepath, fileModifier); err != nil {
			// ???
			log.Panic(err)
		}
	}
}
