package tpl

const (
	DaoHeader = `package %s

import (
	"context"
	"time"
	"errors"
	
	"gorm.io/gorm"
	"%s/%s"
	"%s/%s"
)
`
)

var (
	TmplIDao = `
var _ I{{.StructName.UpperS}}Dao = (*{{.StructName.UpperS}}Dao)(nil)

type I{{.StructName.UpperS}}Dao interface {
	Insert(ctx context.Context, {{.StructName.LowerP}} ...*model.{{.StructName.UpperS}}) (err error)
	Save(ctx context.Context, {{.StructName.LowerS}} *model.{{.StructName.UpperS}}) (err error)
	FindOne(ctx context.Context, condition *opt.{{.StructName.UpperS}}Option) ({{.StructName.LowerS}} *model.{{.StructName.UpperS}}, err error)
	FindList(ctx context.Context, condition *opt.{{.StructName.UpperS}}Option) ({{.StructName.LowerS}}s []*model.{{.StructName.UpperS}}, total int64, err error)
	Count(ctx context.Context, condition *opt.{{.StructName.UpperS}}Option) (count int64, err error)
	Update(ctx context.Context, {{.StructName.LowerS}} *model.{{.StructName.UpperS}}, condition *opt.{{.StructName.UpperS}}Option) (err error)
	Delete(ctx context.Context, condition *opt.{{.StructName.UpperS}}Option) (err error)

	// write you method here
}

`

	TmplDao = `
type {{.StructName.UpperS}}Dao struct {
	db *gorm.DB
}

func New{{.StructName.UpperS}}Dao(db *gorm.DB) *{{.StructName.UpperS}}Dao {
	return &{{.StructName.UpperS}}Dao{db: db}
}

func (m *{{.StructName.UpperS}}Dao) Error(db *gorm.DB) error {
	if db.Error != gorm.ErrRecordNotFound {
		return db.Error
	}
	return nil
}

// TableName Get table name from context, key is TABLE_NAME.
func (m *{{.StructName.UpperS}}Dao) TableName(ctx context.Context) string {
	if ctx != nil {
		val := ctx.Value("TABLE_NAME")
		switch val.(type) {
		case string:
			if tableName := val.(string); tableName != "" {
				return tableName
			}
		}
	}
	return model.{{.StructName.UpperS}}TableName
}
`

	TmplInsert = `
func (m *{{.StructName.UpperS}}Dao) Insert(ctx context.Context, {{.StructName.LowerP}} ...*model.{{.StructName.UpperS}}) (err error) {
	if {{.StructName.LowerP}} == nil {
		return errors.New("insert must include {{.StructName.LowerP}} model")
	}
{{ if and (ne .FieldCreateTime "") (ne .FieldUpdateTime "") }}
	for i := range {{.StructName.LowerP}} {
	{{ if  .IsTimestamp }}  if {{.StructName.LowerP}}[i].{{.FieldCreateTime}} == 0 {
			{{.StructName.LowerP}}[i].{{.FieldCreateTime}} = time.Now().UTC().UnixMilli()
		}
		if {{.StructName.LowerP}}[i].{{.FieldUpdateTime}} == 0 {
			{{.StructName.LowerP}}[i].{{.FieldUpdateTime}} = time.Now().UTC().UnixMilli()
		}
	{{ else }}  if {{.StructName.LowerP}}[i].{{.FieldCreateTime}}.IsZero() {
			{{.StructName.LowerP}}[i].{{.FieldCreateTime}} = time.Now()
		}
		if {{.StructName.LowerP}}[i].{{.FieldUpdateTime}}.IsZero() {
			{{.StructName.LowerP}}[i].{{.FieldUpdateTime}} = time.Now()
		}
	{{ end }}  }
{{ else }}
	{{if ne .FieldCreateTime "" }}
	for i := range {{.StructName.LowerP}} {
	{{ if  .IsTimestamp }}
		if {{.StructName.LowerP}}[i].{{.FieldCreateTime}} == 0 {
			{{.StructName.LowerP}}[i].{{.FieldCreateTime}} = time.Now().UTC().UnixMilli()
		}
	{{ else }}
        if {{.StructName.LowerP}}[i].{{.FieldCreateTime}}.IsZero() {
			{{.StructName.LowerP}}[i].{{.FieldCreateTime}} = time.Now()
		}
	{{ end }}
	}
	{{end}}
	
	{{if ne .FieldUpdateTime "" }}
	for i := range {{.StructName.LowerP}} {
	{{ if  .IsTimestamp }}
		if {{.StructName.LowerP}}[i].{{.FieldUpdateTime}} == 0 {
			{{.StructName.LowerP}}[i].{{.FieldUpdateTime}} = time.Now().UTC().UnixMilli()
		}
	{{ else }}
        if {{.StructName.LowerP}}[i].{{.FieldUpdateTime}}.IsZero() {
			{{.StructName.LowerP}}[i].{{.FieldUpdateTime}} = time.Now()
		}
	{{ end }}
	}
	{{end}}
{{ end }}	db := m.db.WithContext(ctx).Table(m.TableName(ctx)).Create({{.StructName.LowerP}})
	err = m.Error(db)
	return
}
`

	TmplSave = `
func (m *{{.StructName.UpperS}}Dao) Save(ctx context.Context, {{.StructName.LowerS}} *model.{{.StructName.UpperS}}) (err error) {
    if {{.StructName.LowerS}} == nil {
		return errors.New("save must include {{.StructName.LowerS}} model")
	}
{{if ne .FieldCreateTime "" }}
	{{ if  .IsTimestamp }}
	if {{.StructName.LowerS}}.{{.FieldCreateTime}} == 0 {
		{{.StructName.LowerS}}.{{.FieldCreateTime}} = time.Now().UTC().UnixMilli()
	}
	{{ else }}
	if {{.StructName.LowerS}}.{{.FieldCreateTime}}.IsZero() {
		{{.StructName.LowerS}}.{{.FieldCreateTime}} = time.Now()
	}
	{{ end }}
{{end}}
{{if ne .FieldUpdateTime "" }}
	{{ if  .IsTimestamp }}
	if {{.StructName.LowerS}}.{{.FieldUpdateTime}} == 0 {
		{{.StructName.LowerS}}.{{.FieldUpdateTime}} = time.Now().UTC().UnixMilli()
	}
	{{ else }}
    if {{.StructName.LowerS}}.{{.FieldUpdateTime}}.IsZero() {
		{{.StructName.LowerS}}.{{.FieldUpdateTime}} = time.Now()
	}
	{{ end }}
{{end}}
	db := m.db.WithContext(ctx).Table(m.TableName(ctx)).Save({{.StructName.LowerS}})
	err = m.Error(db)
	return
}
`

	TmplFindOne = `
func (m *{{.StructName.UpperS}}Dao) FindOne(ctx context.Context, condition *opt.{{.StructName.UpperS}}Option) ({{.StructName.LowerS}} *model.{{.StructName.UpperS}}, err error) {
	db := m.db.WithContext(ctx)
	if condition != nil {
		db = condition.BuildQuery(db)
	}
{{if ne .FieldSoftDeleteKey "" }}	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db = db.Table(m.TableName(ctx)).First(&{{.StructName.LowerS}})
	err = m.Error(db)
	return
}
`

	TmplFindList = `
func (m *{{.StructName.UpperS}}Dao) FindList(ctx context.Context, condition *opt.{{.StructName.UpperS}}Option) ({{.StructName.LowerP}} []*model.{{.StructName.UpperS}}, total int64, err error) {
	db := m.db.WithContext(ctx).Table(m.TableName(ctx))
	if condition != nil {
		db = condition.BuildQuery(db)
	}
{{if ne .FieldSoftDeleteKey "" }}	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})	{{ end }}
	if err = db.Count(&total).Error; total == 0 {
		return
	}

	if condition != nil {
		db = condition.BuildPage(db)
	}
	db = db.Find(&{{.StructName.LowerP}})
	err = m.Error(db)
	return
}
`

	TmplCount = `
func (m *{{.StructName.UpperS}}Dao) Count(ctx context.Context, condition *opt.{{.StructName.UpperS}}Option) (count int64, err error) {
	db := m.db.WithContext(ctx)
	if condition != nil {
		db = condition.BuildQuery(db)
	}
{{if ne .FieldSoftDeleteKey "" }}	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db = db.Table(m.TableName(ctx)).Count(&count)
	err = m.Error(db)
	return
}
`

	TmplUpdate = `
func (m *{{.StructName.UpperS}}Dao) Update(ctx context.Context, {{.StructName.LowerS}} *model.{{.StructName.UpperS}}, condition *opt.{{.StructName.UpperS}}Option) (err error) {
	if {{.StructName.LowerS}} == nil {
		return errors.New("update must include {{.StructName.LowerS}} model")
	} else if condition == nil {
		return errors.New("update must include where condition")
	}
{{if ne .FieldUpdateTime "" }}
	{{ if  .IsTimestamp }}
	if {{.StructName.LowerS}}.{{.FieldUpdateTime}} == 0 {
		{{.StructName.LowerS}}.{{.FieldUpdateTime}} = time.Now().UTC().UnixMilli()
	}
	{{ else }}
    if {{.StructName.LowerS}}.{{.FieldUpdateTime}}.IsZero() {
		{{.StructName.LowerS}}.{{.FieldUpdateTime}} = time.Now()
	}
	{{ end }}
{{end}}
	db := m.db.WithContext(ctx)
	db = condition.BuildQuery(db)
{{if ne .FieldSoftDeleteKey "" }}	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db = db.Table(m.TableName(ctx)).Updates({{.StructName.LowerS}})
	err = m.Error(db)
	return
}
`

	TmplDelete = `
func (m *{{.StructName.UpperS}}Dao) Delete(ctx context.Context, condition *opt.{{.StructName.UpperS}}Option) (err error) {
	if condition == nil {
		return errors.New("delete must include where condition")
	}

	db := m.db.WithContext(ctx)
	db = condition.BuildQuery(db)
{{if ne .FieldSoftDeleteKey "" }}	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db = db.Table(m.TableName(ctx)).
{{if eq .FieldSoftDeleteKey "" }} Delete(&model.{{.StructName.UpperS}}{})
{{ else }}  {{if eq .FieldUpdateTime "" }}
				Select("{{.TableSoftDeleteKey}}").
				Updates(&model.{{.StructName.UpperS}}{
					{{.FieldSoftDeleteKey}}:{{.TableSoftDeleteValue}},
				})
            {{ else }}
                Select("{{.TableSoftDeleteKey}}","{{.TableUpdateTime}}").
				Updates(&model.{{.StructName.UpperS}}{
					{{.FieldSoftDeleteKey}}:{{.TableSoftDeleteValue}},
				{{ if  .IsTimestamp }}  {{.FieldUpdateTime}} : time.Now().UTC().UnixMilli(),
				{{ else }}  {{.FieldUpdateTime}} : time.Now(),
				{{ end }}
				})
            {{ end }}
{{ end }}	err = m.Error(db)
	return
}
`

	DaoImplTemplate = TmplIDao + TmplDao + TmplInsert + TmplSave + TmplFindOne + TmplFindList + TmplCount + TmplUpdate + TmplDelete
)
