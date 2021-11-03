package jsonbind

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Binder struct {
	node   *BindNode
	schema *jsonschema.Schema
}

func (binder *Binder) Bind(ctx context.Context, b []byte) (interface{}, bool, error) {
	parsed := NewParsed()
	parsed.jsonBytes = b

	o, d, e := binder.node.Bind(ctx, parsed)
	if e != nil || !d {
		return o, d, e
	}
	e = binder.schema.Validate(o)
	return o, d, e
}

func Compile(b []byte) (*Binder, error) {
	schema, err := NewSchema(b)
	if err != nil {
		return nil, err
	}

	raw, err := unmarshal(b)
	if err != nil {
		return nil, err
	}
	path := []string{"$"}
	node, err := compile(raw, path)
	if err != nil {
		return nil, err
	}

	if node == nil {
		node = &BindNode{"$", NewNilGetter(), nil}
	}

	return &Binder{node, schema}, nil
}

func compile(schema interface{}, path []string) (*BindNode, error) {
	pl := len(path)

	switch m := schema.(type) {
	case map[string]interface{}:
		spec, binded := m["bind"]
		name := "default"
		if binded {
			s, ok := spec.(string)
			if !ok {
				return nil, fmt.Errorf("%w %s", ErrCast, "can not cast to string")
			}
			name, spec = SplitSpec(s)
		} else {
			if _, ok := m["properties"]; ok {
				name = "map"
			} else if _, ok := m["items"]; ok {
				name = "slice"
			}
		}

		var next interface{}

		if props, ok := m["properties"]; ok {
			switch kv := props.(type) {
			case map[string]interface{}:
				nodes := map[string]*BindNode{}
				for k, v := range kv {
					path = append(path, k)
					n, e := compile(v, path)
					if e != nil {
						return nil, e
					}
					if n != nil {
						nodes[k] = n
						binded = true
					}
					path = path[0:pl]
				}
				if len(nodes) > 0 {
					next = nodes
				}
			}
		} else if items, ok := m["items"]; ok {
			switch m := items.(type) {
			case []interface{}:
				lm := len(m)
				nodes := make([]interface{}, lm)
				if name == "slice" {
					spec = fmt.Sprintf("%d", lm)
				}
				for i, v := range m {
					path = append(path, fmt.Sprintf("[%d]", i))
					n, e := compile(v, path)
					if e != nil {
						return nil, e
					}
					if n != nil {
						nodes[i] = n
						binded = true
					}
					path = path[0:pl]
				}
				if len(nodes) > 0 {
					next = nodes
				}
			default:
			}
		}

		if next == nil && !binded {
			return nil, nil
		}

		newGetter, ok := newGetters[name]
		if !ok {
			return nil, fmt.Errorf("%w %s", ErrNotSupported, name)
		}
		if spec == nil {
			spec = ""
		}
		getter, err := newGetter(spec.(string))
		if err != nil {
			return nil, err
		}

		return &BindNode{JSONPath(path), getter, next}, nil
	}

	return nil, ErrNotSupported
}

type BindNode struct {
	Path string
	Getter
	Next interface{}
}

func (node *BindNode) wrapError(err error, msg string) error {
	return fmt.Errorf("[%w] %s", BindPathError(node.Path), fmt.Errorf("[%w] %s", err, msg))
}

func (node *BindNode) Bind(ctx context.Context, p *Parsed) (interface{}, bool, error) {
	if node.Next == nil {
		o, d, e := node.Get(ctx, p)
		if e != nil {
			e = node.wrapError(e, "bind error")
		}
		return o, d, e
	}

	bbj, bdd, err := node.Get(ctx, p)
	if err != nil || (bbj == nil && !bdd) {
		if err != nil {
			var _e BindPathError
			if !errors.As(err, &_e) {
				err = node.wrapError(err, "bind error")
			}
		}
		return bbj, bdd, err
	}

	switch next := node.Next.(type) {
	case map[string]*BindNode:
		obj, ok := bbj.(map[string]interface{})
		if !ok {
			return nil, false, node.wrapError(ErrCast, "can not cast to map[string]interface{}")
		}
		for k, v := range next {
			o, d, e := v.Bind(ctx, p)
			if e != nil {
				return nil, false, e
			}
			if d {
				bdd = d
				obj[k] = o
			}
		}
	case []*BindNode:
		obj, ok := bbj.([]interface{})
		if !ok {
			return nil, false, node.wrapError(ErrCast, "can not cast to []interface{}")
		}
		if len(obj) < len(next) {
			return nil, false, node.wrapError(ErrPreallocate, "not enough slice space")
		}
		for i, v := range next {
			o, d, e := v.Bind(ctx, p)
			if e != nil {
				return nil, false, e
			}
			if d {
				bdd = d
				obj[i] = o
			}
		}
	}

	return bbj, bdd, nil
}

func NewSchema(b []byte) (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()
	Register(compiler)

	uri := "schema.json"
	e := compiler.AddResource(uri, bytes.NewReader(b))
	if e != nil {
		return nil, e
	}

	return compiler.Compile(uri)
}

func unmarshal(b []byte) (map[string]interface{}, error) {
	raw, err := Unmarshal(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w %s", ErrCast, "can not cast to map[string]interface{}")
	}

	return m, nil
}
