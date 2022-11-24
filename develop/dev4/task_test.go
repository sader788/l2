package main

import (
	"reflect"
	"testing"
)

func TestAnagram(t *testing.T) {
	tests := []struct {
		test           []string
		expectedResult map[string][]string
	}{
		{
			[]string{`пятак`, `пятка`, `тяпка`, `листок`, `слиток`, `столик`, `рофлинка`, `калинроф`},
			map[string][]string{`калинроф`: {`калинроф`, `рофлинка`}, `листок`: {`листок`, `слиток`, `столик`}, `пятак`: {`пятак`, `пятка`, `тяпка`}},
		},
	}

	for _, test := range tests {
		anagramMap := NewAnagramMap(test.test)

		if !reflect.DeepEqual(*anagramMap, test.expectedResult) {
			t.Error("result for test: ", test.test, " is not correct")
		}

	}
}
