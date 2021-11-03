package jsonbind

import "github.com/santhosh-tekuri/jsonschema/v5"

var bindMeta = jsonschema.MustCompileString("bind.json", `
{
	"properties": {
		"bind": {
			"type": "string"
		}
	}
}
`)

type nilCompiler struct{}

func (nilCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	return nil, nil
}

func Register(c *jsonschema.Compiler) {
	c.RegisterExtension("bind", bindMeta, nilCompiler{})
}
