package structenv

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structtag"
	"github.com/juju/errors"
	"github.com/mkmik/multierror"
)

var (
	falseMatcher = regexp.MustCompile(`(0|[Ff][Aa][Ll][Ss][Ee]|[Nn][Oo])`)
)

func convError(err error, key, value, kind string) error {
	return errors.Annotatef(err, "unable to convert %s=%s to %s", key, value, kind)
}

func envMap(env []string) map[string]string {
	set := make(map[string]string)
	for _, e := range env {
		elems := strings.SplitN(e, "=", 2)
		if len(elems) == 2 {
			set[elems[0]] = elems[1]
		}
		if len(elems) == 1 {
			set[elems[0]] = ""
		}
	}
	return set
}

// Parse parses the environment variables to assign values in the provided interface.
func Parse(i interface{}) error {
	return ParseEnv(os.Environ(), i)
}

// ParseEnv parses the provided environment variables list in the form key=value
// to assign values in the provided interface.
func ParseEnv(env []string, i interface{}) error {
	envs := envMap(env)

	return walkFields(i, func(n int, f reflect.StructField, v reflect.Value) error {
		tags, err := structtag.Parse(string(f.Tag))
		if err != nil {
			return errors.Trace(err)
		}

		// Parse nested objects
		if f.Type.Kind() == reflect.Struct {
			if err := ParseEnv(env, v.Field(n).Addr().Interface()); err != nil {
				return err
			}
		}

		// Get `env:"K"` tags (ignore errors) and parse env
		if tag, _ := tags.Get("env"); tag != nil {
			key := tag.Name

			if key == "" {
				return nil
			}

			value, ok := envs[key]
			if !ok {
				return nil
			}

			kind := f.Type.String()
			if err := setValue(kind, key, value, v.Field(n)); err != nil {
				return convError(err, key, value, kind)
			}
		}

		return nil
	})
}

func walkFields(i interface{}, fn func(n int, f reflect.StructField, v reflect.Value) error) error {
	v := reflect.ValueOf(i)

	if k := reflect.TypeOf(i).Kind(); k != reflect.Ptr {
		return fmt.Errorf("cannot parse non pointers structs (%v was provided)", k)
	}

	e := v.Elem()
	if k := e.Kind(); k != reflect.Struct {
		return fmt.Errorf("cannot parse non structs (%v was provided)", k)
	}

	// Loop struct
	var errs error
	for ni := 0; ni < e.NumField(); ni++ {
		typeField := e.Type().Field(ni)
		if err := fn(ni, typeField, e); err != nil {
			errs = multierror.Append(errs, errors.Annotatef(err, "unable to parse %s field", typeField.Name))
			continue
		}
	}

	return errs
}

func setValueBool(v reflect.Value, s string) {
	var vv bool
	// Any non-empty value is true unless it matches the falsy matcher
	if s != "" {
		vv = true
	}

	if falseMatcher.MatchString(s) {
		vv = false
	}

	pv := v.Addr().Interface().(*bool)
	*pv = vv
}

func setValueDuration(v reflect.Value, s string) error {
	vv, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	pv := v.Addr().Interface().(*time.Duration)
	*pv = vv
	return nil
}

func setValueFloat64(v reflect.Value, s string) error {
	vv, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	pv := v.Addr().Interface().(*float64)
	*pv = vv
	return nil
}

func setValueFloat32(v reflect.Value, s string) error {
	vv, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}
	pv := v.Addr().Interface().(*float32)
	*pv = float32(vv)
	return nil
}

func setValueInt(v reflect.Value, s string) error {
	vv, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}
	pv := v.Addr().Interface().(*int)
	*pv = int(vv)
	return nil
}

func setValueInt64(v reflect.Value, s string) error {
	vv, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	pv := v.Addr().Interface().(*int64)
	*pv = vv
	return nil
}

func setValueUint(v reflect.Value, s string) error {
	vv, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return err
	}
	pv := v.Addr().Interface().(*uint)
	*pv = uint(vv)
	return nil
}

func setValueUint64(v reflect.Value, s string) error {
	vv, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	pv := v.Addr().Interface().(*uint64)
	*pv = vv
	return nil
}

func setValueString(v reflect.Value, s string) {
	pv := v.Addr().Interface().(*string)
	*pv = s
}

func setValue(kind string, key, s string, v reflect.Value) error {
	switch kind {
	case "bool":
		setValueBool(v, s)
		return nil
	case "time.Duration":
		return setValueDuration(v, s)
	case "float64":
		return setValueFloat64(v, s)
	case "float32":
		return setValueFloat32(v, s)
	case "int":
		return setValueInt(v, s)
	case "int64":
		return setValueInt64(v, s)
	case "uint":
		return setValueUint(v, s)
	case "uint64":
		return setValueUint64(v, s)
	case "string":
		setValueString(v, s)
		return nil
	default:
		return fmt.Errorf("unsupported %s type", kind)
	}
}
