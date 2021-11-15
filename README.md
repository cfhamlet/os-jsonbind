# os-jsonbind
Extract data with JSON schema

This library extend json schema, you can predefine standard json schema with extend tags, which can be used for extract data from source data(bytes) and create new data for json

## Install

```
go get github.com/cfhamlet/os-jsonbind
```

## Usage

* source data

    ```
    srcData := []byte(`{"k1": 1, "k2": "v2", "k3": [2,3,4], "k4": {"i1": 5, "i2": true}}`) // json bytes
    ```

* define a json schema of the target object with "bind" key, value is path string of the source data


    ```
    spec := []byte(`{"properties": {"K": {"type": "number", "bind": "k1"}}}`) // "k1" gjson style
    ```

* create a binder with the schema

    ```
    import "github.com/cfhamlet/os-jsonbind"
    
    binder, err := jsonbind.Compile(spec)
    ```
* extract data from source data, return target result, binded flag and error

    ```
    result, binded, err := binder.Bind(context.Background(), srcData)
    // result is map[string]interface{}{"K": 1}
    // binded is true
    // err is null
    ```

### Rule Syntax

The default rule type is [gjson](https://github.com/tidwall/gjson), you can specify other syntax with rule name as prefix

```
"bind": "ajson:$.k1"
```

Supported syntax rule types:
* ojg: https://github.com/oliveagle/jsonpath
* gval: https://github.com/PaesslerAG/gval
* gjson: https://github.com/tidwall/gjson
* ojson: https://github.com/oliveagle/jsonpath
* ajson: https://github.com/spyzhov/ajson



## License

MIT licensed.
