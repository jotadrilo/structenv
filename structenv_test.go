package structenv_test

import (
	"flag"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jotadrilo/structenv"
)

// Nested is an object for testing purposes
type Nested struct {
	Bool1     bool          `flag:"nested-bool1" env:"NESTED_BOOL1"`
	Duration1 time.Duration `flag:"nested-duration1" env:"NESTED_DURATION1"`
	Float1    float64       `flag:"nested-float1" env:"NESTED_FLOAT1"`
	Float2    float32       `flag:"nested-float2" env:"NESTED_FLOAT2"`
	Int1      int           `flag:"nested-int1" env:"NESTED_INT1"`
	Int2      int64         `flag:"nested-int2" env:"NESTED_INT2"`
	String1   string        `flag:"nested-string1" env:"NESTED_STRING1"`
	Uint1     uint          `flag:"nested-uint1" env:"NESTED_UINT1"`
	Uint2     uint64        `flag:"nested-uint2" env:"NESTED_UINT2"`
}

// Test is an object for testing purposes
type Test struct {
	Bool1     bool          `flag:"bool1" env:"BOOL1"`
	Duration1 time.Duration `flag:"duration1" env:"DURATION1"`
	Float1    float64       `flag:"float1" env:"FLOAT1"`
	Float2    float32       `flag:"float2" env:"FLOAT2"`
	Int1      int           `flag:"int1" env:"INT1"`
	Int2      int64         `flag:"int2" env:"INT2"`
	String1   string        `flag:"string1" env:"STRING1"`
	Uint1     uint          `flag:"uint1" env:"UINT1"`
	Uint2     uint64        `flag:"uint2" env:"UINT2"`

	Nested *Nested
}

var (
	_ = flag.Bool("bool1", false, "")
	_ = flag.Duration("duration1", 0, "")
	_ = flag.Float64("float1", 0, "")
	_ = flag.Float64("float2", 0, "")
	_ = flag.Int("int1", 0, "")
	_ = flag.Int64("int2", 0, "")
	_ = flag.String("string1", "s", "")
	_ = flag.Uint("uint1", 0, "")
	_ = flag.Uint64("uint2", 0, "")

	_ = flag.Bool("nested-bool1", false, "")
	_ = flag.Duration("nested-duration1", 0, "")
	_ = flag.Float64("nested-float1", 0, "")
	_ = flag.Float64("nested-float2", 0, "")
	_ = flag.Int("nested-int1", 0, "")
	_ = flag.Int64("nested-int2", 0, "")
	_ = flag.String("nested-string1", "s", "")
	_ = flag.Uint("nested-uint1", 0, "")
	_ = flag.Uint64("nested-uint2", 0, "")
)

func TestParseEnv(t *testing.T) {
	testCases := []struct {
		desc string
		env  []string
		want *Test
	}{
		{
			desc: "parses bool env",
			env:  []string{"BOOL1=1"},
			want: &Test{Bool1: true},
		},
		// {
		// 	desc: "parses nested bool env",
		// 	env:  []string{"BOOL1=1", "NESTED_BOOL1=1"},
		// 	want: &Test{Bool1: true, Nested: &Nested{Bool1: true}},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := &Test{}
			if err := structenv.ParseEnv(tc.env, got); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("want vs got diff:\n %+v", diff)
			}
		})
	}
}
