package tpl

type StructLevel struct {
	// table -> struct
	TableName string

	StructName StructName

	// table column -> struct field
	Columns []FieldLevel

	// create time
	TableCreateTime string
	FieldCreateTime string

	// update time
	TableUpdateTime string
	FieldUpdateTime string

	// is timestamp
	IsTimestamp bool

	// soft delete
	TableSoftDeleteKey   string
	TableSoftDeleteValue int
	FieldSoftDeleteKey   string
}

// 复数Plural
// 单数Singular
type StructName struct {
	// first letter upper/lower
	UpperS, UpperP string // Order,Orders
	LowerS, LowerP string // order,orders
}

type FieldLevel struct {
	FieldName  string
	FieldType  string
	PrimaryKey string
	// gorm tag for field
	GormName string
	// json tag for field
	JsonName string
	// comment from create table sql
	Comment string
}
