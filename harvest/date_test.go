package harvest

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDate_UnmarshalJSONParse(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want Date
		err  error
	}{
		{
			name: "2nd of January",
			args: args{"2006-01-02"},
			want: Date{Time: time.Date(2006, time.January, 2, 0, 0, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "31th of December",
			args: args{"2006-12-31"},
			want: Date{Time: time.Date(2006, time.December, 31, 0, 0, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "With quotes",
			args: args{"\"2006-04-05\""},
			want: Date{Time: time.Date(2006, time.April, 5, 0, 0, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "empty",
			args: args{""},
			err:  DateParseError,
		},
		{
			name: "Totally invalid",
			args: args{"gibberish"},
			err:  DateParseError,
		},
	}
	for _, tt := range tests {
		tm := Date{}
		if gotErr := tm.UnmarshalJSON([]byte(tt.args.str)); gotErr != tt.err {
			t.Errorf("%q. UnmarshalJSON() error = %v, want %v", tt.name, gotErr, tt.err)
		} else if !tm.Equal(tt.want) {
			t.Errorf("%q. UnmarshalJSON() = %v, want %v", tt.name, tm, tt.want)
		}
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	type foo struct {
		ID   *int64 `json:"id"`
		Date *Date  `json:"date"`
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
				ID:   Int64(123),
				Date: &Date{Time: time.Date(2019, time.January, 2, 0, 0, 0, 0, time.Local)},
			},
			err: nil,
		},
		{
			name: "null time",
			args: args{`{"id": 123, "date": null}`},
			want: foo{
				ID:   Int64(123),
				Date: nil,
			},
			err: nil,
		},
		{
			name: "Totally invalid",
			args: args{`{"id": 123, "date": "gibberish"}`},
			err:  DateParseError,
		},
	}

	for _, tt := range tests {
		var f foo
		if gotErr := json.Unmarshal([]byte(tt.args.jsonStr), &f); gotErr != tt.err {
			t.Errorf("%q. UnmarshalJSON() error = %v, want %v", tt.name, gotErr, tt.err)
		} else if tt.err == nil {
			if f.ID == nil || *f.ID != *tt.want.ID {
				t.Errorf("%q. UnmarshalJSON() = %v, want %v - ID messed up", tt.name, f, tt.want)
			} else if tt.want.Date == nil && f.Date != nil {
				t.Errorf("%q. UnmarshalJSON() = %v, want %v - unexpected time", tt.name, f, tt.want)
			} else if tt.want.Date != nil && (f.Date == nil || !(*tt.want.Date).Equal(*f.Date)) {
				t.Errorf("%q. UnmarshalJSON() = %v, want %v", tt.name, f, tt.want)
			}
		}
	}
}
