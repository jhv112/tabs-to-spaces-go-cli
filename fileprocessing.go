package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
			return recurseThroughDirs(filepath.Join(dirPath, entry.Name()), filter, outChan)
		}

		if filter.MatchString(entry.Name()) {
			outChan <- entry.Name()
		}

		// else ignore
	}

	return nil
}

func readFileContents(file *os.File) (string, error) {
	// from: https://stackoverflow.com/questions/48596338/a/58508195
	var buf strings.Builder

	_, err := io.Copy(&buf, file)

	if err != nil {
		return "", err
	}

	// assuming utf8
	return buf.String(), nil
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
	stopSignal <-chan struct{},
) {
	for {
		_, stop := <-stopSignal

		if stop {
			break
		}

		for {
			filepath, ok := <-filepathChannel

			if !ok {
				break
			}

			if err := processFile(filepath, fileModifier); err != nil {
				// ???
				log.Panic(err)
			}
		}
	}
}
