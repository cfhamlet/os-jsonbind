package jsonbind

import (
	"bytes"

	"github.com/ohler55/ojg/oj"
	"github.com/spyzhov/ajson"
	"github.com/tidwall/gjson"
)

type Parsed struct {
	jsonBytes     []byte
	jsonCompiled  interface{}
	ojgCompiled   interface{}
	gjsonCompiled *gjson.Result
	ajsonCompiled *ajson.Node
}

func NewParsed() *Parsed {
	return &Parsed{}
}

func (p *Parsed) EnsureGJSON() error {
	if p.gjsonCompiled == nil {
		parsed := gjson.ParseBytes(p.jsonBytes)
		p.gjsonCompiled = &parsed
	}

	return nil
}

func (p *Parsed) EnsureJSON() error {
	if p.jsonCompiled == nil {
		parsed, err := Unmarshal(bytes.NewReader(p.jsonBytes))
		if err != nil {
			return err
		}
		p.jsonCompiled = parsed
	}

	return nil
}

func (p *Parsed) EnsureOJG() error {
	if p.ojgCompiled == nil {
		parsed, err := oj.Parse(p.jsonBytes)
		if err != nil {
			return err
		}
		p.ojgCompiled = parsed
	}

	return nil
}

func (p *Parsed) EnsureAJSON() error {
	if p.ajsonCompiled == nil {
		parsed, err := ajson.Unmarshal(p.jsonBytes)
		if err != nil {
			return err
		}
		p.ajsonCompiled = parsed
	}

	return nil
}
