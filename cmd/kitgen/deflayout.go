package main

import "path/filepath"

type deflayout struct {
	targetDir string
}

func (l deflayout) packagePath(sub string) string {
	return filepath.Join(l.targetDir, sub)
}

func (l deflayout) transformAST(ctx *sourceContext) (files, error) {
	out := make(outputTree)

	endpoints := out.addFile("endpoints/endpoints.go", "endpoints")
	http := out.addFile("http/http.go", "http")
	service := out.addFile("service/service.go", "service")

	addImports(endpoints, ctx)
	addImports(http, ctx)
	addImports(service, ctx)

	for _, typ := range ctx.types {
		addType(service, typ)
	}

	for _, iface := range ctx.interfaces { //only one...
		addStubStruct(service, iface)

		for _, meth := range iface.methods {
			addMethod(service, iface, meth)
			addRequestStruct(endpoints, meth)
			addResponseStruct(endpoints, meth)
			addEndpointMaker(endpoints, iface, meth)
		}

		addEndpointsStruct(endpoints, iface)
		addHTTPHandler(http, iface)

		for _, meth := range iface.methods {
			addDecoder(http, meth)
			addEncoder(http, meth)
		}

		for _, file := range out {
			selectify(file, "service", iface.stubName().Name, l.packagePath("service"))
			for _, meth := range iface.methods {
				selectify(file, "endpoints", meth.requestStructName().Name, l.packagePath("endpoints"))
			}
		}
	}

	for _, file := range out {
		selectify(file, "endpoints", "Endpoints", l.packagePath("endpoints"))
		for _, typ := range ctx.types {
			selectify(file, "service", typ.Name.Name, l.packagePath("service"))
		}
	}

	return formatNodes(out)
}
