package harvest

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshalJSONParse(t *testing.T) {

	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want Time
		err  error
	}{
		{
			name: "24-hour clock after 12:00",
			args: args{"15:04"},
			want: Time{Time: time.Date(0, time.January, 1, 15, 4, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "24-hour clock before 12:00",
			args: args{"07:34"},
			want: Time{Time: time.Date(0, time.January, 1, 7, 34, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "24-hour clock no leading zero",
			args: args{"7:34"},
			want: Time{Time: time.Date(0, time.January, 1, 7, 34, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "Kitchen clock upper case",
			args: args{"9:13AM"},
			want: Time{Time: time.Date(0, time.January, 1, 9, 13, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "Kitchen clock in the morning",
			args: args{"8:13am"},
			want: Time{Time: time.Date(0, time.January, 1, 8, 13, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "Kitchen clock with leading zero",
			args: args{"09:39am"},
			want: Time{Time: time.Date(0, time.January, 1, 9, 39, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "Kitchen clock in the afternoon",
			args: args{"10:13pm"},
			want: Time{Time: time.Date(0, time.January, 1, 22, 13, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "With quotes",
			args: args{"\"10:13pm\""},
			want: Time{Time: time.Date(0, time.January, 1, 22, 13, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "empty",
			args: args{""},
			err:  TimeParseError,
		},
		{
			name: "24-hour with pm",
			args: args{"13:03pm"},
			err:  TimeParseError,
		},
		{
			name: "24-hour with am",
			args: args{"13:03am"},
			err:  TimeParseError,
		},
		{
			name: "Totally invalid",
			args: args{"gibberish"},
			err:  TimeParseError,
		},
	}
	for _, tt := range tests {
		tm := Time{}
		if gotErr := tm.UnmarshalJSON([]byte(tt.args.str)); gotErr != tt.err {
			t.Errorf("%q. UnmarshalJSON() error = %v, want %v", tt.name, gotErr, tt.err)
		} else if !tm.Equal(tt.want) {
			t.Errorf("%q. UnmarshalJSON() = %v, want %v", tt.name, tm, tt.want)
		}
	}
}

func intPtr(v int64) *int64 {
	return &v
}

func TestUnmarshalJSON(t *testing.T) {

	type foo struct {
		ID   *int64 `json:"id"`
		Time *Time  `json:"time"`
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
			args: args{`{"id": 123, "time": "3:13am"}`},
			want: foo{
				ID:   intPtr(123),
				Time: &Time{Time: time.Date(0, time.January, 1, 3, 13, 0, 0, time.Local)},
			},
			err: nil,
		},
		{
			name: "null time",
			args: args{`{"id": 123, "time": null}`},
			want: foo{
				ID:   intPtr(123),
				Time: nil,
			},
			err: nil,
		},
		{
			name: "Totally invalid",
			args: args{`{"id": 123, "time": "gibberish"}`},
			err:  TimeParseError,
		},
	}
	for _, tt := range tests {
		var f foo
		if gotErr := json.Unmarshal([]byte(tt.args.jsonStr), &f); gotErr != tt.err {
			t.Errorf("%q. UnmarshalJSON() error = %v, want %v", tt.name, gotErr, tt.err)
		} else if tt.err == nil {
			if f.ID == nil || *f.ID != *tt.want.ID {
				t.Errorf("%q. UnmarshalJSON() = %v, want %v - ID messed up", tt.name, f, tt.want)
			} else if tt.want.Time == nil && f.Time != nil {
				t.Errorf("%q. UnmarshalJSON() = %v, want %v - unexpected time", tt.name, f, tt.want)
			} else if tt.want.Time != nil && (f.Time == nil || !(*tt.want.Time).Equal(*f.Time)) {
				t.Errorf("%q. UnmarshalJSON() = %v, want %v", tt.name, f, tt.want)
			}
		}
	}
}
