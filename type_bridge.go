package main

import "io"
import "fmt"
import "reflect"

type TypeBridge struct {
	Structs   []StructInfo
	Queue     []reflect.Type
	SeenTypes map[reflect.Type]bool
}

func NewTypeBridge() *TypeBridge {
	var bridge = new(TypeBridge)
	bridge.SeenTypes = make(map[reflect.Type]bool)
	return bridge
}

type StructInfo struct {
	Type   reflect.Type
	Fields []StructFieldInfo
}

type StructFieldInfo struct {
	Name       string
	Type       reflect.Type
	CustomType string
}

func TSType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Struct:
		return t.Name()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Ptr:
		return TSType(t.Elem())
	case reflect.Array:
		return TSType(t.Elem()) + "[]"
	case reflect.Slice:
		return TSType(t.Elem()) + "[] | null"
	case reflect.Map:
		return fmt.Sprintf("{ [key: %s]: %s } | null", TSType(t.Key()), TSType(t.Elem()))
	default:
		fmt.Println("WARNING: don't know what to output for type:", t.Name())
		return "unknown"
	}
}

func AddFieldsToStruct(bridge *TypeBridge, sinfo *StructInfo, t reflect.Type) {
	numFields := t.NumField()
	for index := 0; index < numFields; index++ {
		field := t.Field(index)
		// fmt.Printf("%#v\n", field)
		if field.Anonymous {
			AddFieldsToStruct(bridge, sinfo, field.Type)
			continue
		}
		var sfield StructFieldInfo
		sfield.Name = field.Name
		sfield.Type = field.Type
		sfield.CustomType = field.Tag.Get("ts")

		if sfield.CustomType == "" {
			QueueType(bridge, sfield.Type)
		}
		sinfo.Fields = append(sinfo.Fields, sfield)
	}
}

func ProcessType(bridge *TypeBridge, t reflect.Type) {
	var sinfo StructInfo
	sinfo.Type = t
	AddFieldsToStruct(bridge, &sinfo, t)
	bridge.Structs = append(bridge.Structs, sinfo)
}

func DescribeStruct(w io.Writer, s StructInfo) {
	fmt.Fprintf(w, "interface %s { \n", s.Type.Name())
	for _, field := range s.Fields {
		var tstype = field.CustomType
		if tstype == "" {
			tstype = TSType(field.Type)
		}
		fmt.Fprintf(w, "    %s: %s;\n", field.Name, tstype)
	}
	fmt.Fprintf(w, "}\n")
}

func DescribeTypes(bridge *TypeBridge, w io.Writer) {
	for _, sinfo := range bridge.Structs {
		DescribeStruct(w, sinfo)
		fmt.Fprintln(w)
	}
}

func QueueInstance(bridge *TypeBridge, inst interface{}) {
	QueueType(bridge, reflect.TypeOf(inst))
}

func QueueType(bridge *TypeBridge, t reflect.Type) {
	// see if we need to process another type referenced here directly or indirectly
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array:
		QueueType(bridge, t.Elem())
	case reflect.Map:
		QueueType(bridge, t.Key())
		QueueType(bridge, t.Elem())
	case reflect.Struct:
		if bridge.SeenTypes[t] {
			return
		}
		bridge.SeenTypes[t] = true
		bridge.Queue = append(bridge.Queue, t)
	}
}

func Process(bridge *TypeBridge) {
	for len(bridge.Queue) > 0 {
		t := bridge.Queue[0]
		bridge.Queue = bridge.Queue[1:]
		ProcessType(bridge, t)
	}
}
