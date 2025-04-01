package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"testing"
)

/*
test dir structure:

	                root
	                  |
	 +---+---+---+----+----------------+
	1.a 2.a 3.b  B                     C
	             |                     |
	         +---+---+    +---+---+---++-----------+
	        1.b 2.b 3.c  1.c 2.c 3.d  D            E
	                                  |            |
	                              +---+---+    +---+---+
	                             1.d 2.d 3.e  1.e 2.e 3.f
*/
func setupTestFileSystem() (string, error) {
	rootDirPath, err := os.MkdirTemp("", "")

	if err != nil {
		return "", err
	}

	subdirs := [4]string{
		path.Join(rootDirPath, "B"),
		path.Join(rootDirPath, "C"),
		path.Join(rootDirPath, "C", "D"),
		path.Join(rootDirPath, "C", "E"),
	}

	possibleErrs := [19]error{}

	for i := range subdirs {
		possibleErrs[i] = os.MkdirAll(subdirs[i], 0777)
	}

	possibleErrs[4] = os.WriteFile(path.Join(rootDirPath, "1.a"), []byte{}, 0777)
	possibleErrs[5] = os.WriteFile(path.Join(rootDirPath, "2.a"), []byte{}, 0777)
	possibleErrs[6] = os.WriteFile(path.Join(rootDirPath, "3.b"), []byte{}, 0777)

	for i := range subdirs {
		possibleErrs[3*i+7] = os.WriteFile(path.Join(subdirs[i], fmt.Sprintf("1.%c", 'b'+i)), []byte{}, 0777)
		possibleErrs[3*i+8] = os.WriteFile(path.Join(subdirs[i], fmt.Sprintf("2.%c", 'b'+i)), []byte{}, 0777)
		possibleErrs[3*i+9] = os.WriteFile(path.Join(subdirs[i], fmt.Sprintf("3.%c", 'b'+i+1)), []byte{}, 0777)
	}

	for _, possibleErr := range possibleErrs {
		if possibleErr != nil {
			return "", possibleErr
		}
	}

	return rootDirPath, nil
}

func setupExpectedFilePaths(rootDirPath string) []string {
	expectedFiles := []string{
		path.Join(rootDirPath, "1.a"),
		path.Join(rootDirPath, "2.a"),
		path.Join(rootDirPath, "3.b"),
		path.Join(rootDirPath, "B", "1.b"),
		path.Join(rootDirPath, "B", "2.b"),
		path.Join(rootDirPath, "B", "3.c"),
		path.Join(rootDirPath, "C", "1.c"),
		path.Join(rootDirPath, "C", "2.c"),
		path.Join(rootDirPath, "C", "3.d"),
		path.Join(rootDirPath, "C", "D", "1.d"),
		path.Join(rootDirPath, "C", "D", "2.d"),
		path.Join(rootDirPath, "C", "D", "3.e"),
		path.Join(rootDirPath, "C", "E", "1.e"),
		path.Join(rootDirPath, "C", "E", "2.e"),
		path.Join(rootDirPath, "C", "E", "3.f"),
	}

	// windows exception
	if os.PathSeparator == '\\' {
		for i := range expectedFiles {
			expectedFiles[i] = strings.ReplaceAll(expectedFiles[i], "/", `\`)
		}
	}

	return expectedFiles
}

func tearDownTestFileSystem(rootDirPath string) error {
	return os.RemoveAll(rootDirPath)
}

func TestRecurseThroughDirs(t *testing.T) {
	rootDirPath, err := setupTestFileSystem()

	defer func() {
		err = tearDownTestFileSystem(rootDirPath)

		if err != nil {
			t.Fatal(err)
		}
	}()

	if err != nil {
		t.Fatal(err)
	}

	expectedFiles := setupExpectedFilePaths(rootDirPath)

	fileChan := make(chan string)

	// because of unbuffered channel use has to run asynchronously
	go func() {
		recurseThroughDirs(rootDirPath, regexp.MustCompile(`\d.[a-f]$`), fileChan)

		close(fileChan)
	}()

	observedFiles := make([]string, 0, 15)

	for f := range fileChan {
		observedFiles = append(observedFiles, f)
	}

	if !slices.Equal(observedFiles, expectedFiles) {
		t.Fatalf("want %s files, have %s files", strings.Join(expectedFiles, ";"), strings.Join(observedFiles, ";"))
	}
}
