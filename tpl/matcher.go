package tpl

var DBTableInfo []TableInfo

type TableInfo struct {
	Name            string `yaml:"name"`
	GoStruct        string `yaml:"goStruct"`
	CreateTime      string `yaml:"createTime"`
	UpdateTime      string `yaml:"updateTime"`
	IsTimestamp     bool   `yaml:"isTimestamp"`
	SoftDeleteKey   string `yaml:"softDeleteKey"`
	SoftDeleteValue int    `yaml:"softDeleteValue"`
}

type ArgInfo struct {
	Name string
	Type string
}

func GetTableMatcher() map[string]TableInfo {
	var tMatcher = make(map[string]TableInfo)
	for _, matcher := range DBTableInfo {
		tMatcher[matcher.Name] = matcher
	}
	return tMatcher
}
