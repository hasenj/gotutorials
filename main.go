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

func addTypeFields(linker *TypeLinker, sinfo *StructInfo, t reflect.Type) {
	numFields := t.NumField()
	for index := 0; index < numFields; index++ {
		field := t.Field(index)
		// fmt.Printf("%#v\n", field)
		if field.Anonymous {
			addTypeFields(linker, sinfo, field.Type)
			continue
		}
		var sfield StructFieldInfo
		sfield.Name = field.Name
		sfield.Type = field.Type
		sfield.CustomType = field.Tag.Get("ts")

		if sfield.CustomType == "" {
			QueueType(linker, sfield.Type)
		}
		sinfo.Fields = append(sinfo.Fields, sfield)
	}
}

func ProcessType(linker *TypeLinker, t reflect.Type) {
	var sinfo StructInfo
	sinfo.Type = t
	addTypeFields(linker, &sinfo, t)
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
	// see if we need to process another type referenced here directly or indirectly
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array:
		QueueType(linker, t.Elem())
	case reflect.Map:
		QueueType(linker, t.Key())
		QueueType(linker, t.Elem())
	case reflect.Struct:
		if linker.SeenTypes[t] {
			return
		}
		linker.SeenTypes[t] = true
		linker.Queue = append(linker.Queue, t)
	}
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
