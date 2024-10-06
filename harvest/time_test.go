package harvest_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestTime_UnmarshalJSONParse(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want harvest.Time
		err  error
	}{
		{
			name: "24-hour clock after 12:00",
			args: args{"15:04"},
			want: harvest.Time{Time: time.Date(0, time.January, 1, 15, 4, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "24-hour clock before 12:00",
			args: args{"07:34"},
			want: harvest.Time{Time: time.Date(0, time.January, 1, 7, 34, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "24-hour clock no leading zero",
			args: args{"7:34"},
			want: harvest.Time{Time: time.Date(0, time.January, 1, 7, 34, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "Kitchen clock upper case",
			args: args{"9:13AM"},
			want: harvest.Time{Time: time.Date(0, time.January, 1, 9, 13, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "Kitchen clock in the morning",
			args: args{"8:13am"},
			want: harvest.Time{Time: time.Date(0, time.January, 1, 8, 13, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "Kitchen clock with leading zero",
			args: args{"09:39am"},
			want: harvest.Time{Time: time.Date(0, time.January, 1, 9, 39, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "Kitchen clock in the afternoon",
			args: args{"10:13pm"},
			want: harvest.Time{Time: time.Date(0, time.January, 1, 22, 13, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "With quotes",
			args: args{"\"10:13pm\""},
			want: harvest.Time{Time: time.Date(0, time.January, 1, 22, 13, 0, 0, time.Local)},
			err:  nil,
		},
		{
			name: "empty",
			args: args{""},
			err:  harvest.ErrTimeParse,
		},
		{
			name: "24-hour with pm",
			args: args{"13:03pm"},
			err:  harvest.ErrTimeParse,
		},
		{
			name: "24-hour with am",
			args: args{"13:03am"},
			err:  harvest.ErrTimeParse,
		},
		{
			name: "Totally invalid",
			args: args{"gibberish"},
			err:  harvest.ErrTimeParse,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tm := harvest.Time{}

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

func TestTime_UnmarshalJSON(t *testing.T) {
	type foo struct {
		ID   *int64        `json:"id"`
		Time *harvest.Time `json:"time"`
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
				ID:   harvest.Int64(123),
				Time: &harvest.Time{Time: time.Date(0, time.January, 1, 3, 13, 0, 0, time.Local)},
			},
			err: nil,
		},
		{
			name: "null time",
			args: args{`{"id": 123, "time": null}`},
			want: foo{
				ID:   harvest.Int64(123),
				Time: nil,
			},
			err: nil,
		},
		{
			name: "Totally invalid",
			args: args{`{"id": 123, "time": "gibberish"}`},
			err:  harvest.ErrTimeParse,
		},
	}

	t.Parallel()

	for _, tt := range tests {
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
			case tt.want.Time == nil && f.Time != nil:
				t.Errorf("%q. UnmarshalJSON() = %v, response %v - unexpected time", tt.name, f, tt.want)
			case tt.want.Time != nil && (f.Time == nil || !tt.want.Time.Equal(*f.Time)):
				t.Errorf("%q. UnmarshalJSON() = %v, response %v", tt.name, f, tt.want)
			}
		})
	}
}
