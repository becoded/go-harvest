package harvest

import (
	"strings"
	"time"
)

type Date struct {
	time.Time
}

func (t Date) String() string {
	return t.Time.Format("2006-01-02")
}

func (t *Date) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	str = strings.Trim(str, "\"")
	const shortForm = "2006-01-02"
	(*t).Time, err = time.Parse(shortForm, str)
	return
}

// Equal reports whether t and u are equal based on time.Equal
func (t Date) Equal(u Date) bool {
	return t.Time.Equal(u.Time)
}
