package jsonbind

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeepGet(t *testing.T) {
	data := []byte(`
{"a": {"b": 7}}
`)
	type Proper struct {
		B map[string]string `json:"B"`
	}
	type Shm struct {
		Properties Proper `json:"properties"`
	}
	raw := Shm{Proper{map[string]string{}}}
	for _, s := range []string{
		"a.b",
		"ojg:a.b",
		"ojg:$.a.b",
		"gjson:a.b",
		"gval:$.a.b",
		"ojson:$.a.b",
		"ajson:$.a.b",
	} {
		raw.Properties.B["bind"] = s
		j, _ := json.Marshal(raw)

		binder, err := Compile(j)
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
	type Proper struct {
		A map[string]string `json:"A"`
	}
	type Shm struct {
		Properties Proper `json:"properties"`
	}

	raw := Shm{Proper{map[string]string{}}}
	for _, s := range []string{
		"a",
		"ojg:a",
		"ojg:$.a",
		"gjson:a",
		"gval:$.a",
		"ojson:$.a",
		"ajson:$.a",
	} {
		raw.Properties.A["bind"] = s
		j, _ := json.Marshal(raw)

		binder, err := Compile(j)
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
