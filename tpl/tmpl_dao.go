package tpl

const (
	DaoImplHeader = `package %s

import (
	"context"
	"time"
	"errors"
	
	"gorm.io/gorm"
	"%s/%s"
)
`
)

var (
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
`

	TmplInsert = `
func (m *{{.StructName.UpperS}}Dao) Insert(ctx context.Context, {{.StructName.LowerP}} ...*model.{{.StructName.UpperS}}) (err error) {
	if {{.StructName.LowerP}} == nil {
		return errors.New("insert must include {{.StructName.LowerP}} model")
	}
{{ if and (ne .FieldCreateTime "") (ne .FieldUpdateTime "") }}
	for i := range {{.StructName.LowerP}} {
        if {{.StructName.LowerP}}[i].{{.FieldCreateTime}}.IsZero() {
			{{.StructName.LowerP}}[i].{{.FieldCreateTime}} = time.Now()
		}
		if {{.StructName.LowerP}}[i].{{.FieldUpdateTime}}.IsZero() {
			{{.StructName.LowerP}}[i].{{.FieldUpdateTime}} = time.Now()
		}
	}
{{ else }}
	{{if ne .FieldCreateTime "" }}
	for i := range {{.StructName.LowerP}} {
        if {{.StructName.LowerP}}[i].{{.FieldCreateTime}}.IsZero() {
			{{.StructName.LowerP}}[i].{{.FieldCreateTime}} = time.Now()
		}
	}
	{{end}}
	
	{{if ne .FieldUpdateTime "" }}
	for i := range {{.StructName.LowerP}} {
        if {{.StructName.LowerP}}[i].{{.FieldUpdateTime}}.IsZero() {
			{{.StructName.LowerP}}[i].{{.FieldUpdateTime}} = time.Now()
		}
	}
	{{end}}
{{ end }}	db := m.db.WithContext(ctx).Create({{.StructName.LowerP}})
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
	if {{.StructName.LowerS}}.{{.FieldCreateTime}}.IsZero() {
		{{.StructName.LowerS}}.{{.FieldCreateTime}} = time.Now()
	}
{{end}}
{{if ne .FieldUpdateTime "" }}
    if {{.StructName.LowerS}}.{{.FieldUpdateTime}}.IsZero() {
		{{.StructName.LowerS}}.{{.FieldUpdateTime}} = time.Now()
	}
{{end}}
	db := m.db.WithContext(ctx).Save({{.StructName.LowerS}})
	err = m.Error(db)
	return
}
`

	TmplFindOne = `
func (m *{{.StructName.UpperS}}Dao) FindOne(ctx context.Context, condition *model.{{.StructName.UpperS}}Option) ({{.StructName.LowerS}} *model.{{.StructName.UpperS}}, err error) {
	db := m.db.WithContext(ctx)
	if condition != nil {
		db = condition.BuildQuery(db)
	}
{{if ne .FieldSoftDeleteKey "" }}	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db = db.First(&{{.StructName.LowerS}})
	err = m.Error(db)
	return
}
`

	TmplFindList = `
func (m *{{.StructName.UpperS}}Dao) FindList(ctx context.Context, condition *model.{{.StructName.UpperS}}Option) ({{.StructName.LowerP}} []model.{{.StructName.UpperS}}, total int64, err error) {
	db := m.db.WithContext(ctx)
	if condition != nil {
		db = condition.BuildQuery(db)
	}
{{if ne .FieldSoftDeleteKey "" }}	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db = db.Table(model.{{.StructName.UpperS}}TableName)
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
func (m *{{.StructName.UpperS}}Dao) Count(ctx context.Context, condition *model.{{.StructName.UpperS}}Option) (count int64, err error) {
	db := m.db.WithContext(ctx)
	if condition != nil {
		db = condition.BuildQuery(db)
	}
{{if ne .FieldSoftDeleteKey "" }}	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db = db.Table(model.{{.StructName.UpperS}}TableName).Count(&count)
	err = m.Error(db)
	return
}
`

	TmplUpdate = `
func (m *{{.StructName.UpperS}}Dao) Update(ctx context.Context, {{.StructName.LowerS}} *model.{{.StructName.UpperS}}, condition *model.{{.StructName.UpperS}}Option) (err error) {
	if {{.StructName.LowerS}} == nil {
		return errors.New("update must include {{.StructName.LowerS}} model")
	} else if condition == nil {
		return errors.New("update must include where condition")
	}
{{if ne .FieldUpdateTime "" }}
    if {{.StructName.LowerS}}.{{.FieldUpdateTime}}.IsZero() {
		{{.StructName.LowerS}}.{{.FieldUpdateTime}} = time.Now()
	}
{{end}}
	db := m.db.WithContext(ctx)
	db = condition.BuildQuery(db)
{{if ne .FieldSoftDeleteKey "" }}	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db = db.Table(model.{{.StructName.UpperS}}TableName).Updates({{.StructName.LowerS}})
	err = m.Error(db)
	return
}
`

	TmplDelete = `
func (m *{{.StructName.UpperS}}Dao) Delete(ctx context.Context, condition *model.{{.StructName.UpperS}}Option) (err error) {
	if condition == nil {
		return errors.New("delete must include where condition")
	}

	db := m.db.WithContext(ctx)
	db = condition.BuildQuery(db)
{{if ne .FieldSoftDeleteKey "" }}	db = db.Where("{{.TableSoftDeleteKey}} != ?", {{.TableSoftDeleteValue}})
{{ end }}
	db = db.Table(model.{{.StructName.UpperS}}TableName).
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
					{{.FieldUpdateTime}} : time.Now(),
				})
            {{ end }}
{{ end }}	err = m.Error(db)
	return
}
`

	DaoImplTemplate = TmplDao + TmplInsert + TmplSave + TmplFindOne + TmplFindList + TmplCount + TmplUpdate + TmplDelete
)
