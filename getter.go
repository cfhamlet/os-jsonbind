package jsonbind

import (
	"context"

	"github.com/PaesslerAG/gval"
	gjp "github.com/PaesslerAG/jsonpath"
	"github.com/ohler55/ojg/jp"
	ojp "github.com/oliveagle/jsonpath"
	"github.com/spyzhov/ajson"
)

type Getter interface {
	Get(context.Context, *Parsed) (interface{}, bool, error)
}

type MapGetter struct{}

var mapGetter = &MapGetter{}

func NewMapGetter() *MapGetter {
	return mapGetter
}

func (g *MapGetter) Get(context.Context, *Parsed) (interface{}, bool, error) {
	return make(map[string]interface{}), false, nil
}

type SliceGetter struct {
	length int
}

func NewSliceGetter(length int) *SliceGetter {
	return &SliceGetter{length: length}
}

func (g *SliceGetter) Get(context.Context, *Parsed) (interface{}, bool, error) {
	return make([]interface{}, g.length), false, nil
}

type NilGetter struct{}

var nilGetter *NilGetter = &NilGetter{}

func NewNilGetter() *NilGetter {
	return nilGetter
}

func (g *NilGetter) Get(context.Context, *Parsed) (interface{}, bool, error) {
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

func (g *OJGGetter) Get(ctx context.Context, p *Parsed) (interface{}, bool, error) {
	err := p.EnsureOJG()
	if err != nil {
		return nil, false, err
	}

	result := g.expr.Get(p.ojgCompiled)
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

func (g *GJSONGetter) Get(ctx context.Context, p *Parsed) (interface{}, bool, error) {
	err := p.EnsureGJSON()
	if err != nil {
		return nil, false, err
	}

	result := p.gjsonCompiled.Get(g.spec)
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

func (g *GvalJSONPathGetter) Get(ctx context.Context, p *Parsed) (interface{}, bool, error) {
	err := p.EnsureJSON()
	if err != nil {
		return nil, false, err
	}

	result, err := g.eval(ctx, p.jsonCompiled)
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

func (g *OJSONPathGetter) Get(ctx context.Context, p *Parsed) (interface{}, bool, error) {
	err := p.EnsureJSON()
	if err != nil {
		return nil, false, err
	}

	result, err := g.pat.Lookup(p.jsonCompiled)
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

func (g *AJSONGetter) Get(ctx context.Context, p *Parsed) (interface{}, bool, error) {
	err := p.EnsureAJSON()
	if err != nil {
		return nil, false, err
	}

	results, err := ajson.ApplyJSONPath(p.ajsonCompiled, g.cmds)
	if err != nil {
		return nil, false, err
	}

	var out interface{} = results
	binded := false
	if len(results) > 0 {
		binded = true
		if len(results) == 1 {
			out, err = results[0].Value()
			if err != nil {
				return nil, false, err
			}
		} else {
			o := make([]interface{}, len(results))
			for i, r := range results {
				n, e := r.Value()
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

type NewGetter func(string) (Getter, error)

var newGetters map[string]NewGetter = map[string]NewGetter{
	"nil":     func(spec string) (Getter, error) { return NewNilGetter(), nil },
	"map":     func(spec string) (Getter, error) { return NewMapGetter(), nil },
	"slice":   func(spec string) (Getter, error) { return NewSliceGetter(MustAtoi(spec)), nil },
	"ojg":     func(spec string) (Getter, error) { return NewOJGGetter(spec) },
	"gjson":   func(spec string) (Getter, error) { return NewGJSONGetter(spec), nil },
	"default": func(spec string) (Getter, error) { return NewGJSONGetter(spec), nil },
	"gval":    func(spec string) (Getter, error) { return NewGvalJSONPathGetter(spec) },
	"ojson":   func(spec string) (Getter, error) { return NewOJSONPathGetter(spec) },
	"ajson":   func(spec string) (Getter, error) { return NewAJSONGetter(spec) },
}
