package harvest

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

var dateType = reflect.TypeOf(Date{}) //nolint: gochecknoglobals

var dateTimeType = reflect.TypeOf(time.Time{}) //nolint: gochecknoglobals

// Stringify attempts to create a reasonable string representation of types in
// the Harvest library. It does things like resolve pointers to their values
// and omits struct fields with nil values.
func Stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v)

	return buf.String()
}

// stringifyValue was heavily inspired by the goprotobuf library.

func stringifyValue(w io.Writer, val reflect.Value) { //nolint:funlen,gocognit
	if val.Kind() == reflect.Ptr && val.IsNil() {
		if _, err := w.Write([]byte("<nil>")); err != nil {
			logrus.Error(err)
		}

		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() { //nolint:exhaustive
	case reflect.String:
		fmt.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		if _, err := w.Write([]byte{'['}); err != nil {
			logrus.Error(err)
		}

		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				if _, err := w.Write([]byte{' '}); err != nil {
					logrus.Error(err)
				}
			}

			stringifyValue(w, v.Index(i))
		}

		if _, err := w.Write([]byte{']'}); err != nil {
			logrus.Error(err)
		}

		return
	case reflect.Struct:
		if v.Type().Name() != "" {
			if _, err := w.Write([]byte(v.Type().String())); err != nil {
				logrus.Error(err)
			}
		}

		// special handling of date values
		switch v.Type() {
		case dateType:
			fmt.Fprintf(w, "{%s}", v.Interface())

			return
		case dateTimeType:
			fmt.Fprintf(w, "{%s}", v.Interface())

			return
		}

		if _, err := w.Write([]byte{'{'}); err != nil {
			logrus.Error(err)
		}

		var sep bool

		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}

			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}

			if sep {
				if _, err := w.Write([]byte(", ")); err != nil {
					logrus.Error(err)
				}
			} else {
				sep = true
			}

			if _, err := w.Write([]byte(v.Type().Field(i).Name)); err != nil {
				logrus.Error(err)
			}

			if _, err := w.Write([]byte{':'}); err != nil {
				logrus.Error(err)
			}

			stringifyValue(w, fv)
		}

		if _, err := w.Write([]byte{'}'}); err != nil {
			logrus.Error(err)
		}
	default:
		if v.CanInterface() {
			fmt.Fprint(w, v.Interface())
		}
	}
}
