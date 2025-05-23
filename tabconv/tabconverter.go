package tabconv

import (
	"fmt"
	"regexp"
	"strings"
)

// Standard implementation for 4 tabs
//
// scan for tabs
// while tabs exist:
//
//	/(^(?:[^\t]{4})*[^\t]{0})\t/ug → "$1    "
//	/(^(?:[^\t]{4})*[^\t]{1})\t/ug → "$1   "
//	/(^(?:[^\t]{4})*[^\t]{2})\t/ug → "$1  "
//	/(^(?:[^\t]{4})*[^\t]{3})\t/ug → "$1 "
type TabConverter struct {
	matcher *regexp.Regexp
	replace string
}

func NewTabConverters(tabsize int) []TabConverter {
	// special case
	if tabsize == 0 {
		return []TabConverter{{regexp.MustCompile(`\t+`), ""}}
	}

	// special case
	if tabsize == 1 {
		return []TabConverter{{regexp.MustCompile(`\t`), " "}}
	}

	// assures that tabsize is nonnegative
	converters := make([]TabConverter, tabsize)

	for i := 0; i < tabsize; i++ {
		converters[i].matcher = regexp.MustCompile(fmt.Sprintf(`((?:\r\n|\r|\n|^)(?:[^\t\r\n]{%d})*[^\t\r\n]{%d})\t`, tabsize, i))
		// assurance from: https://stackoverflow.com/questions/43586091/a/43586154
		converters[i].replace = "$1" + strings.Repeat(" ", tabsize-i)
	}

	return converters
}

// While tabs exist in text, modify text in order of tabConverters occurence
func ConvertTabsIn(text string, tabConverters []TabConverter) string {
	for strings.ContainsRune(text, '\t') {
		for _, converter := range tabConverters {
			text = converter.matcher.ReplaceAllString(text, converter.replace)
		}
	}

	return text
}
