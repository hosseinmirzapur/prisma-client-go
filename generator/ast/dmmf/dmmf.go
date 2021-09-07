package dmmf

import (
	"github.com/prisma/prisma-client-go/generator/types"
)

// FieldKind describes a scalar, object or enum.
type FieldKind string

// FieldKind values
const (
	FieldKindScalar FieldKind = "scalar"
	FieldKindObject FieldKind = "object"
	FieldKindEnum   FieldKind = "enum"
)

// IncludeInStruct shows whether to include a field in a model struct.
func (v FieldKind) IncludeInStruct() bool {
	return v == FieldKindScalar || v == FieldKindEnum
}

// IsRelation returns whether field is a relation
func (v FieldKind) IsRelation() bool {
	return v == FieldKindObject
}

// DatamodelFieldKind describes a scalar, object or enum.
type DatamodelFieldKind string

// DatamodelFieldKind values
const (
	DatamodelFieldKindScalar   DatamodelFieldKind = "scalar"
	DatamodelFieldKindRelation DatamodelFieldKind = "relation"
	DatamodelFieldKindEnum     DatamodelFieldKind = "enum"
)

// IncludeInStruct shows whether to include a field in a model struct.
func (v DatamodelFieldKind) IncludeInStruct() bool {
	return v == DatamodelFieldKindScalar || v == DatamodelFieldKindEnum
}

// IsRelation returns whether field is a relation
func (v DatamodelFieldKind) IsRelation() bool {
	return v == DatamodelFieldKindRelation
}

// Document describes the root of the AST.
type Document struct {
	Datamodel Datamodel `json:"datamodel"`
	Schema    Schema    `json:"schema"`
	Mappings  Mappings  `json:"mappings"`
}

type Mappings struct {
	ModelOperations []ModelOperation `json:"modelOperations"`
	OtherOperations struct {
		Read  []string `json:"read"`
		Write []string `json:"write"`
	} `json:"otherOperations"`
}

type ModelOperation struct {
	Model      types.String `json:"model"`
	Aggregate  types.String `json:"aggregate"`
	CreateOne  types.String `json:"createOne"`
	DeleteMany types.String `json:"deleteMany"`
	DeleteOne  types.String `json:"deleteOne"`
	FindFirst  types.String `json:"findFirst"`
	FindMany   types.String `json:"findMany"`
	FindUnique types.String `json:"findUnique"`
	GroupBy    types.String `json:"groupBy"`
	UpdateMany types.String `json:"updateMany"`
	UpdateOne  types.String `json:"updateOne"`
	UpsertOne  types.String `json:"upsertOne"`
}

func (m *ModelOperation) Namespace() string {
	return m.Model.GoCase() + "Namespace"
}

// Operator describes a query operator such as NOT, OR, etc.
type Operator struct {
	Name   string
	Action string
}

// Operators returns a list of all query operators such as NOT, OR, etc.
func (Document) Operators() []Operator {
	return []Operator{{
		Name:   "Not",
		Action: "NOT",
	}, {
		Name:   "Or",
		Action: "OR",
	}, {
		Name:   "And",
		Action: "AND",
	}}
}

// Action describes a CRUD operation.
type Action struct {
	// Type describes a query or a mutation
	Type string
	Name types.String
}

// ActionType describes a CRUD operation type.
type ActionType struct {
	Name       types.String
	InnerName  types.String
	List       bool
	ReturnList bool
}

func (Document) Types() []string {
	return []string{"Unique", "Many"}
}

// Variations contains different query capabilities such as Unique, First and Many
func (Document) Variations() []ActionType {
	return []ActionType{{
		Name:      "Unique",
		InnerName: "One",
	}, {
		Name:      "First",
		List:      true,
		InnerName: "One",
	}, {
		Name:       "Many",
		List:       true,
		ReturnList: true,
		InnerName:  "Many",
	}}
}

