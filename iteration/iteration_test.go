package iteration

import (
	"fmt"
	"strings"
	"testing"
)

func TestRepeat(t *testing.T) {
	repeated := Repeat("a", 5)
	expected := "aaaaa"

	if repeated != expected {
		t.Errorf("expected %q, got %q", expected, repeated)
	}
}

func TestCompare(t *testing.T) {
	assertEqual := func(t *testing.T, expected, actual int) {
		if expected != actual {
			t.Errorf("actual: %d, but got %d", actual, expected)
		}
	}
	assertEqual(t, -1, strings.Compare("a", "b"))
	assertEqual(t, 0, strings.Compare("a", "a"))
	assertEqual(t, 1, strings.Compare("b", "a"))
}

func ExampleRepeat() {
	repeated := Repeat("a", 5)
	fmt.Println(repeated)
	// output: aaaaa
}

func BenchmarkRepeat(b *testing.B) {
	for i:=0;i<b.N;i++ {
		Repeat("a", 5)
	}
}