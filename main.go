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
	return t.Name()
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

		if sfield.CustomType == "" && sfield.Type.Kind() == reflect.Struct {
			QueueType(linker, sfield.Type)
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