// Actions returns all possible CRUD operations.
func (Document) Actions() []Action {
	return []Action{{
		"query",
		"Find",
	}, {
		"mutation",
		"Create",
	}, {
		"mutation",
		"Update",
	}, {
		"mutation",
		"Delete",
	}}
}

// Method defines the method for the virtual types method
type Method struct {
	Name   string
	Action string
}

// Type defines the data struct for the virtual types method
type Type struct {
	Name    string
	Methods []Method
}

func (Document) WriteTypes() []Type {
	number := []Method{{
		Name:   "Increment",
		Action: "increment",
	}, {
		Name:   "Decrement",
		Action: "decrement",
	}, {
		Name:   "Multiply",
		Action: "multiply",
	}, {
		Name:   "Divide",
		Action: "divide",
	}}
	return []Type{{
		Name:    "Int",
		Methods: number,
	}, {
		Name:    "Float",
		Methods: number,
	}}
}

// SchemaEnum describes an enumerated type.
type SchemaEnum struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
	// DBName (optional)
	DBName types.String `json:"dBName"`
}

// Enum describes an enumerated type.
type Enum struct {
	Name   types.String `json:"name"`
	Values []EnumValue  `json:"values"`
	// DBName (optional)
	DBName types.String `json:"dBName"`
}

// EnumValue contains detailed information about an enum type.
type EnumValue struct {
	Name types.String `json:"name"`
	// DBName (optional)
	DBName types.String `json:"dBName"`
}

// Datamodel contains all types of the Prisma Datamodel.
type Datamodel struct {
	Models []Model `json:"models"`
	Enums  []Enum  `json:"enums"`
}

type UniqueIndex struct {
	InternalName types.String   `json:"name"`
	Fields       []types.String `json:"fields"`
}

func (m *UniqueIndex) Name() types.String {
	if m.InternalName != "" {
		return m.InternalName
	}
	var name string
	for _, f := range m.Fields {
		name += f.GoCase()
	}
	return types.String(name)
}

func (m *UniqueIndex) ASTName() types.String {
	if m.InternalName != "" {
		return m.InternalName
	}
	return concatFieldsToName(m.Fields)
}

// Model describes a Prisma type model, which usually maps to a database table or collection.
type Model struct {
	// Name describes the singular name of the model.
	Name       types.String `json:"name"`
	IsEmbedded bool         `json:"isEmbedded"`
	// DBName (optional)
	DBName        types.String  `json:"dbName"`
	Fields        []Field       `json:"fields"`
	UniqueIndexes []UniqueIndex `json:"uniqueIndexes"`
	PrimaryKey    PrimaryKey    `json:"primaryKey"`
}

type PrimaryKey struct {
	Name   types.String   `json:"name"`
	Fields []types.String `json:"fields"`
}

func (m Model) Actions() []string {
	return []string{"Set", "Equals"}
}

func (m Model) CompositeIndexes() []UniqueIndex {
	var indexes []UniqueIndex
	indexes = append(indexes, m.UniqueIndexes...)
	if len(m.PrimaryKey.Fields) > 0 {
		indexes = append(indexes, UniqueIndex{
			InternalName: concatFieldsToName(m.PrimaryKey.Fields),
			Fields:       m.PrimaryKey.Fields,
		})
	}
	return indexes
}

// RelationFieldsPlusOne returns all fields plus an empty one, so it's easier to iterate through it in some gotpl files
func (m Model) RelationFieldsPlusOne() []Field {
	var fields []Field
	for _, field := range m.Fields {
		if field.Kind.IsRelation() {
			fields = append(fields, field)
		}
	}
	fields = append(fields, Field{})
	return fields
}

