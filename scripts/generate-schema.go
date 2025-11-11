package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/config"
	"github.com/invopop/jsonschema"
)

func getScriptsDirectory() string {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		return path.Dir(file)
	}
	return ""
}

func main() {
	r := new(jsonschema.Reflector)

	r.FieldNameTag = "yaml"

	// taken from https://github.com/megaease/easeprobe/blob/8a29940850fe335fe91fd085835c039bd6c745a1/conf/conf.go#L154-L170
	// Apache License 2.0

	// The Struct name could be same, but the package name is different
	// This would cause the json schema to be wrong `$ref` to the same name.
	// the following code is to fix this issue by adding the package name to the struct name
	// p.s. this issue has been reported in: https://github.com/invopop/jsonschema/issues/42
	r.Namer = func(t reflect.Type) string {
		name := t.Name()
		if t.Kind() == reflect.Struct {
			v := reflect.New(t)
			vt := v.Elem().Type()
			if vt.PkgPath() != "github.com/TwiN/gatus/v5/config" {
				name = vt.PkgPath() + "/" + vt.Name()
				name = strings.TrimPrefix(name, "github.com/TwiN/gatus/v5/")
				name = strings.ReplaceAll(name, "/", "_")
			}
		}
		return name
	}
	///////////////

	durationType := reflect.TypeOf(time.Duration(0))

	r.Mapper = func(t reflect.Type) *jsonschema.Schema {
		if t == durationType {
			return &jsonschema.Schema{
				Type: "string",
			}
		}
		return nil
	}

	s := r.Reflect(&config.Config{})
	data, _ := json.MarshalIndent(s, "", "  ")

	f, _ := os.Create(path.Join(getScriptsDirectory(), "../.schema", "gatus-config-schema.json"))
	defer f.Close()
	_, _ = f.WriteString(string(data))
	fmt.Println(string(data))
}
