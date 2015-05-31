package tango

import (
	"errors"
	"net/http"
	"strconv"
)

type Forms http.Request

var _ Set = &Forms{}

func (f *Forms) String(key string) (string, error) {
	(*http.Request)(f).ParseForm()
	if v, ok := (*http.Request)(f).Form[key]; ok {
		return v[0], nil
	}
	return "", errors.New("not exist")
}

func (f *Forms) Int(key string) (int, error) {
	return strconv.Atoi((*http.Request)(f).FormValue(key))
}

func (f *Forms) Int32(key string) (int32, error) {
	v, err := strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 32)
	return int32(v), err
}

func (f *Forms) Int64(key string) (int64, error) {
	return strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 64)
}

func (f *Forms) Uint(key string) (uint, error) {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
	return uint(v), err
}

func (f *Forms) Uint32(key string) (uint32, error) {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 32)
	return uint32(v), err
}

func (f *Forms) Uint64(key string) (uint64, error) {
	return strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
}

func (f *Forms) Bool(key string) (bool, error) {
	return strconv.ParseBool((*http.Request)(f).FormValue(key))
}

func (f *Forms) Float32(key string) (float32, error) {
	v, err := strconv.ParseFloat((*http.Request)(f).FormValue(key), 64)
	return float32(v), err
}

func (f *Forms) Float64(key string) (float64, error) {
	return strconv.ParseFloat((*http.Request)(f).FormValue(key), 64)
}

func (f *Forms) MustString(key string, defs ...string) string {
	(*http.Request)(f).ParseForm()
	if v, ok := (*http.Request)(f).Form[key]; ok {
		return v[0]
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return ""
}

func (f *Forms) MustInt(key string, defs ...int) int {
	v, err := strconv.Atoi((*http.Request)(f).FormValue(key))
	if len(defs) > 0 && err != nil {
		return defs[0]
	}
	return v
}

func (f *Forms) MustInt32(key string, defs ...int32) int32 {
	v, err := strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 32)
	if len(defs) > 0 && err != nil {
		return defs[0]
	}
	return int32(v)
}

func (f *Forms) MustInt64(key string, defs ...int64) int64 {
	v, err := strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 64)
	if len(defs) > 0 && err != nil {
		return defs[0]
	}
	return v
}

func (f *Forms) MustUint(key string, defs ...uint) uint {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
	if len(defs) > 0 && err != nil {
		return defs[0]
	}
	return uint(v)
}

func (f *Forms) MustUint32(key string, defs ...uint32) uint32 {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 32)
	if len(defs) > 0 && err != nil {
		return defs[0]
	}
	return uint32(v)
}

func (f *Forms) MustUint64(key string, defs ...uint64) uint64 {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
	if len(defs) > 0 && err != nil {
		return defs[0]
	}
	return v
}

func (f *Forms) MustFloat32(key string, defs ...float32) float32 {
	v, err := strconv.ParseFloat((*http.Request)(f).FormValue(key), 32)
	if len(defs) > 0 && err != nil {
		return defs[0]
	}
	return float32(v)
}

func (f *Forms) MustFloat64(key string, defs ...float64) float64 {
	v, err := strconv.ParseFloat((*http.Request)(f).FormValue(key), 64)
	if len(defs) > 0 && err != nil {
		return defs[0]
	}
	return v
}

func (f *Forms) MustBool(key string, defs ...bool) bool {
	v, err := strconv.ParseBool((*http.Request)(f).FormValue(key))
	if len(defs) > 0 && err != nil {
		return defs[0]
	}
	return v
}