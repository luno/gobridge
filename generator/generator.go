package generator

import (
	"gobridge/ioeasy"
	"gobridge/reader"
	"gobridge/templates"
	"os"
	"strings"
)

func TSClient(tsPath string, rawTypes map[string]map[string]string, fs map[string][]reader.FunctionSignature) error {
	err := ioeasy.CreateFileFromPath(tsPath)
	if err != nil {
		return err
	}

	// Reset file
	err = os.Truncate(tsPath, 0)
	if err != nil {
		return err
	}

	// Apply correct permissions
	file, err := os.OpenFile(tsPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	tsdata := make(map[string]map[string]string)
	for name, fields := range rawTypes {
		tsdata[name] = fields
		f := tsdata[name]
		for key, value := range f {
			f[key] = switchToTypescriptType(value)
		}
	}

	tsi := new(templates.TSService)
	for serviceName, methods := range fs {
		tsi.Name = serviceName
		for _, m := range methods {
			tsi.MethodNames = append(tsi.MethodNames, m.Name)

			req := templates.TSType{
				Name:   "Request" + m.Name,
				Fields: m.Params,
			}
			tsi.Types = append(tsi.Types, req)

			resp := templates.TSType{
				Name:   "Response" + m.Name,
				Fields: m.Results,
			}

			tsi.Types = append(tsi.Types, resp)
		}
	}

	for name, fields := range tsdata {
		tst := templates.TSType{
			Name:   name,
			Fields: fields,
		}

		tsi.Types = append(tsi.Types, tst)
	}

	err = tsi.AddTo(file)
	if err != nil {
		return err
	}

	return nil
}

func GoClient(data map[string]map[string]string) (fContents string, err error) {
	return
}

func Server() (fContents string, err error) {
	return
}

func isBuiltInType(typ string) bool {
	switch typ {
	case "bool", "byte", "complex128", "complex64", "error":
	case "float32", "float64":
	case "int", "int16", "int32", "int64", "int8":
	case "rune", "string":
	case "uint", "uint16", "uint32", "uint64", "uint8", "uintptr":
	default:
		return false
	}
	return true
}

func switchToTypescriptType(typ string) string {
	switch typ {
	case "byte", "complex128", "complex64", "error":
		return "string"
	case "bool":
		return "boolean"
	case "float32", "float64":
		return "number"
	case "int", "int16", "int32", "int64", "int8":
		return "number"
	case "rune", "string":
		return "string"
	case "uint", "uint16", "uint32", "uint64", "uint8", "uintptr":
		return "number"
	case "Time":
		return "Date"
	default:
		// Consider and treat it as a non-primitive type

		// Trade Go slices for TS arrays
		if strings.HasPrefix(typ, "[]") {
			typ = strings.TrimPrefix(typ, "[]")
			typ += "[]"
		}

		// Trade Go maps for TS any as there is no easy swap
		if strings.HasPrefix(typ, "map") {
			typ = "any"
		}
		return typ
	}
}