// Field describes properties of a single model field.
type Field struct {
	Kind       FieldKind    `json:"kind"`
	Name       types.String `json:"name"`
	IsRequired bool         `json:"isRequired"`
	IsList     bool         `json:"isList"`
	IsUnique   bool         `json:"isUnique"`
	IsReadOnly bool         `json:"isReadOnly"`
	IsID       bool         `json:"isId"`
	Type       types.Type   `json:"type"`
	// DBName (optional)
	DBName      types.String `json:"dBName"`
	IsGenerated bool         `json:"isGenerated"`
	IsUpdatedAt bool         `json:"isUpdatedAt"`
	// RelationToFields (optional)
	RelationToFields []interface{} `json:"relationToFields"`
	// RelationOnDelete (optional)
	RelationOnDelete types.String `json:"relationOnDelete"`
	// RelationName (optional)
	RelationName types.String `json:"relationName"`
	// HasDefaultValue
	HasDefaultValue bool `json:"hasDefaultValue"`
}

func (f Field) RequiredOnCreate() bool {
	if !f.IsRequired || f.IsUpdatedAt || f.HasDefaultValue || f.IsReadOnly || f.IsList {
		return false
	}

	if f.RelationName != "" && f.IsList {
		return false
	}

	return true
}

// Schema provides the GraphQL/PQL AST.
type Schema struct {
	// RootQueryType (optional)
	RootQueryType types.String `json:"rootQueryType"`
	// RootMutationType (optional)
	RootMutationType  types.String    `json:"rootMutationType"`
	InputObjectTypes  InputObjectType `json:"inputObjectTypes"`
	OutputObjectTypes OutputObject    `json:"outputObjectTypes"`
	EnumTypes         EnumTypes       `json:"enumTypes"`
}

type EnumTypes struct {
	Enums []SchemaEnum `json:"model"`
}

type InputObjectType struct {
	Prisma []CoreType `json:"prisma"`
}

type OutputObject struct {
	Prisma []OutputType `json:"prisma"`
}

// OuterInputType provides the arguments of a given field.
type OuterInputType struct {
	Name       types.String      `json:"name"`
	InputTypes []SchemaInputType `json:"inputTypes"`
	// IsRelationFilter (optional)
	IsRelationFilter bool `json:"isRelationFilter"`
}

// SchemaInputType describes an input type of a given field.
type SchemaInputType struct {
	IsRequired bool         `json:"isRequired"`
	IsList     bool         `json:"isList"`
	Type       types.Type   `json:"type"` // this was declared as ArgType
	Kind       FieldKind    `json:"kind"`
	Namespace  types.String `json:"namespace"`
	// can be "scalar", "inputObjectTypes" or "outputObjectTypes"
	Location string `json:"location"`
}

// OutputType describes a GraphQL/PQL return type.
type OutputType struct {
	Name   types.String  `json:"name"`
	Fields []SchemaField `json:"fields"`
	// IsEmbedded (optional)
	IsEmbedded bool `json:"isEmbedded"`
}

// SchemaField describes the information of an output type field.
type SchemaField struct {
	Name       types.String     `json:"name"`
	OutputType SchemaOutputType `json:"outputType"`
	Args       []OuterInputType `json:"args"`
}

// SchemaOutputType describes an output type of a given field.
type SchemaOutputType struct {
	Type       types.String `json:"type"`
	IsList     bool         `json:"isList"`
	IsRequired bool         `json:"isRequired"`
	Kind       FieldKind    `json:"kind"`
}

// CoreType describes a GraphQL/PQL input type.
type CoreType struct {
	Name types.String `json:"name"`
	// IsWhereType (optional)
	IsWhereType bool `json:"isWhereType"`
	// IsOrderType (optional)
	IsOrderType bool `json:"isOrderType"`
	// AtLeastOne (optional)
	AtLeastOne bool `json:"atLeastOne"`
	// AtMostOne (optional)
	AtMostOne bool             `json:"atMostOne"`
	Fields    []OuterInputType `json:"fields"`
}

func concatFieldsToName(fields []types.String) types.String {
	var name types.String
	for i, f := range fields {
		if i > 0 {
			name += "_"
		}
		name += f
	}
	return name
}