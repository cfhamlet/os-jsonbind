package jsonbind

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var replaceMe = "--replace-me--"

func TestGetInMap(t *testing.T) {
	data := []byte(`
{"a": {"b": 7}}
`)
	tmpl := `
{"properties": {"A": {"properties": {"B": {"bind": "--replace-me--"}}}}}
`
	for _, s := range []string{
		"a.b",
		"ojg:a.b",
		"ojg:$.a.b",
		"gjson:a.b",
		"gval:$.a.b",
		"ojson:$.a.b",
		"ajson:$.a.b",
	} {
		j := strings.Replace(tmpl, replaceMe, s, -1)

		binder, err := Compile([]byte(j))
		require.Nil(t, err)
		require.NotNil(t, binder)
		result, binded, err := binder.Bind(context.Background(), data)
		require.Nil(t, err)
		require.True(t, binded)
		m := result.(map[string]interface{})
		a := m["A"]
		b := a.(map[string]interface{})["B"]
		switch i := b.(type) {
		case json.Number:
			x, _ := i.Int64()
			require.EqualValues(t, 7, x)
		default:
			require.EqualValues(t, 7, b)
		}
	}
}

func TestDeepGet(t *testing.T) {
	data := []byte(`
{"a": {"b": 7}}
`)
	tmpl := `
{"properties": {"B": {"bind": "--replace-me--"}}}
`
	for _, s := range []string{
		"a.b",
		"ojg:a.b",
		"ojg:$.a.b",
		"gjson:a.b",
		"gval:$.a.b",
		"ojson:$.a.b",
		"ajson:$.a.b",
	} {
		j := strings.Replace(tmpl, replaceMe, s, -1)

		binder, err := Compile([]byte(j))
		require.Nil(t, err)
		require.NotNil(t, binder)
		result, binded, err := binder.Bind(context.Background(), data)
		require.Nil(t, err)
		require.True(t, binded)
		m := result.(map[string]interface{})
		a := m["B"]
		switch i := a.(type) {
		case json.Number:
			x, _ := i.Int64()
			require.EqualValues(t, 7, x)
		default:
			require.EqualValues(t, 7, a)
		}
	}
}

func TestSimpleGet(t *testing.T) {
	data := []byte(`
{"a":1, "b": true}
`)
	tmpl := `
{"properties": {"A": {"bind": "--replace-me--"}}}
`

	for _, s := range []string{
		"a",
		"ojg:a",
		"ojg:$.a",
		"gjson:a",
		"gval:$.a",
		"ojson:$.a",
		"ajson:$.a",
	} {
		j := strings.Replace(tmpl, replaceMe, s, -1)

		binder, err := Compile([]byte(j))
		require.Nil(t, err)
		require.NotNil(t, binder)
		result, binded, err := binder.Bind(context.Background(), data)
		require.Nil(t, err)
		require.True(t, binded)
		m := result.(map[string]interface{})
		a := m["A"]
		switch i := a.(type) {
		case json.Number:
			x, _ := i.Int64()
			require.EqualValues(t, 1, x)
		default:
			require.EqualValues(t, 1, a)
		}
	}
}
