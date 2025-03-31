package tabconv

import "testing"

type expected struct{ regex, replace string }

func TestNewTabConverters(t *testing.T) {
	testCases := []struct {
		tabsize  int
		expected []expected
	}{
		// just remove all tabs
		{0, []expected{{`\t`, "$1"}}},
		{1, []expected{
			{`(^(?:[^\t]{1})*[^\t]{0})\t`, "$1 "},
		}},
		{4, []expected{
			{`(^(?:[^\t]{4})*[^\t]{0})\t`, "$1    "},
			{`(^(?:[^\t]{4})*[^\t]{1})\t`, "$1   "},
			{`(^(?:[^\t]{4})*[^\t]{2})\t`, "$1  "},
			{`(^(?:[^\t]{4})*[^\t]{3})\t`, "$1 "},
		}},
	}

	for _, testCase := range testCases {
		observed := NewTabConverters(testCase.tabsize)

		if len(observed) != len(testCase.expected) {
			t.Fatalf("want %d converters, have %d", len(testCase.expected), len(observed))
		}

		for i := range observed {
			if observed[i].Matcher.String() != testCase.expected[i].regex {
				t.Fatalf("want /%s/, have /%s/", testCase.expected[i].regex, observed[i].Matcher.String())
			}

			if observed[i].Replace != testCase.expected[i].replace {
				t.Fatalf("want /%s/, have /%s/", testCase.expected[i].regex, observed[i].Matcher.String())
			}
		}
	}
}
