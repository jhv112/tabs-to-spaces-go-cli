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
	Matcher *regexp.Regexp
	Replace string
}

func NewTabConverters(tabsize int) []TabConverter {
	// special case
	if tabsize == 0 {
		return []TabConverter{{regexp.MustCompile(`\t`), "$1"}}
	}

	// assures that tabsize is nonnegative
	converters := make([]TabConverter, tabsize)

	for i := 0; i < tabsize; i++ {
		converters[i].Matcher = regexp.MustCompile(fmt.Sprintf(`(^(?:[^\t]{%d})*[^\t]{%d})\t`, tabsize, i))
		// assurance from: https://stackoverflow.com/questions/43586091/a/43586154
		converters[i].Replace = "$1" + strings.Repeat(" ", tabsize-i)
	}

	return converters
}

// While tabs exist in text, modify text in order of tabConverters occurence
func ConvertTabsIn(text string, tabConverters []TabConverter) string {
	for strings.ContainsRune(text, '\t') {
		for _, converter := range tabConverters {
			text = converter.Matcher.ReplaceAllString(text, converter.Replace)
		}
	}

	return text
}
