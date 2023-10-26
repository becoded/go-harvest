package harvest_test

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/becoded/go-harvest/harvest"
	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
)

func TestDate_String(t *testing.T) {
	loc, err := time.LoadLocation("Local")
	assert.NoError(t, err)

	time.UTC = loc

	tests := []struct {
		name  string
		input harvest.Date
		want  string
	}{
		{
			name:  "2nd of January",
			input: harvest.Date{Time: time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)},
			want:  "2006-01-02",
		},
		{
			name:  "31th of December",
			input: harvest.Date{Time: time.Date(2006, time.December, 31, 0, 0, 0, 0, time.UTC)},
			want:  "2006-12-31",
		},
		{
			name:  "6th of May",
			input: harvest.Date{Time: time.Date(2016, time.May, 6, 0, 0, 0, 0, time.UTC)},
			want:  "2016-05-06",
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.input.String())
		})
	}
}

func TestDate_UnmarshalJSONParse(t *testing.T) {
	loc, err := time.LoadLocation("Local")
	assert.NoError(t, err)

	time.UTC = loc

	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want harvest.Date
		err  error
	}{
		{
			name: "2nd of January",
			args: args{"2006-01-02"},
			want: harvest.Date{Time: time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)},
			err:  nil,
		},
		{
			name: "31th of December",
			args: args{"2006-12-31"},
			want: harvest.Date{Time: time.Date(2006, time.December, 31, 0, 0, 0, 0, time.UTC)},
			err:  nil,
		},
		{
			name: "With quotes",
			args: args{"\"2006-04-05\""},
			want: harvest.Date{Time: time.Date(2006, time.April, 5, 0, 0, 0, 0, time.UTC)},
			err:  nil,
		},
		{
			name: "empty",
			args: args{""},
			err:  harvest.ErrDateParse,
		},
		{
			name: "Totally invalid",
			args: args{"gibberish"},
			err:  harvest.ErrDateParse,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tm := harvest.Date{}

			gotErr := tm.UnmarshalJSON([]byte(tt.args.str))
			if tt.err != nil {
				assert.EqualError(t, gotErr, tt.err.Error())

				return
			}
			assert.NoError(t, gotErr)

			if !tm.Equal(tt.want) {
				t.Errorf("%q. UnmarshalJSON() = %v, response %v", tt.name, tm, tt.want)
			}
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	type foo struct {
		ID   *int64        `json:"id"`
		Date *harvest.Date `json:"date"`
	}

	type args struct {
		jsonStr string
	}

	tests := []struct {
		name string
		args args
		want foo
		err  error
	}{
		{
			name: "Typical json",
			args: args{`{"id": 123, "date": "2019-01-02"}`},
			want: foo{
				ID:   harvest.Int64(123),
				Date: &harvest.Date{Time: time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC)},
			},
			err: nil,
		},
		{
			name: "null time",
			args: args{`{"id": 123, "date": null}`},
			want: foo{
				ID:   harvest.Int64(123),
				Date: nil,
			},
			err: nil,
		},
		{
			name: "Totally invalid",
			args: args{`{"id": 123, "date": "gibberish"}`},
			err:  harvest.ErrDateParse,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var f foo

			gotErr := json.Unmarshal([]byte(tt.args.jsonStr), &f)
			if tt.err != nil {
				assert.EqualError(t, gotErr, tt.err.Error())

				return
			}
			assert.NoError(t, gotErr)

			switch {
			case f.ID == nil || *f.ID != *tt.want.ID:
				t.Errorf("%q. UnmarshalJSON() = %v, response %v - ID messed up", tt.name, f, tt.want)
			case tt.want.Date == nil && f.Date != nil:
				t.Errorf("%q. UnmarshalJSON() = %v, response %v - unexpected time", tt.name, f, tt.want)
			case tt.want.Date != nil && (f.Date == nil || !tt.want.Date.Equal(*f.Date)):
				t.Errorf("%q. UnmarshalJSON() = %v, response %v", tt.name, f, tt.want)
			}
		})
	}
}

func TestDate_EncodeValues(t *testing.T) {
	type foo struct {
		Query *string       `url:"query,omitempty"`
		Date  *harvest.Date `url:"date,omitempty"`
	}

	tests := []struct {
		name string
		args *foo
		want url.Values
	}{
		{
			name: "All fields filled in",
			args: &foo{
				Query: harvest.String("foo"),
				Date:  &harvest.Date{Time: time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC)},
			},
			want: url.Values{
				"query": []string{"foo"},
				"date":  []string{"2019-01-02"},
			},
		},
		{
			name: "No query",
			args: &foo{
				Date: &harvest.Date{Time: time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC)},
			},
			want: url.Values{
				"date": []string{"2019-01-02"},
			},
		},
		{
			name: "No date",
			args: &foo{
				Query: harvest.String("foo"),
			},
			want: url.Values{
				"query": []string{"foo"},
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			qs, err := query.Values(tt.args)
			if err != nil {
				t.Errorf("%q. EncodeValues() error = %v", tt.name, err)
			}
			if !reflect.DeepEqual(qs, tt.want) {
				t.Errorf("%q. EncodeValues() = %v, response %v", tt.name, qs, tt.want)
			}
		})
	}
}
