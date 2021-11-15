package jsonbind

import (
	"bytes"
	"context"

	"github.com/PaesslerAG/gval"
	gjp "github.com/PaesslerAG/jsonpath"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	ojp "github.com/oliveagle/jsonpath"
	"github.com/spyzhov/ajson"
	"github.com/tidwall/gjson"
)

type Getter interface {
	Get(context.Context, interface{}) (interface{}, bool, error)
	New([]byte) (interface{}, error)
}

type MapGetter struct{}

var mapGetter = &MapGetter{}

func NewMapGetter() *MapGetter {
	return mapGetter
}

func (g *MapGetter) New([]byte) (interface{}, error) {
	return nil, nil
}

func (g *MapGetter) Get(context.Context, interface{}) (interface{}, bool, error) {
	return make(map[string]interface{}), false, nil
}

type SliceGetter struct {
	length int
}

func NewSliceGetter(length int) *SliceGetter {
	return &SliceGetter{length: length}
}

func (g *SliceGetter) New([]byte) (interface{}, error) {
	return nil, nil
}

func (g *SliceGetter) Get(context.Context, interface{}) (interface{}, bool, error) {
	return make([]interface{}, g.length), false, nil
}

type NilGetter struct{}

var nilGetter *NilGetter = &NilGetter{}

func NewNilGetter() *NilGetter {
	return nilGetter
}

func (g *NilGetter) New([]byte) (interface{}, error) {
	return nil, nil
}

func (g *NilGetter) Get(context.Context, interface{}) (interface{}, bool, error) {
	return nil, false, nil
}

type OJGGetter struct {
	spec string
	expr jp.Expr
}

func NewOJGGetter(spec string) (*OJGGetter, error) {
	x, e := jp.ParseString(spec)
	if e != nil {
		return nil, e
	}

	return &OJGGetter{spec, x}, nil
}

func (g *OJGGetter) New(b []byte) (interface{}, error) {
	return oj.Parse(b)
}

func (g *OJGGetter) Get(ctx context.Context, obj interface{}) (interface{}, bool, error) {
	result := g.expr.Get(obj)
	binded := false
	var out interface{} = result
	if len(result) > 0 {
		binded = true
		if len(result) == 1 {
			out = result[0]
		}
	}

	return out, binded, nil
}

type GJSONGetter struct {
	spec string
}

func NewGJSONGetter(spec string) *GJSONGetter {
	return &GJSONGetter{spec}
}

func (g *GJSONGetter) New(b []byte) (interface{}, error) {
	parsed := gjson.ParseBytes(b)
	return &parsed, nil
}

func (g *GJSONGetter) Get(ctx context.Context, obj interface{}) (interface{}, bool, error) {
	parsed := obj.(*gjson.Result)
	result := parsed.Get(g.spec)
	binded := false
	if result.Exists() {
		binded = true
	}

	return result.Value(), binded, nil
}

type GvalJSONPathGetter struct {
	spec string
	eval gval.Evaluable
}

func NewGvalJSONPathGetter(spec string) (*GvalJSONPathGetter, error) {
	eval, err := gjp.New(spec)
	if err != nil {
		return nil, err
	}

	return &GvalJSONPathGetter{spec, eval}, nil
}

func (g *GvalJSONPathGetter) New(b []byte) (interface{}, error) {
	return Unmarshal(bytes.NewReader(b))
}

func (g *GvalJSONPathGetter) Get(ctx context.Context, obj interface{}) (interface{}, bool, error) {
	result, err := g.eval(ctx, obj)
	if err != nil {
		return nil, false, err
	}

	return result, true, nil
}

type OJSONPathGetter struct {
	spec string
	pat  *ojp.Compiled
}

func NewOJSONPathGetter(spec string) (*OJSONPathGetter, error) {
	pat, err := ojp.Compile(spec)
	if err != nil {
		return nil, err
	}

	return &OJSONPathGetter{spec, pat}, nil
}

func (g *OJSONPathGetter) New(b []byte) (interface{}, error) {
	return Unmarshal(bytes.NewReader(b))
}

func (g *OJSONPathGetter) Get(ctx context.Context, obj interface{}) (interface{}, bool, error) {
	result, err := g.pat.Lookup(obj)
	if err != nil {
		return nil, false, err
	}

	return result, true, err
}

type AJSONGetter struct {
	spec string
	cmds []string
}

func NewAJSONGetter(spec string) (*AJSONGetter, error) {
	cmds, err := ajson.ParseJSONPath(spec)
	if err != nil {
		return nil, err
	}

	return &AJSONGetter{spec, cmds}, nil
}

func (g *AJSONGetter) New(b []byte) (interface{}, error) {
	return ajson.Unmarshal(b)
}

func (g *AJSONGetter) Get(ctx context.Context, obj interface{}) (interface{}, bool, error) {
	parsed := obj.(*ajson.Node)
	results, err := ajson.ApplyJSONPath(parsed, g.cmds)
	if err != nil {
		return nil, false, err
	}

	var out interface{} = results
	binded := false
	if len(results) > 0 {
		binded = true
		if len(results) == 1 {
			out, err = results[0].Unpack()
			if err != nil {
				return nil, false, err
			}
		} else {
			o := make([]interface{}, len(results))
			for i, r := range results {
				n, e := r.Unpack()
				if e != nil {
					return nil, false, e
				}
				o[i] = n
			}
			out = o
		}
	}

	return out, binded, nil
}

type autoCacheGetter struct {
	CacheName string
	Getter
}

type autoCacheGetterCreator struct {
	cacheName string
	new       NewGetter
}

func (c autoCacheGetterCreator) New(spec string) (*autoCacheGetter, error) {
	g, e := c.new(spec)
	if e != nil {
		return nil, e
	}
	return &autoCacheGetter{c.cacheName, g}, nil
}

type NewGetter func(string) (Getter, error)

var autoCacheGetterCreators = map[string]autoCacheGetterCreator{
	"nil":   {"", func(string) (Getter, error) { return NewNilGetter(), nil }},
	"map":   {"", func(string) (Getter, error) { return NewMapGetter(), nil }},
	"slice": {"", func(spec string) (Getter, error) { return NewSliceGetter(MustAtoi(spec)), nil }},
	"ojg":   {"ojg", func(spec string) (Getter, error) { return NewOJGGetter(spec) }},
	"gjson": {"gjson", func(spec string) (Getter, error) { return NewGJSONGetter(spec), nil }},
	"dft":   {"gjson", func(spec string) (Getter, error) { return NewGJSONGetter(spec), nil }},
	"gval":  {"json", func(spec string) (Getter, error) { return NewGvalJSONPathGetter(spec) }},
	"ojson": {"json", func(spec string) (Getter, error) { return NewOJSONPathGetter(spec) }},
	"ajson": {"ajson", func(spec string) (Getter, error) { return NewAJSONGetter(spec) }},
}
