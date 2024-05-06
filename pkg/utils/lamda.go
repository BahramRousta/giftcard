package utils

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func Map[T1, T2 any](in []T1, f func(T1) T2) []T2 {
	res := make([]T2, len(in))
	for i, item := range in {
		res[i] = f(item)
	}
	return res
}

func Filter[T any](in []T, f func(T) bool) []T {
	res := make([]T, 0, len(in))
	for _, item := range in {
		if f(item) {
			res = append(res, item)
		}
	}
	return res
}

func Reduce[T1, T2 any](in []T1, init T2, reducer func(T2, T1) T2) T2 {
	for _, item := range in {
		init = reducer(init, item)
	}
	return init
}

func Any[T any](in []T, f func(T) bool) bool {
	res := false
	for i := 0; i < len(in); i++ {
		res = res || f(in[i])
	}
	return res
}

func All[T any](in []T, f func(T) bool) bool {
	res := true
	for i := 0; i < len(in); i++ {
		res = res && f(in[i])
	}
	return res
}

func Marshal(v any) string {
	b, err := json.Marshal(v)

	if err != nil {
		return err.Error()
	}

	return string(b)
}

func ToStringData(data interface{}) string {
	switch v := data.(type) {
	case int:
		// v is an int here, so e.g. v + 1 is possible.
		return fmt.Sprintf("%v", v)
	case float64:
		// v is a float64 here, so e.g. v + 1.0 is possible.
		return fmt.Sprintf("%v", v)
	case string:
		// v is a string here, so e.g. v + " Yeah!" is possible.
		return data.(string)
	default:
		// And here I'm feeling dumb. ;)
		res, _ := json.Marshal(data)
		return string(res)
	}
}

func ValueExistsInArray(value string, array []string) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}
	return false
}

func MapToUuid(str []string) []uuid.UUID {
	var ids []uuid.UUID
	for _, s := range str {
		id, err := uuid.Parse(s)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	return ids
}
