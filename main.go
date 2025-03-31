package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/jhv112/tabs-to-spaces-go-cli/tabconv"
)

const usage = `

Description:

    Command-line utility that replaces tabs in certain files with spaces

Command-line arguments:

    | flag       | description        | implementation | default |
    |:-----------|:-------------------|:--------------:|:-------:|
    | --filter   | file filter        |  regex filter  |    .*   |
    | --startdir | start directory    |  relative dir  |    .    |
    | --tabsize  | tab size in spaces |  range [0,127] |    4    |

`

type Args struct {
	filter   *regexp.Regexp
	startdir string
	tabsize  int8
}

// I know: Not the go way, to do it.
func parseArgs(argsStrings []string) (args Args, err error) {
	// defaults
	args = Args{regexp.MustCompile(".*"), ".", 4}

	for _, arg := range argsStrings {
		if strings.HasPrefix(arg, "--filter=") {
			args.filter, err = regexp.Compile(arg[len("--filter="):])

			if err != nil {
				return Args{}, err
			}
		}

		if strings.HasPrefix(arg, "--startdir=") {
			// from: https://stackoverflow.com/questions/45941821/a/64499397
			args.startdir, err = filepath.Abs(arg[len("--startdir="):])

			if err != nil {
				return Args{}, err
			}
		}

		if strings.HasPrefix(arg, "--tabsize=") {
			parsedUint, err := strconv.ParseUint(arg[len("--tabsize="):], 10, 7)

			if err != nil {
				return Args{}, err
			}

			// guaranteed to fit and be nonnegative
			args.tabsize = int8(parsedUint)
		}

		if slices.Contains([]string{"-h", "--help"}, arg) {
			return Args{}, fmt.Errorf("help")
		}
	}

	return
}

func main() {
	args, err := parseArgs(os.Args)

	if err != nil {
		log.Fatal(err.Error() + usage)
	}

	tabConverters := tabconv.NewTabConverters(int(args.tabsize))
	filepathChannel := make(chan string)
	stopSignal := make(chan struct{})

	defer func() {
		stopSignal <- struct{}{}
	}()

	// from: https://stackoverflow.com/questions/24073697/a/24073875
	processorCount := runtime.NumCPU()

	fileModifier := func(input string) string {
		return tabconv.ConvertTabsIn(input, tabConverters)
	}

	for i := 0; i < processorCount; i++ {
		go processFiles(fileModifier, filepathChannel, stopSignal)
	}

	err = recurseThroughDirs(args.startdir, args.filter, filepathChannel)

	if err != nil {
		log.Fatal(err)
	}
}
