package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/shinxiang/gormgen/example/opt"
	"log"
	"os"
	"time"

	"github.com/shinxiang/gormgen/example/dao"
	"github.com/shinxiang/gormgen/example/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Viper = viper.New()

func main() {
	var configFilePath = flag.String("c", "./", "config file path")
	flag.Parse()

	Viper.SetConfigName("config")        // config file name without file type
	Viper.SetConfigType("yaml")          // config file type
	Viper.AddConfigPath(*configFilePath) // config file path
	if err := Viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Connect mysql
	dsn := Viper.GetString("database.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalln("connect mysql error:", err)
	}

	TestInsert(db)
	TestSave(db)
	TestFindOne(db)
	TestFindList(db)
	TestFindListAll(db)
	TestCount(db)
	TestUpdate(db)
	TestDelete(db)
}

func TestInsert(db *gorm.DB) {
	var dao dao.IUserDao = dao.NewUserDao(db)

	var user = model.User{
		Username: "admin",
		Nickname: "admin",
		Status:   "Y",
	}
	err := dao.Insert(context.Background(), &user)
	if err != nil {
		fmt.Printf("Insert error %+v", err)
		os.Exit(1)
	}
}

func TestSave(db *gorm.DB) {
	var dao dao.IUserDao = dao.NewUserDao(db)

	user := model.User{
		Id:       1,
		Username: "root",
		Nickname: "root",
		Status:   "N",
	}
	err := dao.Save(context.Background(), &user)
	if err != nil {
		fmt.Printf("Save error %+v", err)
		os.Exit(1)
	}
}

func TestFindOne(db *gorm.DB) {
	var dao dao.IUserDao = dao.NewUserDao(db)

	option := opt.NewUserOption().
		Where("username = ?", "admin").
		OrderBy("id desc")

	data, err := dao.FindOne(context.Background(), option)
	if err != nil {
		fmt.Printf("FindOne error %+v", err)
		os.Exit(1)
	}
	fmt.Println("FindOne:", data)
}

func TestFindList(db *gorm.DB) {
	var dao dao.IUserDao = dao.NewUserDao(db)

	option := opt.NewUserOption().
		SetPage(0, 2).
		Where("id > ? and username like ?", 1, "%"+"admin"+"%")

	list, total, err := dao.FindList(context.Background(), option)
	if err != nil {
		fmt.Printf("FindList error %+v", err)
		os.Exit(1)
	}
	fmt.Println("Total:", total, "Page List:", list)
}

func TestFindListAll(db *gorm.DB) {
	var dao dao.IUserDao = dao.NewUserDao(db)

	list, total, err := dao.FindList(context.Background(), nil)
	if err != nil {
		fmt.Printf("FindListAll error %+v", err)
		os.Exit(1)
	}
	fmt.Println("Total:", total, "List:", list)
}

func TestCount(db *gorm.DB) {
	var dao dao.IUserDao = dao.NewUserDao(db)

	option := opt.NewUserOption()
	option.Where("username = ?", "admin")

	count, err := dao.Count(context.Background(), option)
	if err != nil {
		fmt.Printf("Count error %+v", err)
		os.Exit(1)
	}
	fmt.Println("count =", count)
}

func TestUpdate(db *gorm.DB) {
	var dao dao.IUserDao = dao.NewUserDao(db)

	user := model.User{
		Id:         1,
		Username:   "root",
		Nickname:   "root",
		Status:     "Y",
		CreateTime: time.Now(),
	}

	// GORM when updating with struct it will only update non-zero fields by default.
	// Use option.SelectAll() can solve this problem.
	option := opt.NewUserOption().SelectAll()

	err := dao.Update(context.Background(), &user, option)
	if err != nil {
		fmt.Printf("Update error %+v", err)
		os.Exit(1)
	}
}

func TestDelete(db *gorm.DB) {
	var dao dao.IUserDao = dao.NewUserDao(db)

	option := opt.NewUserOption()
	option.Where("id = ?", 2)

	err := dao.Delete(context.Background(), option)
	if err != nil {
		fmt.Printf("Delete error %+v", err)
		os.Exit(1)
	}
}
