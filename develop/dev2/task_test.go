package main

import "testing"

func TestSeq(t *testing.T) {
	tests := []struct {
		test           string
		expectedResult string
	}{
		{"a4", "aaaa"},
		{"", ""},
		{"a2b2c2", "aabbcc"},
		{"aaa3bbb2", "aaaaabbbb"},
	}

	for _, test := range tests {
		extract, _ := Extract(test.test)

		if extract != test.expectedResult {
			t.Error("Expected result not equal: ", extract, "!=", test.expectedResult)
		}
	}
}
