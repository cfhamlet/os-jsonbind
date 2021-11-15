package jsonbind

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Unmarshal(r io.Reader) (interface{}, error) {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	var doc interface{}
	if err := decoder.Decode(&doc); err != nil {
		return nil, err
	}
	if t, _ := decoder.Token(); t != nil {
		return nil, fmt.Errorf("invalid character %v after top-level value", t)
	}

	return doc, nil
}

func MustAtoi(i string) int {
	o, e := strconv.Atoi(i)
	if e != nil {
		panic(e)
	}

	return o
}

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func SplitSpec(s string) (string, string) {
	idx := strings.Index(s[0:minInt(8, len(s))], ":")
	if idx <= 0 {
		return "dft", s
	}
	k := s[0:idx]
	if _, ok := autoCacheGetterCreators[k]; ok {
		return k, s[idx+1:]
	}

	return "dft", s
}

func JSONPath(path []string) string {
	if len(path) <= 0 {
		return ""
	}
	var builder strings.Builder
	for i, s := range path {
		if !strings.HasPrefix(s, "[") {
			if i != 0 {
				builder.WriteByte('.')
			}
		}
		builder.WriteString(s)
	}

	return builder.String()
}
