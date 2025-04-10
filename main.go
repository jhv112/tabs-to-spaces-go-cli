package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jhv112/tabs-to-spaces-go-cli/tabconv"
)

//go:embed USAGE
var usage string

type Args struct {
	filter   *regexp.Regexp
	startdir string
	tabsize  int8
}

// I know: Not the go way, to do it.
func parseArgs(argsStrings []string) (args Args, err error) {
	// defaults
	args = Args{regexp.MustCompile(".*"), ".", 4}

	for _, arg := range argsStrings[1:] {
		if strings.HasPrefix(arg, "--filter=") {
			args.filter, err = regexp.Compile(arg[len("--filter="):])

			if err != nil {
				return Args{}, err
			}
		} else if strings.HasPrefix(arg, "--startdir=") {
			// from: https://stackoverflow.com/questions/45941821/a/64499397
			args.startdir, err = filepath.Abs(arg[len("--startdir="):])

			if err != nil {
				return Args{}, err
			}
		} else if strings.HasPrefix(arg, "--tabsize=") {
			parsedUint, err := strconv.ParseUint(arg[len("--tabsize="):], 10, 7)

			if err != nil {
				return Args{}, err
			}

			// guaranteed to fit and be nonnegative
			args.tabsize = int8(parsedUint)
		} else if slices.Contains([]string{"-h", "--help"}, arg) {
			return Args{}, fmt.Errorf("help")
		} else {
			return Args{}, fmt.Errorf("unrecognised argument: \"%s\"", arg)
		}
	}

	return
}

func convertTabsToSpacesInDirectory(
	filter *regexp.Regexp,
	startdir string,
	tabsize int,
) error {
	tabConverters := tabconv.NewTabConverters(tabsize)

	fileModifier := func(input string) string {
		return tabconv.ConvertTabsIn(input, tabConverters)
	}

	return produceSyncConsumeAsync(
		func(ch chan<- string) error {
			return recurseThroughDirs(startdir, filter, ch)
		},
		func(ch <-chan string) {
			processFiles(fileModifier, ch)
		},
	)
}

func main() {
	args, err := parseArgs(os.Args)

	if err != nil {
		log.Fatal(err.Error() + usage)
	}

	err = convertTabsToSpacesInDirectory(
		args.filter,
		args.startdir,
		int(args.tabsize),
	)

	if err != nil {
		log.Fatal(err)
	}
}
