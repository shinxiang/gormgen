package tpl

const (
	DaoInterfaceHeader = `package %s

import (
	"context"

	"%s/%s"
)
`
)

var (
	DaoInterfaceTemplate = `
var _ I{{.StructName.UpperS}}Dao = (*{{.StructName.UpperS}}Dao)(nil)

type I{{.StructName.UpperS}}Dao interface {
	Insert(ctx context.Context, {{.StructName.LowerP}} ...*model.{{.StructName.UpperS}}) (err error)
	Save(ctx context.Context, {{.StructName.LowerS}} *model.{{.StructName.UpperS}}) (err error)
	FindOne(ctx context.Context, condition *model.{{.StructName.UpperS}}Option) ({{.StructName.LowerS}} *model.{{.StructName.UpperS}}, err error)
	FindList(ctx context.Context, condition *model.{{.StructName.UpperS}}Option) ({{.StructName.LowerS}}s []model.{{.StructName.UpperS}}, total int64, err error)
	Count(ctx context.Context, condition *model.{{.StructName.UpperS}}Option) (count int64, err error)
	Update(ctx context.Context, {{.StructName.LowerS}} *model.{{.StructName.UpperS}}, condition *model.{{.StructName.UpperS}}Option) (err error)
	Delete(ctx context.Context, condition *model.{{.StructName.UpperS}}Option) (err error)

	// write you method here
}

`
)
