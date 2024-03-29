/*
 * @Author: Jerry.Yang
 * @Date: 2023-07-17 16:50:14
 * @LastEditors: Jerry.Yang
 * @LastEditTime: 2023-07-17 16:52:59
 * @Description:
 */
package pkg

import (
	"flag"
	"fmt"
	"strings"

	"github.com/yangjerry110/protoc-gen-go/gengo"
	"github.com/yangjerry110/protoc-gen-go/gengogrpc"
	"google.golang.org/protobuf/compiler/protogen"
)

func Action() {

	var (
		flags        flag.FlagSet
		plugins      = flags.String("plugins", "", "list of plugins to enable (supported values: grpc)")
		importPrefix = flags.String("import_prefix", "", "prefix to prepend to import paths")
	)
	importRewriteFunc := func(importPath protogen.GoImportPath) protogen.GoImportPath {
		switch importPath {
		case "context", "fmt", "math":
			return importPath
		}
		if *importPrefix != "" {
			return protogen.GoImportPath(*importPrefix) + importPath
		}
		return importPath
	}
	protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: importRewriteFunc,
	}.Run(func(gen *protogen.Plugin) error {
		grpc := false
		for _, plugin := range strings.Split(*plugins, ",") {
			switch plugin {
			case "grpc":
				grpc = true
			case "":
			default:
				return fmt.Errorf("protoc-gen-go: unknown plugin %q", plugin)
			}
		}
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			g := gengo.GenerateFile(gen, f)
			if grpc {
				gengogrpc.GenerateFileContent(gen, f, g)
			}
		}
		gen.SupportedFeatures = gengo.SupportedFeatures
		return nil
	})
}
