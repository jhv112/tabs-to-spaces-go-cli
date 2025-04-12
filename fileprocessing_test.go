package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
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

func adaptPathsForOS(unixPaths []string) []string {
	if len(unixPaths) == 0 {
		return unixPaths
	}

	// windows exception
	if os.PathSeparator == '\\' {
		winPaths := make([]string, len(unixPaths))

		for i := range unixPaths {
			winPaths[i] = strings.ReplaceAll(unixPaths[i], "/", `\`)
		}

		return winPaths
	}

	return unixPaths
}

func setupExpectedFilePaths(rootDirPath string) []string {
	return adaptPathsForOS([]string{
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
	})
}

func tearDownTestFileSystem(rootDirPath string) error {
	return os.RemoveAll(rootDirPath)
}

func getAllFilesIn(dirPath string, filter *regexp.Regexp) []string {
	fileChan := make(chan string)

	// because of unbuffered channel use has to run asynchronously
	go func() {
		recurseThroughDirs(dirPath, filter, fileChan)

		close(fileChan)
	}()

	observedFiles := make([]string, 0, 15)

	for f := range fileChan {
		observedFiles = append(observedFiles, f)
	}

	return observedFiles
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
	observedFiles := getAllFilesIn(rootDirPath, regexp.MustCompile(`\d.[a-f]$`))

	if !slices.Equal(observedFiles, expectedFiles) {
		t.Fatalf("want %s files, have %s files", strings.Join(expectedFiles, ";"), strings.Join(observedFiles, ";"))
	}
}

func TestRecurseThroughDirsOnProjectDirectory(t *testing.T) {
	currentDirAbsolutePath, _ := filepath.Abs(".")
	expectedFiles := adaptPathsForOS([]string{
		path.Join(currentDirAbsolutePath, "asyncprocessing.go"),
		path.Join(currentDirAbsolutePath, "asyncprocessing_test.go"),
		path.Join(currentDirAbsolutePath, "fileprocessing.go"),
		path.Join(currentDirAbsolutePath, "fileprocessing_test.go"),
		path.Join(currentDirAbsolutePath, "main.go"),
		path.Join(currentDirAbsolutePath, "tabconv/tabconverter.go"),
		path.Join(currentDirAbsolutePath, "tabconv/tabconverter_test.go"),
	})

	observedFiles := getAllFilesIn(currentDirAbsolutePath, regexp.MustCompile(`\.go$`))

	if !slices.Equal(observedFiles, expectedFiles) {
		t.Fatalf("want %s files, have %s files", strings.Join(expectedFiles, ";"), strings.Join(observedFiles, ";"))
	}
}

func setupTempFileWithContent(content string) (*os.File, error) {
	file, err := os.CreateTemp(os.TempDir(), "")

	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(content)

	if err != nil {
		return nil, err
	}

	// to read from beginning and not from last append
	file.Seek(0, 0)

	return file, nil
}

const utf8String = "Hello, 世界"

func TestReadFileContents(t *testing.T) {
	file, err := setupTempFileWithContent(utf8String)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		// required; otherwise file cannot be removed
		file.Close()
		os.Remove(file.Name())
	}()

	observed, err := readFileContents(file)

	if err != nil {
		t.Fatalf("want no error, have %v", err)
	}

	if observed != utf8String {
		t.Fatalf("want \"%s\", have \"%s\"", utf8String, observed)
	}
}

func TestProcessFile(t *testing.T) {
	const initial = utf8String
	file, err := setupTempFileWithContent(initial)

	if err != nil {
		t.Fatal(err)
	}

	file.Close()

	defer os.Remove(file.Name())

	const expected = initial + " 2"

	err = processFile(file.Name(), func(_ string) string { return expected })

	if err != nil {
		t.Fatalf("want no error, have %v", err)
	}

	fileBytes, err := os.ReadFile(file.Name())

	if err != nil {
		t.Fatalf("want no error, have %v", err)
	}

	observed := string(fileBytes)

	if observed != expected {
		t.Fatalf("want \"%s\", have \"%s\"", expected, observed)
	}
}
