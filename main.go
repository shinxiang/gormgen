package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/shinxiang/gormgen/tpl"
	"github.com/spf13/viper"
)

var Viper = viper.New()

func main() {
	var configFilePath = flag.String("f", "./", "config file path")
	flag.Parse()

	Viper.SetConfigName("config")        // config file name without file type
	Viper.SetConfigType("yaml")          // config file type
	Viper.AddConfigPath(*configFilePath) // config file path
	if err := Viper.ReadInConfig(); err != nil {
		panic(err)
	} else if err = Viper.UnmarshalKey("database.tables", &tpl.DBTableInfo); err != nil {
		panic(err)
	}

	for _, dbTableInfo := range tpl.DBTableInfo {
		fmt.Printf("table: %s, %+v\n", dbTableInfo.Name, dbTableInfo)
	}
	fmt.Println()

	var (
		dsn         = Viper.GetString("database.dsn")
		goMod       = Viper.GetString("project.go_mod")
		projectPath = Viper.GetString("project.base")
		modelPath   = "model/" // model实例
		daoPath     = "dao/"   // dao接口层
		optPath     = "opt/"   // opt封装层
	)

	if dsn == "" || goMod == "" || projectPath == "" {
		fmt.Println("dsn,goMod,projectPath 为必填参数，请检查")
		os.Exit(1)
	}

	// 创建文件夹（如果已存在会报错，不影响）
	for _, path := range []string{projectPath + daoPath, projectPath + modelPath, projectPath + optPath, projectPath} {
		os.MkdirAll(path, os.ModePerm)
	}

	// Check go mod path
	if strings.HasPrefix(projectPath, "./") {
		p := strings.Trim(projectPath, "./")
		if p != "" {
			goMod += "/" + p
		}
	}

	// 检查文件路径
	if daoPath[len(daoPath)-1] == '/' {
		daoPath = daoPath[:len(daoPath)-1]
	}
	if modelPath[len(modelPath)-1] == '/' {
		modelPath = modelPath[:len(modelPath)-1]
	}
	if optPath[len(optPath)-1] == '/' {
		optPath = optPath[:len(optPath)-1]
	}

	var (
		tables   []string
		tMatcher map[string]tpl.TableInfo
		db       *sql.DB
		err      error
	)
	if dsn != "" {
		// 连接mysql
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		fmt.Println("start to generate gorm structs")

		// 读取数据库中的表
		tables, err = tpl.GetAllTables(db)
		if err != nil {
			fmt.Printf("getAllTables error %+v", err)
			os.Exit(1)
		}
	} else {
		for _, v := range tpl.DBTableInfo {
			tables = append(tables, v.Name)
		}
	}
	tMatcher = tpl.GetTableMatcher()

	for _, table := range tables {
		// 不存在的表直接过滤
		if tMatcher != nil {
			if _, ok := tMatcher[table]; !ok {
				fmt.Printf("table %s ignored\n", table)
				continue
			}
		}

		// 1.生成结构
		structResult, err := tpl.Generate(db, table, tMatcher[table])
		if err != nil {
			fmt.Printf("Generate table %s error %+v\n", table, err)
			os.Exit(1)
		}

		// 2.生成model file
		dirs := strings.Split(modelPath, "/")
		header := fmt.Sprintf(tpl.ModelHeader, dirs[len(dirs)-1])
		err = parseToFile(projectPath+modelPath, tMatcher[table].GoStruct, header, structResult, parseToTmpl, tpl.ModelTemplate)
		if err != nil {
			fmt.Printf("parseToFile error %+v\n", err)
			os.Exit(1)
		}

		// 3.生成model opt file
		dirs = strings.Split(optPath, "/")
		header = fmt.Sprintf(tpl.ModelOptHeader, dirs[len(dirs)-1])
		err = parseToFile(projectPath+optPath, tMatcher[table].GoStruct+".opt", header, structResult, parseToTmpl, tpl.ModelOptTemplate)
		if err != nil {
			fmt.Printf("parseToFile error %+v\n", err)
			os.Exit(1)
		}

		// 4.生成dao file
		dirs = strings.Split(daoPath, "/")
		header = fmt.Sprintf(tpl.DaoHeader, dirs[len(dirs)-1], goMod, modelPath, goMod, optPath)
		err = parseToFile(projectPath+daoPath, tMatcher[table].GoStruct, header, structResult, parseToTmpl, tpl.DaoImplTemplate)
		if err != nil {
			fmt.Printf("parseToFile error %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Generate Table %s Finished\n", table)
	}

}

func parseToFile(filePath string, fileName string, fileHeader string, structResult tpl.StructLevel, parseFunc func(tpl.StructLevel, string) (string, error), text string) error {
	result, err := parseFunc(structResult, text)
	if err != nil {
		return errors.Wrapf(err, "parseToDaoTmpl structResult %v", structResult)
	}
	path := fmt.Sprintf("%s/%s.go", filePath, fileName)
	file, err := os.OpenFile(path, os.O_WRONLY+os.O_CREATE+os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "OpenFile path %s", path)
	}
	defer file.Close()

	_, err = file.WriteString(fileHeader + result)
	if err != nil {
		return errors.Wrap(err, "WriteString to file")
	}

	// go fmt files
	exec.Command("go", "fmt", path).Run()
	return nil
}

func parseToTmpl(structData tpl.StructLevel, text string) (string, error) {
	tmpl, err := template.New("t").Funcs(template.FuncMap{"counter": counter}).Parse(text)
	// tmpl, err := template.New("t").Parse(text)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, structData)
	return buf.String(), nil
}

// counter 是为了去除一个数组最后一个分隔符的问题，如 1,2,3 不填最后的逗号
func counter() func() int {
	i := -1
	return func() int {
		i++
		return i
	}
}
