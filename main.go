package main

import "os"
import "io"
import "fmt"
import "reflect"

type TypeLinker struct {
	Structs []StructInfo
	Queue []reflect.Type
	SeenTypes map[reflect.Type]bool
}

func NewTypeLinker() *TypeLinker {
	var linker = new(TypeLinker)
	linker.SeenTypes = make(map[reflect.Type]bool)
	return linker
}

type StructInfo struct {
	Type reflect.Type
	Fields []StructFieldInfo
}

type StructFieldInfo struct {
	Name string
	Type reflect.Type
	CustomType string
}

func TypescriptTypeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Struct:
		return t.Name()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "number"
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Ptr:
		return TypescriptTypeName(t.Elem())
	case reflect.Slice, reflect.Array:
		return TypescriptTypeName(t.Elem()) + "[]"
	case reflect.Map:
		return fmt.Sprintf("{ [key: %s]: %s }", TypescriptTypeName(t.Key()), TypescriptTypeName(t.Elem()))
	default:
		fmt.Println("WARNING: don't know what to output for type:", t.Name())
		return "unknown"
	}
}

func ProcessType(linker *TypeLinker, t reflect.Type) {
	var sinfo StructInfo
	sinfo.Type = t
	numFields := sinfo.Type.NumField()
	for index := 0; index < numFields; index++ {
		field := sinfo.Type.Field(index)
		// fmt.Printf("%#v\n", field)
		var sfield StructFieldInfo
		sfield.Name = field.Name
		sfield.Type = field.Type
		sfield.CustomType = field.Tag.Get("ts")

		if sfield.CustomType == "" {
			// see if we need to process another type referenced here directly or indirectly
			// FIXME: maybe move this logic to QueueType
			switch sfield.Type.Kind() {
			case reflect.Struct:
				QueueType(linker, sfield.Type)
			case reflect.Ptr, reflect.Slice, reflect.Array:
				QueueType(linker, sfield.Type.Elem())
			case reflect.Map:
				QueueType(linker, sfield.Type.Key())
				QueueType(linker, sfield.Type.Elem())
			}
		}
		sinfo.Fields = append(sinfo.Fields, sfield)
	}
	linker.Structs = append(linker.Structs, sinfo)
}

func DescribeStruct(w io.Writer, s StructInfo) {
	fmt.Fprintf(w, "interface %s { \n", TypescriptTypeName(s.Type))
	for _, field := range s.Fields {
		var tstype = field.CustomType
		if tstype == "" {
			tstype = TypescriptTypeName(field.Type)
		}
		fmt.Fprintf(w, "    %s: %s;\n", field.Name, tstype)
	}
	fmt.Fprintf(w, "}")
}

func DescribeTypes(w io.Writer, linker *TypeLinker) {
	for _, sinfo := range linker.Structs {
		DescribeStruct(w, sinfo)
		fmt.Fprintln(w, "\n")
	}
}

func QueueInstance(linker *TypeLinker, inst interface{}) {
	QueueType(linker, reflect.TypeOf(inst))
}

func QueueType(linker *TypeLinker, t reflect.Type) {
	if t.Kind() != reflect.Struct {
		return
	}
	if linker.SeenTypes[t] {
		return
	}
	linker.SeenTypes[t] = true
	linker.Queue = append(linker.Queue, t)
}

func Process(linker *TypeLinker) {
	for len(linker.Queue) > 0 {
		t := linker.Queue[0]
		linker.Queue = linker.Queue[1:]
		ProcessType(linker, t)
	}
}

func main() {
	var linker = NewTypeLinker()
	var inst UserProfile
	QueueInstance(linker, inst)
	Process(linker)
	DescribeTypes(os.Stdout, linker)
}
