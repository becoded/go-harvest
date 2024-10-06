package harvest_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestStringify(t *testing.T) {
	var nilPointer *string

	tests := []struct {
		name string
		in   interface{}
		out  string
	}{
		// basic types
		{"case 1", "foo", `"foo"`},
		{"case 2", 123, `123`},
		{"case 3", 1.5, `1.5`},
		{"case 4", false, `false`},
		{
			"case 5",
			[]string{"a", "b"},
			`["a" "b"]`,
		},
		{
			"case 6",
			struct {
				A []string
			}{nil},
			// nil slice is skipped
			`{}`,
		},
		{
			"case 7",
			struct {
				A string
			}{"foo"},
			// structs not of a named type get no prefix
			`{A:"foo"}`,
		},

		// pointers
		{"case 8", nilPointer, `<nil>`},
		{"case 9", harvest.String("foo"), `"foo"`},
		{"case 10", harvest.Int(123), `123`},
		{"case 11", harvest.Bool(false), `false`},
		{
			"case 12",
			[]*string{harvest.String("a"), harvest.String("b")},
			`["a" "b"]`,
		},
		{
			"case 13",
			harvest.User{ID: harvest.Int64(123), FirstName: harvest.String("n")},
			`harvest.User{ID:123, FirstName:"n"}`,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := harvest.Stringify(tt.in)

			assert.Equal(t, tt.out, s)
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
		out  string
	}{
		{
			"happy",
			harvest.User{ID: harvest.Int64(1)},
			`harvest.User{ID:1}`,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := tt.in.(fmt.Stringer).String()

			assert.Equal(t, tt.out, s)
		})
	}
}
