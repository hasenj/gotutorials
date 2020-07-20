package main

import "os"
import "io"
import "fmt"
import "reflect"

type StructInfo struct {
	Type reflect.Type
	Fields []StructFieldInfo
}

type StructFieldInfo struct {
	Name string
	Type reflect.Type
}

func TypescriptTypeName(t reflect.Type) string {
	return t.Name()
}

func GenerateTypeInfo(inst interface{}) *StructInfo {
	var result = new(StructInfo)
	result.Type = reflect.TypeOf(inst)
	numFields := result.Type.NumField()
	for index := 0; index < numFields; index++ {
		field := result.Type.Field(index)
		// fmt.Printf("%#v\n", field)
		sfield := StructFieldInfo {
			Name: field.Name,
			Type: field.Type,
		}
		result.Fields = append(result.Fields, sfield)
	}
	return result
}

func DescribeStruct(w io.Writer, s *StructInfo) {
	fmt.Fprintf(w, "interface %s { \n", TypescriptTypeName(s.Type))
	for _, field := range s.Fields {
		fmt.Fprintf(w, "    %s: %s;\n", field.Name, TypescriptTypeName(field.Type))
	}
	fmt.Fprintf(w, "}")
}

func main() {
	var inst UserLoginInfo
	var typeInfo = GenerateTypeInfo(inst)
	DescribeStruct(os.Stdout, typeInfo)
}
