package store

import (
	"reflect"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	var tests = []struct {
		s map[int]string
	}{
		{
			map[int]string{
				123: "CS",
				122: "Witcher 3",
			},
		},
	}

	for caseid, c := range tests {
		var obj map[int]string
		bin, err := Encode(c.s)
		if err != nil {
			t.Errorf("case #%d, encode err: %v", caseid+1, err)
		}
		if err := Decode(bin, &obj); err != nil {
			t.Errorf("case #%d, decode err: %v", caseid+1, err)
		}
		if !reflect.DeepEqual(obj, c.s) {
			t.Errorf("case #%d, got: %v, expected: %v", caseid+1, obj, c.s)
		}
		t.Logf("Result: %v", bin)
	}
}
