package searcher

import (
	"testing"
)

type testCase struct {
	Arg  []string
	Want []string
}

var cases = []testCase{
	{[]string{}, []string{}},
	{[]string{"hello"}, []string{"hello"}},
	{[]string{"hello", "hello"}, []string{"hello"}},
	{[]string{"hello", "world", "hello", "hello"}, []string{"hello", "world"}},
	{[]string{"hello", "world", "world", "hello"}, []string{"hello", "world"}},
}

func TestSliceUnique(t *testing.T) {
	for _, tc := range cases {
		got := sliceUnique(tc.Arg)

		if len(got) != len(tc.Want) {
			t.Errorf("Wrong object received:\ngot:\t%#v\nwant:\t%#v", got, tc.Want)
			continue
		}

		for i, it := range tc.Want {
			if got[i] != it {
				t.Errorf("Objects don't match at index %d: got -> %s, want -> %s\n", i, got[i], it)
				t.Errorf("Wrong object received:\ngot:\t%#v\nwant:\t%#v", got, tc.Want)
				break
			}
		}
	}
}
