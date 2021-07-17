package parser

import "fmt"

type Identifier string
type BoolOp string
type CompareOp string
type Ref string
type None int8
type Boolean bool
type StringObj string
type Number float64

type Type struct {
	Name  string
	TEnum TypeEnumValue

	//  if struct
	Field []*FieldItem

	// if enum
	Enum []*EnumItem

	KeyType     *Type // If map
	ValueType   *Type // If map, list, or set
	Annotations []*Annotation
}

func (t *Type) String() string {
	switch t.Name {
	case "map":
		return fmt.Sprintf("map<%s,%s>", t.KeyType.String(), t.ValueType.String())
	case "list":
		return fmt.Sprintf("list<%s>", t.ValueType.String())
	case "set":
		return fmt.Sprintf("set<%s>", t.ValueType.String())
	}
	return t.Name
}

func (t *Type) GetField(fieldName string) *FieldItem {
	for _, f := range t.Field {
		if f.Name == fieldName {
			return f
		}
	}
	return nil
}

type FieldItem struct {
	Name        string
	Type        *Type
	Annotations []*Annotation
}

type EnumItem struct {
	Name    string
	Value   int
	Comment string
}

type Typedef struct {
	*Type

	Alias       string
	Annotations []*Annotation
}

type EnumValue struct {
	Name        string
	Value       int
	Annotations []*Annotation
	Comment     string
}

type Enum struct {
	Name        string
	Values      map[string]*EnumValue
	Annotations []*Annotation
	Comment     string
}

type Constant struct {
	Name  string
	Type  *Type
	Value interface{}
}

type Field struct {
	ID          int
	Name        string
	Optional    bool
	Type        *Type
	Default     interface{}
	Annotations []*Annotation
	Comment     string
}

type Struct struct {
	Name        string
	Fields      []*Field
	Annotations []*Annotation
	Comment     string
}

type Annotation struct {
	Name  string
	Value string
}

type TypeEnumValue string

const (
	Bool         TypeEnumValue = "bool"
	Byte         TypeEnumValue = "byte"
	I16          TypeEnumValue = "i16"
	I32          TypeEnumValue = "i32"
	I64          TypeEnumValue = "i64"
	Double       TypeEnumValue = "double"
	String       TypeEnumValue = "string"
	Binary       TypeEnumValue = "binary"
	Map          TypeEnumValue = "map"
	Set          TypeEnumValue = "set"
	List         TypeEnumValue = "list"
	CustomEnum   TypeEnumValue = "enum"
	CustomStruct TypeEnumValue = "struct"
	Unknown      TypeEnumValue = "unknown"
)
