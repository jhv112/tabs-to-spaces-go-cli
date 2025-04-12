package tabconv

import "testing"

type expected struct{ regex, replace string }

func TestNewTabConverters(t *testing.T) {
	testCases := []struct {
		tabsize  int
		expected []expected
	}{
		// just remove all tabs
		{0, []expected{{`\t+`, ""}}},
		// replace each tab one-to-one with a space
		{1, []expected{{`\t`, " "}}},
		{2, []expected{
			{`((?:\r\n|\r|\n|^)(?:[^\t\r\n]{2})*[^\t\r\n]{0})\t`, "$1  "},
			{`((?:\r\n|\r|\n|^)(?:[^\t\r\n]{2})*[^\t\r\n]{1})\t`, "$1 "},
		}},
		{4, []expected{
			{`((?:\r\n|\r|\n|^)(?:[^\t\r\n]{4})*[^\t\r\n]{0})\t`, "$1    "},
			{`((?:\r\n|\r|\n|^)(?:[^\t\r\n]{4})*[^\t\r\n]{1})\t`, "$1   "},
			{`((?:\r\n|\r|\n|^)(?:[^\t\r\n]{4})*[^\t\r\n]{2})\t`, "$1  "},
			{`((?:\r\n|\r|\n|^)(?:[^\t\r\n]{4})*[^\t\r\n]{3})\t`, "$1 "},
		}},
	}

	for _, testCase := range testCases {
		observed := NewTabConverters(testCase.tabsize)

		if len(observed) != len(testCase.expected) {
			t.Fatalf("want %d converters, have %d", len(testCase.expected), len(observed))
		}

		for i := range observed {
			if observed[i].matcher.String() != testCase.expected[i].regex {
				t.Fatalf("want /%s/, have /%s/", testCase.expected[i].regex, observed[i].matcher.String())
			}

			if observed[i].replace != testCase.expected[i].replace {
				t.Fatalf("want /%s/, have /%s/", testCase.expected[i].regex, observed[i].matcher.String())
			}
		}
	}
}

func TestNewTabConvertersFailOnNegativeTabsize(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fail()
		}
	}()

	NewTabConverters(-1)
}

func newConverterRange(maxTabSize int) [][]TabConverter {
	converters := make([][]TabConverter, maxTabSize+1)

	for i := 0; i <= maxTabSize; i++ {
		converters[i] = NewTabConverters(i)
	}

	return converters
}

func TestConvertTabsInText(t *testing.T) {
	testCases := []struct {
		tabsize           int
		initial, expected string
	}{
		{4, `/*	Text 1`, `/*  Text 1`},
		{4, `
/*	Text 1`, `
/*  Text 1`},
		{4, `
/*	Text 1
**	Text 2
*/	Text 3`, `
/*  Text 1
**  Text 2
*/  Text 3`},
		{4, `	char const *field;		// Comment`, `    char const *field;      // Comment`},
		{3, `	char const *field;		// Comment`, `   char const *field;      // Comment`},
		{2, `	char const *field;		// Comment`, `  char const *field;    // Comment`},
		{1, `	char const *field;		// Comment`, ` char const *field;  // Comment`},
		{0, `	char const *field;		// Comment`, `char const *field;// Comment`},
	}

	// preinit converters
	converters := newConverterRange(4)

	for _, testCase := range testCases {
		observed := ConvertTabsIn(testCase.initial, converters[testCase.tabsize])

		if observed != testCase.expected {
			t.Fatalf("want \"%s\", have \"%s\"", testCase.expected, observed)
		}
	}
}
